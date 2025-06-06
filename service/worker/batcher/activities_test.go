package batcher

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"
	"unicode"

	"github.com/stretchr/testify/suite"
	commonpb "go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	historypb "go.temporal.io/api/history/v1"
	workflowpb "go.temporal.io/api/workflow/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/primitives/timestamp"
	"go.temporal.io/server/common/testing/mockapi/workflowservicemock/v1"
	"go.uber.org/mock/gomock"
)

type activitiesSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	controller *gomock.Controller

	mockFrontendClient *workflowservicemock.MockWorkflowServiceClient
}

func (s *activitiesSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())

	s.mockFrontendClient = workflowservicemock.NewMockWorkflowServiceClient(s.controller)
}

func TestActivitiesSuite(t *testing.T) {
	suite.Run(t, new(activitiesSuite))
}

const NumTotalEvents = 10

// Pattern contains either c or f representing completed or failed task.
// Schedule events for each task has id of NumTotalEvents*i + 1 where i is the index of the character
// EventId for each task has id of NumTotalEvents*i+NumTotalEvents where i is the index of the character
func generateEventHistory(pattern string) *historypb.History {
	events := make([]*historypb.HistoryEvent, 0)
	for i, char := range pattern {
		// add a Schedule event independent of type of event
		scheduledEventId := int64(NumTotalEvents*i + 1)
		scheduledEvent := historypb.HistoryEvent{EventId: scheduledEventId, EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_SCHEDULED}
		events = append(events, &scheduledEvent)

		event := historypb.HistoryEvent{EventId: int64(NumTotalEvents*i + NumTotalEvents)}
		switch unicode.ToLower(char) {
		case 'c':
			event.EventType = enumspb.EVENT_TYPE_WORKFLOW_TASK_COMPLETED
			event.Attributes = &historypb.HistoryEvent_WorkflowTaskCompletedEventAttributes{
				WorkflowTaskCompletedEventAttributes: &historypb.WorkflowTaskCompletedEventAttributes{ScheduledEventId: scheduledEventId},
			}
		case 'f':
			event.EventType = enumspb.EVENT_TYPE_WORKFLOW_TASK_FAILED
		}
		events = append(events, &event)
	}

	return &historypb.History{Events: events}
}

func (s *activitiesSuite) TestGetLastWorkflowTaskEventID() {
	namespaceStr := "test-namespace"
	tests := []struct {
		name                    string
		history                 *historypb.History
		wantWorkflowTaskEventID int64
		wantErr                 bool
	}{
		{
			name:                    "Test history with all completed task event history",
			history:                 generateEventHistory("ccccc"),
			wantWorkflowTaskEventID: NumTotalEvents*4 + NumTotalEvents,
		},
		{
			name:                    "Test history with last task failing",
			history:                 generateEventHistory("ccccf"),
			wantWorkflowTaskEventID: NumTotalEvents*3 + NumTotalEvents,
		},
		{
			name:                    "Test history with all tasks failing",
			history:                 generateEventHistory("fffff"),
			wantWorkflowTaskEventID: 2,
		},
		{
			name:                    "Test history with some tasks failing in the middle",
			history:                 generateEventHistory("cfffc"),
			wantWorkflowTaskEventID: NumTotalEvents*4 + NumTotalEvents,
		},
		{
			name:    "Test history with empty history should error",
			history: generateEventHistory(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			slices.Reverse(tt.history.Events)
			workflowExecution := &commonpb.WorkflowExecution{}
			s.mockFrontendClient.EXPECT().GetWorkflowExecutionHistoryReverse(ctx, gomock.Any()).Return(
				&workflowservice.GetWorkflowExecutionHistoryReverseResponse{History: tt.history, NextPageToken: nil}, nil)
			gotWorkflowTaskEventID, err := getLastWorkflowTaskEventID(ctx, namespaceStr, workflowExecution, s.mockFrontendClient, log.NewTestLogger())
			s.Equal(tt.wantErr, err != nil)
			s.Equal(tt.wantWorkflowTaskEventID, gotWorkflowTaskEventID)
		})
	}
}

func (s *activitiesSuite) TestGetFirstWorkflowTaskEventID() {
	namespaceStr := "test-namespace"
	workflowExecution := commonpb.WorkflowExecution{}
	tests := []struct {
		name                    string
		history                 *historypb.History
		wantWorkflowTaskEventID int64
		wantErr                 bool
	}{
		{
			name:                    "Test history with all completed task event history",
			history:                 generateEventHistory("ccccc"),
			wantWorkflowTaskEventID: NumTotalEvents,
		},
		{
			name:                    "Test history with last task failing",
			history:                 generateEventHistory("ccccf"),
			wantWorkflowTaskEventID: NumTotalEvents,
		},
		{
			name:                    "Test history with first task failing",
			history:                 generateEventHistory("fcccc"),
			wantWorkflowTaskEventID: NumTotalEvents*1 + NumTotalEvents,
		},
		{
			name:                    "Test history with all tasks failing",
			history:                 generateEventHistory("fffff"),
			wantWorkflowTaskEventID: 2,
		},
		{
			name:                    "Test history with some tasks failing in the middle",
			history:                 generateEventHistory("cfffc"),
			wantWorkflowTaskEventID: NumTotalEvents,
		},
		{
			name:    "Test history with empty history should error",
			history: generateEventHistory(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.mockFrontendClient.EXPECT().GetWorkflowExecutionHistory(ctx, gomock.Any()).Return(
				&workflowservice.GetWorkflowExecutionHistoryResponse{History: tt.history, NextPageToken: nil}, nil)
			gotWorkflowTaskEventID, err := getFirstWorkflowTaskEventID(ctx, namespaceStr, &workflowExecution, s.mockFrontendClient, log.NewTestLogger())
			s.Equal(tt.wantErr, err != nil)
			s.Equal(tt.wantWorkflowTaskEventID, gotWorkflowTaskEventID)
		})
	}
}

func (s *activitiesSuite) TestGetResetPoint() {
	ctx := context.Background()
	ns := "namespacename"
	tests := []struct {
		name                    string
		points                  []*workflowpb.ResetPointInfo
		buildId                 string
		currentRunOnly          bool
		wantWorkflowTaskEventID int64
		wantErr                 bool
		wantSetRunId            string
	}{
		{
			name: "not found",
			points: []*workflowpb.ResetPointInfo{
				{
					BuildId:                      "build1",
					RunId:                        "run1",
					FirstWorkflowTaskCompletedId: 123,
					Resettable:                   true,
				},
			},
			buildId: "otherbuild",
			wantErr: true,
		},
		{
			name: "found",
			points: []*workflowpb.ResetPointInfo{
				{
					BuildId:                      "build1",
					RunId:                        "run1",
					FirstWorkflowTaskCompletedId: 123,
					Resettable:                   true,
				},
			},
			buildId:                 "build1",
			wantWorkflowTaskEventID: 123,
		},
		{
			name: "not resettable",
			points: []*workflowpb.ResetPointInfo{
				{
					BuildId:                      "build1",
					RunId:                        "run1",
					FirstWorkflowTaskCompletedId: 123,
					Resettable:                   false,
				},
			},
			buildId: "build1",
			wantErr: true,
		},
		{
			name: "from another run",
			points: []*workflowpb.ResetPointInfo{
				{
					BuildId:                      "build1",
					RunId:                        "run0",
					FirstWorkflowTaskCompletedId: 34,
					Resettable:                   true,
				},
			},
			buildId:                 "build1",
			wantWorkflowTaskEventID: 34,
			wantSetRunId:            "run0",
		},
		{
			name: "from another run but not allowed",
			points: []*workflowpb.ResetPointInfo{
				{
					BuildId:                      "build1",
					RunId:                        "run0",
					FirstWorkflowTaskCompletedId: 34,
					Resettable:                   true,
				},
			},
			buildId:        "build1",
			currentRunOnly: true,
			wantErr:        true,
		},
		{
			name: "expired",
			points: []*workflowpb.ResetPointInfo{
				{
					BuildId:                      "build1",
					RunId:                        "run1",
					FirstWorkflowTaskCompletedId: 123,
					Resettable:                   true,
					ExpireTime:                   timestamp.TimePtr(time.Now().Add(-1 * time.Hour)),
				},
			},
			buildId: "build1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mockFrontendClient.EXPECT().DescribeWorkflowExecution(gomock.Any(), gomock.Any()).Return(
				&workflowservice.DescribeWorkflowExecutionResponse{
					WorkflowExecutionInfo: &workflowpb.WorkflowExecutionInfo{
						AutoResetPoints: &workflowpb.ResetPoints{
							Points: tt.points,
						},
					},
				},
				nil,
			)
			execution := &commonpb.WorkflowExecution{
				WorkflowId: "wfid",
				RunId:      "run1",
			}
			id, err := getResetPoint(ctx, ns, execution, s.mockFrontendClient, tt.buildId, tt.currentRunOnly)
			s.Equal(tt.wantErr, err != nil)
			s.Equal(tt.wantWorkflowTaskEventID, id)
			if tt.wantSetRunId != "" {
				s.Equal(tt.wantSetRunId, execution.RunId)
			}
		})
	}
}

func (s *activitiesSuite) TestAdjustQuery() {
	tests := []struct {
		name           string
		query          string
		expectedResult string
		batchType      string
	}{
		{
			name:           "Empty query",
			query:          "",
			expectedResult: "",
			batchType:      BatchTypeTerminate,
		},
		{
			name:           "Acceptance",
			query:          "A=B",
			expectedResult: fmt.Sprintf("(A=B) AND (%s)", statusRunningQueryFilter),
			batchType:      BatchTypeTerminate,
		},
		{
			name:           "Acceptance with parenthesis",
			query:          "(A=B)",
			expectedResult: fmt.Sprintf("((A=B)) AND (%s)", statusRunningQueryFilter),
			batchType:      BatchTypeTerminate,
		},
		{
			name:           "Acceptance with multiple conditions",
			query:          "(A=B) OR C=D",
			expectedResult: fmt.Sprintf("((A=B) OR C=D) AND (%s)", statusRunningQueryFilter),
			batchType:      BatchTypeTerminate,
		},
		{
			name:           "Contains status - 1",
			query:          "ExecutionStatus=Completed",
			expectedResult: fmt.Sprintf("(ExecutionStatus=Completed) AND (%s)", statusRunningQueryFilter),
			batchType:      BatchTypeTerminate,
		},
		{
			name:           "Contains status - 2",
			query:          "A=B OR ExecutionStatus='Completed'",
			expectedResult: fmt.Sprintf("(A=B OR ExecutionStatus='Completed') AND (%s)", statusRunningQueryFilter),
			batchType:      BatchTypeTerminate,
		},
		{
			name:           "Not supported batch type",
			query:          "A=B",
			expectedResult: "A=B",
			batchType:      "NotSupported",
		},
	}
	for _, testRun := range tests {
		s.Run(testRun.name, func() {
			a := activities{}
			batchParams := BatchParams{Query: testRun.query, BatchType: testRun.batchType}
			adjustedQuery := a.adjustQuery(batchParams)
			s.Equal(testRun.expectedResult, adjustedQuery)
		})
	}
}
