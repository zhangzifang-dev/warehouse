package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"warehouse/internal/model"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
)

type mockAuditLogService struct {
	getByIDFunc func(ctx context.Context, id int64) (*model.AuditLog, error)
	listFunc    func(ctx context.Context, filter *service.AuditLogQueryFilter) (*service.AuditLogListResult, error)
}

func (m *mockAuditLogService) GetByID(ctx context.Context, id int64) (*model.AuditLog, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuditLogService) List(ctx context.Context, filter *service.AuditLogQueryFilter) (*service.AuditLogListResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, errors.New("not implemented")
}

func TestAuditLogHandler_List_Success(t *testing.T) {
	mockSvc := &mockAuditLogService{
		listFunc: func(ctx context.Context, filter *service.AuditLogQueryFilter) (*service.AuditLogListResult, error) {
			return &service.AuditLogListResult{
				Items: []model.AuditLog{
					{ID: 1, TableName: "users", RecordID: 1, Action: model.AuditActionCreate},
					{ID: 2, TableName: "products", RecordID: 2, Action: model.AuditActionUpdate},
				},
				Total: 2,
			}, nil
		},
	}

	handler := NewAuditLogHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/audit-logs", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/audit-logs?page=1&size=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in response")
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		t.Fatal("expected items array in response")
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestAuditLogHandler_List_WithFilters(t *testing.T) {
	var capturedFilter *service.AuditLogQueryFilter
	mockSvc := &mockAuditLogService{
		listFunc: func(ctx context.Context, filter *service.AuditLogQueryFilter) (*service.AuditLogListResult, error) {
			capturedFilter = filter
			return &service.AuditLogListResult{Items: []model.AuditLog{}, Total: 0}, nil
		},
	}

	handler := NewAuditLogHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/audit-logs", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/audit-logs?table_name=users&record_id=1&operated_by=2", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if capturedFilter.TableName != "users" {
		t.Errorf("expected table_name 'users', got '%s'", capturedFilter.TableName)
	}
	if capturedFilter.RecordID == nil || *capturedFilter.RecordID != 1 {
		t.Errorf("expected record_id 1, got %v", capturedFilter.RecordID)
	}
	if capturedFilter.OperatedBy == nil || *capturedFilter.OperatedBy != 2 {
		t.Errorf("expected operated_by 2, got %v", capturedFilter.OperatedBy)
	}
}

func TestAuditLogHandler_List_TimeFilters(t *testing.T) {
	var capturedFilter *service.AuditLogQueryFilter
	mockSvc := &mockAuditLogService{
		listFunc: func(ctx context.Context, filter *service.AuditLogQueryFilter) (*service.AuditLogListResult, error) {
			capturedFilter = filter
			return &service.AuditLogListResult{Items: []model.AuditLog{}, Total: 0}, nil
		},
	}

	handler := NewAuditLogHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/audit-logs", handler.List)

	startTime := "2024-01-01T00:00:00Z"
	endTime := "2024-12-31T23:59:59Z"
	req := httptest.NewRequest(http.MethodGet, "/audit-logs?start_time="+startTime+"&end_time="+endTime, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if capturedFilter.StartTime == nil {
		t.Error("expected start_time to be set")
	}
	if capturedFilter.EndTime == nil {
		t.Error("expected end_time to be set")
	}
}

func TestAuditLogHandler_List_DefaultPagination(t *testing.T) {
	mockSvc := &mockAuditLogService{
		listFunc: func(ctx context.Context, filter *service.AuditLogQueryFilter) (*service.AuditLogListResult, error) {
			if filter.Page != 1 {
				t.Errorf("expected page 1, got %d", filter.Page)
			}
			if filter.PageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", filter.PageSize)
			}
			return &service.AuditLogListResult{Items: []model.AuditLog{}, Total: 0}, nil
		},
	}

	handler := NewAuditLogHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/audit-logs", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/audit-logs", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestAuditLogHandler_GetByID_Success(t *testing.T) {
	mockSvc := &mockAuditLogService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.AuditLog, error) {
			return &model.AuditLog{
				ID:        id,
				TableName: "users",
				RecordID:  1,
				Action:    model.AuditActionCreate,
			}, nil
		},
	}

	handler := NewAuditLogHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/audit-logs/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/audit-logs/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in response")
	}
	if data["table_name"] != "users" {
		t.Errorf("expected table_name 'users', got '%v'", data["table_name"])
	}
}

func TestAuditLogHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockAuditLogService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.AuditLog, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewAuditLogHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/audit-logs/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/audit-logs/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestAuditLogHandler_GetByID_InvalidID(t *testing.T) {
	mockSvc := &mockAuditLogService{}
	handler := NewAuditLogHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/audit-logs/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/audit-logs/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestNewAuditLogHandler(t *testing.T) {
	mockSvc := &mockAuditLogService{}
	handler := NewAuditLogHandler(mockSvc)
	if handler == nil {
		t.Error("NewAuditLogHandler() returned nil")
	}
}

func TestRegisterAuditLogRoutes(t *testing.T) {
	mockSvc := &mockAuditLogService{}
	handler := NewAuditLogHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterAuditLogRoutes(router.Group("/api"), handler)

	routes := router.Routes()
	if len(routes) != 2 {
		t.Errorf("expected 2 routes, got %d", len(routes))
	}
}
