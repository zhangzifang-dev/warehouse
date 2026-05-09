package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"warehouse/internal/model"
)

func TestDashboardRepository_GetOverviewStats(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo := NewDashboardRepository(testDB)
	ctx := context.Background()

	stats, err := repo.GetOverviewStats(ctx)
	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.GreaterOrEqual(t, stats.TotalInventory, float64(0))
	assert.GreaterOrEqual(t, stats.InventoryWarning, 0)
	assert.GreaterOrEqual(t, stats.TodayInbound, 0)
	assert.GreaterOrEqual(t, stats.TodayOutbound, 0)
}

func TestDashboardRepository_GetTrendData(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo := NewDashboardRepository(testDB)
	ctx := context.Background()

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
	}

	trend, err := repo.GetTrendData(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, trend)

	assert.GreaterOrEqual(t, len(trend), 0)
}

func TestDashboardRepository_GetTopProducts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo := NewDashboardRepository(testDB)
	ctx := context.Background()

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     10,
	}

	products, err := repo.GetTopProducts(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, products)

	assert.LessOrEqual(t, len(products), 10)
}

func TestDashboardRepository_GetWarehouseUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo := NewDashboardRepository(testDB)
	ctx := context.Background()

	usage, err := repo.GetWarehouseUsage(ctx)
	require.NoError(t, err)
	require.NotNil(t, usage)
}

func TestDashboardRepository_GetSupplierPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo := NewDashboardRepository(testDB)
	ctx := context.Background()

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     10,
	}

	performance, err := repo.GetSupplierPerformance(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, performance)
}

func TestDashboardRepository_GetPendingOrders(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo := NewDashboardRepository(testDB)
	ctx := context.Background()

	pending, err := repo.GetPendingOrders(ctx)
	require.NoError(t, err)
	require.NotNil(t, pending)

	assert.GreaterOrEqual(t, pending.InboundPending, 0)
	assert.GreaterOrEqual(t, pending.OutboundPending, 0)
	assert.GreaterOrEqual(t, pending.TransferPending, 0)
}
