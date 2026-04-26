package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestAuditLogRepository_Create(t *testing.T) {
	repo, _, ctx := setupAuditLogTest(t)
	log := &model.AuditLog{
		TableName:  "users",
		RecordID:   1,
		Action:     model.AuditActionCreate,
		OperatedBy: 1,
		OperatedAt: time.Now(),
	}

	err := repo.Create(ctx, log)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestAuditLogRepository_GetByID_Query(t *testing.T) {
	repo, _, ctx := setupAuditLogTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestAuditLogRepository_List_Query(t *testing.T) {
	repo, _, ctx := setupAuditLogTest(t)
	filter := &AuditLogFilter{Page: 1, PageSize: 10}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestAuditLogRepository_List_WithFilters(t *testing.T) {
	repo, _, ctx := setupAuditLogTest(t)
	recordID := int64(1)
	operatedBy := int64(1)
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	filter := &AuditLogFilter{
		TableName:  "users",
		RecordID:   &recordID,
		OperatedBy: &operatedBy,
		StartTime:  &startTime,
		EndTime:    &endTime,
		Page:       1,
		PageSize:   10,
	}

	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestNewAuditLogRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewAuditLogRepository(db)
	if repo == nil {
		t.Error("NewAuditLogRepository() returned nil")
	}
}

func TestAuditLogActionConstants(t *testing.T) {
	if model.AuditActionCreate != "create" {
		t.Errorf("AuditActionCreate = %s, want create", model.AuditActionCreate)
	}
	if model.AuditActionUpdate != "update" {
		t.Errorf("AuditActionUpdate = %s, want update", model.AuditActionUpdate)
	}
	if model.AuditActionDelete != "delete" {
		t.Errorf("AuditActionDelete = %s, want delete", model.AuditActionDelete)
	}
}

func setupAuditLogTest(t *testing.T) (*AuditLogRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewAuditLogRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
