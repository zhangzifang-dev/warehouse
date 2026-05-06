package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockInboundOrderService struct {
	listFunc             func(ctx context.Context, page, pageSize int, warehouseID, status int) (*service.ListInboundOrdersResult, error)
	listWithFilterFunc   func(ctx context.Context, filter *model.InboundOrderQueryFilter) (*service.ListInboundOrdersResult, error)
	getByIDFunc          func(ctx context.Context, id int64) (*model.InboundOrder, error)
	createFunc           func(ctx context.Context, input *service.CreateInboundOrderInput) (*model.InboundOrder, error)
	updateFunc           func(ctx context.Context, id int64, input *service.UpdateInboundOrderInput) (*model.InboundOrder, error)
	deleteFunc           func(ctx context.Context, id int64) error
	confirmFunc          func(ctx context.Context, id int64) (*model.InboundOrder, error)
}

func (m *mockInboundOrderService) List(ctx context.Context, page, pageSize int, warehouseID, status int) (*service.ListInboundOrdersResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, warehouseID, status)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInboundOrderService) ListWithFilter(ctx context.Context, filter *model.InboundOrderQueryFilter) (*service.ListInboundOrdersResult, error) {
	if m.listWithFilterFunc != nil {
		return m.listWithFilterFunc(ctx, filter)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInboundOrderService) GetByID(ctx context.Context, id int64) (*model.InboundOrder, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInboundOrderService) Create(ctx context.Context, input *service.CreateInboundOrderInput) (*model.InboundOrder, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInboundOrderService) Update(ctx context.Context, id int64, input *service.UpdateInboundOrderInput) (*model.InboundOrder, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInboundOrderService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockInboundOrderService) Confirm(ctx context.Context, id int64) (*model.InboundOrder, error) {
	if m.confirmFunc != nil {
		return m.confirmFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func setupInboundOrderHandlerTest(t *testing.T) (*gin.Engine, *InboundOrderHandler, *mockInboundOrderService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockInboundOrderService{}
	handler := NewInboundOrderHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestInboundOrderHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		mockOrders     []model.InboundOrder
		mockTotal      int
		mockError      error
		queryWarehouse string
		queryStatus    string
		queryPage      string
		querySize      string
		wantStatus     int
		wantTotal      int
	}{
		{
			name: "success with default pagination",
			mockOrders: []model.InboundOrder{
				{BaseModel: model.BaseModel{ID: 1}, OrderNo: "PO-2024-001", WarehouseID: 1, TotalQuantity: 100},
				{BaseModel: model.BaseModel{ID: 2}, OrderNo: "PO-2024-002", WarehouseID: 1, TotalQuantity: 200},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with warehouse filter",
			mockOrders: []model.InboundOrder{
				{BaseModel: model.BaseModel{ID: 1}},
			},
			mockTotal:      1,
			queryWarehouse: "1",
			wantStatus:     http.StatusOK,
			wantTotal:       1,
		},
		{
			name: "success with status filter",
			mockOrders: []model.InboundOrder{
				{BaseModel: model.BaseModel{ID: 1}},
			},
			mockTotal:   1,
			queryStatus: "1",
			wantStatus:  http.StatusOK,
			wantTotal:   1,
		},
		{
			name:       "empty list",
			mockOrders: []model.InboundOrder{},
			mockTotal:  0,
			wantStatus: http.StatusOK,
			wantTotal:  0,
		},
		{
			name:       "service error",
			mockError:  apperrors.NewAppError(apperrors.CodeInternalError, "database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInboundOrderHandlerTest(t)
			mockSvc.listWithFilterFunc = func(ctx context.Context, filter *model.InboundOrderQueryFilter) (*service.ListInboundOrdersResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				if tt.queryWarehouse != "" {
					warehouseID := int64(1)
					if filter.WarehouseID == nil || *filter.WarehouseID != warehouseID {
						t.Errorf("expected WarehouseID %d", warehouseID)
					}
				}
				return &service.ListInboundOrdersResult{
					Orders: tt.mockOrders,
					Total:  tt.mockTotal,
				}, nil
			}

			router.GET("/inbound-orders", handler.List)

			req := httptest.NewRequest("GET", "/inbound-orders?warehouse_id="+tt.queryWarehouse+"&status="+tt.queryStatus+"&page="+tt.queryPage+"&size="+tt.querySize, nil)
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

func TestInboundOrderHandler_ListWithFilter(t *testing.T) {
	now := time.Now()
	supplierID := int64(1)
	warehouseID := int64(2)

	tests := []struct {
		name       string
		query      string
		mockResult *service.ListInboundOrdersResult
		mockError  error
		wantStatus int
		wantTotal  int
	}{
		{
			name:  "success with all filters",
			query: "?order_no=PO-2024&supplier_id=1&warehouse_id=2&quantity_min=10&quantity_max=100&created_at_start=" + now.Format(time.RFC3339) + "&created_at_end=" + now.Add(24*time.Hour).Format(time.RFC3339),
			mockResult: &service.ListInboundOrdersResult{
				Orders: []model.InboundOrder{
					{BaseModel: model.BaseModel{ID: 1}, OrderNo: "PO-2024-001"},
				},
				Total: 1,
			},
			wantStatus: http.StatusOK,
			wantTotal:  1,
		},
		{
			name:  "success with partial filters",
			query: "?order_no=PO-2024&warehouse_id=2",
			mockResult: &service.ListInboundOrdersResult{
				Orders: []model.InboundOrder{
					{BaseModel: model.BaseModel{ID: 1}, OrderNo: "PO-2024-001", WarehouseID: 2},
				},
				Total: 1,
			},
			wantStatus: http.StatusOK,
			wantTotal:  1,
		},
		{
			name:       "service error",
			query:      "?order_no=PO-2024",
			mockError:  apperrors.NewAppError(apperrors.CodeInternalError, "database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInboundOrderHandlerTest(t)
			mockSvc.listWithFilterFunc = func(ctx context.Context, filter *model.InboundOrderQueryFilter) (*service.ListInboundOrdersResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				if filter.OrderNo != "" && filter.OrderNo != "PO-2024" {
					t.Errorf("expected OrderNo 'PO-2024', got '%s'", filter.OrderNo)
				}
				if filter.SupplierID != nil && *filter.SupplierID != supplierID {
					t.Errorf("expected SupplierID %d, got %d", supplierID, *filter.SupplierID)
				}
				if filter.WarehouseID != nil && *filter.WarehouseID != warehouseID {
					t.Errorf("expected WarehouseID %d, got %d", warehouseID, *filter.WarehouseID)
				}
				return tt.mockResult, nil
			}

			router.GET("/inbound-orders", handler.List)

			req := httptest.NewRequest("GET", "/inbound-orders"+tt.query, nil)
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

func TestInboundOrderHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		orderID    string
		mockOrder  *model.InboundOrder
		mockError  error
		wantStatus int
	}{
		{
			name:       "success",
			orderID:    "1",
			mockOrder:  &model.InboundOrder{BaseModel: model.BaseModel{ID: 1}, OrderNo: "PO-2024-001", WarehouseID: 1, TotalQuantity: 100},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			orderID:    "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			orderID:    "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "inbound order not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInboundOrderHandlerTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.InboundOrder, error) {
				return tt.mockOrder, tt.mockError
			}

			router.GET("/inbound-orders/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/inbound-orders/"+tt.orderID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestInboundOrderHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		mockOrder  *model.InboundOrder
		mockError  error
		wantStatus int
	}{
		{
			name: "success",
			body: CreateInboundOrderRequest{
				OrderNo:       "PO-2024-001",
				WarehouseID:   1,
				TotalQuantity: 100,
			},
			mockOrder:  &model.InboundOrder{BaseModel: model.BaseModel{ID: 1}, OrderNo: "PO-2024-001", WarehouseID: 1, TotalQuantity: 100},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateInboundOrderRequest{},
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
			router, handler, mockSvc := setupInboundOrderHandlerTest(t)
			mockSvc.createFunc = func(ctx context.Context, input *service.CreateInboundOrderInput) (*model.InboundOrder, error) {
				return tt.mockOrder, tt.mockError
			}

			router.POST("/inbound-orders", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/inbound-orders", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestInboundOrderHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		orderID    string
		body       interface{}
		mockOrder  *model.InboundOrder
		mockError  error
		wantStatus int
	}{
		{
			name:    "success",
			orderID: "1",
			body: UpdateInboundOrderRequest{
				TotalQuantity: floatPtrIO(200),
			},
			mockOrder:  &model.InboundOrder{BaseModel: model.BaseModel{ID: 1}, TotalQuantity: 200},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			orderID:    "invalid",
			body:       UpdateInboundOrderRequest{TotalQuantity: floatPtrIO(200)},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			orderID:    "999",
			body:       UpdateInboundOrderRequest{TotalQuantity: floatPtrIO(200)},
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "inbound order not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			orderID:    "1",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInboundOrderHandlerTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, input *service.UpdateInboundOrderInput) (*model.InboundOrder, error) {
				return tt.mockOrder, tt.mockError
			}

			router.PUT("/inbound-orders/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/inbound-orders/"+tt.orderID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestInboundOrderHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		orderID    string
		mockError  error
		wantStatus int
	}{
		{
			name:       "success",
			orderID:    "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			orderID:    "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			orderID:    "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "inbound order not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInboundOrderHandlerTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/inbound-orders/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/inbound-orders/"+tt.orderID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestInboundOrderHandler_Confirm(t *testing.T) {
	tests := []struct {
		name       string
		orderID    string
		mockOrder  *model.InboundOrder
		mockError  error
		wantStatus int
	}{
		{
			name:       "success",
			orderID:    "1",
			mockOrder:  &model.InboundOrder{BaseModel: model.BaseModel{ID: 1}, Status: 1},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			orderID:    "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "already completed",
			orderID:    "1",
			mockError:  apperrors.NewAppError(apperrors.CodeBadRequest, "order already completed"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			orderID:    "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "inbound order not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInboundOrderHandlerTest(t)
			mockSvc.confirmFunc = func(ctx context.Context, id int64) (*model.InboundOrder, error) {
				return tt.mockOrder, tt.mockError
			}

			router.POST("/inbound-orders/:id/confirm", handler.Confirm)

			req := httptest.NewRequest("POST", "/inbound-orders/"+tt.orderID+"/confirm", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func floatPtrIO(f float64) *float64 {
	return &f
}
