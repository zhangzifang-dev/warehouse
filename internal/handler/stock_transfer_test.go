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

type mockStockTransferService struct {
	listFunc    func(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) (*service.ListStockTransfersResult, error)
	getByIDFunc func(ctx context.Context, id int64) (*model.StockTransfer, error)
	createFunc  func(ctx context.Context, input *service.CreateStockTransferInput) (*model.StockTransfer, error)
	updateFunc  func(ctx context.Context, id int64, input *service.UpdateStockTransferInput) (*model.StockTransfer, error)
	deleteFunc  func(ctx context.Context, id int64) error
	confirmFunc func(ctx context.Context, id int64) (*model.StockTransfer, error)
}

func (m *mockStockTransferService) List(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) (*service.ListStockTransfersResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, fromWarehouseID, toWarehouseID, status)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStockTransferService) GetByID(ctx context.Context, id int64) (*model.StockTransfer, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStockTransferService) Create(ctx context.Context, input *service.CreateStockTransferInput) (*model.StockTransfer, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStockTransferService) Update(ctx context.Context, id int64, input *service.UpdateStockTransferInput) (*model.StockTransfer, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStockTransferService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockStockTransferService) Confirm(ctx context.Context, id int64) (*model.StockTransfer, error) {
	if m.confirmFunc != nil {
		return m.confirmFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func setupStockTransferHandlerTest(t *testing.T) (*gin.Engine, *StockTransferHandler, *mockStockTransferService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockStockTransferService{}
	handler := NewStockTransferHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestStockTransferHandler_List(t *testing.T) {
	tests := []struct {
		name            string
		mockTransfers   []model.StockTransfer
		mockTotal       int
		mockError       error
		queryFromWH     string
		queryToWH       string
		queryStatus     string
		queryPage       string
		querySize       string
		wantStatus      int
		wantTotal       int
	}{
		{
			name: "success with default pagination",
			mockTransfers: []model.StockTransfer{
				{BaseModel: model.BaseModel{ID: 1}, OrderNo: "ST-2024-001", FromWarehouseID: 1, ToWarehouseID: 2, TotalQuantity: 100},
				{BaseModel: model.BaseModel{ID: 2}, OrderNo: "ST-2024-002", FromWarehouseID: 1, ToWarehouseID: 3, TotalQuantity: 200},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with warehouse filters",
			mockTransfers: []model.StockTransfer{
				{BaseModel: model.BaseModel{ID: 1}},
			},
			mockTotal:   1,
			queryFromWH: "1",
			queryToWH:   "2",
			wantStatus:  http.StatusOK,
			wantTotal:   1,
		},
		{
			name: "success with status filter",
			mockTransfers: []model.StockTransfer{
				{BaseModel: model.BaseModel{ID: 1}},
			},
			mockTotal:   1,
			queryStatus: "1",
			wantStatus:  http.StatusOK,
			wantTotal:   1,
		},
		{
			name:         "empty list",
			mockTransfers: []model.StockTransfer{},
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
			router, handler, mockSvc := setupStockTransferHandlerTest(t)
			mockSvc.listFunc = func(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) (*service.ListStockTransfersResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				return &service.ListStockTransfersResult{
					Transfers: tt.mockTransfers,
					Total:     tt.mockTotal,
				}, nil
			}

			router.GET("/stock-transfers", handler.List)

			req := httptest.NewRequest("GET", "/stock-transfers?from_warehouse_id="+tt.queryFromWH+"&to_warehouse_id="+tt.queryToWH+"&status="+tt.queryStatus+"&page="+tt.queryPage+"&size="+tt.querySize, nil)
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

func TestStockTransferHandler_GetByID(t *testing.T) {
	tests := []struct {
		name         string
		transferID   string
		mockTransfer *model.StockTransfer
		mockError    error
		wantStatus   int
	}{
		{
			name:         "success",
			transferID:   "1",
			mockTransfer: &model.StockTransfer{BaseModel: model.BaseModel{ID: 1}, OrderNo: "ST-2024-001", FromWarehouseID: 1, ToWarehouseID: 2, TotalQuantity: 100},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "invalid id",
			transferID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			transferID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "stock transfer not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupStockTransferHandlerTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.StockTransfer, error) {
				return tt.mockTransfer, tt.mockError
			}

			router.GET("/stock-transfers/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/stock-transfers/"+tt.transferID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestStockTransferHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		mockTransfer *model.StockTransfer
		mockError    error
		wantStatus   int
	}{
		{
			name: "success",
			body: CreateStockTransferRequest{
				OrderNo:         "ST-2024-001",
				FromWarehouseID: 1,
				ToWarehouseID:   2,
				TotalQuantity:   100,
			},
			mockTransfer: &model.StockTransfer{BaseModel: model.BaseModel{ID: 1}, OrderNo: "ST-2024-001", FromWarehouseID: 1, ToWarehouseID: 2, TotalQuantity: 100},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateStockTransferRequest{},
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
			router, handler, mockSvc := setupStockTransferHandlerTest(t)
			mockSvc.createFunc = func(ctx context.Context, input *service.CreateStockTransferInput) (*model.StockTransfer, error) {
				return tt.mockTransfer, tt.mockError
			}

			router.POST("/stock-transfers", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/stock-transfers", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestStockTransferHandler_Update(t *testing.T) {
	tests := []struct {
		name         string
		transferID   string
		body         interface{}
		mockTransfer *model.StockTransfer
		mockError    error
		wantStatus   int
	}{
		{
			name:       "success",
			transferID: "1",
			body: UpdateStockTransferRequest{
				TotalQuantity: floatPtrSTH(200),
			},
			mockTransfer: &model.StockTransfer{BaseModel: model.BaseModel{ID: 1}, TotalQuantity: 200},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "invalid id",
			transferID: "invalid",
			body:       UpdateStockTransferRequest{TotalQuantity: floatPtrSTH(200)},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			transferID: "999",
			body:       UpdateStockTransferRequest{TotalQuantity: floatPtrSTH(200)},
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "stock transfer not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			transferID: "1",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupStockTransferHandlerTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, input *service.UpdateStockTransferInput) (*model.StockTransfer, error) {
				return tt.mockTransfer, tt.mockError
			}

			router.PUT("/stock-transfers/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/stock-transfers/"+tt.transferID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestStockTransferHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		transferID string
		mockError  error
		wantStatus int
	}{
		{
			name:       "success",
			transferID: "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			transferID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			transferID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "stock transfer not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupStockTransferHandlerTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/stock-transfers/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/stock-transfers/"+tt.transferID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestStockTransferHandler_Confirm(t *testing.T) {
	tests := []struct {
		name         string
		transferID   string
		mockTransfer *model.StockTransfer
		mockError    error
		wantStatus   int
	}{
		{
			name:         "success",
			transferID:   "1",
			mockTransfer: &model.StockTransfer{BaseModel: model.BaseModel{ID: 1}, Status: 1},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "invalid id",
			transferID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "already completed",
			transferID: "1",
			mockError:  apperrors.NewAppError(apperrors.CodeBadRequest, "transfer already completed"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			transferID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "stock transfer not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupStockTransferHandlerTest(t)
			mockSvc.confirmFunc = func(ctx context.Context, id int64) (*model.StockTransfer, error) {
				return tt.mockTransfer, tt.mockError
			}

			router.POST("/stock-transfers/:id/confirm", handler.Confirm)

			req := httptest.NewRequest("POST", "/stock-transfers/"+tt.transferID+"/confirm", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func floatPtrSTH(f float64) *float64 {
	return &f
}
