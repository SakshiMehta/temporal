package filestore

import (
	"context"
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	commonpb "go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	workflowpb "go.temporal.io/api/workflow/v1"
	archiverspb "go.temporal.io/server/api/archiver/v1"
	"go.temporal.io/server/common/archiver"
	"go.temporal.io/server/common/codec"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/payload"
	"go.temporal.io/server/common/primitives/timestamp"
	"go.temporal.io/server/common/searchattribute"
	"go.temporal.io/server/common/util"
	"go.temporal.io/server/tests/testutils"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	testWorkflowTypeName = "test-workflow-type"
)

type visibilityArchiverSuite struct {
	*require.Assertions
	suite.Suite

	container          *archiver.VisibilityBootstrapContainer
	testArchivalURI    archiver.URI
	testQueryDirectory string
	visibilityRecords  []*archiverspb.VisibilityRecord

	controller *gomock.Controller
}

func TestVisibilityArchiverSuite(t *testing.T) {
	suite.Run(t, new(visibilityArchiverSuite))
}

func (s *visibilityArchiverSuite) SetupSuite() {
	var err error
	s.testQueryDirectory, err = os.MkdirTemp("", "TestQuery")
	s.Require().NoError(err)
	s.setupVisibilityDirectory()
	s.testArchivalURI, err = archiver.NewURI("file:///a/b/c")
	s.Require().NoError(err)
}

func (s *visibilityArchiverSuite) TearDownSuite() {
	if err := os.RemoveAll(s.testQueryDirectory); err != nil {
		s.Fail("Failed to remove test query directory %v: %v", s.testQueryDirectory, err)
	}
}

func (s *visibilityArchiverSuite) SetupTest() {
	s.Assertions = require.New(s.T())
	s.container = &archiver.VisibilityBootstrapContainer{
		Logger: log.NewNoopLogger(),
	}
	s.controller = gomock.NewController(s.T())
}

func (s *visibilityArchiverSuite) TearDownTest() {
	s.controller.Finish()
}

func (s *visibilityArchiverSuite) TestValidateURI() {
	testCases := []struct {
		URI         string
		expectedErr error
	}{
		{
			URI:         "wrongscheme:///a/b/c",
			expectedErr: archiver.ErrURISchemeMismatch,
		},
		{
			URI:         "file://",
			expectedErr: errEmptyDirectoryPath,
		},
		{
			URI:         "file:///a/b/c",
			expectedErr: nil,
		},
	}

	visibilityArchiver := s.newTestVisibilityArchiver()
	for _, tc := range testCases {
		URI, err := archiver.NewURI(tc.URI)
		s.NoError(err)
		s.Equal(tc.expectedErr, visibilityArchiver.ValidateURI(URI))
	}
}

func (s *visibilityArchiverSuite) TestArchive_Fail_InvalidURI() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	URI, err := archiver.NewURI("wrongscheme://")
	s.NoError(err)
	request := &archiverspb.VisibilityRecord{
		Namespace:        testNamespace,
		NamespaceId:      testNamespaceID,
		WorkflowId:       testWorkflowID,
		RunId:            testRunID,
		WorkflowTypeName: testWorkflowTypeName,
		StartTime:        timestamp.TimeNowPtrUtc(),
		ExecutionTime:    nil, // workflow without backoff
		CloseTime:        timestamp.TimeNowPtrUtc(),
		Status:           enumspb.WORKFLOW_EXECUTION_STATUS_FAILED,
		HistoryLength:    int64(101),
	}
	err = visibilityArchiver.Archive(context.Background(), URI, request)
	s.Error(err)
}

func (s *visibilityArchiverSuite) TestArchive_Fail_InvalidRequest() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	err := visibilityArchiver.Archive(context.Background(), s.testArchivalURI, &archiverspb.VisibilityRecord{})
	s.Error(err)
}

func (s *visibilityArchiverSuite) TestArchive_Fail_NonRetryableErrorOption() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	nonRetryableErr := errors.New("some non-retryable error")
	err := visibilityArchiver.Archive(
		context.Background(),
		s.testArchivalURI,
		&archiverspb.VisibilityRecord{},
		archiver.GetNonRetryableErrorOption(nonRetryableErr),
	)
	s.Equal(nonRetryableErr, err)
}

func (s *visibilityArchiverSuite) TestArchive_Success() {
	dir := testutils.MkdirTemp(s.T(), "", "TestVisibilityArchive")

	visibilityArchiver := s.newTestVisibilityArchiver()
	closeTimestamp := timestamp.TimeNowPtrUtc()
	request := &archiverspb.VisibilityRecord{
		NamespaceId:      testNamespaceID,
		Namespace:        testNamespace,
		WorkflowId:       testWorkflowID,
		RunId:            testRunID,
		WorkflowTypeName: testWorkflowTypeName,
		StartTime:        timestamppb.New(closeTimestamp.AsTime().Add(-time.Hour)),
		ExecutionTime:    nil, // workflow without backoff
		CloseTime:        closeTimestamp,
		Status:           enumspb.WORKFLOW_EXECUTION_STATUS_FAILED,
		HistoryLength:    int64(101),
		Memo: &commonpb.Memo{
			Fields: map[string]*commonpb.Payload{
				"testFields": payload.EncodeBytes([]byte{1, 2, 3}),
			},
		},
		SearchAttributes: map[string]string{
			"testAttribute": "456",
		},
	}
	URI, err := archiver.NewURI("file://" + dir)
	s.NoError(err)
	err = visibilityArchiver.Archive(context.Background(), URI, request)
	s.NoError(err)

	expectedFilename := constructVisibilityFilename(closeTimestamp.AsTime(), testRunID)
	filepath := path.Join(dir, testNamespaceID, expectedFilename)
	s.assertFileExists(filepath)

	data, err := readFile(filepath)
	s.NoError(err)

	archivedRecord := &archiverspb.VisibilityRecord{}
	encoder := codec.NewJSONPBEncoder()
	err = encoder.Decode(data, archivedRecord)
	s.NoError(err)
	s.Equal(request, archivedRecord)
}

func (s *visibilityArchiverSuite) TestMatchQuery() {
	testCases := []struct {
		query       *parsedQuery
		record      *archiverspb.VisibilityRecord
		shouldMatch bool
	}{
		{
			query: &parsedQuery{
				earliestCloseTime: time.Unix(0, 1000),
				latestCloseTime:   time.Unix(0, 12345),
			},
			record: &archiverspb.VisibilityRecord{
				CloseTime: timestamp.UnixOrZeroTimePtr(1999),
			},
			shouldMatch: true,
		},
		{
			query: &parsedQuery{
				earliestCloseTime: time.Unix(0, 1000),
				latestCloseTime:   time.Unix(0, 12345),
			},
			record: &archiverspb.VisibilityRecord{
				CloseTime: timestamp.UnixOrZeroTimePtr(999),
			},
			shouldMatch: false,
		},
		{
			query: &parsedQuery{
				earliestCloseTime: time.Unix(0, 1000),
				latestCloseTime:   time.Unix(0, 12345),
				workflowID:        util.Ptr("random workflowID"),
			},
			record: &archiverspb.VisibilityRecord{
				CloseTime: timestamp.UnixOrZeroTimePtr(2000),
			},
			shouldMatch: false,
		},
		{
			query: &parsedQuery{
				earliestCloseTime: time.Unix(0, 1000),
				latestCloseTime:   time.Unix(0, 12345),
				workflowID:        util.Ptr("random workflowID"),
				runID:             util.Ptr("random runID"),
			},
			record: &archiverspb.VisibilityRecord{
				CloseTime:        timestamp.UnixOrZeroTimePtr(12345),
				WorkflowId:       "random workflowID",
				RunId:            "random runID",
				WorkflowTypeName: "random type name",
			},
			shouldMatch: true,
		},
		{
			query: &parsedQuery{
				earliestCloseTime: time.Unix(0, 1000),
				latestCloseTime:   time.Unix(0, 12345),
				workflowTypeName:  util.Ptr("some random type name"),
			},
			record: &archiverspb.VisibilityRecord{
				CloseTime: timestamp.UnixOrZeroTimePtr(12345),
			},
			shouldMatch: false,
		},
		{
			query: &parsedQuery{
				earliestCloseTime: time.Unix(0, 1000),
				latestCloseTime:   time.Unix(0, 12345),
				workflowTypeName:  util.Ptr("some random type name"),
				status:            toWorkflowExecutionStatusPtr(enumspb.WORKFLOW_EXECUTION_STATUS_CONTINUED_AS_NEW),
			},
			record: &archiverspb.VisibilityRecord{
				CloseTime:        timestamp.UnixOrZeroTimePtr(12345),
				Status:           enumspb.WORKFLOW_EXECUTION_STATUS_CONTINUED_AS_NEW,
				WorkflowTypeName: "some random type name",
			},
			shouldMatch: true,
		},
	}

	for _, tc := range testCases {
		s.Equal(tc.shouldMatch, matchQuery(tc.record, tc.query))
	}
}

func (s *visibilityArchiverSuite) TestSortAndFilterFiles() {
	testCases := []struct {
		filenames      []string
		token          *queryVisibilityToken
		expectedResult []string
	}{
		{
			filenames:      []string{"9_12345.vis", "5_0.vis", "9_54321.vis", "1000_654.vis", "1000_78.vis"},
			expectedResult: []string{"1000_78.vis", "1000_654.vis", "9_54321.vis", "9_12345.vis", "5_0.vis"},
		},
		{
			filenames: []string{"9_12345.vis", "5_0.vis", "9_54321.vis", "1000_654.vis", "1000_78.vis"},
			token: &queryVisibilityToken{
				LastCloseTime: time.Unix(0, 3),
			},
			expectedResult: []string{},
		},
		{
			filenames: []string{"9_12345.vis", "5_0.vis", "9_54321.vis", "1000_654.vis", "1000_78.vis"},
			token: &queryVisibilityToken{
				LastCloseTime: time.Unix(0, 999),
			},
			expectedResult: []string{"9_54321.vis", "9_12345.vis", "5_0.vis"},
		},
		{
			filenames: []string{"9_12345.vis", "5_0.vis", "9_54321.vis", "1000_654.vis", "1000_78.vis"},
			token: &queryVisibilityToken{
				LastCloseTime: time.Unix(0, 5).UTC(),
			},
			expectedResult: []string{"5_0.vis"},
		},
	}

	for i, tc := range testCases {
		result, err := sortAndFilterFiles(tc.filenames, tc.token)
		s.NoError(err, "case %d", i)
		s.Equal(tc.expectedResult, result, "case %d", i)
	}
}

func (s *visibilityArchiverSuite) TestQuery_Fail_InvalidURI() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	URI, err := archiver.NewURI("wrongscheme://")
	s.NoError(err)
	request := &archiver.QueryVisibilityRequest{
		NamespaceID: testNamespaceID,
		PageSize:    1,
	}
	response, err := visibilityArchiver.Query(context.Background(), URI, request, searchattribute.TestNameTypeMap)
	s.Error(err)
	s.Nil(response)
}

func (s *visibilityArchiverSuite) TestQuery_Fail_InvalidRequest() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	response, err := visibilityArchiver.Query(context.Background(), s.testArchivalURI, &archiver.QueryVisibilityRequest{}, searchattribute.TestNameTypeMap)
	s.Error(err)
	s.Nil(response)
}

func (s *visibilityArchiverSuite) TestQuery_Fail_InvalidQuery() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	mockParser := NewMockQueryParser(s.controller)
	mockParser.EXPECT().Parse(gomock.Any()).Return(nil, errors.New("invalid query"))
	visibilityArchiver.queryParser = mockParser
	response, err := visibilityArchiver.Query(context.Background(), s.testArchivalURI, &archiver.QueryVisibilityRequest{
		NamespaceID: "some random namespaceID",
		PageSize:    10,
		Query:       "some invalid query",
	}, searchattribute.TestNameTypeMap)
	s.Error(err)
	s.Nil(response)
}

func (s *visibilityArchiverSuite) TestQuery_Success_DirectoryNotExist() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	mockParser := NewMockQueryParser(s.controller)
	mockParser.EXPECT().Parse(gomock.Any()).Return(&parsedQuery{
		earliestCloseTime: time.Unix(0, 1),
		latestCloseTime:   time.Unix(0, 101),
	}, nil)
	visibilityArchiver.queryParser = mockParser
	request := &archiver.QueryVisibilityRequest{
		NamespaceID: testNamespaceID,
		Query:       "parsed by mockParser",
		PageSize:    1,
	}
	response, err := visibilityArchiver.Query(context.Background(), s.testArchivalURI, request, searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.NotNil(response)
	s.Empty(response.Executions)
	s.Empty(response.NextPageToken)
}

func (s *visibilityArchiverSuite) TestQuery_Fail_InvalidToken() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	mockParser := NewMockQueryParser(s.controller)
	mockParser.EXPECT().Parse(gomock.Any()).Return(&parsedQuery{
		earliestCloseTime: time.Unix(0, 1),
		latestCloseTime:   time.Unix(0, 101),
	}, nil)
	visibilityArchiver.queryParser = mockParser
	request := &archiver.QueryVisibilityRequest{
		NamespaceID:   testNamespaceID,
		Query:         "parsed by mockParser",
		PageSize:      1,
		NextPageToken: []byte{1, 2, 3},
	}
	response, err := visibilityArchiver.Query(context.Background(), s.testArchivalURI, request, searchattribute.TestNameTypeMap)
	s.Error(err)
	s.Nil(response)
}

func (s *visibilityArchiverSuite) TestQuery_Success_NoNextPageToken() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	mockParser := NewMockQueryParser(s.controller)
	mockParser.EXPECT().Parse(gomock.Any()).Return(&parsedQuery{
		earliestCloseTime: time.Unix(0, 1),
		latestCloseTime:   time.Unix(0, 10001),
		workflowID:        util.Ptr(testWorkflowID),
	}, nil)
	visibilityArchiver.queryParser = mockParser
	request := &archiver.QueryVisibilityRequest{
		NamespaceID: testNamespaceID,
		PageSize:    10,
		Query:       "parsed by mockParser",
	}
	URI, err := archiver.NewURI("file://" + s.testQueryDirectory)
	s.NoError(err)
	response, err := visibilityArchiver.Query(context.Background(), URI, request, searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.NotNil(response)
	s.Nil(response.NextPageToken)
	s.Len(response.Executions, 1)
	ei, err := convertToExecutionInfo(s.visibilityRecords[0], searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.Equal(ei, response.Executions[0])
}

func (s *visibilityArchiverSuite) TestQuery_Success_SmallPageSize() {
	visibilityArchiver := s.newTestVisibilityArchiver()
	mockParser := NewMockQueryParser(s.controller)
	mockParser.EXPECT().Parse(gomock.Any()).Return(&parsedQuery{
		earliestCloseTime: time.Unix(0, 1),
		latestCloseTime:   time.Unix(0, 10001),
		status:            toWorkflowExecutionStatusPtr(enumspb.WORKFLOW_EXECUTION_STATUS_FAILED),
	}, nil).AnyTimes()
	visibilityArchiver.queryParser = mockParser
	request := &archiver.QueryVisibilityRequest{
		NamespaceID: testNamespaceID,
		PageSize:    2,
		Query:       "parsed by mockParser",
	}
	URI, err := archiver.NewURI("file://" + s.testQueryDirectory)
	s.NoError(err)
	response, err := visibilityArchiver.Query(context.Background(), URI, request, searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.NotNil(response)
	s.NotNil(response.NextPageToken)
	s.Len(response.Executions, 2)
	ei, err := convertToExecutionInfo(s.visibilityRecords[0], searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.Equal(ei, response.Executions[0])
	ei, err = convertToExecutionInfo(s.visibilityRecords[1], searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.Equal(ei, response.Executions[1])

	request.NextPageToken = response.NextPageToken
	response, err = visibilityArchiver.Query(context.Background(), URI, request, searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.NotNil(response)
	s.Nil(response.NextPageToken)
	s.Len(response.Executions, 1)
	ei, err = convertToExecutionInfo(s.visibilityRecords[3], searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.Equal(ei, response.Executions[0])
}

func (s *visibilityArchiverSuite) TestArchiveAndQuery() {
	dir := testutils.MkdirTemp(s.T(), "", "TestArchiveAndQuery")

	visibilityArchiver := s.newTestVisibilityArchiver()
	mockParser := NewMockQueryParser(s.controller)
	mockParser.EXPECT().Parse(gomock.Any()).Return(&parsedQuery{
		earliestCloseTime: time.Unix(0, 10),
		latestCloseTime:   time.Unix(0, 10001),
		status:            toWorkflowExecutionStatusPtr(enumspb.WORKFLOW_EXECUTION_STATUS_FAILED),
	}, nil).AnyTimes()
	visibilityArchiver.queryParser = mockParser
	URI, err := archiver.NewURI("file://" + dir)
	s.NoError(err)
	for _, record := range s.visibilityRecords {
		err := visibilityArchiver.Archive(context.Background(), URI, (*archiverspb.VisibilityRecord)(record))
		s.NoError(err)
	}

	request := &archiver.QueryVisibilityRequest{
		NamespaceID: testNamespaceID,
		PageSize:    1,
		Query:       "parsed by mockParser",
	}
	executions := []*workflowpb.WorkflowExecutionInfo{}
	for len(executions) == 0 || request.NextPageToken != nil {
		response, err := visibilityArchiver.Query(context.Background(), URI, request, searchattribute.TestNameTypeMap)
		s.NoError(err)
		s.NotNil(response)
		executions = append(executions, response.Executions...)
		request.NextPageToken = response.NextPageToken
	}
	s.Len(executions, 2)
	ei, err := convertToExecutionInfo(s.visibilityRecords[0], searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.Equal(ei, executions[0])
	ei, err = convertToExecutionInfo(s.visibilityRecords[1], searchattribute.TestNameTypeMap)
	s.NoError(err)
	s.Equal(ei, executions[1])
}

func (s *visibilityArchiverSuite) TestQuery_EmptyQuery_InvalidNamespace() {
	URI := s.testArchivalURI

	visibilityArchiver := s.newTestVisibilityArchiver()
	mockParser := NewMockQueryParser(s.controller)
	mockParser.EXPECT().Parse(gomock.Any()).Return(&parsedQuery{
		earliestCloseTime: time.Unix(0, 10),
		latestCloseTime:   time.Unix(0, 10001),
		status:            toWorkflowExecutionStatusPtr(enumspb.WORKFLOW_EXECUTION_STATUS_FAILED),
	}, nil).AnyTimes()
	visibilityArchiver.queryParser = mockParser
	req := &archiver.QueryVisibilityRequest{
		NamespaceID:   "",
		PageSize:      1,
		NextPageToken: nil,
		Query:         "",
	}
	_, err := visibilityArchiver.Query(context.Background(), URI, req, searchattribute.TestNameTypeMap)

	var svcErr *serviceerror.InvalidArgument

	s.ErrorAs(err, &svcErr)
}

func (s *visibilityArchiverSuite) TestQuery_EmptyQuery_ZeroPageSize() {
	visibilityArchiver := s.newTestVisibilityArchiver()

	req := &archiver.QueryVisibilityRequest{
		NamespaceID:   testNamespaceID,
		PageSize:      0,
		NextPageToken: nil,
		Query:         "",
	}
	_, err := visibilityArchiver.Query(context.Background(), s.testArchivalURI, req, searchattribute.TestNameTypeMap)

	var svcErr *serviceerror.InvalidArgument

	s.ErrorAs(err, &svcErr)
}

func (s *visibilityArchiverSuite) TestQuery_EmptyQuery_Pagination() {
	dir := testutils.MkdirTemp(s.T(), "", "TestQuery_EmptyQuery_Pagination")

	visibilityArchiver := s.newTestVisibilityArchiver()
	URI, err := archiver.NewURI("file://" + dir)
	s.NoError(err)
	for _, record := range s.visibilityRecords {
		err := visibilityArchiver.Archive(context.Background(), URI, record)
		s.NoError(err)
	}

	request := &archiver.QueryVisibilityRequest{
		NamespaceID: testNamespaceID,
		PageSize:    1,
		Query:       "",
	}
	var executions []*workflowpb.WorkflowExecutionInfo
	for len(executions) == 0 || request.NextPageToken != nil {
		response, err := visibilityArchiver.Query(context.Background(), URI, request, searchattribute.TestNameTypeMap)
		s.NoError(err)
		s.NotNil(response)
		executions = append(executions, response.Executions...)
		request.NextPageToken = response.NextPageToken
	}
	s.Len(executions, 4)
}

func (s *visibilityArchiverSuite) newTestVisibilityArchiver() *visibilityArchiver {
	config := &config.FilestoreArchiver{
		FileMode: testFileModeStr,
		DirMode:  testDirModeStr,
	}
	archiver, err := NewVisibilityArchiver(s.container, config)
	s.NoError(err)
	return archiver.(*visibilityArchiver)
}

func (s *visibilityArchiverSuite) setupVisibilityDirectory() {
	s.visibilityRecords = []*archiverspb.VisibilityRecord{
		{
			NamespaceId:      testNamespaceID,
			Namespace:        testNamespace,
			WorkflowId:       testWorkflowID,
			RunId:            testRunID,
			WorkflowTypeName: testWorkflowTypeName,
			StartTime:        timestamp.UnixOrZeroTimePtr(1),
			CloseTime:        timestamp.UnixOrZeroTimePtr(10000),
			Status:           enumspb.WORKFLOW_EXECUTION_STATUS_FAILED,
			HistoryLength:    101,
		},
		{
			NamespaceId:      testNamespaceID,
			Namespace:        testNamespace,
			WorkflowId:       "some random workflow ID",
			RunId:            "some random run ID",
			WorkflowTypeName: testWorkflowTypeName,
			StartTime:        timestamp.UnixOrZeroTimePtr(2),
			ExecutionTime:    nil,
			CloseTime:        timestamp.UnixOrZeroTimePtr(1000),
			Status:           enumspb.WORKFLOW_EXECUTION_STATUS_FAILED,
			HistoryLength:    123,
		},
		{
			NamespaceId:      testNamespaceID,
			Namespace:        testNamespace,
			WorkflowId:       "another workflow ID",
			RunId:            "another run ID",
			WorkflowTypeName: testWorkflowTypeName,
			StartTime:        timestamp.UnixOrZeroTimePtr(3),
			ExecutionTime:    nil,
			CloseTime:        timestamp.UnixOrZeroTimePtr(10),
			Status:           enumspb.WORKFLOW_EXECUTION_STATUS_CONTINUED_AS_NEW,
			HistoryLength:    456,
		},
		{
			NamespaceId:      testNamespaceID,
			Namespace:        testNamespace,
			WorkflowId:       "and another workflow ID",
			RunId:            "and another run ID",
			WorkflowTypeName: testWorkflowTypeName,
			StartTime:        timestamp.UnixOrZeroTimePtr(3),
			ExecutionTime:    nil,
			CloseTime:        timestamp.UnixOrZeroTimePtr(5),
			Status:           enumspb.WORKFLOW_EXECUTION_STATUS_FAILED,
			HistoryLength:    456,
		},
		{
			NamespaceId:      "some random namespace ID",
			Namespace:        "some random namespace name",
			WorkflowId:       "another workflow ID",
			RunId:            "another run ID",
			WorkflowTypeName: testWorkflowTypeName,
			StartTime:        timestamp.UnixOrZeroTimePtr(3),
			ExecutionTime:    nil,
			CloseTime:        timestamp.UnixOrZeroTimePtr(10000),
			Status:           enumspb.WORKFLOW_EXECUTION_STATUS_CONTINUED_AS_NEW,
			HistoryLength:    456,
		},
	}

	for _, record := range s.visibilityRecords {
		s.writeVisibilityRecordForQueryTest(record)
	}
}

func (s *visibilityArchiverSuite) writeVisibilityRecordForQueryTest(record *archiverspb.VisibilityRecord) {
	data, err := encode(record)
	s.Require().NoError(err)
	filename := constructVisibilityFilename(record.CloseTime.AsTime(), record.GetRunId())
	s.Require().NoError(os.MkdirAll(path.Join(s.testQueryDirectory, record.GetNamespaceId()), testDirMode))
	err = writeFile(path.Join(s.testQueryDirectory, record.GetNamespaceId(), filename), data, testFileMode)
	s.Require().NoError(err)
}

func (s *visibilityArchiverSuite) assertFileExists(filepath string) {
	exists, err := fileExists(filepath)
	s.NoError(err)
	s.True(exists)
}
