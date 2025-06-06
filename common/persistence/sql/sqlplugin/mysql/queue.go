package mysql

import (
	"context"
	"database/sql"

	"go.temporal.io/server/common/persistence"
	"go.temporal.io/server/common/persistence/sql/sqlplugin"
)

const (
	templateEnqueueMessageQuery      = `INSERT INTO queue (queue_type, message_id, message_payload, message_encoding) VALUES(:queue_type, :message_id, :message_payload, :message_encoding)`
	templateGetMessageQuery          = `SELECT message_id, message_payload, message_encoding FROM queue WHERE queue_type = ? and message_id = ?`
	templateGetMessagesQuery         = `SELECT message_id, message_payload, message_encoding FROM queue WHERE queue_type = ? and message_id > ? and message_id <= ? ORDER BY message_id ASC LIMIT ?`
	templateDeleteMessageQuery       = `DELETE FROM queue WHERE queue_type = ? and message_id = ?`
	templateRangeDeleteMessagesQuery = `DELETE FROM queue WHERE queue_type = ? and message_id > ? and message_id <= ?`

	// Note that even though this query takes a range lock that serializes all writes, it will return multiple rows
	// whenever more than one enqueue-er blocks. This is why we max().
	templateGetLastMessageIDQuery = `SELECT MAX(message_id) FROM queue WHERE message_id >= (SELECT message_id FROM queue WHERE queue_type=? ORDER BY message_id DESC LIMIT 1) FOR UPDATE`

	templateCreateQueueMetadataQuery = `INSERT INTO queue_metadata (queue_type, data, data_encoding, version) VALUES(:queue_type, :data, :data_encoding, :version)`
	templateUpdateQueueMetadataQuery = `UPDATE queue_metadata SET data = :data, data_encoding = :data_encoding, version= :version+1 WHERE queue_type = :queue_type and version = :version`
	templateGetQueueMetadataQuery    = `SELECT data, data_encoding, version from queue_metadata WHERE queue_type = ?`
	templateLockQueueMetadataQuery   = templateGetQueueMetadataQuery + " FOR UPDATE"
)

// InsertIntoMessages inserts a new row into queue table
func (mdb *db) InsertIntoMessages(
	ctx context.Context,
	row []sqlplugin.QueueMessageRow,
) (sql.Result, error) {
	return mdb.NamedExecContext(ctx,
		templateEnqueueMessageQuery,
		row,
	)
}

func (mdb *db) SelectFromMessages(
	ctx context.Context,
	filter sqlplugin.QueueMessagesFilter,
) ([]sqlplugin.QueueMessageRow, error) {
	var rows []sqlplugin.QueueMessageRow
	err := mdb.SelectContext(ctx,
		&rows,
		templateGetMessageQuery,
		filter.QueueType,
		filter.MessageID,
	)
	return rows, err
}

func (mdb *db) RangeSelectFromMessages(
	ctx context.Context,
	filter sqlplugin.QueueMessagesRangeFilter,
) ([]sqlplugin.QueueMessageRow, error) {
	var rows []sqlplugin.QueueMessageRow
	err := mdb.SelectContext(ctx,
		&rows,
		templateGetMessagesQuery,
		filter.QueueType,
		filter.MinMessageID,
		filter.MaxMessageID,
		filter.PageSize,
	)
	return rows, err
}

// DeleteFromMessages deletes message with a messageID from the queue
func (mdb *db) DeleteFromMessages(
	ctx context.Context,
	filter sqlplugin.QueueMessagesFilter,
) (sql.Result, error) {
	return mdb.ExecContext(ctx,
		templateDeleteMessageQuery,
		filter.QueueType,
		filter.MessageID,
	)
}

// RangeDeleteFromMessages deletes messages before messageID from the queue
func (mdb *db) RangeDeleteFromMessages(
	ctx context.Context,
	filter sqlplugin.QueueMessagesRangeFilter,
) (sql.Result, error) {
	return mdb.ExecContext(ctx,
		templateRangeDeleteMessagesQuery,
		filter.QueueType,
		filter.MinMessageID,
		filter.MaxMessageID,
	)
}

// GetLastEnqueuedMessageIDForUpdate returns the last enqueued message ID
func (mdb *db) GetLastEnqueuedMessageIDForUpdate(
	ctx context.Context,
	queueType persistence.QueueType,
) (int64, error) {
	var lastMessageID *int64
	err := mdb.GetContext(ctx,
		&lastMessageID,
		templateGetLastMessageIDQuery,
		queueType,
	)
	if lastMessageID == nil {
		// The layer of code above us expects ErrNoRows when the queue is empty. MAX() yields
		// null when the queue is empty, so we need to turn that into the correct error.
		return 0, sql.ErrNoRows
	} else {
		return *lastMessageID, err
	}
}

func (mdb *db) InsertIntoQueueMetadata(
	ctx context.Context,
	row *sqlplugin.QueueMetadataRow,
) (sql.Result, error) {
	return mdb.NamedExecContext(ctx,
		templateCreateQueueMetadataQuery,
		row,
	)
}

func (mdb *db) UpdateQueueMetadata(
	ctx context.Context,
	row *sqlplugin.QueueMetadataRow,
) (sql.Result, error) {
	return mdb.NamedExecContext(ctx,
		templateUpdateQueueMetadataQuery,
		row,
	)
}

func (mdb *db) SelectFromQueueMetadata(
	ctx context.Context,
	filter sqlplugin.QueueMetadataFilter,
) (*sqlplugin.QueueMetadataRow, error) {
	var row sqlplugin.QueueMetadataRow
	err := mdb.GetContext(ctx,
		&row,
		templateGetQueueMetadataQuery,
		filter.QueueType,
	)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (mdb *db) LockQueueMetadata(
	ctx context.Context,
	filter sqlplugin.QueueMetadataFilter,
) (*sqlplugin.QueueMetadataRow, error) {
	var row sqlplugin.QueueMetadataRow
	err := mdb.GetContext(ctx,
		&row,
		templateLockQueueMetadataQuery,
		filter.QueueType,
	)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
