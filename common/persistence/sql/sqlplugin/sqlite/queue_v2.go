package sqlite

import (
	"context"
	"database/sql"

	"go.temporal.io/server/common/persistence/sql/sqlplugin"
)

const (
	templateEnqueueMessageQueryV2      = `INSERT INTO queue_messages (queue_type, queue_name, queue_partition, message_id, message_payload, message_encoding) VALUES(:queue_type, :queue_name, :queue_partition, :message_id, :message_payload, :message_encoding)`
	templateGetMessagesQueryV2         = `SELECT message_id, message_payload, message_encoding FROM queue_messages WHERE queue_type = ? and queue_name = ? and queue_partition = ? and message_id >= ? ORDER BY message_id ASC LIMIT ?`
	templateRangeDeleteMessagesQueryV2 = `DELETE FROM queue_messages WHERE queue_type = ? and queue_name = ? and queue_partition = ? and message_id >= ? and message_id <= ?`
	templateGetLastMessageIDQueryV2    = `SELECT message_id FROM queue_messages WHERE queue_type=? and queue_name=? and queue_partition=? ORDER BY message_id DESC LIMIT 1`
	templateCreateQueueMetadataQueryV2 = `INSERT INTO queues (queue_type, queue_name, metadata_payload, metadata_encoding) VALUES(:queue_type, :queue_name, :metadata_payload, :metadata_encoding)`
	templateUpdateQueueMetadataQueryV2 = `UPDATE queues SET metadata_payload = :metadata_payload, metadata_encoding = :metadata_encoding WHERE queue_type = :queue_type and queue_name = :queue_name`
	templateGetQueueMetadataQueryV2    = `SELECT metadata_payload, metadata_encoding from queues WHERE queue_type=? and queue_name=?`
	templateGetNameFromQueueMetadataV2 = `SELECT queue_type, queue_name, metadata_payload, metadata_encoding from queues WHERE queue_type=? LIMIT ? OFFSET ?`
)

func (sdb *db) InsertIntoQueueV2Metadata(ctx context.Context, row *sqlplugin.QueueV2MetadataRow) (sql.Result, error) {
	return sdb.conn.NamedExecContext(ctx,
		templateCreateQueueMetadataQueryV2,
		row,
	)
}
func (sdb *db) UpdateQueueV2Metadata(ctx context.Context, row *sqlplugin.QueueV2MetadataRow) (sql.Result, error) {
	return sdb.conn.NamedExecContext(ctx,
		templateUpdateQueueMetadataQueryV2,
		row,
	)
}
func (sdb *db) SelectFromQueueV2Metadata(ctx context.Context, filter sqlplugin.QueueV2MetadataFilter) (*sqlplugin.QueueV2MetadataRow, error) {
	var row sqlplugin.QueueV2MetadataRow
	err := sdb.conn.GetContext(ctx,
		&row,
		templateGetQueueMetadataQueryV2,
		filter.QueueType,
		filter.QueueName,
	)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
func (sdb *db) SelectFromQueueV2MetadataForUpdate(ctx context.Context, filter sqlplugin.QueueV2MetadataFilter) (*sqlplugin.QueueV2MetadataRow, error) {
	// sqlite does not have FOR UPDATE clause. Calling SelectFromQueueV2Metadata() itself.
	return sdb.SelectFromQueueV2Metadata(ctx, filter)
}
func (sdb *db) SelectNameFromQueueV2Metadata(ctx context.Context, filter sqlplugin.QueueV2MetadataTypeFilter) ([]sqlplugin.QueueV2MetadataRow, error) {
	var rows []sqlplugin.QueueV2MetadataRow
	err := sdb.conn.SelectContext(ctx,
		&rows,
		templateGetNameFromQueueMetadataV2,
		filter.QueueType,
		filter.PageSize,
		filter.PageOffset,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
func (sdb *db) InsertIntoQueueV2Messages(ctx context.Context, row []sqlplugin.QueueV2MessageRow) (sql.Result, error) {
	return sdb.conn.NamedExecContext(ctx,
		templateEnqueueMessageQueryV2,
		row,
	)
}
func (sdb *db) RangeSelectFromQueueV2Messages(ctx context.Context, filter sqlplugin.QueueV2MessagesFilter) ([]sqlplugin.QueueV2MessageRow, error) {
	var rows []sqlplugin.QueueV2MessageRow
	err := sdb.conn.SelectContext(ctx,
		&rows,
		templateGetMessagesQueryV2,
		filter.QueueType,
		filter.QueueName,
		filter.Partition,
		filter.MinMessageID,
		filter.PageSize,
	)
	return rows, err
}
func (sdb *db) RangeDeleteFromQueueV2Messages(ctx context.Context, filter sqlplugin.QueueV2MessagesFilter) (sql.Result, error) {
	return sdb.conn.ExecContext(ctx,
		templateRangeDeleteMessagesQueryV2,
		filter.QueueType,
		filter.QueueName,
		filter.Partition,
		filter.MinMessageID,
		filter.MaxMessageID,
	)
}

func (sdb *db) GetLastEnqueuedMessageIDForUpdateV2(ctx context.Context, filter sqlplugin.QueueV2Filter) (int64, error) {
	var lastMessageID int64
	err := sdb.conn.GetContext(ctx,
		&lastMessageID,
		templateGetLastMessageIDQueryV2,
		filter.QueueType,
		filter.QueueName,
		filter.Partition,
	)
	return lastMessageID, err
}
