package model

import (
	"time"
)

type AuditLog struct {
	bun             struct{}       `bun:"table:audit_logs"`
	ID              int64          `bun:"id,pk,autoincrement" json:"id"`
	TableName       string         `bun:"table_name,notnull" json:"table_name"`
	RecordID        int64          `bun:"record_id,notnull" json:"record_id"`
	Action          string         `bun:"action,notnull" json:"action"`
	OldValue        map[string]any `bun:"old_value,type:json" json:"old_value"`
	NewValue        map[string]any `bun:"new_value,type:json" json:"new_value"`
	OperatedBy      int64          `bun:"operated_by,notnull" json:"operated_by"`
	OperatedByName  string         `bun:"operated_by_name,type:varchar" json:"operated_by_name"`
	OperatedAt      time.Time      `bun:"operated_at,notnull" json:"operated_at"`
	IPAddress       string         `bun:"ip_address" json:"ip_address"`
}

const (
	AuditActionCreate = "create"
	AuditActionUpdate = "update"
	AuditActionDelete = "delete"
)
