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

type mockSupplierService struct {
	listFunc    func(ctx context.Context, page, pageSize int, keyword string) (*service.ListSuppliersResult, error)
	getByIDFunc func(ctx context.Context, id int64) (*model.Supplier, error)
	createFunc  func(ctx context.Context, input *service.CreateSupplierInput) (*model.Supplier, error)
	updateFunc  func(ctx context.Context, id int64, input *service.UpdateSupplierInput) (*model.Supplier, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockSupplierService) List(ctx context.Context, page, pageSize int, keyword string) (*service.ListSuppliersResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, keyword)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSupplierService) GetByID(ctx context.Context, id int64) (*model.Supplier, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSupplierService) Create(ctx context.Context, input *service.CreateSupplierInput) (*model.Supplier, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSupplierService) Update(ctx context.Context, id int64, input *service.UpdateSupplierInput) (*model.Supplier, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSupplierService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func setupSupplierHandlerTest(t *testing.T) (*gin.Engine, *SupplierHandler, *mockSupplierService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockSupplierService{}
	handler := NewSupplierHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestSupplierHandler_List(t *testing.T) {
	tests := []struct {
		name          string
		mockSuppliers []model.Supplier
		mockTotal     int
		mockError     error
		queryKeyword  string
		queryPage     string
		querySize     string
		wantStatus    int
		wantTotal     int
	}{
		{
			name: "success with default pagination",
			mockSuppliers: []model.Supplier{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Supplier A"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "Supplier B"},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with keyword filter",
			mockSuppliers: []model.Supplier{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Test Supplier"},
			},
			mockTotal:    1,
			queryKeyword: "Test",
			wantStatus:   http.StatusOK,
			wantTotal:    1,
		},
		{
			name:          "empty list",
			mockSuppliers: []model.Supplier{},
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
			router, handler, mockSvc := setupSupplierHandlerTest(t)
			mockSvc.listFunc = func(ctx context.Context, page, pageSize int, keyword string) (*service.ListSuppliersResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				return &service.ListSuppliersResult{
					Suppliers: tt.mockSuppliers,
					Total:     tt.mockTotal,
				}, nil
			}

			router.GET("/suppliers", handler.List)

			req := httptest.NewRequest("GET", "/suppliers?keyword="+tt.queryKeyword+"&page="+tt.queryPage+"&size="+tt.querySize, nil)
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

func TestSupplierHandler_GetByID(t *testing.T) {
	tests := []struct {
		name         string
		supplierID   string
		mockSupplier *model.Supplier
		mockError    error
		wantStatus   int
	}{
		{
			name:       "success",
			supplierID: "1",
			mockSupplier: &model.Supplier{BaseModel: model.BaseModel{ID: 1}, Name: "Test Supplier"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			supplierID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			supplierID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "supplier not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupSupplierHandlerTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.Supplier, error) {
				return tt.mockSupplier, tt.mockError
			}

			router.GET("/suppliers/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/suppliers/"+tt.supplierID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSupplierHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		mockSupplier *model.Supplier
		mockError    error
		wantStatus   int
	}{
		{
			name: "success",
			body: CreateSupplierRequest{
				Name:    "Test Supplier",
				Code:    strPtrSupplierHandler("SUP001"),
				Contact: strPtrSupplierHandler("John Doe"),
			},
			mockSupplier: &model.Supplier{BaseModel: model.BaseModel{ID: 1}, Name: "Test Supplier"},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateSupplierRequest{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "duplicate code",
			body:       CreateSupplierRequest{Name: "Test", Code: strPtrSupplierHandler("DUP")},
			mockError:  apperrors.NewAppError(apperrors.CodeDuplicateEntry, "supplier code already exists"),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupSupplierHandlerTest(t)
			mockSvc.createFunc = func(ctx context.Context, input *service.CreateSupplierInput) (*model.Supplier, error) {
				return tt.mockSupplier, tt.mockError
			}

			router.POST("/suppliers", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/suppliers", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSupplierHandler_Update(t *testing.T) {
	tests := []struct {
		name         string
		supplierID   string
		body         interface{}
		mockSupplier *model.Supplier
		mockError    error
		wantStatus   int
	}{
		{
			name:       "success",
			supplierID: "1",
			body: UpdateSupplierRequest{
				Name:    strPtrSupplierHandler("Updated Supplier"),
				Contact: strPtrSupplierHandler("Jane Doe"),
			},
			mockSupplier: &model.Supplier{BaseModel: model.BaseModel{ID: 1}, Name: "Updated Supplier"},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "invalid id",
			supplierID: "invalid",
			body:       UpdateSupplierRequest{Name: strPtrSupplierHandler("Updated")},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			supplierID: "999",
			body:       UpdateSupplierRequest{Name: strPtrSupplierHandler("Updated")},
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "supplier not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			supplierID: "1",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupSupplierHandlerTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, input *service.UpdateSupplierInput) (*model.Supplier, error) {
				return tt.mockSupplier, tt.mockError
			}

			router.PUT("/suppliers/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/suppliers/"+tt.supplierID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSupplierHandler_Delete(t *testing.T) {
	tests := []struct {
		name        string
		supplierID  string
		mockError   error
		wantStatus  int
	}{
		{
			name:       "success",
			supplierID: "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			supplierID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			supplierID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "supplier not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupSupplierHandlerTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/suppliers/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/suppliers/"+tt.supplierID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func strPtrSupplierHandler(s string) *string {
	return &s
}
