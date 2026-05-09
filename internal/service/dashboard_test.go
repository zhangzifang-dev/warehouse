package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"warehouse/internal/model"
)

type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) GetOverviewStats(ctx context.Context) (*model.OverviewStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OverviewStats), args.Error(1)
}

func (m *MockDashboardRepository) GetTrendData(ctx context.Context, params *model.DashboardQueryParams) ([]model.TrendData, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TrendData), args.Error(1)
}

func (m *MockDashboardRepository) GetTopProducts(ctx context.Context, params *model.DashboardQueryParams) ([]model.TopProduct, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TopProduct), args.Error(1)
}

func (m *MockDashboardRepository) GetWarehouseUsage(ctx context.Context) ([]model.WarehouseUsage, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.WarehouseUsage), args.Error(1)
}

func (m *MockDashboardRepository) GetSupplierPerformance(ctx context.Context, params *model.DashboardQueryParams) ([]model.SupplierPerformance, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SupplierPerformance), args.Error(1)
}

func (m *MockDashboardRepository) GetPendingOrders(ctx context.Context) (*model.PendingOrders, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PendingOrders), args.Error(1)
}

func TestDashboardService_GetOverview(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo)

	ctx := context.Background()
	expectedStats := &model.OverviewStats{
		TotalInventory:   1000.0,
		InventoryWarning: 5,
		TodayInbound:     10,
		TodayOutbound:    15,
	}

	mockRepo.On("GetOverviewStats", ctx).Return(expectedStats, nil)

	stats, err := service.GetOverview(ctx)
	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, expectedStats.TotalInventory, stats.TotalInventory)
	assert.Equal(t, expectedStats.InventoryWarning, stats.InventoryWarning)
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetTrendData(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo)

	ctx := context.Background()
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
	}

	expectedTrend := []model.TrendData{
		{Date: "2026-04-08", InboundQty: 100, OutboundQty: 80},
		{Date: "2026-04-09", InboundQty: 120, OutboundQty: 90},
	}

	mockRepo.On("GetTrendData", ctx, params).Return(expectedTrend, nil)

	trend, err := service.GetTrendData(ctx, startDate, endDate)
	require.NoError(t, err)
	require.NotNil(t, trend)

	assert.Len(t, trend, 2)
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetTopProducts(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo)

	ctx := context.Background()
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     10,
	}

	expectedProducts := []model.TopProduct{
		{ProductID: 1, ProductName: "Product A", TotalQty: 500, OrderCount: 20},
	}

	mockRepo.On("GetTopProducts", ctx, params).Return(expectedProducts, nil)

	products, err := service.GetTopProducts(ctx, startDate, endDate, 10)
	require.NoError(t, err)
	require.NotNil(t, products)

	assert.Len(t, products, 1)
	mockRepo.AssertExpectations(t)
}
