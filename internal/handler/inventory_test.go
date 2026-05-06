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

type mockInventoryService struct {
	listFunc            func(ctx context.Context, filter *model.InventoryQueryFilter) (*service.ListInventoriesResult, error)
	getByIDFunc         func(ctx context.Context, id int64) (*model.Inventory, error)
	createFunc          func(ctx context.Context, input *service.CreateInventoryInput) (*model.Inventory, error)
	updateFunc          func(ctx context.Context, id int64, input *service.UpdateInventoryInput) (*model.Inventory, error)
	deleteFunc          func(ctx context.Context, id int64) error
	adjustQuantityFunc  func(ctx context.Context, input *service.AdjustQuantityInput) (*model.Inventory, error)
	checkStockFunc      func(ctx context.Context, input *service.CheckStockInput) (*service.CheckStockResult, error)
}

func (m *mockInventoryService) List(ctx context.Context, filter *model.InventoryQueryFilter) (*service.ListInventoriesResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryService) GetByID(ctx context.Context, id int64) (*model.Inventory, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryService) Create(ctx context.Context, input *service.CreateInventoryInput) (*model.Inventory, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryService) Update(ctx context.Context, id int64, input *service.UpdateInventoryInput) (*model.Inventory, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockInventoryService) AdjustQuantity(ctx context.Context, input *service.AdjustQuantityInput) (*model.Inventory, error) {
	if m.adjustQuantityFunc != nil {
		return m.adjustQuantityFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryService) CheckStock(ctx context.Context, input *service.CheckStockInput) (*service.CheckStockResult, error) {
	if m.checkStockFunc != nil {
		return m.checkStockFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func setupInventoryHandlerTest(t *testing.T) (*gin.Engine, *InventoryHandler, *mockInventoryService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockInventoryService{}
	handler := NewInventoryHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestInventoryHandler_List(t *testing.T) {
	tests := []struct {
		name              string
		mockInventories   []model.Inventory
		mockTotal         int
		mockError         error
		queryProductName  string
		queryQuantityMin  string
		queryQuantityMax  string
		queryBatchNo      string
		queryProductID    string
		queryWarehouseID  string
		queryPage         string
		querySize         string
		wantStatus        int
		wantTotal         int
	}{
		{
			name: "success with default pagination",
			mockInventories: []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, WarehouseID: 1, ProductID: 1, Quantity: 100},
				{BaseModel: model.BaseModel{ID: 2}, WarehouseID: 1, ProductID: 2, Quantity: 200},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with product name filter",
			mockInventories: []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, ProductID: 1, Quantity: 100},
			},
			mockTotal:        1,
			queryProductName: "测试商品",
			wantStatus:       http.StatusOK,
			wantTotal:        1,
		},
		{
			name: "success with quantity min filter",
			mockInventories: []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, Quantity: 150},
			},
			mockTotal:        1,
			queryQuantityMin: "100",
			wantStatus:       http.StatusOK,
			wantTotal:        1,
		},
		{
			name: "success with quantity max filter",
			mockInventories: []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, Quantity: 50},
			},
			mockTotal:        1,
			queryQuantityMax: "100",
			wantStatus:       http.StatusOK,
			wantTotal:        1,
		},
		{
			name: "success with quantity range filter",
			mockInventories: []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, Quantity: 100},
				{BaseModel: model.BaseModel{ID: 2}, Quantity: 150},
			},
			mockTotal:        2,
			queryQuantityMin: "50",
			queryQuantityMax: "200",
			wantStatus:       http.StatusOK,
			wantTotal:        2,
		},
		{
			name: "success with batch no filter",
			mockInventories: []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, BatchNo: "BATCH001", Quantity: 100},
			},
			mockTotal:     1,
			queryBatchNo:  "BATCH001",
			wantStatus:    http.StatusOK,
			wantTotal:     1,
		},
		{
			name: "success with product id filter",
			mockInventories: []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, ProductID: 1, Quantity: 100},
			},
			mockTotal:       1,
			queryProductID:  "1",
			wantStatus:      http.StatusOK,
			wantTotal:       1,
		},
		{
			name: "success with warehouse id filter",
			mockInventories: []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, WarehouseID: 1, Quantity: 100},
			},
			mockTotal:        1,
			queryWarehouseID: "1",
			wantStatus:       http.StatusOK,
			wantTotal:        1,
		},
		{
			name: "success with all filters",
			mockInventories: []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, ProductID: 1, BatchNo: "BATCH001", Quantity: 100},
			},
			mockTotal:        1,
			queryProductName: "测试商品",
			queryQuantityMin: "50",
			queryQuantityMax: "200",
			queryBatchNo:     "BATCH001",
			wantStatus:       http.StatusOK,
			wantTotal:        1,
		},
		{
			name:            "empty list",
			mockInventories: []model.Inventory{},
			mockTotal:       0,
			wantStatus:      http.StatusOK,
			wantTotal:       0,
		},
		{
			name:       "service error",
			mockError:  apperrors.NewAppError(apperrors.CodeInternalError, "database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInventoryHandlerTest(t)
			mockSvc.listFunc = func(ctx context.Context, filter *model.InventoryQueryFilter) (*service.ListInventoriesResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				return &service.ListInventoriesResult{
					Inventories: tt.mockInventories,
					Total:       tt.mockTotal,
				}, nil
			}

			router.GET("/inventory", handler.List)

			req := httptest.NewRequest("GET", "/inventory?product_name="+tt.queryProductName+"&product_id="+tt.queryProductID+"&warehouse_id="+tt.queryWarehouseID+"&quantity_min="+tt.queryQuantityMin+"&quantity_max="+tt.queryQuantityMax+"&batch_no="+tt.queryBatchNo+"&page="+tt.queryPage+"&size="+tt.querySize, nil)
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

func TestInventoryHandler_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		inventoryID   string
		mockInventory *model.Inventory
		mockError     error
		wantStatus    int
	}{
		{
			name:          "success",
			inventoryID:   "1",
			mockInventory: &model.Inventory{BaseModel: model.BaseModel{ID: 1}, WarehouseID: 1, ProductID: 1, Quantity: 100},
			wantStatus:    http.StatusOK,
		},
		{
			name:        "invalid id",
			inventoryID: "invalid",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "not found",
			inventoryID: "999",
			mockError:   apperrors.NewAppError(apperrors.CodeNotFound, "inventory not found"),
			wantStatus:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInventoryHandlerTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.Inventory, error) {
				return tt.mockInventory, tt.mockError
			}

			router.GET("/inventory/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/inventory/"+tt.inventoryID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestInventoryHandler_Create(t *testing.T) {
	tests := []struct {
		name          string
		body          interface{}
		mockInventory *model.Inventory
		mockError     error
		wantStatus    int
	}{
		{
			name: "success",
			body: CreateInventoryRequest{
				WarehouseID: 1,
				ProductID:   1,
				Quantity:    invFloatPtrH(100),
			},
			mockInventory: &model.Inventory{BaseModel: model.BaseModel{ID: 1}, WarehouseID: 1, ProductID: 1, Quantity: 100},
			wantStatus:    http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateInventoryRequest{},
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
			router, handler, mockSvc := setupInventoryHandlerTest(t)
			mockSvc.createFunc = func(ctx context.Context, input *service.CreateInventoryInput) (*model.Inventory, error) {
				return tt.mockInventory, tt.mockError
			}

			router.POST("/inventory", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/inventory", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestInventoryHandler_Update(t *testing.T) {
	tests := []struct {
		name          string
		inventoryID   string
		body          interface{}
		mockInventory *model.Inventory
		mockError     error
		wantStatus    int
	}{
		{
			name:        "success",
			inventoryID: "1",
			body: UpdateInventoryRequest{
				Quantity: invFloatPtrH(200),
			},
			mockInventory: &model.Inventory{BaseModel: model.BaseModel{ID: 1}, Quantity: 200},
			wantStatus:    http.StatusOK,
		},
		{
			name:        "invalid id",
			inventoryID: "invalid",
			body:        UpdateInventoryRequest{Quantity: invFloatPtrH(200)},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "not found",
			inventoryID: "999",
			body:        UpdateInventoryRequest{Quantity: invFloatPtrH(200)},
			mockError:   apperrors.NewAppError(apperrors.CodeNotFound, "inventory not found"),
			wantStatus:  http.StatusNotFound,
		},
		{
			name:        "invalid json",
			inventoryID: "1",
			body:        "invalid",
			wantStatus:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInventoryHandlerTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, input *service.UpdateInventoryInput) (*model.Inventory, error) {
				return tt.mockInventory, tt.mockError
			}

			router.PUT("/inventory/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/inventory/"+tt.inventoryID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestInventoryHandler_Delete(t *testing.T) {
	tests := []struct {
		name        string
		inventoryID string
		mockError   error
		wantStatus  int
	}{
		{
			name:        "success",
			inventoryID: "1",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "invalid id",
			inventoryID: "invalid",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "not found",
			inventoryID: "999",
			mockError:   apperrors.NewAppError(apperrors.CodeNotFound, "inventory not found"),
			wantStatus:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInventoryHandlerTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/inventory/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/inventory/"+tt.inventoryID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestInventoryHandler_AdjustQuantity(t *testing.T) {
	tests := []struct {
		name          string
		body          interface{}
		mockInventory *model.Inventory
		mockError     error
		wantStatus    int
	}{
		{
			name: "increase success",
			body: AdjustQuantityRequest{
				InventoryID: 1,
				Quantity:    50,
			},
			mockInventory: &model.Inventory{BaseModel: model.BaseModel{ID: 1}, Quantity: 150},
			wantStatus:    http.StatusOK,
		},
		{
			name: "decrease success",
			body: AdjustQuantityRequest{
				InventoryID: 1,
				Quantity:    -50,
			},
			mockInventory: &model.Inventory{BaseModel: model.BaseModel{ID: 1}, Quantity: 50},
			wantStatus:    http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       AdjustQuantityRequest{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "insufficient stock",
			body: AdjustQuantityRequest{
				InventoryID: 1,
				Quantity:    -100,
			},
			mockError:  apperrors.NewAppError(apperrors.CodeInsufficientStock, "insufficient stock"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			body: AdjustQuantityRequest{
				InventoryID: 999,
				Quantity:    50,
			},
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "inventory not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInventoryHandlerTest(t)
			mockSvc.adjustQuantityFunc = func(ctx context.Context, input *service.AdjustQuantityInput) (*model.Inventory, error) {
				return tt.mockInventory, tt.mockError
			}

			router.POST("/inventory/adjust", handler.AdjustQuantity)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/inventory/adjust", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestInventoryHandler_CheckStock(t *testing.T) {
	tests := []struct {
		name        string
		body        interface{}
		mockResult  *service.CheckStockResult
		mockError   error
		wantStatus  int
		wantAvail   bool
	}{
		{
			name: "available",
			body: CheckStockRequest{
				WarehouseID: 1,
				ProductID:   1,
				BatchNo:     "BATCH001",
				Quantity:    50,
			},
			mockResult: &service.CheckStockResult{
				Available:    true,
				CurrentStock: 100,
				Requested:    50,
			},
			wantStatus: http.StatusOK,
			wantAvail:  true,
		},
		{
			name: "not available",
			body: CheckStockRequest{
				WarehouseID: 1,
				ProductID:   1,
				Quantity:    150,
			},
			mockResult: &service.CheckStockResult{
				Available:    false,
				CurrentStock: 100,
				Requested:    150,
			},
			wantStatus: http.StatusOK,
			wantAvail:  false,
		},
		{
			name:       "missing required fields",
			body:       CheckStockRequest{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupInventoryHandlerTest(t)
			mockSvc.checkStockFunc = func(ctx context.Context, input *service.CheckStockInput) (*service.CheckStockResult, error) {
				return tt.mockResult, tt.mockError
			}

			router.POST("/inventory/check", handler.CheckStock)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/inventory/check", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, tt.wantAvail, data["available"])
			}
		})
	}
}

func invFloatPtrH(f float64) *float64 {
	return &f
}

func invIntPtrH(i int64) *int64 {
	return &i
}

func invStrPtrH(s string) *string {
	return &s
}
