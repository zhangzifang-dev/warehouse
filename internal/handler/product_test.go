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

type mockProductService struct {
	listFunc    func(ctx context.Context, page, pageSize int, categoryID int64, keyword string) (*service.ListProductsResult, error)
	getByIDFunc func(ctx context.Context, id int64) (*model.Product, error)
	createFunc  func(ctx context.Context, input *service.CreateProductInput) (*model.Product, error)
	updateFunc  func(ctx context.Context, id int64, input *service.UpdateProductInput) (*model.Product, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockProductService) List(ctx context.Context, page, pageSize int, categoryID int64, keyword string) (*service.ListProductsResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, categoryID, keyword)
	}
	return nil, errors.New("not implemented")
}

func (m *mockProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockProductService) Create(ctx context.Context, input *service.CreateProductInput) (*model.Product, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockProductService) Update(ctx context.Context, id int64, input *service.UpdateProductInput) (*model.Product, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockProductService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func setupProductHandlerTest(t *testing.T) (*gin.Engine, *ProductHandler, *mockProductService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockProductService{}
	handler := NewProductHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestProductHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		mockProducts   []model.Product
		mockTotal      int
		mockError      error
		queryCategory  string
		queryKeyword   string
		queryPage      string
		querySize      string
		wantStatus     int
		wantTotal      int
	}{
		{
			name: "success with default pagination",
			mockProducts: []model.Product{
				{BaseModel: model.BaseModel{ID: 1}, SKU: "SKU001", Name: "Product 1"},
				{BaseModel: model.BaseModel{ID: 2}, SKU: "SKU002", Name: "Product 2"},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with category filter",
			mockProducts: []model.Product{
				{BaseModel: model.BaseModel{ID: 1}, CategoryID: 1},
			},
			mockTotal:     1,
			queryCategory: "1",
			wantStatus:    http.StatusOK,
			wantTotal:     1,
		},
		{
			name: "success with keyword filter",
			mockProducts: []model.Product{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Test Product"},
			},
			mockTotal:     1,
			queryKeyword:  "test",
			wantStatus:    http.StatusOK,
			wantTotal:     1,
		},
		{
			name:         "empty list",
			mockProducts: []model.Product{},
			mockTotal:    0,
			wantStatus:   http.StatusOK,
			wantTotal:    0,
		},
		{
			name:       "service error",
			mockError:  apperrors.NewAppError(apperrors.CodeInternalError, "database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupProductHandlerTest(t)
			mockSvc.listFunc = func(ctx context.Context, page, pageSize int, categoryID int64, keyword string) (*service.ListProductsResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				return &service.ListProductsResult{
					Products: tt.mockProducts,
					Total:    tt.mockTotal,
				}, nil
			}

			router.GET("/products", handler.List)

			req := httptest.NewRequest("GET", "/products?category_id="+tt.queryCategory+"&keyword="+tt.queryKeyword+"&page="+tt.queryPage+"&size="+tt.querySize, nil)
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

func TestProductHandler_GetByID(t *testing.T) {
	tests := []struct {
		name         string
		productID    string
		mockProduct  *model.Product
		mockError    error
		wantStatus   int
	}{
		{
			name:        "success",
			productID:   "1",
			mockProduct: &model.Product{BaseModel: model.BaseModel{ID: 1}, SKU: "SKU001", Name: "Test Product"},
			wantStatus:  http.StatusOK,
		},
		{
			name:       "invalid id",
			productID:  "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			productID:  "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "product not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupProductHandlerTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.Product, error) {
				return tt.mockProduct, tt.mockError
			}

			router.GET("/products/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/products/"+tt.productID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestProductHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		mockProduct  *model.Product
		mockError    error
		wantStatus   int
	}{
		{
			name: "success",
			body: CreateProductRequest{
				SKU:   "SKU001",
				Name:  "Test Product",
				Price: productFloatPtrHandler(99.99),
			},
			mockProduct: &model.Product{BaseModel: model.BaseModel{ID: 1}, SKU: "SKU001", Name: "Test Product"},
			wantStatus:  http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateProductRequest{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "duplicate SKU",
			body:       CreateProductRequest{SKU: "SKU001", Name: "Test Product"},
			mockError:  apperrors.NewAppError(apperrors.CodeDuplicateEntry, "SKU already exists"),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupProductHandlerTest(t)
			mockSvc.createFunc = func(ctx context.Context, input *service.CreateProductInput) (*model.Product, error) {
				return tt.mockProduct, tt.mockError
			}

			router.POST("/products", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/products", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestProductHandler_Update(t *testing.T) {
	tests := []struct {
		name         string
		productID    string
		body         interface{}
		mockProduct  *model.Product
		mockError    error
		wantStatus   int
	}{
		{
			name:      "success",
			productID: "1",
			body: UpdateProductRequest{
				Name:  productStrPtrHandler("Updated Product"),
				Price: productFloatPtrHandler(199.99),
			},
			mockProduct: &model.Product{BaseModel: model.BaseModel{ID: 1}, Name: "Updated Product"},
			wantStatus:  http.StatusOK,
		},
		{
			name:       "invalid id",
			productID:  "invalid",
			body:       UpdateProductRequest{Name: productStrPtrHandler("Updated")},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			productID:  "999",
			body:       UpdateProductRequest{Name: productStrPtrHandler("Updated")},
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "product not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			productID:  "1",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "duplicate SKU",
			productID:  "1",
			body:       UpdateProductRequest{SKU: productStrPtrHandler("SKU002")},
			mockError:  apperrors.NewAppError(apperrors.CodeDuplicateEntry, "SKU already exists"),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupProductHandlerTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, input *service.UpdateProductInput) (*model.Product, error) {
				return tt.mockProduct, tt.mockError
			}

			router.PUT("/products/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/products/"+tt.productID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestProductHandler_Delete(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		mockError   error
		wantStatus  int
	}{
		{
			name:       "success",
			productID:  "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			productID:  "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			productID:  "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "product not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupProductHandlerTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/products/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/products/"+tt.productID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func productStrPtrHandler(s string) *string {
	return &s
}

func productFloatPtrHandler(f float64) *float64 {
	return &f
}
