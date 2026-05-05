package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockWarehouseService struct {
	listFunc    func(ctx context.Context, filter *service.WarehouseFilter) (*service.ListWarehousesResult, error)
	getByIDFunc func(ctx context.Context, id int64) (*model.Warehouse, error)
	createFunc  func(ctx context.Context, input *service.CreateWarehouseInput) (*model.Warehouse, error)
	updateFunc  func(ctx context.Context, id int64, input *service.UpdateWarehouseInput) (*model.Warehouse, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockWarehouseService) List(ctx context.Context, filter *service.WarehouseFilter) (*service.ListWarehousesResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWarehouseService) GetByID(ctx context.Context, id int64) (*model.Warehouse, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWarehouseService) Create(ctx context.Context, input *service.CreateWarehouseInput) (*model.Warehouse, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWarehouseService) Update(ctx context.Context, id int64, input *service.UpdateWarehouseInput) (*model.Warehouse, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWarehouseService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func setupWarehouseHandlerTest(t *testing.T) (*gin.Engine, *WarehouseHandler, *mockWarehouseService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockWarehouseService{}
	handler := NewWarehouseHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestWarehouseHandler_List(t *testing.T) {
	tests := []struct {
		name          string
		mockWarehouses []model.Warehouse
		mockTotal     int
		mockError     error
		queryPage     string
		querySize     string
		wantStatus    int
		wantTotal     int
	}{
		{
			name: "success with default pagination",
			mockWarehouses: []model.Warehouse{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Warehouse 1", Code: "WH001"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "Warehouse 2", Code: "WH002"},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with custom pagination",
			mockWarehouses: []model.Warehouse{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Warehouse 1", Code: "WH001"},
			},
			mockTotal:  10,
			queryPage:  "2",
			querySize:  "5",
			wantStatus: http.StatusOK,
			wantTotal:  10,
		},
		{
			name:          "empty list",
			mockWarehouses: []model.Warehouse{},
			mockTotal:     0,
			wantStatus:    http.StatusOK,
			wantTotal:     0,
		},
		{
			name:       "service error",
			mockError:  apperrors.NewAppError(apperrors.CodeInternalError, "database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupWarehouseHandlerTest(t)
			mockSvc.listFunc = func(ctx context.Context, filter *service.WarehouseFilter) (*service.ListWarehousesResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				return &service.ListWarehousesResult{
					Warehouses: tt.mockWarehouses,
					Total:      tt.mockTotal,
				}, nil
			}

			router.GET("/warehouses", handler.List)

			req := httptest.NewRequest("GET", "/warehouses?page="+tt.queryPage+"&size="+tt.querySize, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, float64(tt.wantTotal), data["total"])
			}
		})
	}
}

func TestWarehouseHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		warehouseID string
		mockWarehouse *model.Warehouse
		mockError  error
		wantStatus int
	}{
		{
			name:       "success",
			warehouseID: "1",
			mockWarehouse: &model.Warehouse{BaseModel: model.BaseModel{ID: 1}, Name: "Main Warehouse", Code: "WH001"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			warehouseID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			warehouseID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupWarehouseHandlerTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.Warehouse, error) {
				return tt.mockWarehouse, tt.mockError
			}

			router.GET("/warehouses/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/warehouses/"+tt.warehouseID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestWarehouseHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		mockWarehouse *model.Warehouse
		mockError    error
		wantStatus   int
	}{
		{
			name: "success",
			body: CreateWarehouseRequest{
				Name:    "Main Warehouse",
				Code:    "WH001",
				Address: "123 Main St",
				Contact: "John Doe",
				Phone:   "1234567890",
			},
			mockWarehouse: &model.Warehouse{BaseModel: model.BaseModel{ID: 1}, Name: "Main Warehouse", Code: "WH001"},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateWarehouseRequest{Name: ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate code",
			body: CreateWarehouseRequest{
				Name: "Main Warehouse",
				Code: "WH001",
			},
			mockError:  apperrors.NewAppError(apperrors.CodeDuplicateEntry, "warehouse code already exists"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupWarehouseHandlerTest(t)
			mockSvc.createFunc = func(ctx context.Context, input *service.CreateWarehouseInput) (*model.Warehouse, error) {
				return tt.mockWarehouse, tt.mockError
			}

			router.POST("/warehouses", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/warehouses", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestWarehouseHandler_Update(t *testing.T) {
	tests := []struct {
		name         string
		warehouseID  string
		body         interface{}
		mockWarehouse *model.Warehouse
		mockError    error
		wantStatus   int
	}{
		{
			name:        "success",
			warehouseID: "1",
			body: UpdateWarehouseRequest{
				Name:    "Updated Warehouse",
				Address: "New Address",
			},
			mockWarehouse: &model.Warehouse{BaseModel: model.BaseModel{ID: 1}, Name: "Updated Warehouse"},
			wantStatus:   http.StatusOK,
		},
		{
			name:        "invalid id",
			warehouseID: "invalid",
			body:        UpdateWarehouseRequest{Name: "Test"},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "not found",
			warehouseID: "999",
			body:        UpdateWarehouseRequest{Name: "Test"},
			mockError:   apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found"),
			wantStatus:  http.StatusNotFound,
		},
		{
			name:        "invalid json",
			warehouseID: "1",
			body:        "invalid",
			wantStatus:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupWarehouseHandlerTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, input *service.UpdateWarehouseInput) (*model.Warehouse, error) {
				return tt.mockWarehouse, tt.mockError
			}

			router.PUT("/warehouses/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/warehouses/"+tt.warehouseID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestWarehouseHandler_Delete(t *testing.T) {
	tests := []struct {
		name        string
		warehouseID string
		mockError   error
		wantStatus  int
	}{
		{
			name:        "success",
			warehouseID: "1",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "invalid id",
			warehouseID: "invalid",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "not found",
			warehouseID: "999",
			mockError:   apperrors.NewAppError(apperrors.CodeNotFound, "warehouse not found"),
			wantStatus:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupWarehouseHandlerTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/warehouses/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/warehouses/"+tt.warehouseID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
