// Code generated by cmd/tools/genrpcwrappers. DO NOT EDIT.

package history

import (
	"context"

	"go.temporal.io/server/api/historyservice/v1"
	"google.golang.org/grpc"

	"go.temporal.io/server/common/backoff"
)

func (c *retryableClient) AddTasks(
	ctx context.Context,
	request *historyservice.AddTasksRequest,
	opts ...grpc.CallOption,
) (*historyservice.AddTasksResponse, error) {
	var resp *historyservice.AddTasksResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.AddTasks(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) CloseShard(
	ctx context.Context,
	request *historyservice.CloseShardRequest,
	opts ...grpc.CallOption,
) (*historyservice.CloseShardResponse, error) {
	var resp *historyservice.CloseShardResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.CloseShard(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) CompleteNexusOperation(
	ctx context.Context,
	request *historyservice.CompleteNexusOperationRequest,
	opts ...grpc.CallOption,
) (*historyservice.CompleteNexusOperationResponse, error) {
	var resp *historyservice.CompleteNexusOperationResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.CompleteNexusOperation(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) DeepHealthCheck(
	ctx context.Context,
	request *historyservice.DeepHealthCheckRequest,
	opts ...grpc.CallOption,
) (*historyservice.DeepHealthCheckResponse, error) {
	var resp *historyservice.DeepHealthCheckResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.DeepHealthCheck(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) DeleteDLQTasks(
	ctx context.Context,
	request *historyservice.DeleteDLQTasksRequest,
	opts ...grpc.CallOption,
) (*historyservice.DeleteDLQTasksResponse, error) {
	var resp *historyservice.DeleteDLQTasksResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.DeleteDLQTasks(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) DeleteWorkflowExecution(
	ctx context.Context,
	request *historyservice.DeleteWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.DeleteWorkflowExecutionResponse, error) {
	var resp *historyservice.DeleteWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.DeleteWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) DeleteWorkflowVisibilityRecord(
	ctx context.Context,
	request *historyservice.DeleteWorkflowVisibilityRecordRequest,
	opts ...grpc.CallOption,
) (*historyservice.DeleteWorkflowVisibilityRecordResponse, error) {
	var resp *historyservice.DeleteWorkflowVisibilityRecordResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.DeleteWorkflowVisibilityRecord(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) DescribeHistoryHost(
	ctx context.Context,
	request *historyservice.DescribeHistoryHostRequest,
	opts ...grpc.CallOption,
) (*historyservice.DescribeHistoryHostResponse, error) {
	var resp *historyservice.DescribeHistoryHostResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.DescribeHistoryHost(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) DescribeMutableState(
	ctx context.Context,
	request *historyservice.DescribeMutableStateRequest,
	opts ...grpc.CallOption,
) (*historyservice.DescribeMutableStateResponse, error) {
	var resp *historyservice.DescribeMutableStateResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.DescribeMutableState(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) DescribeWorkflowExecution(
	ctx context.Context,
	request *historyservice.DescribeWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.DescribeWorkflowExecutionResponse, error) {
	var resp *historyservice.DescribeWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.DescribeWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ExecuteMultiOperation(
	ctx context.Context,
	request *historyservice.ExecuteMultiOperationRequest,
	opts ...grpc.CallOption,
) (*historyservice.ExecuteMultiOperationResponse, error) {
	var resp *historyservice.ExecuteMultiOperationResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ExecuteMultiOperation(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ForceDeleteWorkflowExecution(
	ctx context.Context,
	request *historyservice.ForceDeleteWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.ForceDeleteWorkflowExecutionResponse, error) {
	var resp *historyservice.ForceDeleteWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ForceDeleteWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GenerateLastHistoryReplicationTasks(
	ctx context.Context,
	request *historyservice.GenerateLastHistoryReplicationTasksRequest,
	opts ...grpc.CallOption,
) (*historyservice.GenerateLastHistoryReplicationTasksResponse, error) {
	var resp *historyservice.GenerateLastHistoryReplicationTasksResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GenerateLastHistoryReplicationTasks(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetDLQMessages(
	ctx context.Context,
	request *historyservice.GetDLQMessagesRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetDLQMessagesResponse, error) {
	var resp *historyservice.GetDLQMessagesResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetDLQMessages(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetDLQReplicationMessages(
	ctx context.Context,
	request *historyservice.GetDLQReplicationMessagesRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetDLQReplicationMessagesResponse, error) {
	var resp *historyservice.GetDLQReplicationMessagesResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetDLQReplicationMessages(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetDLQTasks(
	ctx context.Context,
	request *historyservice.GetDLQTasksRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetDLQTasksResponse, error) {
	var resp *historyservice.GetDLQTasksResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetDLQTasks(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetMutableState(
	ctx context.Context,
	request *historyservice.GetMutableStateRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetMutableStateResponse, error) {
	var resp *historyservice.GetMutableStateResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetMutableState(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetReplicationMessages(
	ctx context.Context,
	request *historyservice.GetReplicationMessagesRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetReplicationMessagesResponse, error) {
	var resp *historyservice.GetReplicationMessagesResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetReplicationMessages(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetReplicationStatus(
	ctx context.Context,
	request *historyservice.GetReplicationStatusRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetReplicationStatusResponse, error) {
	var resp *historyservice.GetReplicationStatusResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetReplicationStatus(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetShard(
	ctx context.Context,
	request *historyservice.GetShardRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetShardResponse, error) {
	var resp *historyservice.GetShardResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetShard(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetWorkflowExecutionHistory(
	ctx context.Context,
	request *historyservice.GetWorkflowExecutionHistoryRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetWorkflowExecutionHistoryResponse, error) {
	var resp *historyservice.GetWorkflowExecutionHistoryResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetWorkflowExecutionHistory(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetWorkflowExecutionHistoryReverse(
	ctx context.Context,
	request *historyservice.GetWorkflowExecutionHistoryReverseRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetWorkflowExecutionHistoryReverseResponse, error) {
	var resp *historyservice.GetWorkflowExecutionHistoryReverseResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetWorkflowExecutionHistoryReverse(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetWorkflowExecutionRawHistory(
	ctx context.Context,
	request *historyservice.GetWorkflowExecutionRawHistoryRequest,
	opts ...grpc.CallOption,
) (*historyservice.GetWorkflowExecutionRawHistoryResponse, error) {
	var resp *historyservice.GetWorkflowExecutionRawHistoryResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetWorkflowExecutionRawHistory(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) GetWorkflowExecutionRawHistoryV2(
	ctx context.Context,
	request *historyservice.GetWorkflowExecutionRawHistoryV2Request,
	opts ...grpc.CallOption,
) (*historyservice.GetWorkflowExecutionRawHistoryV2Response, error) {
	var resp *historyservice.GetWorkflowExecutionRawHistoryV2Response
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.GetWorkflowExecutionRawHistoryV2(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ImportWorkflowExecution(
	ctx context.Context,
	request *historyservice.ImportWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.ImportWorkflowExecutionResponse, error) {
	var resp *historyservice.ImportWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ImportWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) InvokeStateMachineMethod(
	ctx context.Context,
	request *historyservice.InvokeStateMachineMethodRequest,
	opts ...grpc.CallOption,
) (*historyservice.InvokeStateMachineMethodResponse, error) {
	var resp *historyservice.InvokeStateMachineMethodResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.InvokeStateMachineMethod(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) IsActivityTaskValid(
	ctx context.Context,
	request *historyservice.IsActivityTaskValidRequest,
	opts ...grpc.CallOption,
) (*historyservice.IsActivityTaskValidResponse, error) {
	var resp *historyservice.IsActivityTaskValidResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.IsActivityTaskValid(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) IsWorkflowTaskValid(
	ctx context.Context,
	request *historyservice.IsWorkflowTaskValidRequest,
	opts ...grpc.CallOption,
) (*historyservice.IsWorkflowTaskValidResponse, error) {
	var resp *historyservice.IsWorkflowTaskValidResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.IsWorkflowTaskValid(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ListQueues(
	ctx context.Context,
	request *historyservice.ListQueuesRequest,
	opts ...grpc.CallOption,
) (*historyservice.ListQueuesResponse, error) {
	var resp *historyservice.ListQueuesResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ListQueues(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ListTasks(
	ctx context.Context,
	request *historyservice.ListTasksRequest,
	opts ...grpc.CallOption,
) (*historyservice.ListTasksResponse, error) {
	var resp *historyservice.ListTasksResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ListTasks(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) MergeDLQMessages(
	ctx context.Context,
	request *historyservice.MergeDLQMessagesRequest,
	opts ...grpc.CallOption,
) (*historyservice.MergeDLQMessagesResponse, error) {
	var resp *historyservice.MergeDLQMessagesResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.MergeDLQMessages(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) PauseActivity(
	ctx context.Context,
	request *historyservice.PauseActivityRequest,
	opts ...grpc.CallOption,
) (*historyservice.PauseActivityResponse, error) {
	var resp *historyservice.PauseActivityResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.PauseActivity(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) PollMutableState(
	ctx context.Context,
	request *historyservice.PollMutableStateRequest,
	opts ...grpc.CallOption,
) (*historyservice.PollMutableStateResponse, error) {
	var resp *historyservice.PollMutableStateResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.PollMutableState(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) PollWorkflowExecutionUpdate(
	ctx context.Context,
	request *historyservice.PollWorkflowExecutionUpdateRequest,
	opts ...grpc.CallOption,
) (*historyservice.PollWorkflowExecutionUpdateResponse, error) {
	var resp *historyservice.PollWorkflowExecutionUpdateResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.PollWorkflowExecutionUpdate(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) PurgeDLQMessages(
	ctx context.Context,
	request *historyservice.PurgeDLQMessagesRequest,
	opts ...grpc.CallOption,
) (*historyservice.PurgeDLQMessagesResponse, error) {
	var resp *historyservice.PurgeDLQMessagesResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.PurgeDLQMessages(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) QueryWorkflow(
	ctx context.Context,
	request *historyservice.QueryWorkflowRequest,
	opts ...grpc.CallOption,
) (*historyservice.QueryWorkflowResponse, error) {
	var resp *historyservice.QueryWorkflowResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.QueryWorkflow(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ReapplyEvents(
	ctx context.Context,
	request *historyservice.ReapplyEventsRequest,
	opts ...grpc.CallOption,
) (*historyservice.ReapplyEventsResponse, error) {
	var resp *historyservice.ReapplyEventsResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ReapplyEvents(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RebuildMutableState(
	ctx context.Context,
	request *historyservice.RebuildMutableStateRequest,
	opts ...grpc.CallOption,
) (*historyservice.RebuildMutableStateResponse, error) {
	var resp *historyservice.RebuildMutableStateResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RebuildMutableState(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RecordActivityTaskHeartbeat(
	ctx context.Context,
	request *historyservice.RecordActivityTaskHeartbeatRequest,
	opts ...grpc.CallOption,
) (*historyservice.RecordActivityTaskHeartbeatResponse, error) {
	var resp *historyservice.RecordActivityTaskHeartbeatResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RecordActivityTaskHeartbeat(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RecordActivityTaskStarted(
	ctx context.Context,
	request *historyservice.RecordActivityTaskStartedRequest,
	opts ...grpc.CallOption,
) (*historyservice.RecordActivityTaskStartedResponse, error) {
	var resp *historyservice.RecordActivityTaskStartedResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RecordActivityTaskStarted(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RecordChildExecutionCompleted(
	ctx context.Context,
	request *historyservice.RecordChildExecutionCompletedRequest,
	opts ...grpc.CallOption,
) (*historyservice.RecordChildExecutionCompletedResponse, error) {
	var resp *historyservice.RecordChildExecutionCompletedResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RecordChildExecutionCompleted(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RecordWorkflowTaskStarted(
	ctx context.Context,
	request *historyservice.RecordWorkflowTaskStartedRequest,
	opts ...grpc.CallOption,
) (*historyservice.RecordWorkflowTaskStartedResponse, error) {
	var resp *historyservice.RecordWorkflowTaskStartedResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RecordWorkflowTaskStarted(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RefreshWorkflowTasks(
	ctx context.Context,
	request *historyservice.RefreshWorkflowTasksRequest,
	opts ...grpc.CallOption,
) (*historyservice.RefreshWorkflowTasksResponse, error) {
	var resp *historyservice.RefreshWorkflowTasksResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RefreshWorkflowTasks(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RemoveSignalMutableState(
	ctx context.Context,
	request *historyservice.RemoveSignalMutableStateRequest,
	opts ...grpc.CallOption,
) (*historyservice.RemoveSignalMutableStateResponse, error) {
	var resp *historyservice.RemoveSignalMutableStateResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RemoveSignalMutableState(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RemoveTask(
	ctx context.Context,
	request *historyservice.RemoveTaskRequest,
	opts ...grpc.CallOption,
) (*historyservice.RemoveTaskResponse, error) {
	var resp *historyservice.RemoveTaskResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RemoveTask(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ReplicateEventsV2(
	ctx context.Context,
	request *historyservice.ReplicateEventsV2Request,
	opts ...grpc.CallOption,
) (*historyservice.ReplicateEventsV2Response, error) {
	var resp *historyservice.ReplicateEventsV2Response
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ReplicateEventsV2(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ReplicateWorkflowState(
	ctx context.Context,
	request *historyservice.ReplicateWorkflowStateRequest,
	opts ...grpc.CallOption,
) (*historyservice.ReplicateWorkflowStateResponse, error) {
	var resp *historyservice.ReplicateWorkflowStateResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ReplicateWorkflowState(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RequestCancelWorkflowExecution(
	ctx context.Context,
	request *historyservice.RequestCancelWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.RequestCancelWorkflowExecutionResponse, error) {
	var resp *historyservice.RequestCancelWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RequestCancelWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ResetActivity(
	ctx context.Context,
	request *historyservice.ResetActivityRequest,
	opts ...grpc.CallOption,
) (*historyservice.ResetActivityResponse, error) {
	var resp *historyservice.ResetActivityResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ResetActivity(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ResetStickyTaskQueue(
	ctx context.Context,
	request *historyservice.ResetStickyTaskQueueRequest,
	opts ...grpc.CallOption,
) (*historyservice.ResetStickyTaskQueueResponse, error) {
	var resp *historyservice.ResetStickyTaskQueueResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ResetStickyTaskQueue(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ResetWorkflowExecution(
	ctx context.Context,
	request *historyservice.ResetWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.ResetWorkflowExecutionResponse, error) {
	var resp *historyservice.ResetWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ResetWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RespondActivityTaskCanceled(
	ctx context.Context,
	request *historyservice.RespondActivityTaskCanceledRequest,
	opts ...grpc.CallOption,
) (*historyservice.RespondActivityTaskCanceledResponse, error) {
	var resp *historyservice.RespondActivityTaskCanceledResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RespondActivityTaskCanceled(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RespondActivityTaskCompleted(
	ctx context.Context,
	request *historyservice.RespondActivityTaskCompletedRequest,
	opts ...grpc.CallOption,
) (*historyservice.RespondActivityTaskCompletedResponse, error) {
	var resp *historyservice.RespondActivityTaskCompletedResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RespondActivityTaskCompleted(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RespondActivityTaskFailed(
	ctx context.Context,
	request *historyservice.RespondActivityTaskFailedRequest,
	opts ...grpc.CallOption,
) (*historyservice.RespondActivityTaskFailedResponse, error) {
	var resp *historyservice.RespondActivityTaskFailedResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RespondActivityTaskFailed(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RespondWorkflowTaskCompleted(
	ctx context.Context,
	request *historyservice.RespondWorkflowTaskCompletedRequest,
	opts ...grpc.CallOption,
) (*historyservice.RespondWorkflowTaskCompletedResponse, error) {
	var resp *historyservice.RespondWorkflowTaskCompletedResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RespondWorkflowTaskCompleted(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) RespondWorkflowTaskFailed(
	ctx context.Context,
	request *historyservice.RespondWorkflowTaskFailedRequest,
	opts ...grpc.CallOption,
) (*historyservice.RespondWorkflowTaskFailedResponse, error) {
	var resp *historyservice.RespondWorkflowTaskFailedResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.RespondWorkflowTaskFailed(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) ScheduleWorkflowTask(
	ctx context.Context,
	request *historyservice.ScheduleWorkflowTaskRequest,
	opts ...grpc.CallOption,
) (*historyservice.ScheduleWorkflowTaskResponse, error) {
	var resp *historyservice.ScheduleWorkflowTaskResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.ScheduleWorkflowTask(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) SignalWithStartWorkflowExecution(
	ctx context.Context,
	request *historyservice.SignalWithStartWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.SignalWithStartWorkflowExecutionResponse, error) {
	var resp *historyservice.SignalWithStartWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.SignalWithStartWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) SignalWorkflowExecution(
	ctx context.Context,
	request *historyservice.SignalWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.SignalWorkflowExecutionResponse, error) {
	var resp *historyservice.SignalWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.SignalWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) StartWorkflowExecution(
	ctx context.Context,
	request *historyservice.StartWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.StartWorkflowExecutionResponse, error) {
	var resp *historyservice.StartWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.StartWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) SyncActivity(
	ctx context.Context,
	request *historyservice.SyncActivityRequest,
	opts ...grpc.CallOption,
) (*historyservice.SyncActivityResponse, error) {
	var resp *historyservice.SyncActivityResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.SyncActivity(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) SyncShardStatus(
	ctx context.Context,
	request *historyservice.SyncShardStatusRequest,
	opts ...grpc.CallOption,
) (*historyservice.SyncShardStatusResponse, error) {
	var resp *historyservice.SyncShardStatusResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.SyncShardStatus(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) SyncWorkflowState(
	ctx context.Context,
	request *historyservice.SyncWorkflowStateRequest,
	opts ...grpc.CallOption,
) (*historyservice.SyncWorkflowStateResponse, error) {
	var resp *historyservice.SyncWorkflowStateResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.SyncWorkflowState(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) TerminateWorkflowExecution(
	ctx context.Context,
	request *historyservice.TerminateWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.TerminateWorkflowExecutionResponse, error) {
	var resp *historyservice.TerminateWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.TerminateWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) UnpauseActivity(
	ctx context.Context,
	request *historyservice.UnpauseActivityRequest,
	opts ...grpc.CallOption,
) (*historyservice.UnpauseActivityResponse, error) {
	var resp *historyservice.UnpauseActivityResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.UnpauseActivity(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) UpdateActivityOptions(
	ctx context.Context,
	request *historyservice.UpdateActivityOptionsRequest,
	opts ...grpc.CallOption,
) (*historyservice.UpdateActivityOptionsResponse, error) {
	var resp *historyservice.UpdateActivityOptionsResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.UpdateActivityOptions(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) UpdateWorkflowExecution(
	ctx context.Context,
	request *historyservice.UpdateWorkflowExecutionRequest,
	opts ...grpc.CallOption,
) (*historyservice.UpdateWorkflowExecutionResponse, error) {
	var resp *historyservice.UpdateWorkflowExecutionResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.UpdateWorkflowExecution(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) UpdateWorkflowExecutionOptions(
	ctx context.Context,
	request *historyservice.UpdateWorkflowExecutionOptionsRequest,
	opts ...grpc.CallOption,
) (*historyservice.UpdateWorkflowExecutionOptionsResponse, error) {
	var resp *historyservice.UpdateWorkflowExecutionOptionsResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.UpdateWorkflowExecutionOptions(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) VerifyChildExecutionCompletionRecorded(
	ctx context.Context,
	request *historyservice.VerifyChildExecutionCompletionRecordedRequest,
	opts ...grpc.CallOption,
) (*historyservice.VerifyChildExecutionCompletionRecordedResponse, error) {
	var resp *historyservice.VerifyChildExecutionCompletionRecordedResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.VerifyChildExecutionCompletionRecorded(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}

func (c *retryableClient) VerifyFirstWorkflowTaskScheduled(
	ctx context.Context,
	request *historyservice.VerifyFirstWorkflowTaskScheduledRequest,
	opts ...grpc.CallOption,
) (*historyservice.VerifyFirstWorkflowTaskScheduledResponse, error) {
	var resp *historyservice.VerifyFirstWorkflowTaskScheduledResponse
	op := func(ctx context.Context) error {
		var err error
		resp, err = c.client.VerifyFirstWorkflowTaskScheduled(ctx, request, opts...)
		return err
	}
	err := backoff.ThrottleRetryContext(ctx, op, c.policy, c.isRetryable)
	return resp, err
}
