package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"warehouse/internal/model"
	"warehouse/internal/repository"
)

type mockAuditLogRepository struct {
	createFunc  func(ctx context.Context, log *model.AuditLog) error
	getByIDFunc func(ctx context.Context, id int64) (*model.AuditLog, error)
	listFunc    func(ctx context.Context, filter *repository.AuditLogFilter) ([]model.AuditLog, int, error)
}

func (m *mockAuditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, log)
	}
	return errors.New("not implemented")
}

func (m *mockAuditLogRepository) GetByID(ctx context.Context, id int64) (*model.AuditLog, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuditLogRepository) List(ctx context.Context, filter *repository.AuditLogFilter) ([]model.AuditLog, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, 0, errors.New("not implemented")
}

func TestAuditLogService_Log_Success(t *testing.T) {
	var createdLog *model.AuditLog
	mockRepo := &mockAuditLogRepository{
		createFunc: func(ctx context.Context, log *model.AuditLog) error {
			createdLog = log
			return nil
		},
	}

	svc := NewAuditLogService(mockRepo)
	input := &CreateAuditLogInput{
		TableName:  "users",
		RecordID:   1,
		Action:     model.AuditActionCreate,
		OperatedBy: 1,
		IPAddress:  "127.0.0.1",
	}

	err := svc.Log(context.Background(), input)

	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}
	if createdLog == nil {
		t.Fatal("expected log to be created")
	}
	if createdLog.TableName != "users" {
		t.Errorf("expected table_name 'users', got '%s'", createdLog.TableName)
	}
	if createdLog.Action != model.AuditActionCreate {
		t.Errorf("expected action 'create', got '%s'", createdLog.Action)
	}
	if createdLog.OperatedAt.IsZero() {
		t.Error("expected operated_at to be set")
	}
}

func TestAuditLogService_Log_WithValues(t *testing.T) {
	var createdLog *model.AuditLog
	mockRepo := &mockAuditLogRepository{
		createFunc: func(ctx context.Context, log *model.AuditLog) error {
			createdLog = log
			return nil
		},
	}

	svc := NewAuditLogService(mockRepo)
	input := &CreateAuditLogInput{
		TableName:  "products",
		RecordID:   42,
		Action:     model.AuditActionUpdate,
		OldValue:   map[string]any{"name": "old"},
		NewValue:   map[string]any{"name": "new"},
		OperatedBy: 2,
		IPAddress:  "192.168.1.1",
	}

	err := svc.Log(context.Background(), input)

	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}
	if createdLog.OldValue["name"] != "old" {
		t.Errorf("expected old_value name 'old', got '%v'", createdLog.OldValue["name"])
	}
	if createdLog.NewValue["name"] != "new" {
		t.Errorf("expected new_value name 'new', got '%v'", createdLog.NewValue["name"])
	}
}

func TestAuditLogService_GetByID_Success(t *testing.T) {
	mockRepo := &mockAuditLogRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.AuditLog, error) {
			return &model.AuditLog{
				ID:        id,
				TableName: "users",
				RecordID:  1,
				Action:    model.AuditActionCreate,
			}, nil
		},
	}

	svc := NewAuditLogService(mockRepo)

	log, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if log == nil {
		t.Fatal("expected log, got nil")
	}
	if log.TableName != "users" {
		t.Errorf("expected table_name 'users', got '%s'", log.TableName)
	}
}

func TestAuditLogService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockAuditLogRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.AuditLog, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewAuditLogService(mockRepo)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent log, got nil")
	}
}

func TestAuditLogService_List_Success(t *testing.T) {
	mockRepo := &mockAuditLogRepository{
		listFunc: func(ctx context.Context, filter *repository.AuditLogFilter) ([]model.AuditLog, int, error) {
			return []model.AuditLog{
				{ID: 1, TableName: "users", Action: model.AuditActionCreate},
				{ID: 2, TableName: "products", Action: model.AuditActionUpdate},
			}, 2, nil
		},
	}

	svc := NewAuditLogService(mockRepo)

	result, err := svc.List(context.Background(), &AuditLogQueryFilter{Page: 1, PageSize: 10})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Items))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestAuditLogService_List_DefaultPagination(t *testing.T) {
	mockRepo := &mockAuditLogRepository{
		listFunc: func(ctx context.Context, filter *repository.AuditLogFilter) ([]model.AuditLog, int, error) {
			if filter.Page != 1 {
				t.Errorf("expected page 1, got %d", filter.Page)
			}
			if filter.PageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", filter.PageSize)
			}
			return []model.AuditLog{}, 0, nil
		},
	}

	svc := NewAuditLogService(mockRepo)

	_, err := svc.List(context.Background(), &AuditLogQueryFilter{})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestAuditLogService_List_MaxPageSize(t *testing.T) {
	mockRepo := &mockAuditLogRepository{
		listFunc: func(ctx context.Context, filter *repository.AuditLogFilter) ([]model.AuditLog, int, error) {
			if filter.PageSize > 100 {
				t.Errorf("expected pageSize <= 100, got %d", filter.PageSize)
			}
			return []model.AuditLog{}, 0, nil
		},
	}

	svc := NewAuditLogService(mockRepo)

	_, err := svc.List(context.Background(), &AuditLogQueryFilter{Page: 1, PageSize: 200})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestAuditLogService_List_WithFilters(t *testing.T) {
	var capturedFilter *repository.AuditLogFilter
	mockRepo := &mockAuditLogRepository{
		listFunc: func(ctx context.Context, filter *repository.AuditLogFilter) ([]model.AuditLog, int, error) {
			capturedFilter = filter
			return []model.AuditLog{}, 0, nil
		},
	}

	svc := NewAuditLogService(mockRepo)
	recordID := int64(1)
	operatedBy := int64(2)
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	_, err := svc.List(context.Background(), &AuditLogQueryFilter{
		TableName:  "users",
		RecordID:   &recordID,
		OperatedBy: &operatedBy,
		StartTime:  &startTime,
		EndTime:    &endTime,
		Page:       1,
		PageSize:   10,
	})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if capturedFilter.TableName != "users" {
		t.Errorf("expected table_name 'users', got '%s'", capturedFilter.TableName)
	}
	if *capturedFilter.RecordID != recordID {
		t.Errorf("expected record_id %d, got %d", recordID, *capturedFilter.RecordID)
	}
	if *capturedFilter.OperatedBy != operatedBy {
		t.Errorf("expected operated_by %d, got %d", operatedBy, *capturedFilter.OperatedBy)
	}
}
