package repository

import (
	"time"
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type AuditLogRepository struct {
	db *bun.DB
}

func NewAuditLogRepository(db *bun.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	_, err := r.db.NewInsert().Model(log).ExcludeColumn("operated_by_name").Exec(ctx)
	return err
}

func (r *AuditLogRepository) GetByID(ctx context.Context, id int64) (*model.AuditLog, error) {
	log := new(model.AuditLog)
	err := r.db.NewSelect().
		Model(log).
		ColumnExpr("audit_log.*").
		ColumnExpr("u.username AS operated_by_name").
		Join("LEFT JOIN users u ON u.id = audit_log.operated_by").
		Where("audit_log.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return log, nil
}

type AuditLogFilter struct {
	TableName  string
	RecordID   *int64
	OperatedBy *int64
	OperatedByName string
	StartTime  *time.Time
	EndTime    *time.Time
	Page       int
	PageSize   int
}

func (r *AuditLogRepository) List(ctx context.Context, filter *AuditLogFilter) ([]model.AuditLog, int, error) {
	var logs []model.AuditLog

	q := r.db.NewSelect().Model(&logs).
		ColumnExpr("audit_log.*").
		ColumnExpr("u.username AS operated_by_name").
		Join("LEFT JOIN users u ON u.id = audit_log.operated_by")

	if filter.TableName != "" {
		q = q.Where("audit_log.table_name = ?", filter.TableName)
	}
	if filter.RecordID != nil {
		q = q.Where("audit_log.record_id = ?", *filter.RecordID)
	}
	if filter.OperatedBy != nil {
	if filter.OperatedByName != "" {
		q = q.Where("u.username LIKE ?", "%"+filter.OperatedByName+"%")
	}
		q = q.Where("audit_log.operated_by = ?", *filter.OperatedBy)
	if filter.OperatedByName != "" {
		q = q.Where("u.username LIKE ?", "%"+filter.OperatedByName+"%")
	}
	}
	if filter.OperatedByName != "" {
		q = q.Where("u.username LIKE ?", "%"+filter.OperatedByName+"%")
	}
	if filter.StartTime != nil {
		q = q.Where("audit_log.operated_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		q = q.Where("audit_log.operated_at <= ?", *filter.EndTime)
	}

	total, err := q.
		Order("audit_log.operated_at DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
