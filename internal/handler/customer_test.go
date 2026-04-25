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

type mockCustomerService struct {
	listFunc    func(ctx context.Context, page, pageSize int, keyword string) (*service.ListCustomersResult, error)
	getByIDFunc func(ctx context.Context, id int64) (*model.Customer, error)
	createFunc  func(ctx context.Context, input *service.CreateCustomerInput) (*model.Customer, error)
	updateFunc  func(ctx context.Context, id int64, input *service.UpdateCustomerInput) (*model.Customer, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockCustomerService) List(ctx context.Context, page, pageSize int, keyword string) (*service.ListCustomersResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, keyword)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCustomerService) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCustomerService) Create(ctx context.Context, input *service.CreateCustomerInput) (*model.Customer, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCustomerService) Update(ctx context.Context, id int64, input *service.UpdateCustomerInput) (*model.Customer, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCustomerService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func setupCustomerHandlerTest(t *testing.T) (*gin.Engine, *CustomerHandler, *mockCustomerService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockCustomerService{}
	handler := NewCustomerHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestCustomerHandler_List(t *testing.T) {
	tests := []struct {
		name          string
		mockCustomers []model.Customer
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
			mockCustomers: []model.Customer{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Customer A"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "Customer B"},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with keyword filter",
			mockCustomers: []model.Customer{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Test Customer"},
			},
			mockTotal:    1,
			queryKeyword: "Test",
			wantStatus:   http.StatusOK,
			wantTotal:    1,
		},
		{
			name:          "empty list",
			mockCustomers: []model.Customer{},
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
			router, handler, mockSvc := setupCustomerHandlerTest(t)
			mockSvc.listFunc = func(ctx context.Context, page, pageSize int, keyword string) (*service.ListCustomersResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				return &service.ListCustomersResult{
					Customers: tt.mockCustomers,
					Total:     tt.mockTotal,
				}, nil
			}

			router.GET("/customers", handler.List)

			req := httptest.NewRequest("GET", "/customers?keyword="+tt.queryKeyword+"&page="+tt.queryPage+"&size="+tt.querySize, nil)
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

func TestCustomerHandler_GetByID(t *testing.T) {
	tests := []struct {
		name         string
		customerID   string
		mockCustomer *model.Customer
		mockError    error
		wantStatus   int
	}{
		{
			name:       "success",
			customerID: "1",
			mockCustomer: &model.Customer{BaseModel: model.BaseModel{ID: 1}, Name: "Test Customer"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			customerID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			customerID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "customer not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupCustomerHandlerTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.Customer, error) {
				return tt.mockCustomer, tt.mockError
			}

			router.GET("/customers/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/customers/"+tt.customerID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCustomerHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		mockCustomer *model.Customer
		mockError    error
		wantStatus   int
	}{
		{
			name: "success",
			body: CreateCustomerRequest{
				Name:    "Test Customer",
				Code:    strPtrCustomerHandler("CUS001"),
				Contact: strPtrCustomerHandler("John Doe"),
			},
			mockCustomer: &model.Customer{BaseModel: model.BaseModel{ID: 1}, Name: "Test Customer"},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateCustomerRequest{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "duplicate code",
			body:       CreateCustomerRequest{Name: "Test", Code: strPtrCustomerHandler("DUP")},
			mockError:  apperrors.NewAppError(apperrors.CodeDuplicateEntry, "customer code already exists"),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupCustomerHandlerTest(t)
			mockSvc.createFunc = func(ctx context.Context, input *service.CreateCustomerInput) (*model.Customer, error) {
				return tt.mockCustomer, tt.mockError
			}

			router.POST("/customers", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/customers", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCustomerHandler_Update(t *testing.T) {
	tests := []struct {
		name         string
		customerID   string
		body         interface{}
		mockCustomer *model.Customer
		mockError    error
		wantStatus   int
	}{
		{
			name:       "success",
			customerID: "1",
			body: UpdateCustomerRequest{
				Name:    strPtrCustomerHandler("Updated Customer"),
				Contact: strPtrCustomerHandler("Jane Doe"),
			},
			mockCustomer: &model.Customer{BaseModel: model.BaseModel{ID: 1}, Name: "Updated Customer"},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "invalid id",
			customerID: "invalid",
			body:       UpdateCustomerRequest{Name: strPtrCustomerHandler("Updated")},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			customerID: "999",
			body:       UpdateCustomerRequest{Name: strPtrCustomerHandler("Updated")},
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "customer not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			customerID: "1",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupCustomerHandlerTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, input *service.UpdateCustomerInput) (*model.Customer, error) {
				return tt.mockCustomer, tt.mockError
			}

			router.PUT("/customers/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/customers/"+tt.customerID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCustomerHandler_Delete(t *testing.T) {
	tests := []struct {
		name        string
		customerID  string
		mockError   error
		wantStatus  int
	}{
		{
			name:       "success",
			customerID: "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			customerID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			customerID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "customer not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupCustomerHandlerTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/customers/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/customers/"+tt.customerID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func strPtrCustomerHandler(s string) *string {
	return &s
}
