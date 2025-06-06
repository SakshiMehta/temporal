package searchattribute

import (
	"fmt"
	"strings"

	enumspb "go.temporal.io/api/enums/v1"
	persistencespb "go.temporal.io/server/api/persistence/v1"
)

const (
	NamespaceID           = "NamespaceId"
	WorkflowID            = "WorkflowId"
	RunID                 = "RunId"
	WorkflowType          = "WorkflowType"
	StartTime             = "StartTime"
	ExecutionTime         = "ExecutionTime"
	CloseTime             = "CloseTime"
	ExecutionStatus       = "ExecutionStatus"
	TaskQueue             = "TaskQueue"
	HistoryLength         = "HistoryLength"
	ExecutionDuration     = "ExecutionDuration"
	StateTransitionCount  = "StateTransitionCount"
	TemporalChangeVersion = "TemporalChangeVersion"
	BinaryChecksums       = "BinaryChecksums"
	BatcherNamespace      = "BatcherNamespace"
	BatcherUser           = "BatcherUser"
	HistorySizeBytes      = "HistorySizeBytes"
	ParentWorkflowID      = "ParentWorkflowId"
	ParentRunID           = "ParentRunId"
	RootWorkflowID        = "RootWorkflowId"
	RootRunID             = "RootRunId"

	TemporalNamespaceDivision = "TemporalNamespaceDivision"

	// These fields are not in Elasticsearch mappings definition and therefore are not indexed.
	MemoEncoding      = "MemoEncoding"
	Memo              = "Memo"
	VisibilityTaskKey = "VisibilityTaskKey"

	// Added to workflows started by a schedule.
	TemporalScheduledStartTime = "TemporalScheduledStartTime"
	TemporalScheduledById      = "TemporalScheduledById"

	// Used by scheduler workflow.
	TemporalSchedulePaused = "TemporalSchedulePaused"

	ReservedPrefix = "Temporal"

	// Query clause that mentions TemporalNamespaceDivision to disable special handling of that
	// search attribute in visibility.
	matchAnyNamespaceDivision = TemporalNamespaceDivision + ` IS NULL OR ` + TemporalNamespaceDivision + ` IS NOT NULL`

	// A user may specify a ScheduleID in a query even if a ScheduleId search attribute isn't defined for the namespace.
	// In such a case, ScheduleId is effectively a fake search attribute. Of course, a user may optionally choose to
	// define a custom ScheduleId search attribute, in which case the query using the ScheduleId would operate just like
	// any other custom search attribute.
	ScheduleID = "ScheduleId"

	// TemporalPauseInfo is a search attribute that stores the information about paused entities in the workflow.
	// Format of a single paused entity: "<key>:<value>".
	//  * <key> is something that can be used to identify the filtering condition
	//  * <value> is the value of the corresponding filtering condition.
	// examples:
	//   - for paused activities, manual pause, we may have 2 <key>:<value> pairs:
	//     * "Activity:MyCoolActivityType"
	//     * "Reason:ManualActivityPause"
	//     * or
	//     * "Policy:<some policy id>"
	//   - for paused workflows, we may have the following <key>:<value> pairs:
	//     * "Workflow:WorkflowID"
	//     * "Reason:ManualWorkflowPause"
	TemporalPauseInfo = "TemporalPauseInfo"

	// BuildIds is a KeywordList that holds information about current and past build ids
	// used by the workflow. Used for Worker Versioning
	BuildIds = "BuildIds"

	// TemporalWorkerDeploymentVersion stores the current Worker Deployment Version
	// associated with the execution. It is updated at workflow task completion when
	// the SDK says what Version completed the workflow task. It can have a value for
	// unversioned workflows, if they are processed by an unversioned deployment.
	TemporalWorkerDeploymentVersion = "TemporalWorkerDeploymentVersion"

	// TemporalWorkerDeployment stores the current Worker Deployment associated with
	// the execution. It is updated at workflow task completion when the SDK says what
	// Worker Deployment completed the workflow task. It can have a value for
	// unversioned workflows, if they are processed by an unversioned deployment.
	TemporalWorkerDeployment = "TemporalWorkerDeployment"

	// TemporalWorkflowVersioningBehavior stores the current Versioning Behavior of the
	// execution. It is updated at workflow task completion when the server gets the
	// behavior (`auto_upgrade` or `pinned`) from the SDK. Empty for unversioned workflows.
	TemporalWorkflowVersioningBehavior = "TemporalWorkflowVersioningBehavior"
)

var (
	// system are internal search attributes which are passed and stored as separate fields.
	system = map[string]enumspb.IndexedValueType{
		WorkflowID:           enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		RunID:                enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		WorkflowType:         enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		StartTime:            enumspb.INDEXED_VALUE_TYPE_DATETIME,
		ExecutionTime:        enumspb.INDEXED_VALUE_TYPE_DATETIME,
		CloseTime:            enumspb.INDEXED_VALUE_TYPE_DATETIME,
		ExecutionStatus:      enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		TaskQueue:            enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		HistoryLength:        enumspb.INDEXED_VALUE_TYPE_INT,
		ExecutionDuration:    enumspb.INDEXED_VALUE_TYPE_INT,
		StateTransitionCount: enumspb.INDEXED_VALUE_TYPE_INT,
		HistorySizeBytes:     enumspb.INDEXED_VALUE_TYPE_INT,
		ParentWorkflowID:     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		ParentRunID:          enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		RootWorkflowID:       enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		RootRunID:            enumspb.INDEXED_VALUE_TYPE_KEYWORD,
	}

	// predefined are internal search attributes which are passed and stored in SearchAttributes object together with custom search attributes.
	predefined = map[string]enumspb.IndexedValueType{
		TemporalChangeVersion:              enumspb.INDEXED_VALUE_TYPE_KEYWORD_LIST,
		BinaryChecksums:                    enumspb.INDEXED_VALUE_TYPE_KEYWORD_LIST,
		BuildIds:                           enumspb.INDEXED_VALUE_TYPE_KEYWORD_LIST,
		BatcherNamespace:                   enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		BatcherUser:                        enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		TemporalScheduledStartTime:         enumspb.INDEXED_VALUE_TYPE_DATETIME,
		TemporalScheduledById:              enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		TemporalSchedulePaused:             enumspb.INDEXED_VALUE_TYPE_BOOL,
		TemporalNamespaceDivision:          enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		TemporalPauseInfo:                  enumspb.INDEXED_VALUE_TYPE_KEYWORD_LIST,
		TemporalWorkerDeploymentVersion:    enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		TemporalWorkflowVersioningBehavior: enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		TemporalWorkerDeployment:           enumspb.INDEXED_VALUE_TYPE_KEYWORD,
	}

	// reserved are internal field names that can't be used as search attribute names.
	reserved = map[string]struct{}{
		NamespaceID:  {},
		MemoEncoding: {},
		Memo:         {},
		// Used in the Elasticsearch bulk processor, not needed in SQL databases.
		VisibilityTaskKey: {},
	}

	sqlDbSystemNameToColName = map[string]string{
		NamespaceID:          "namespace_id",
		WorkflowID:           "workflow_id",
		RunID:                "run_id",
		WorkflowType:         "workflow_type_name",
		StartTime:            "start_time",
		ExecutionTime:        "execution_time",
		CloseTime:            "close_time",
		ExecutionStatus:      "status",
		TaskQueue:            "task_queue",
		HistoryLength:        "history_length",
		HistorySizeBytes:     "history_size_bytes",
		ExecutionDuration:    "execution_duration",
		StateTransitionCount: "state_transition_count",
		Memo:                 "memo",
		MemoEncoding:         "encoding",
		ParentWorkflowID:     "parent_workflow_id",
		ParentRunID:          "parent_run_id",
		RootWorkflowID:       "root_workflow_id",
		RootRunID:            "root_run_id",
	}

	sqlDbCustomSearchAttributes = map[string]enumspb.IndexedValueType{
		"Bool01":        enumspb.INDEXED_VALUE_TYPE_BOOL,
		"Bool02":        enumspb.INDEXED_VALUE_TYPE_BOOL,
		"Bool03":        enumspb.INDEXED_VALUE_TYPE_BOOL,
		"Datetime01":    enumspb.INDEXED_VALUE_TYPE_DATETIME,
		"Datetime02":    enumspb.INDEXED_VALUE_TYPE_DATETIME,
		"Datetime03":    enumspb.INDEXED_VALUE_TYPE_DATETIME,
		"Double01":      enumspb.INDEXED_VALUE_TYPE_DOUBLE,
		"Double02":      enumspb.INDEXED_VALUE_TYPE_DOUBLE,
		"Double03":      enumspb.INDEXED_VALUE_TYPE_DOUBLE,
		"Int01":         enumspb.INDEXED_VALUE_TYPE_INT,
		"Int02":         enumspb.INDEXED_VALUE_TYPE_INT,
		"Int03":         enumspb.INDEXED_VALUE_TYPE_INT,
		"Keyword01":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Keyword02":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Keyword03":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Keyword04":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Keyword05":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Keyword06":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Keyword07":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Keyword08":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Keyword09":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Keyword10":     enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		"Text01":        enumspb.INDEXED_VALUE_TYPE_TEXT,
		"Text02":        enumspb.INDEXED_VALUE_TYPE_TEXT,
		"Text03":        enumspb.INDEXED_VALUE_TYPE_TEXT,
		"KeywordList01": enumspb.INDEXED_VALUE_TYPE_KEYWORD_LIST,
		"KeywordList02": enumspb.INDEXED_VALUE_TYPE_KEYWORD_LIST,
		"KeywordList03": enumspb.INDEXED_VALUE_TYPE_KEYWORD_LIST,
	}
)

// IsSystem returns true if name is system search attribute
func IsSystem(name string) bool {
	_, ok := system[name]
	return ok
}

// IsReserved returns true if name is system reserved and can't be used as custom search attribute name.
func IsReserved(name string) bool {
	if _, ok := system[name]; ok {
		return true
	}
	if _, ok := predefined[name]; ok {
		return true
	}
	if _, ok := reserved[name]; ok {
		return true
	}
	return strings.HasPrefix(name, ReservedPrefix)
}

// IsMappable returns true if name can have be mapped tho the alias.
func IsMappable(name string) bool {
	if _, ok := system[name]; ok {
		return false
	}
	if _, ok := predefined[name]; ok {
		return false
	}
	return true
}

// GetSqlDbColName maps system and reserved search attributes to column names for SQL tables.
// If the input is not a system or reserved search attribute, then it returns the input.
func GetSqlDbColName(name string) string {
	if fieldName, ok := sqlDbSystemNameToColName[name]; ok {
		return fieldName
	}
	return name
}

func GetSqlDbIndexSearchAttributes() *persistencespb.IndexSearchAttributes {
	return &persistencespb.IndexSearchAttributes{
		CustomSearchAttributes: sqlDbCustomSearchAttributes,
	}
}

// QueryWithAnyNamespaceDivision returns a modified workflow visibility query that disables
// special handling of namespace division and so matches workflows in all namespace divisions.
// Normally a query that didn't explicitly mention TemporalNamespaceDivision would be limited
// to the default (empty string) namespace division.
func QueryWithAnyNamespaceDivision(query string) string {
	if strings.TrimSpace(query) == "" {
		return matchAnyNamespaceDivision
	}
	return fmt.Sprintf(`(%s) AND (%s)`, query, matchAnyNamespaceDivision)
}
