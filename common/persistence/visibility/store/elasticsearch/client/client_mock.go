// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -package client -source client.go -destination client_mock.go
//

// Package client is a generated GoMock package.
package client

import (
	context "context"
	reflect "reflect"

	elastic "github.com/olivere/elastic/v7"
	enums "go.temporal.io/api/enums/v1"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
	isgomock struct{}
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CatIndices mocks base method.
func (m *MockClient) CatIndices(ctx context.Context, target string) (elastic.CatIndicesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CatIndices", ctx, target)
	ret0, _ := ret[0].(elastic.CatIndicesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CatIndices indicates an expected call of CatIndices.
func (mr *MockClientMockRecorder) CatIndices(ctx, target any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CatIndices", reflect.TypeOf((*MockClient)(nil).CatIndices), ctx, target)
}

// ClosePointInTime mocks base method.
func (m *MockClient) ClosePointInTime(ctx context.Context, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClosePointInTime", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClosePointInTime indicates an expected call of ClosePointInTime.
func (mr *MockClientMockRecorder) ClosePointInTime(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClosePointInTime", reflect.TypeOf((*MockClient)(nil).ClosePointInTime), ctx, id)
}

// CloseScroll mocks base method.
func (m *MockClient) CloseScroll(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseScroll", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseScroll indicates an expected call of CloseScroll.
func (mr *MockClientMockRecorder) CloseScroll(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseScroll", reflect.TypeOf((*MockClient)(nil).CloseScroll), ctx, id)
}

// Count mocks base method.
func (m *MockClient) Count(ctx context.Context, index string, query elastic.Query) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, index, query)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockClientMockRecorder) Count(ctx, index, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockClient)(nil).Count), ctx, index, query)
}

// CountGroupBy mocks base method.
func (m *MockClient) CountGroupBy(ctx context.Context, index string, query elastic.Query, aggName string, agg elastic.Aggregation) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountGroupBy", ctx, index, query, aggName, agg)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountGroupBy indicates an expected call of CountGroupBy.
func (mr *MockClientMockRecorder) CountGroupBy(ctx, index, query, aggName, agg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountGroupBy", reflect.TypeOf((*MockClient)(nil).CountGroupBy), ctx, index, query, aggName, agg)
}

// CreateIndex mocks base method.
func (m *MockClient) CreateIndex(ctx context.Context, index string, body map[string]any) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndex", ctx, index, body)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIndex indicates an expected call of CreateIndex.
func (mr *MockClientMockRecorder) CreateIndex(ctx, index, body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndex", reflect.TypeOf((*MockClient)(nil).CreateIndex), ctx, index, body)
}

// DeleteIndex mocks base method.
func (m *MockClient) DeleteIndex(ctx context.Context, indexName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIndex", ctx, indexName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteIndex indicates an expected call of DeleteIndex.
func (mr *MockClientMockRecorder) DeleteIndex(ctx, indexName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIndex", reflect.TypeOf((*MockClient)(nil).DeleteIndex), ctx, indexName)
}

// Get mocks base method.
func (m *MockClient) Get(ctx context.Context, index, docID string) (*elastic.GetResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, index, docID)
	ret0, _ := ret[0].(*elastic.GetResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockClientMockRecorder) Get(ctx, index, docID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClient)(nil).Get), ctx, index, docID)
}

// GetMapping mocks base method.
func (m *MockClient) GetMapping(ctx context.Context, index string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMapping", ctx, index)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMapping indicates an expected call of GetMapping.
func (mr *MockClientMockRecorder) GetMapping(ctx, index any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMapping", reflect.TypeOf((*MockClient)(nil).GetMapping), ctx, index)
}

// IndexExists mocks base method.
func (m *MockClient) IndexExists(ctx context.Context, indexName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexExists", ctx, indexName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexExists indicates an expected call of IndexExists.
func (mr *MockClientMockRecorder) IndexExists(ctx, indexName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexExists", reflect.TypeOf((*MockClient)(nil).IndexExists), ctx, indexName)
}

// IsPointInTimeSupported mocks base method.
func (m *MockClient) IsPointInTimeSupported(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPointInTimeSupported", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPointInTimeSupported indicates an expected call of IsPointInTimeSupported.
func (mr *MockClientMockRecorder) IsPointInTimeSupported(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPointInTimeSupported", reflect.TypeOf((*MockClient)(nil).IsPointInTimeSupported), ctx)
}

// OpenPointInTime mocks base method.
func (m *MockClient) OpenPointInTime(ctx context.Context, index, keepAliveInterval string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenPointInTime", ctx, index, keepAliveInterval)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenPointInTime indicates an expected call of OpenPointInTime.
func (mr *MockClientMockRecorder) OpenPointInTime(ctx, index, keepAliveInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenPointInTime", reflect.TypeOf((*MockClient)(nil).OpenPointInTime), ctx, index, keepAliveInterval)
}

// OpenScroll mocks base method.
func (m *MockClient) OpenScroll(ctx context.Context, p *SearchParameters, keepAliveInterval string) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenScroll", ctx, p, keepAliveInterval)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenScroll indicates an expected call of OpenScroll.
func (mr *MockClientMockRecorder) OpenScroll(ctx, p, keepAliveInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenScroll", reflect.TypeOf((*MockClient)(nil).OpenScroll), ctx, p, keepAliveInterval)
}

// PutMapping mocks base method.
func (m *MockClient) PutMapping(ctx context.Context, index string, mapping map[string]enums.IndexedValueType) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutMapping", ctx, index, mapping)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutMapping indicates an expected call of PutMapping.
func (mr *MockClientMockRecorder) PutMapping(ctx, index, mapping any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutMapping", reflect.TypeOf((*MockClient)(nil).PutMapping), ctx, index, mapping)
}

// RunBulkProcessor mocks base method.
func (m *MockClient) RunBulkProcessor(ctx context.Context, p *BulkProcessorParameters) (BulkProcessor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunBulkProcessor", ctx, p)
	ret0, _ := ret[0].(BulkProcessor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunBulkProcessor indicates an expected call of RunBulkProcessor.
func (mr *MockClientMockRecorder) RunBulkProcessor(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunBulkProcessor", reflect.TypeOf((*MockClient)(nil).RunBulkProcessor), ctx, p)
}

// Scroll mocks base method.
func (m *MockClient) Scroll(ctx context.Context, id, keepAliveInterval string) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scroll", ctx, id, keepAliveInterval)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scroll indicates an expected call of Scroll.
func (mr *MockClientMockRecorder) Scroll(ctx, id, keepAliveInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scroll", reflect.TypeOf((*MockClient)(nil).Scroll), ctx, id, keepAliveInterval)
}

// Search mocks base method.
func (m *MockClient) Search(ctx context.Context, p *SearchParameters) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, p)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockClientMockRecorder) Search(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockClient)(nil).Search), ctx, p)
}

// WaitForYellowStatus mocks base method.
func (m *MockClient) WaitForYellowStatus(ctx context.Context, index string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForYellowStatus", ctx, index)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WaitForYellowStatus indicates an expected call of WaitForYellowStatus.
func (mr *MockClientMockRecorder) WaitForYellowStatus(ctx, index any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForYellowStatus", reflect.TypeOf((*MockClient)(nil).WaitForYellowStatus), ctx, index)
}

// MockCLIClient is a mock of CLIClient interface.
type MockCLIClient struct {
	ctrl     *gomock.Controller
	recorder *MockCLIClientMockRecorder
	isgomock struct{}
}

// MockCLIClientMockRecorder is the mock recorder for MockCLIClient.
type MockCLIClientMockRecorder struct {
	mock *MockCLIClient
}

// NewMockCLIClient creates a new mock instance.
func NewMockCLIClient(ctrl *gomock.Controller) *MockCLIClient {
	mock := &MockCLIClient{ctrl: ctrl}
	mock.recorder = &MockCLIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCLIClient) EXPECT() *MockCLIClientMockRecorder {
	return m.recorder
}

// CatIndices mocks base method.
func (m *MockCLIClient) CatIndices(ctx context.Context, target string) (elastic.CatIndicesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CatIndices", ctx, target)
	ret0, _ := ret[0].(elastic.CatIndicesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CatIndices indicates an expected call of CatIndices.
func (mr *MockCLIClientMockRecorder) CatIndices(ctx, target any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CatIndices", reflect.TypeOf((*MockCLIClient)(nil).CatIndices), ctx, target)
}

// ClosePointInTime mocks base method.
func (m *MockCLIClient) ClosePointInTime(ctx context.Context, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClosePointInTime", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClosePointInTime indicates an expected call of ClosePointInTime.
func (mr *MockCLIClientMockRecorder) ClosePointInTime(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClosePointInTime", reflect.TypeOf((*MockCLIClient)(nil).ClosePointInTime), ctx, id)
}

// CloseScroll mocks base method.
func (m *MockCLIClient) CloseScroll(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseScroll", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseScroll indicates an expected call of CloseScroll.
func (mr *MockCLIClientMockRecorder) CloseScroll(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseScroll", reflect.TypeOf((*MockCLIClient)(nil).CloseScroll), ctx, id)
}

// Count mocks base method.
func (m *MockCLIClient) Count(ctx context.Context, index string, query elastic.Query) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, index, query)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockCLIClientMockRecorder) Count(ctx, index, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockCLIClient)(nil).Count), ctx, index, query)
}

// CountGroupBy mocks base method.
func (m *MockCLIClient) CountGroupBy(ctx context.Context, index string, query elastic.Query, aggName string, agg elastic.Aggregation) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountGroupBy", ctx, index, query, aggName, agg)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountGroupBy indicates an expected call of CountGroupBy.
func (mr *MockCLIClientMockRecorder) CountGroupBy(ctx, index, query, aggName, agg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountGroupBy", reflect.TypeOf((*MockCLIClient)(nil).CountGroupBy), ctx, index, query, aggName, agg)
}

// CreateIndex mocks base method.
func (m *MockCLIClient) CreateIndex(ctx context.Context, index string, body map[string]any) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndex", ctx, index, body)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIndex indicates an expected call of CreateIndex.
func (mr *MockCLIClientMockRecorder) CreateIndex(ctx, index, body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndex", reflect.TypeOf((*MockCLIClient)(nil).CreateIndex), ctx, index, body)
}

// Delete mocks base method.
func (m *MockCLIClient) Delete(ctx context.Context, indexName, docID string, version int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, indexName, docID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCLIClientMockRecorder) Delete(ctx, indexName, docID, version any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCLIClient)(nil).Delete), ctx, indexName, docID, version)
}

// DeleteIndex mocks base method.
func (m *MockCLIClient) DeleteIndex(ctx context.Context, indexName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIndex", ctx, indexName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteIndex indicates an expected call of DeleteIndex.
func (mr *MockCLIClientMockRecorder) DeleteIndex(ctx, indexName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIndex", reflect.TypeOf((*MockCLIClient)(nil).DeleteIndex), ctx, indexName)
}

// Get mocks base method.
func (m *MockCLIClient) Get(ctx context.Context, index, docID string) (*elastic.GetResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, index, docID)
	ret0, _ := ret[0].(*elastic.GetResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCLIClientMockRecorder) Get(ctx, index, docID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCLIClient)(nil).Get), ctx, index, docID)
}

// GetMapping mocks base method.
func (m *MockCLIClient) GetMapping(ctx context.Context, index string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMapping", ctx, index)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMapping indicates an expected call of GetMapping.
func (mr *MockCLIClientMockRecorder) GetMapping(ctx, index any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMapping", reflect.TypeOf((*MockCLIClient)(nil).GetMapping), ctx, index)
}

// IndexExists mocks base method.
func (m *MockCLIClient) IndexExists(ctx context.Context, indexName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexExists", ctx, indexName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexExists indicates an expected call of IndexExists.
func (mr *MockCLIClientMockRecorder) IndexExists(ctx, indexName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexExists", reflect.TypeOf((*MockCLIClient)(nil).IndexExists), ctx, indexName)
}

// IsPointInTimeSupported mocks base method.
func (m *MockCLIClient) IsPointInTimeSupported(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPointInTimeSupported", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPointInTimeSupported indicates an expected call of IsPointInTimeSupported.
func (mr *MockCLIClientMockRecorder) IsPointInTimeSupported(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPointInTimeSupported", reflect.TypeOf((*MockCLIClient)(nil).IsPointInTimeSupported), ctx)
}

// OpenPointInTime mocks base method.
func (m *MockCLIClient) OpenPointInTime(ctx context.Context, index, keepAliveInterval string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenPointInTime", ctx, index, keepAliveInterval)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenPointInTime indicates an expected call of OpenPointInTime.
func (mr *MockCLIClientMockRecorder) OpenPointInTime(ctx, index, keepAliveInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenPointInTime", reflect.TypeOf((*MockCLIClient)(nil).OpenPointInTime), ctx, index, keepAliveInterval)
}

// OpenScroll mocks base method.
func (m *MockCLIClient) OpenScroll(ctx context.Context, p *SearchParameters, keepAliveInterval string) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenScroll", ctx, p, keepAliveInterval)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenScroll indicates an expected call of OpenScroll.
func (mr *MockCLIClientMockRecorder) OpenScroll(ctx, p, keepAliveInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenScroll", reflect.TypeOf((*MockCLIClient)(nil).OpenScroll), ctx, p, keepAliveInterval)
}

// PutMapping mocks base method.
func (m *MockCLIClient) PutMapping(ctx context.Context, index string, mapping map[string]enums.IndexedValueType) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutMapping", ctx, index, mapping)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutMapping indicates an expected call of PutMapping.
func (mr *MockCLIClientMockRecorder) PutMapping(ctx, index, mapping any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutMapping", reflect.TypeOf((*MockCLIClient)(nil).PutMapping), ctx, index, mapping)
}

// RunBulkProcessor mocks base method.
func (m *MockCLIClient) RunBulkProcessor(ctx context.Context, p *BulkProcessorParameters) (BulkProcessor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunBulkProcessor", ctx, p)
	ret0, _ := ret[0].(BulkProcessor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunBulkProcessor indicates an expected call of RunBulkProcessor.
func (mr *MockCLIClientMockRecorder) RunBulkProcessor(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunBulkProcessor", reflect.TypeOf((*MockCLIClient)(nil).RunBulkProcessor), ctx, p)
}

// Scroll mocks base method.
func (m *MockCLIClient) Scroll(ctx context.Context, id, keepAliveInterval string) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scroll", ctx, id, keepAliveInterval)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scroll indicates an expected call of Scroll.
func (mr *MockCLIClientMockRecorder) Scroll(ctx, id, keepAliveInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scroll", reflect.TypeOf((*MockCLIClient)(nil).Scroll), ctx, id, keepAliveInterval)
}

// Search mocks base method.
func (m *MockCLIClient) Search(ctx context.Context, p *SearchParameters) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, p)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockCLIClientMockRecorder) Search(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockCLIClient)(nil).Search), ctx, p)
}

// WaitForYellowStatus mocks base method.
func (m *MockCLIClient) WaitForYellowStatus(ctx context.Context, index string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForYellowStatus", ctx, index)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WaitForYellowStatus indicates an expected call of WaitForYellowStatus.
func (mr *MockCLIClientMockRecorder) WaitForYellowStatus(ctx, index any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForYellowStatus", reflect.TypeOf((*MockCLIClient)(nil).WaitForYellowStatus), ctx, index)
}

// MockIntegrationTestsClient is a mock of IntegrationTestsClient interface.
type MockIntegrationTestsClient struct {
	ctrl     *gomock.Controller
	recorder *MockIntegrationTestsClientMockRecorder
	isgomock struct{}
}

// MockIntegrationTestsClientMockRecorder is the mock recorder for MockIntegrationTestsClient.
type MockIntegrationTestsClientMockRecorder struct {
	mock *MockIntegrationTestsClient
}

// NewMockIntegrationTestsClient creates a new mock instance.
func NewMockIntegrationTestsClient(ctrl *gomock.Controller) *MockIntegrationTestsClient {
	mock := &MockIntegrationTestsClient{ctrl: ctrl}
	mock.recorder = &MockIntegrationTestsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIntegrationTestsClient) EXPECT() *MockIntegrationTestsClientMockRecorder {
	return m.recorder
}

// CatIndices mocks base method.
func (m *MockIntegrationTestsClient) CatIndices(ctx context.Context, target string) (elastic.CatIndicesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CatIndices", ctx, target)
	ret0, _ := ret[0].(elastic.CatIndicesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CatIndices indicates an expected call of CatIndices.
func (mr *MockIntegrationTestsClientMockRecorder) CatIndices(ctx, target any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CatIndices", reflect.TypeOf((*MockIntegrationTestsClient)(nil).CatIndices), ctx, target)
}

// ClosePointInTime mocks base method.
func (m *MockIntegrationTestsClient) ClosePointInTime(ctx context.Context, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClosePointInTime", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClosePointInTime indicates an expected call of ClosePointInTime.
func (mr *MockIntegrationTestsClientMockRecorder) ClosePointInTime(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClosePointInTime", reflect.TypeOf((*MockIntegrationTestsClient)(nil).ClosePointInTime), ctx, id)
}

// CloseScroll mocks base method.
func (m *MockIntegrationTestsClient) CloseScroll(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseScroll", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseScroll indicates an expected call of CloseScroll.
func (mr *MockIntegrationTestsClientMockRecorder) CloseScroll(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseScroll", reflect.TypeOf((*MockIntegrationTestsClient)(nil).CloseScroll), ctx, id)
}

// Count mocks base method.
func (m *MockIntegrationTestsClient) Count(ctx context.Context, index string, query elastic.Query) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, index, query)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockIntegrationTestsClientMockRecorder) Count(ctx, index, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIntegrationTestsClient)(nil).Count), ctx, index, query)
}

// CountGroupBy mocks base method.
func (m *MockIntegrationTestsClient) CountGroupBy(ctx context.Context, index string, query elastic.Query, aggName string, agg elastic.Aggregation) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountGroupBy", ctx, index, query, aggName, agg)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountGroupBy indicates an expected call of CountGroupBy.
func (mr *MockIntegrationTestsClientMockRecorder) CountGroupBy(ctx, index, query, aggName, agg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountGroupBy", reflect.TypeOf((*MockIntegrationTestsClient)(nil).CountGroupBy), ctx, index, query, aggName, agg)
}

// CreateIndex mocks base method.
func (m *MockIntegrationTestsClient) CreateIndex(ctx context.Context, index string, body map[string]any) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndex", ctx, index, body)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIndex indicates an expected call of CreateIndex.
func (mr *MockIntegrationTestsClientMockRecorder) CreateIndex(ctx, index, body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndex", reflect.TypeOf((*MockIntegrationTestsClient)(nil).CreateIndex), ctx, index, body)
}

// DeleteIndex mocks base method.
func (m *MockIntegrationTestsClient) DeleteIndex(ctx context.Context, indexName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIndex", ctx, indexName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteIndex indicates an expected call of DeleteIndex.
func (mr *MockIntegrationTestsClientMockRecorder) DeleteIndex(ctx, indexName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIndex", reflect.TypeOf((*MockIntegrationTestsClient)(nil).DeleteIndex), ctx, indexName)
}

// Get mocks base method.
func (m *MockIntegrationTestsClient) Get(ctx context.Context, index, docID string) (*elastic.GetResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, index, docID)
	ret0, _ := ret[0].(*elastic.GetResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIntegrationTestsClientMockRecorder) Get(ctx, index, docID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIntegrationTestsClient)(nil).Get), ctx, index, docID)
}

// GetMapping mocks base method.
func (m *MockIntegrationTestsClient) GetMapping(ctx context.Context, index string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMapping", ctx, index)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMapping indicates an expected call of GetMapping.
func (mr *MockIntegrationTestsClientMockRecorder) GetMapping(ctx, index any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMapping", reflect.TypeOf((*MockIntegrationTestsClient)(nil).GetMapping), ctx, index)
}

// IndexExists mocks base method.
func (m *MockIntegrationTestsClient) IndexExists(ctx context.Context, indexName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexExists", ctx, indexName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexExists indicates an expected call of IndexExists.
func (mr *MockIntegrationTestsClientMockRecorder) IndexExists(ctx, indexName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexExists", reflect.TypeOf((*MockIntegrationTestsClient)(nil).IndexExists), ctx, indexName)
}

// IndexGetSettings mocks base method.
func (m *MockIntegrationTestsClient) IndexGetSettings(ctx context.Context, indexName string) (map[string]*elastic.IndicesGetSettingsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexGetSettings", ctx, indexName)
	ret0, _ := ret[0].(map[string]*elastic.IndicesGetSettingsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexGetSettings indicates an expected call of IndexGetSettings.
func (mr *MockIntegrationTestsClientMockRecorder) IndexGetSettings(ctx, indexName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexGetSettings", reflect.TypeOf((*MockIntegrationTestsClient)(nil).IndexGetSettings), ctx, indexName)
}

// IndexPutSettings mocks base method.
func (m *MockIntegrationTestsClient) IndexPutSettings(ctx context.Context, indexName, bodyString string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexPutSettings", ctx, indexName, bodyString)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexPutSettings indicates an expected call of IndexPutSettings.
func (mr *MockIntegrationTestsClientMockRecorder) IndexPutSettings(ctx, indexName, bodyString any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexPutSettings", reflect.TypeOf((*MockIntegrationTestsClient)(nil).IndexPutSettings), ctx, indexName, bodyString)
}

// IndexPutTemplate mocks base method.
func (m *MockIntegrationTestsClient) IndexPutTemplate(ctx context.Context, templateName, bodyString string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexPutTemplate", ctx, templateName, bodyString)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexPutTemplate indicates an expected call of IndexPutTemplate.
func (mr *MockIntegrationTestsClientMockRecorder) IndexPutTemplate(ctx, templateName, bodyString any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexPutTemplate", reflect.TypeOf((*MockIntegrationTestsClient)(nil).IndexPutTemplate), ctx, templateName, bodyString)
}

// IsPointInTimeSupported mocks base method.
func (m *MockIntegrationTestsClient) IsPointInTimeSupported(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPointInTimeSupported", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPointInTimeSupported indicates an expected call of IsPointInTimeSupported.
func (mr *MockIntegrationTestsClientMockRecorder) IsPointInTimeSupported(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPointInTimeSupported", reflect.TypeOf((*MockIntegrationTestsClient)(nil).IsPointInTimeSupported), ctx)
}

// OpenPointInTime mocks base method.
func (m *MockIntegrationTestsClient) OpenPointInTime(ctx context.Context, index, keepAliveInterval string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenPointInTime", ctx, index, keepAliveInterval)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenPointInTime indicates an expected call of OpenPointInTime.
func (mr *MockIntegrationTestsClientMockRecorder) OpenPointInTime(ctx, index, keepAliveInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenPointInTime", reflect.TypeOf((*MockIntegrationTestsClient)(nil).OpenPointInTime), ctx, index, keepAliveInterval)
}

// OpenScroll mocks base method.
func (m *MockIntegrationTestsClient) OpenScroll(ctx context.Context, p *SearchParameters, keepAliveInterval string) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenScroll", ctx, p, keepAliveInterval)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenScroll indicates an expected call of OpenScroll.
func (mr *MockIntegrationTestsClientMockRecorder) OpenScroll(ctx, p, keepAliveInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenScroll", reflect.TypeOf((*MockIntegrationTestsClient)(nil).OpenScroll), ctx, p, keepAliveInterval)
}

// Ping mocks base method.
func (m *MockIntegrationTestsClient) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockIntegrationTestsClientMockRecorder) Ping(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockIntegrationTestsClient)(nil).Ping), ctx)
}

// PutMapping mocks base method.
func (m *MockIntegrationTestsClient) PutMapping(ctx context.Context, index string, mapping map[string]enums.IndexedValueType) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutMapping", ctx, index, mapping)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutMapping indicates an expected call of PutMapping.
func (mr *MockIntegrationTestsClientMockRecorder) PutMapping(ctx, index, mapping any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutMapping", reflect.TypeOf((*MockIntegrationTestsClient)(nil).PutMapping), ctx, index, mapping)
}

// RunBulkProcessor mocks base method.
func (m *MockIntegrationTestsClient) RunBulkProcessor(ctx context.Context, p *BulkProcessorParameters) (BulkProcessor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunBulkProcessor", ctx, p)
	ret0, _ := ret[0].(BulkProcessor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunBulkProcessor indicates an expected call of RunBulkProcessor.
func (mr *MockIntegrationTestsClientMockRecorder) RunBulkProcessor(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunBulkProcessor", reflect.TypeOf((*MockIntegrationTestsClient)(nil).RunBulkProcessor), ctx, p)
}

// Scroll mocks base method.
func (m *MockIntegrationTestsClient) Scroll(ctx context.Context, id, keepAliveInterval string) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scroll", ctx, id, keepAliveInterval)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scroll indicates an expected call of Scroll.
func (mr *MockIntegrationTestsClientMockRecorder) Scroll(ctx, id, keepAliveInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scroll", reflect.TypeOf((*MockIntegrationTestsClient)(nil).Scroll), ctx, id, keepAliveInterval)
}

// Search mocks base method.
func (m *MockIntegrationTestsClient) Search(ctx context.Context, p *SearchParameters) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, p)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockIntegrationTestsClientMockRecorder) Search(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIntegrationTestsClient)(nil).Search), ctx, p)
}

// WaitForYellowStatus mocks base method.
func (m *MockIntegrationTestsClient) WaitForYellowStatus(ctx context.Context, index string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForYellowStatus", ctx, index)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WaitForYellowStatus indicates an expected call of WaitForYellowStatus.
func (mr *MockIntegrationTestsClientMockRecorder) WaitForYellowStatus(ctx, index any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForYellowStatus", reflect.TypeOf((*MockIntegrationTestsClient)(nil).WaitForYellowStatus), ctx, index)
}
