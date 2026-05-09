package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"warehouse/internal/model"
)

type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetOverview(ctx context.Context) (*model.OverviewStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OverviewStats), args.Error(1)
}

func (m *MockDashboardService) GetTrendData(ctx context.Context, startDate, endDate time.Time) ([]model.TrendData, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TrendData), args.Error(1)
}

func (m *MockDashboardService) GetTopProducts(ctx context.Context, startDate, endDate time.Time, limit int) ([]model.TopProduct, error) {
	args := m.Called(ctx, startDate, endDate, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TopProduct), args.Error(1)
}

func (m *MockDashboardService) GetWarehouseUsage(ctx context.Context) ([]model.WarehouseUsage, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.WarehouseUsage), args.Error(1)
}

func (m *MockDashboardService) GetSupplierPerformance(ctx context.Context, startDate, endDate time.Time, limit int) ([]model.SupplierPerformance, error) {
	args := m.Called(ctx, startDate, endDate, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SupplierPerformance), args.Error(1)
}

func (m *MockDashboardService) GetPendingOrders(ctx context.Context) (*model.PendingOrders, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PendingOrders), args.Error(1)
}

func TestDashboardHandler_GetOverview(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)

	router := gin.New()
	router.GET("/dashboard/overview", handler.GetOverview)

	expectedStats := &model.OverviewStats{
		TotalInventory:   1000.0,
		InventoryWarning: 5,
		TodayInbound:     10,
		TodayOutbound:    15,
	}

	mockService.On("GetOverview", mock.Anything).Return(expectedStats, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/overview", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_GetTrendData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)

	router := gin.New()
	router.GET("/dashboard/trend", handler.GetTrendData)

	expectedTrend := []model.TrendData{
		{Date: "2026-04-08", InboundQty: 100, OutboundQty: 80},
	}

	mockService.On("GetTrendData", mock.Anything, mock.Anything, mock.Anything).Return(expectedTrend, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/trend?start_date=2026-04-01&end_date=2026-05-01", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_GetTopProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)

	router := gin.New()
	router.GET("/dashboard/top-products", handler.GetTopProducts)

	expectedProducts := []model.TopProduct{
		{ProductID: 1, ProductName: "Product A", TotalQty: 500},
	}

	mockService.On("GetTopProducts", mock.Anything, mock.Anything, mock.Anything, 10).Return(expectedProducts, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/top-products?start_date=2026-04-01&end_date=2026-05-01", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}
