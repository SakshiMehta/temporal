package tests

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	activitypb "go.temporal.io/api/activity/v1"
	commonpb "go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	sdkclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.temporal.io/server/common/testing/testvars"
	"go.temporal.io/server/tests/testcore"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type ActivityApiUpdateClientTestSuite struct {
	testcore.FunctionalTestSdkSuite
	tv *testvars.TestVars
}

func TestActivityApiUpdateClientTestSuite(t *testing.T) {
	s := new(ActivityApiUpdateClientTestSuite)
	suite.Run(t, s)
}

func (s *ActivityApiUpdateClientTestSuite) SetupTest() {
	s.FunctionalTestSdkSuite.SetupTest()
	s.tv = testvars.New(s.T()).WithTaskQueue(s.TaskQueue()).WithNamespaceName(s.Namespace())
}

type (
	ActivityFunctions func() (string, error)
	WorkflowFunction  func(context2 workflow.Context) error
)

func (s *ActivityApiUpdateClientTestSuite) makeWorkflowFunc(
	activityFunction ActivityFunctions,
	scheduleToCloseTimeout time.Duration,
	initialRetryInterval time.Duration,
) WorkflowFunction {
	return func(ctx workflow.Context) error {

		activityRetryPolicy := &temporal.RetryPolicy{
			InitialInterval:    initialRetryInterval,
			BackoffCoefficient: 1,
		}

		var ret string
		err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			ActivityID:             "activity-id",
			DisableEagerExecution:  true,
			ScheduleToCloseTimeout: scheduleToCloseTimeout,
			StartToCloseTimeout:    scheduleToCloseTimeout,
			RetryPolicy:            activityRetryPolicy,
		}), activityFunction).Get(ctx, &ret)
		return err
	}
}

func (s *ActivityApiUpdateClientTestSuite) TestActivityUpdateApi_ChangeRetryInterval() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	activityUpdated := make(chan struct{})

	var startedActivityCount atomic.Int32
	activityFunction := func() (string, error) {
		startedActivityCount.Add(1)
		if startedActivityCount.Load() == 1 {
			activityErr := errors.New("bad-luck-please-retry")

			return "", activityErr
		}

		s.WaitForChannel(ctx, activityUpdated)
		return "done!", nil
	}

	scheduleToCloseTimeout := 30 * time.Minute
	retryTimeout := 10 * time.Minute
	workflowFn := s.makeWorkflowFunc(activityFunction, scheduleToCloseTimeout, retryTimeout)

	s.Worker().RegisterWorkflow(workflowFn)
	s.Worker().RegisterActivity(activityFunction)

	workflowOptions := sdkclient.StartWorkflowOptions{
		ID:        s.tv.WorkflowID(),
		TaskQueue: s.TaskQueue(),
	}

	workflowRun, err := s.SdkClient().ExecuteWorkflow(ctx, workflowOptions, workflowFn)
	s.NoError(err)

	s.EventuallyWithT(func(t *assert.CollectT) {
		description, err := s.SdkClient().DescribeWorkflowExecution(ctx, workflowRun.GetID(), workflowRun.GetRunID())
		assert.NoError(t, err)
		if description.GetPendingActivities() != nil {
			assert.Len(t, description.GetPendingActivities(), 1)
			assert.Equal(t, int32(1), startedActivityCount.Load())
		}
	}, 10*time.Second, 500*time.Millisecond)

	updateRequest := &workflowservice.UpdateActivityOptionsRequest{
		Namespace: s.Namespace().String(),
		Execution: &commonpb.WorkflowExecution{
			WorkflowId: workflowRun.GetID(),
		},
		Activity: &workflowservice.UpdateActivityOptionsRequest_Id{Id: "activity-id"},
		ActivityOptions: &activitypb.ActivityOptions{
			RetryPolicy: &commonpb.RetryPolicy{
				InitialInterval: durationpb.New(1 * time.Second),
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"retry_policy.initial_interval"}},
	}
	resp, err := s.FrontendClient().UpdateActivityOptions(ctx, updateRequest)
	s.NoError(err)
	s.NotNil(resp)

	description, err := s.SdkClient().DescribeWorkflowExecution(ctx, workflowRun.GetID(), workflowRun.GetRunID())
	s.NoError(err)
	s.Equal(1, len(description.PendingActivities))

	activityUpdated <- struct{}{}

	s.EventuallyWithT(func(t *assert.CollectT) {
		description, err = s.SdkClient().DescribeWorkflowExecution(ctx, workflowRun.GetID(), workflowRun.GetRunID())
		assert.NoError(t, err)
		if description.GetPendingActivities() != nil {
			assert.Len(t, description.GetPendingActivities(), 0)
			assert.Equal(t, int32(2), startedActivityCount.Load())
		}
	}, 3*time.Second, 100*time.Millisecond)

	var out string
	err = workflowRun.Get(ctx, &out)

	s.NoError(err)
}

func (s *ActivityApiUpdateClientTestSuite) TestActivityUpdateApi_ChangeScheduleToClose() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var startedActivityCount atomic.Int32
	activityFunction := func() (string, error) {
		startedActivityCount.Add(1)
		if startedActivityCount.Load() == 1 {
			activityErr := errors.New("bad-luck-please-retry")
			return "", activityErr
		}
		return "done!", nil
	}

	scheduleToCloseTimeout := 30 * time.Minute
	retryTimeout := 10 * time.Minute

	workflowFn := s.makeWorkflowFunc(activityFunction, scheduleToCloseTimeout, retryTimeout)

	s.Worker().RegisterWorkflow(workflowFn)
	s.Worker().RegisterActivity(activityFunction)

	workflowOptions := sdkclient.StartWorkflowOptions{
		ID:        s.tv.WorkflowID(),
		TaskQueue: s.TaskQueue(),
	}

	workflowRun, err := s.SdkClient().ExecuteWorkflow(ctx, workflowOptions, workflowFn)
	s.NoError(err)

	// wait for activity to start (and fail)
	s.EventuallyWithT(func(t *assert.CollectT) {
		description, err := s.SdkClient().DescribeWorkflowExecution(ctx, workflowRun.GetID(), workflowRun.GetRunID())
		assert.NoError(t, err)
		if description.GetPendingActivities() != nil {
			assert.Len(t, description.GetPendingActivities(), 1)
			assert.Equal(t, int32(1), startedActivityCount.Load())
		}

	}, 2*time.Second, 200*time.Millisecond)

	// update schedule_to_close_timeout
	updateRequest := &workflowservice.UpdateActivityOptionsRequest{
		Namespace: s.Namespace().String(),
		Execution: &commonpb.WorkflowExecution{
			WorkflowId: workflowRun.GetID(),
		},
		Activity: &workflowservice.UpdateActivityOptionsRequest_Id{Id: "activity-id"},
		ActivityOptions: &activitypb.ActivityOptions{
			ScheduleToCloseTimeout: durationpb.New(1 * time.Second),
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"schedule_to_close_timeout"}},
	}
	resp, err := s.FrontendClient().UpdateActivityOptions(ctx, updateRequest)
	s.NoError(err)
	s.NotNil(resp)

	// activity should fail immediately
	s.EventuallyWithT(func(t *assert.CollectT) {
		description, err := s.SdkClient().DescribeWorkflowExecution(ctx, workflowRun.GetID(), workflowRun.GetRunID())
		assert.NoError(t, err)
		if description.GetPendingActivities() != nil {
			assert.Len(t, description.GetPendingActivities(), 0)
			assert.Equal(t, int32(1), startedActivityCount.Load())
		}
	}, 2*time.Second, 200*time.Millisecond)

	var out string
	err = workflowRun.Get(ctx, &out)
	var activityError *temporal.ActivityError
	s.True(errors.As(err, &activityError))
	s.Equal(enumspb.RETRY_STATE_NON_RETRYABLE_FAILURE, activityError.RetryState())
	var timeoutError *temporal.TimeoutError
	s.True(errors.As(activityError.Unwrap(), &timeoutError))
	s.Equal(enumspb.TIMEOUT_TYPE_SCHEDULE_TO_CLOSE, timeoutError.TimeoutType())
	s.Equal(int32(1), startedActivityCount.Load())
}

func (s *ActivityApiUpdateClientTestSuite) TestActivityUpdateApi_ChangeScheduleToCloseAndRetry() {
	// change both schedule to close and retry policy
	// initial values are chosen in such a way that activity will fail due to schedule to close timeout
	// we change schedule to close to a longer value and retry policy to a shorter value
	// after that activity should succeed
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var startedActivityCount atomic.Int32
	activityFunction := func() (string, error) {
		startedActivityCount.Add(1)
		if startedActivityCount.Load() == 1 {
			activityErr := errors.New("bad-luck-please-retry")

			return "", activityErr
		}
		return "done!", nil
	}

	// make scheduleToClose shorter than retry 2nd retry interval
	scheduleToCloseTimeout := 8 * time.Second
	retryInterval := 5 * time.Second

	workflowFn := s.makeWorkflowFunc(
		activityFunction, scheduleToCloseTimeout, retryInterval)

	s.Worker().RegisterWorkflow(workflowFn)
	s.Worker().RegisterActivity(activityFunction)

	workflowOptions := sdkclient.StartWorkflowOptions{
		ID:        s.tv.WorkflowID(),
		TaskQueue: s.TaskQueue(),
	}

	workflowRun, err := s.SdkClient().ExecuteWorkflow(ctx, workflowOptions, workflowFn)
	s.NoError(err)

	// wait for activity to start (and fail)
	s.EventuallyWithT(func(t *assert.CollectT) {
		assert.True(t, startedActivityCount.Load() > 0)
	}, 2*time.Second, 200*time.Millisecond)

	// update schedule_to_close_timeout, make it longer
	// also update retry policy interval, make it shorter
	newScheduleToCloseTimeout := 10 * time.Second
	updateRequest := &workflowservice.UpdateActivityOptionsRequest{
		Namespace: s.Namespace().String(),
		Execution: &commonpb.WorkflowExecution{
			WorkflowId: workflowRun.GetID(),
		},
		Activity: &workflowservice.UpdateActivityOptionsRequest_Id{Id: "activity-id"},
		ActivityOptions: &activitypb.ActivityOptions{
			ScheduleToCloseTimeout: durationpb.New(newScheduleToCloseTimeout),
			RetryPolicy: &commonpb.RetryPolicy{
				InitialInterval: durationpb.New(1 * time.Second),
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"schedule_to_close_timeout", "retry_policy.initial_interval"}},
	}

	resp, err := s.FrontendClient().UpdateActivityOptions(ctx, updateRequest)
	s.NoError(err)
	s.NotNil(resp)
	// check that the update was successful
	s.Equal(int64(newScheduleToCloseTimeout.Seconds()), resp.GetActivityOptions().ScheduleToCloseTimeout.GetSeconds())
	// check that field we didn't update is the same
	s.Equal(int64(scheduleToCloseTimeout.Seconds()), resp.GetActivityOptions().StartToCloseTimeout.GetSeconds())

	// now activity should succeed
	s.EventuallyWithT(func(t *assert.CollectT) {
		description, err := s.SdkClient().DescribeWorkflowExecution(ctx, workflowRun.GetID(), workflowRun.GetRunID())
		assert.NoError(t, err)
		if description.GetPendingActivities() != nil {
			assert.Len(t, description.GetPendingActivities(), 0)
		}
		assert.Equal(t, int32(2), startedActivityCount.Load())
	}, 5*time.Second, 200*time.Millisecond)

	var out string
	err = workflowRun.Get(ctx, &out)
	s.NoError(err)
}
