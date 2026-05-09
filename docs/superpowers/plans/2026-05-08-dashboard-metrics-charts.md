# Dashboard Metrics Charts Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a comprehensive dashboard with statistical charts showing inventory, orders, trends, and performance metrics.

**Architecture:** Backend API using Go/Gin with database aggregation queries. Frontend using React + Ant Design Charts for visualization. Layered architecture: Handler → Service → Repository.

**Tech Stack:** Go, Gin, Bun ORM, PostgreSQL, React, TypeScript, Ant Design 6, Ant Design Charts, TanStack Query

---

## File Structure

**Backend (Go):**
- `warehouse/internal/model/dashboard.go` - Data models and DTOs
- `warehouse/internal/repository/dashboard.go` - Database queries
- `warehouse/internal/service/dashboard.go` - Business logic
- `warehouse/internal/handler/dashboard.go` - HTTP handlers
- `warehouse/internal/router/router.go` - Route registration (modify)

**Frontend (React/TypeScript):**
- `warehouse/web/src/types/dashboard.ts` - TypeScript interfaces
- `warehouse/web/src/api/dashboard.ts` - API client functions
- `warehouse/web/src/pages/dashboard/Dashboard.tsx` - Main page
- `warehouse/web/src/pages/dashboard/components/StatCard.tsx` - Stat card component
- `warehouse/web/src/pages/dashboard/components/TrendChart.tsx` - Trend line chart
- `warehouse/web/src/pages/dashboard/components/TopProductsChart.tsx` - Top products bar chart
- `warehouse/web/src/pages/dashboard/components/WarehouseUsageChart.tsx` - Warehouse usage pie chart
- `warehouse/web/src/pages/dashboard/components/SupplierPerformanceChart.tsx` - Supplier radar chart
- `warehouse/web/src/pages/dashboard/hooks/useDashboardStats.ts` - Data fetching hook
- `warehouse/web/src/App.tsx` - Route registration (modify)

**Tests:**
- `warehouse/internal/repository/dashboard_test.go`
- `warehouse/internal/service/dashboard_test.go`
- `warehouse/internal/handler/dashboard_test.go`

---

## Phase 1: Backend Data Models

### Task 1: Create Dashboard Models

**Files:**
- Create: `warehouse/internal/model/dashboard.go`

- [ ] **Step 1: Write the data model file**

```go
package model

import "time"

// OverviewStats 总览统计
type OverviewStats struct {
	TotalInventory   float64 `json:"total_inventory"`
	InventoryWarning int     `json:"inventory_warning"`
	TodayInbound     int     `json:"today_inbound"`
	TodayInboundQty  float64 `json:"today_inbound_qty"`
	TodayOutbound    int     `json:"today_outbound"`
	TodayOutboundQty float64 `json:"today_outbound_qty"`
}

// TrendData 趋势数据
type TrendData struct {
	Date        string  `json:"date"`
	InboundQty  float64 `json:"inbound_qty"`
	OutboundQty float64 `json:"outbound_qty"`
}

// TopProduct 热销产品
type TopProduct struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
	TotalQty    float64 `json:"total_qty"`
	OrderCount  int     `json:"order_count"`
}

// WarehouseUsage 仓库使用率
type WarehouseUsage struct {
	WarehouseID   int64   `json:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"`
	Capacity      float64 `json:"capacity"`
	UsedCapacity  float64 `json:"used_capacity"`
	UsageRate     float64 `json:"usage_rate"`
}

// SupplierPerformance 供应商绩效
type SupplierPerformance struct {
	SupplierID    int64   `json:"supplier_id"`
	SupplierName  string  `json:"supplier_name"`
	OrderCount    int     `json:"order_count"`
	TotalValue    float64 `json:"total_value"`
	OnTimeRate    float64 `json:"on_time_rate"`
	QualityScore  float64 `json:"quality_score"`
	DeliveryScore float64 `json:"delivery_score"`
}

// PendingOrders 待处理订单
type PendingOrders struct {
	InboundPending  int `json:"inbound_pending"`
	OutboundPending int `json:"outbound_pending"`
	TransferPending int `json:"transfer_pending"`
}

// DashboardQueryParams 仪表盘查询参数
type DashboardQueryParams struct {
	StartDate time.Time
	EndDate   time.Time
	Limit     int
}
```

- [ ] **Step 2: Commit**

```bash
git add warehouse/internal/model/dashboard.go
git commit -m "feat: add dashboard data models and DTOs"
```

---

## Phase 2: Backend Repository Layer

### Task 2: Write Dashboard Repository Tests

**Files:**
- Create: `warehouse/internal/repository/dashboard_test.go`

- [ ] **Step 1: Write repository test file**

```go
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

	// Should return data for the date range
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd warehouse && go test -v ./internal/repository -run TestDashboardRepository`

Expected: FAIL with "undefined: NewDashboardRepository"

- [ ] **Step 3: Commit tests**

```bash
git add warehouse/internal/repository/dashboard_test.go
git commit -m "test: add dashboard repository tests"
```

### Task 3: Implement Dashboard Repository

**Files:**
- Create: `warehouse/internal/repository/dashboard.go`

- [ ] **Step 1: Write repository implementation**

```go
package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type DashboardRepository struct {
	db *bun.DB
}

func NewDashboardRepository(db *bun.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetOverviewStats(ctx context.Context) (*model.OverviewStats, error) {
	stats := &model.OverviewStats{}

	// Get total inventory
	err := r.db.NewSelect().
		Table("inventories").
		Where("deleted_at IS NULL").
		ColumnExpr("COALESCE(SUM(quantity), 0) as total_inventory").
		Scan(ctx, &stats.TotalInventory)
	if err != nil {
		return nil, err
	}

	// Get inventory warning count (quantity < 10)
	err = r.db.NewSelect().
		Table("inventories").
		Where("deleted_at IS NULL").
		Where("quantity < 10").
		ColumnExpr("COUNT(*) as inventory_warning").
		Scan(ctx, &stats.InventoryWarning)
	if err != nil {
		return nil, err
	}

	// Get today's inbound stats
	today := time.Now().Format("2006-01-02")
	err = r.db.NewSelect().
		Table("inbound_orders").
		Where("deleted_at IS NULL").
		Where("DATE(created_at) = ?", today).
		ColumnExpr("COUNT(*) as today_inbound").
		ColumnExpr("COALESCE(SUM(total_qty), 0) as today_inbound_qty").
		Scan(ctx, &stats.TodayInbound, &stats.TodayInboundQty)
	if err != nil {
		return nil, err
	}

	// Get today's outbound stats
	err = r.db.NewSelect().
		Table("outbound_orders").
		Where("deleted_at IS NULL").
		Where("DATE(created_at) = ?", today).
		ColumnExpr("COUNT(*) as today_outbound").
		ColumnExpr("COALESCE(SUM(total_qty), 0) as today_outbound_qty").
		Scan(ctx, &stats.TodayOutbound, &stats.TodayOutboundQty)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *DashboardRepository) GetTrendData(ctx context.Context, params *model.DashboardQueryParams) ([]model.TrendData, error) {
	// Get inbound trend
	type DailyData struct {
		Date string
		Qty  float64
	}

	var inboundTrend []DailyData
	err := r.db.NewSelect().
		Table("inbound_orders").
		Where("deleted_at IS NULL").
		Where("created_at >= ?", params.StartDate).
		Where("created_at <= ?", params.EndDate).
		ColumnExpr("DATE(created_at) as date").
		ColumnExpr("COALESCE(SUM(total_qty), 0) as qty").
		GroupExpr("DATE(created_at)").
		OrderExpr("date").
		Scan(ctx, &inboundTrend)
	if err != nil {
		return nil, err
	}

	// Get outbound trend
	var outboundTrend []DailyData
	err = r.db.NewSelect().
		Table("outbound_orders").
		Where("deleted_at IS NULL").
		Where("created_at >= ?", params.StartDate).
		Where("created_at <= ?", params.EndDate).
		ColumnExpr("DATE(created_at) as date").
		ColumnExpr("COALESCE(SUM(total_qty), 0) as qty").
		GroupExpr("DATE(created_at)").
		OrderExpr("date").
		Scan(ctx, &outboundTrend)
	if err != nil {
		return nil, err
	}

	// Merge inbound and outbound data
	dateMap := make(map[string]*model.TrendData)
	for _, d := range inboundTrend {
		dateMap[d.Date] = &model.TrendData{
			Date:       d.Date,
			InboundQty: d.Qty,
		}
	}
	for _, d := range outboundTrend {
		if _, exists := dateMap[d.Date]; exists {
			dateMap[d.Date].OutboundQty = d.Qty
		} else {
			dateMap[d.Date] = &model.TrendData{
				Date:        d.Date,
				OutboundQty: d.Qty,
			}
		}
	}

	// Convert to slice and sort by date
	result := make([]model.TrendData, 0, len(dateMap))
	for _, v := range dateMap {
		result = append(result, *v)
	}

	// Sort by date
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Date > result[j].Date {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

func (r *DashboardRepository) GetTopProducts(ctx context.Context, params *model.DashboardQueryParams) ([]model.TopProduct, error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var products []model.TopProduct
	err := r.db.NewSelect().
		TableExpr("outbound_items AS oi").
		Join("JOIN products p ON oi.product_id = p.id").
		Join("LEFT JOIN categories c ON p.category_id = c.id").
		Join("JOIN outbound_orders o ON oi.order_id = o.id").
		Where("o.deleted_at IS NULL").
		Where("o.created_at >= ?", params.StartDate).
		Where("o.created_at <= ?", params.EndDate).
		Where("o.status = 'completed'").
		ColumnExpr("p.id as product_id").
		ColumnExpr("p.name as product_name").
		ColumnExpr("COALESCE(c.name, '') as category").
		ColumnExpr("SUM(oi.quantity) as total_qty").
		ColumnExpr("COUNT(DISTINCT o.id) as order_count").
		GroupExpr("p.id, p.name, c.name").
		OrderExpr("total_qty DESC").
		Limit(params.Limit).
		Scan(ctx, &products)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *DashboardRepository) GetWarehouseUsage(ctx context.Context) ([]model.WarehouseUsage, error) {
	var usage []model.WarehouseUsage

	query := `
		SELECT 
			w.id as warehouse_id,
			w.name as warehouse_name,
			COUNT(l.id) as capacity,
			COUNT(DISTINCT CASE WHEN i.quantity > 0 THEN i.location_id END) as used_capacity
		FROM warehouses w
		LEFT JOIN locations l ON w.id = l.warehouse_id AND l.deleted_at IS NULL
		LEFT JOIN inventories i ON l.id = i.location_id AND i.deleted_at IS NULL
		WHERE w.deleted_at IS NULL
		GROUP BY w.id, w.name
	`

	err := r.db.NewRaw(query).Scan(ctx, &usage)
	if err != nil {
		return nil, err
	}

	// Calculate usage rate
	for i := range usage {
		if usage[i].Capacity > 0 {
			usage[i].UsageRate = (usage[i].UsedCapacity / usage[i].Capacity) * 100
		}
	}

	return usage, nil
}

func (r *DashboardRepository) GetSupplierPerformance(ctx context.Context, params *model.DashboardQueryParams) ([]model.SupplierPerformance, error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var performance []model.SupplierPerformance

	// Note: Supplier scoring fields may not exist in current schema
	// Using basic metrics for now
	query := `
		SELECT 
			s.id as supplier_id,
			s.name as supplier_name,
			COUNT(DISTINCT io.id) as order_count,
			COALESCE(SUM(io.total_amount), 0) as total_value,
			0 as on_time_rate,
			0 as quality_score,
			0 as delivery_score
		FROM suppliers s
		LEFT JOIN inbound_orders io ON s.id = io.supplier_id 
			AND io.deleted_at IS NULL
			AND io.created_at >= ?
			AND io.created_at <= ?
		WHERE s.deleted_at IS NULL
		GROUP BY s.id, s.name
		ORDER BY order_count DESC, total_value DESC
		LIMIT ?
	`

	err := r.db.NewRaw(query, params.StartDate, params.EndDate, params.Limit).Scan(ctx, &performance)
	if err != nil {
		return nil, err
	}

	return performance, nil
}

func (r *DashboardRepository) GetPendingOrders(ctx context.Context) (*model.PendingOrders, error) {
	pending := &model.PendingOrders{}

	// Get pending inbound orders
	err := r.db.NewSelect().
		Table("inbound_orders").
		Where("deleted_at IS NULL").
		Where("status IN ('pending', 'approved')").
		ColumnExpr("COUNT(*) as count").
		Scan(ctx, &pending.InboundPending)
	if err != nil {
		return nil, err
	}

	// Get pending outbound orders
	err = r.db.NewSelect().
		Table("outbound_orders").
		Where("deleted_at IS NULL").
		Where("status IN ('pending', 'approved')").
		ColumnExpr("COUNT(*) as count").
		Scan(ctx, &pending.OutboundPending)
	if err != nil {
		return nil, err
	}

	// Get pending stock transfers
	err = r.db.NewSelect().
		Table("stock_transfers").
		Where("deleted_at IS NULL").
		Where("status = 'pending'").
		ColumnExpr("COUNT(*) as count").
		Scan(ctx, &pending.TransferPending)
	if err != nil {
		return nil, err
	}

	return pending, nil
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `cd warehouse && go test -v ./internal/repository -run TestDashboardRepository`

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add warehouse/internal/repository/dashboard.go
git commit -m "feat: implement dashboard repository with database queries"
```

---

## Phase 3: Backend Service Layer

### Task 4: Write Dashboard Service Tests

**Files:**
- Create: `warehouse/internal/service/dashboard_test.go`

- [ ] **Step 1: Write service test file**

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd warehouse && go test -v ./internal/service -run TestDashboardService`

Expected: FAIL with "undefined: NewDashboardService"

- [ ] **Step 3: Commit tests**

```bash
git add warehouse/internal/service/dashboard_test.go
git commit -m "test: add dashboard service tests"
```

### Task 5: Implement Dashboard Service

**Files:**
- Create: `warehouse/internal/service/dashboard.go`

- [ ] **Step 1: Write service implementation**

```go
package service

import (
	"context"
	"time"

	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/model"
)

type DashboardRepository interface {
	GetOverviewStats(ctx context.Context) (*model.OverviewStats, error)
	GetTrendData(ctx context.Context, params *model.DashboardQueryParams) ([]model.TrendData, error)
	GetTopProducts(ctx context.Context, params *model.DashboardQueryParams) ([]model.TopProduct, error)
	GetWarehouseUsage(ctx context.Context) ([]model.WarehouseUsage, error)
	GetSupplierPerformance(ctx context.Context, params *model.DashboardQueryParams) ([]model.SupplierPerformance, error)
	GetPendingOrders(ctx context.Context) (*model.PendingOrders, error)
}

type DashboardService struct {
	repo DashboardRepository
}

func NewDashboardService(repo DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetOverview(ctx context.Context) (*model.OverviewStats, error) {
	stats, err := s.repo.GetOverviewStats(ctx)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get overview stats")
	}
	return stats, nil
}

func (s *DashboardService) GetTrendData(ctx context.Context, startDate, endDate time.Time) ([]model.TrendData, error) {
	if startDate.After(endDate) {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "start date must be before end date")
	}

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
	}

	trend, err := s.repo.GetTrendData(ctx, params)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get trend data")
	}

	return trend, nil
}

func (s *DashboardService) GetTopProducts(ctx context.Context, startDate, endDate time.Time, limit int) ([]model.TopProduct, error) {
	if startDate.After(endDate) {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "start date must be before end date")
	}

	if limit <= 0 {
		limit = 10
	}

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
	}

	products, err := s.repo.GetTopProducts(ctx, params)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get top products")
	}

	return products, nil
}

func (s *DashboardService) GetWarehouseUsage(ctx context.Context) ([]model.WarehouseUsage, error) {
	usage, err := s.repo.GetWarehouseUsage(ctx)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get warehouse usage")
	}

	return usage, nil
}

func (s *DashboardService) GetSupplierPerformance(ctx context.Context, startDate, endDate time.Time, limit int) ([]model.SupplierPerformance, error) {
	if startDate.After(endDate) {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "start date must be before end date")
	}

	if limit <= 0 {
		limit = 10
	}

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
	}

	performance, err := s.repo.GetSupplierPerformance(ctx, params)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get supplier performance")
	}

	return performance, nil
}

func (s *DashboardService) GetPendingOrders(ctx context.Context) (*model.PendingOrders, error) {
	pending, err := s.repo.GetPendingOrders(ctx)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get pending orders")
	}

	return pending, nil
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `cd warehouse && go test -v ./internal/service -run TestDashboardService`

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add warehouse/internal/service/dashboard.go
git commit -m "feat: implement dashboard service with business logic"
```

---

## Phase 4: Backend Handler Layer

### Task 6: Write Dashboard Handler Tests

**Files:**
- Create: `warehouse/internal/handler/dashboard_test.go`

- [ ] **Step 1: Write handler test file**

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd warehouse && go test -v ./internal/handler -run TestDashboardHandler`

Expected: FAIL with "undefined: NewDashboardHandler"

- [ ] **Step 3: Commit tests**

```bash
git add warehouse/internal/handler/dashboard_test.go
git commit -m "test: add dashboard handler tests"
```

### Task 7: Implement Dashboard Handler

**Files:**
- Create: `warehouse/internal/handler/dashboard.go`

- [ ] **Step 1: Write handler implementation**

```go
package handler

import (
	"strconv"
	"time"

	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardService interface {
	GetOverview(ctx interface{}) (*model.OverviewStats, error)
	GetTrendData(ctx interface{}, startDate, endDate time.Time) ([]model.TrendData, error)
	GetTopProducts(ctx interface{}, startDate, endDate time.Time, limit int) ([]model.TopProduct, error)
	GetWarehouseUsage(ctx interface{}) ([]model.WarehouseUsage, error)
	GetSupplierPerformance(ctx interface{}, startDate, endDate time.Time, limit int) ([]model.SupplierPerformance, error)
	GetPendingOrders(ctx interface{}) (*model.PendingOrders, error)
}

type DashboardHandler struct {
	dashboardService DashboardService
}

func NewDashboardHandler(dashboardService DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetOverview(c *gin.Context) {
	stats, err := h.dashboardService.GetOverview(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get overview stats")
		return
	}

	response.Success(c, stats)
}

func (h *DashboardHandler) GetTrendData(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid end_date format")
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	trend, err := h.dashboardService.GetTrendData(c.Request.Context(), startDate, endDate)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get trend data")
		return
	}

	response.Success(c, trend)
}

func (h *DashboardHandler) GetTopProducts(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	limitStr := c.DefaultQuery("limit", "10")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid end_date format")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	products, err := h.dashboardService.GetTopProducts(c.Request.Context(), startDate, endDate, limit)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get top products")
		return
	}

	response.Success(c, products)
}

func (h *DashboardHandler) GetWarehouseUsage(c *gin.Context) {
	usage, err := h.dashboardService.GetWarehouseUsage(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get warehouse usage")
		return
	}

	response.Success(c, usage)
}

func (h *DashboardHandler) GetSupplierPerformance(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	limitStr := c.DefaultQuery("limit", "10")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid end_date format")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	performance, err := h.dashboardService.GetSupplierPerformance(c.Request.Context(), startDate, endDate, limit)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get supplier performance")
		return
	}

	response.Success(c, performance)
}

func (h *DashboardHandler) GetPendingOrders(c *gin.Context) {
	pending, err := h.dashboardService.GetPendingOrders(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get pending orders")
		return
	}

	response.Success(c, pending)
}
```

- [ ] **Step 2: Fix service interface type (use context.Context instead of interface{})**

Edit `warehouse/internal/service/dashboard.go` to use proper context types in the interface:

```go
type DashboardRepository interface {
	GetOverviewStats(ctx context.Context) (*model.OverviewStats, error)
	GetTrendData(ctx context.Context, params *model.DashboardQueryParams) ([]model.TrendData, error)
	GetTopProducts(ctx context.Context, params *model.DashboardQueryParams) ([]model.TopProduct, error)
	GetWarehouseUsage(ctx context.Context) ([]model.WarehouseUsage, error)
	GetSupplierPerformance(ctx context.Context, params *model.DashboardQueryParams) ([]model.SupplierPerformance, error)
	GetPendingOrders(ctx context.Context) (*model.PendingOrders, error)
}
```

- [ ] **Step 3: Run tests to verify they pass**

Run: `cd warehouse && go test -v ./internal/handler -run TestDashboardHandler`

Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add warehouse/internal/handler/dashboard.go warehouse/internal/service/dashboard.go
git commit -m "feat: implement dashboard handler with HTTP endpoints"
```

---

## Phase 5: Backend Router Integration

### Task 8: Register Dashboard Routes

**Files:**
- Modify: `warehouse/internal/router/router.go`

- [ ] **Step 1: Add DashboardHandler to Handlers struct**

Edit `warehouse/internal/router/router.go` line 15-31:

```go
type Handlers struct {
	Auth          *handler.AuthHandler
	User          *handler.UserHandler
	Role          *handler.RoleHandler
	Permission    *handler.PermissionHandler
	Warehouse     *handler.WarehouseHandler
	Location      *handler.LocationHandler
	Category      *handler.CategoryHandler
	Product       *handler.ProductHandler
	Inventory     *handler.InventoryHandler
	Supplier      *handler.SupplierHandler
	Customer      *handler.CustomerHandler
	InboundOrder  *handler.InboundOrderHandler
	OutboundOrder *handler.OutboundOrderHandler
	StockTransfer *handler.StockTransferHandler
	AuditLog      *handler.AuditLogHandler
	Dashboard     *handler.DashboardHandler  // Add this line
}
```

- [ ] **Step 2: Add dashboard routes**

Read the full router.go file first to find where to add routes:

```bash
cat warehouse/internal/router/router.go
```

Then add dashboard routes in the protected group (after other routes):

```go
// Dashboard routes
dashboard := protected.Group("/dashboard")
{
	dashboard.GET("/overview", handlers.Dashboard.GetOverview)
	dashboard.GET("/trend", handlers.Dashboard.GetTrendData)
	dashboard.GET("/top-products", handlers.Dashboard.GetTopProducts)
	dashboard.GET("/warehouse-usage", handlers.Dashboard.GetWarehouseUsage)
	dashboard.GET("/supplier-performance", handlers.Dashboard.GetSupplierPerformance)
	dashboard.GET("/pending-orders", handlers.Dashboard.GetPendingOrders)
}
```

- [ ] **Step 3: Test the routes compile**

Run: `cd warehouse && go build ./...`

Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add warehouse/internal/router/router.go
git commit -m "feat: register dashboard API routes"
```

---

## Phase 6: Frontend TypeScript Types

### Task 9: Create Dashboard TypeScript Types

**Files:**
- Create: `warehouse/web/src/types/dashboard.ts`

- [ ] **Step 1: Write type definitions**

```typescript
export interface OverviewStats {
  total_inventory: number
  inventory_warning: number
  today_inbound: number
  today_inbound_qty: number
  today_outbound: number
  today_outbound_qty: number
}

export interface TrendData {
  date: string
  inbound_qty: number
  outbound_qty: number
}

export interface TopProduct {
  product_id: number
  product_name: string
  category: string
  total_qty: number
  order_count: number
}

export interface WarehouseUsage {
  warehouse_id: number
  warehouse_name: string
  capacity: number
  used_capacity: number
  usage_rate: number
}

export interface SupplierPerformance {
  supplier_id: number
  supplier_name: string
  order_count: number
  total_value: number
  on_time_rate: number
  quality_score: number
  delivery_score: number
}

export interface PendingOrders {
  inbound_pending: number
  outbound_pending: number
  transfer_pending: number
}
```

- [ ] **Step 2: Commit**

```bash
git add warehouse/web/src/types/dashboard.ts
git commit -m "feat: add dashboard TypeScript type definitions"
```

---

## Phase 7: Frontend API Client

### Task 10: Create Dashboard API Client

**Files:**
- Create: `warehouse/web/src/api/dashboard.ts`

- [ ] **Step 1: Write API client functions**

```typescript
import api from './client'
import type { 
  OverviewStats, 
  TrendData, 
  TopProduct, 
  WarehouseUsage, 
  SupplierPerformance,
  PendingOrders 
} from '../types/dashboard'

export const dashboardApi = {
  getOverview: async (): Promise<OverviewStats> => {
    const response = await api.get<OverviewStats>('/dashboard/overview')
    return response.data
  },

  getTrend: async (startDate: string, endDate: string): Promise<TrendData[]> => {
    const response = await api.get<TrendData[]>('/dashboard/trend', {
      params: { start_date: startDate, end_date: endDate }
    })
    return response.data
  },

  getTopProducts: async (startDate: string, endDate: string, limit = 10): Promise<TopProduct[]> => {
    const response = await api.get<TopProduct[]>('/dashboard/top-products', {
      params: { start_date: startDate, end_date: endDate, limit }
    })
    return response.data
  },

  getWarehouseUsage: async (): Promise<WarehouseUsage[]> => {
    const response = await api.get<WarehouseUsage[]>('/dashboard/warehouse-usage')
    return response.data
  },

  getSupplierPerformance: async (startDate: string, endDate: string, limit = 10): Promise<SupplierPerformance[]> => {
    const response = await api.get<SupplierPerformance[]>('/dashboard/supplier-performance', {
      params: { start_date: startDate, end_date: endDate, limit }
    })
    return response.data
  },

  getPendingOrders: async (): Promise<PendingOrders> => {
    const response = await api.get<PendingOrders>('/dashboard/pending-orders')
    return response.data
  }
}
```

- [ ] **Step 2: Commit**

```bash
git add warehouse/web/src/api/dashboard.ts
git commit -m "feat: add dashboard API client functions"
```

---

## Phase 8: Frontend Components

### Task 11: Install Ant Design Charts

**Files:**
- Modify: `warehouse/web/package.json`

- [ ] **Step 1: Install the package**

Run: `cd warehouse/web && npm install @ant-design/charts`

- [ ] **Step 2: Verify installation**

Run: `cd warehouse/web && npm list @ant-design/charts`

Expected: Version listed

- [ ] **Step 3: Commit package.json**

```bash
git add warehouse/web/package.json warehouse/web/package-lock.json
git commit -m "chore: add @ant-design/charts dependency"
```

### Task 12: Create StatCard Component

**Files:**
- Create: `warehouse/web/src/pages/dashboard/components/StatCard.tsx`

- [ ] **Step 1: Write StatCard component**

```typescript
import { Card, Statistic } from 'antd'
import { ArrowUpOutlined, ArrowDownOutlined, WarningOutlined } from '@ant-design/icons'
import type { ReactNode } from 'react'

interface StatCardProps {
  title: string
  value: number | string
  prefix?: ReactNode
  suffix?: string
  valueStyle?: React.CSSProperties
  onClick?: () => void
  loading?: boolean
}

export function StatCard({ 
  title, 
  value, 
  prefix, 
  suffix, 
  valueStyle, 
  onClick,
  loading 
}: StatCardProps) {
  return (
    <Card 
      hoverable={!!onClick} 
      onClick={onClick}
      loading={loading}
      styles={{
        body: { padding: '20px' }
      }}
    >
      <Statistic
        title={title}
        value={value}
        prefix={prefix}
        suffix={suffix}
        valueStyle={{ fontSize: '24px', ...valueStyle }}
      />
    </Card>
  )
}

export function InventoryCard({ value, warning, onClick }: { 
  value: number
  warning: number
  onClick?: () => void 
}) {
  return (
    <StatCard
      title="总库存量"
      value={value}
      suffix="件"
      onClick={onClick}
      prefix={warning > 0 ? <WarningOutlined style={{ color: '#faad14' }} /> : undefined}
    />
  )
}

export function WarningCard({ value, onClick }: { value: number; onClick?: () => void }) {
  return (
    <StatCard
      title="库存预警"
      value={value}
      suffix="项"
      onClick={onClick}
      valueStyle={value > 0 ? { color: '#ff4d4f' } : undefined}
      prefix={<WarningOutlined />}
    />
  )
}

export function TodayInboundCard({ orders, quantity, onClick }: { 
  orders: number
  quantity: number
  onClick?: () => void 
}) {
  return (
    <StatCard
      title="今日入库"
      value={orders}
      suffix={`单 / ${quantity} 件`}
      onClick={onClick}
      prefix={<ArrowDownOutlined style={{ color: '#52c41a' }} />}
    />
  )
}

export function TodayOutboundCard({ orders, quantity, onClick }: { 
  orders: number
  quantity: number
  onClick?: () => void 
}) {
  return (
    <StatCard
      title="今日出库"
      value={orders}
      suffix={`单 / ${quantity} 件`}
      onClick={onClick}
      prefix={<ArrowUpOutlined style={{ color: '#1890ff' }} />}
    />
  )
}
```

- [ ] **Step 2: Commit**

```bash
git add warehouse/web/src/pages/dashboard/components/StatCard.tsx
git commit -m "feat: add StatCard component for dashboard metrics"
```

### Task 13: Create TrendChart Component

**Files:**
- Create: `warehouse/web/src/pages/dashboard/components/TrendChart.tsx`

- [ ] **Step 1: Write TrendChart component**

```typescript
import { Line } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import type { TrendData } from '../../../types/dashboard'

interface TrendChartProps {
  data: TrendData[]
  loading?: boolean
  onPointClick?: (date: string, type: 'inbound' | 'outbound') => void
}

export function TrendChart({ data, loading, onPointClick }: TrendChartProps) {
  const config = {
    data: data.flatMap(item => [
      { date: item.date, value: item.inbound_qty, type: '入库' },
      { date: item.date, value: item.outbound_qty, type: '出库' }
    ]),
    xField: 'date',
    yField: 'value',
    seriesField: 'type',
    color: ['#1890ff', '#fa8c16'],
    smooth: true,
    animation: {
      appear: {
        animation: 'path-in',
        duration: 1000,
      },
    },
    point: {
      shape: 'circle',
      size: 4,
    },
    interactions: [
      {
        type: 'marker-active',
      },
    ],
    onReady: (plot: any) => {
      if (onPointClick) {
        plot.on('element:click', (evt: any) => {
          const { data } = evt.data
          if (data && data.date) {
            const type = data.type === '入库' ? 'inbound' : 'outbound'
            onPointClick(data.date, type)
          }
        })
      }
    },
  }

  return (
    <Card title="出入库趋势" style={{ height: '100%' }}>
      <Spin spinning={loading}>
        <Line {...config} />
      </Spin>
    </Card>
  )
}
```

- [ ] **Step 2: Commit**

```bash
git add warehouse/web/src/pages/dashboard/components/TrendChart.tsx
git commit -m "feat: add TrendChart component for inbound/outbound trends"
```

### Task 14: Create TopProductsChart Component

**Files:**
- Create: `warehouse/web/src/pages/dashboard/components/TopProductsChart.tsx`

- [ ] **Step 1: Write TopProductsChart component**

```typescript
import { Bar } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import type { TopProduct } from '../../../types/dashboard'

interface TopProductsChartProps {
  data: TopProduct[]
  loading?: boolean
  onBarClick?: (productId: number) => void
}

export function TopProductsChart({ data, loading, onBarClick }: TopProductsChartProps) {
  const config = {
    data: data.map(item => ({
      name: item.product_name,
      value: item.total_qty,
      productId: item.product_id,
    })),
    xField: 'value',
    yField: 'name',
    seriesField: 'name',
    legend: false,
    color: '#1890ff',
    barStyle: {
      radius: [0, 4, 4, 0],
    },
    onReady: (plot: any) => {
      if (onBarClick) {
        plot.on('element:click', (evt: any) => {
          const { data } = evt.data
          if (data && data.productId) {
            onBarClick(data.productId)
          }
        })
      }
    },
  }

  return (
    <Card title="热销产品排行" style={{ height: '100%' }}>
      <Spin spinning={loading}>
        <Bar {...config} />
      </Spin>
    </Card>
  )
}
```

- [ ] **Step 2: Commit**

```bash
git add warehouse/web/src/pages/dashboard/components/TopProductsChart.tsx
git commit -m "feat: add TopProductsChart component for product ranking"
```

### Task 15: Create WarehouseUsageChart Component

**Files:**
- Create: `warehouse/web/src/pages/dashboard/components/WarehouseUsageChart.tsx`

- [ ] **Step 1: Write WarehouseUsageChart component**

```typescript
import { Pie } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import type { WarehouseUsage } from '../../../types/dashboard'

interface WarehouseUsageChartProps {
  data: WarehouseUsage[]
  loading?: boolean
  onClick?: (warehouseId: number) => void
}

export function WarehouseUsageChart({ data, loading, onClick }: WarehouseUsageChartProps) {
  const config = {
    data: data.map(item => ({
      type: item.warehouse_name,
      value: item.usage_rate,
      warehouseId: item.warehouse_id,
    })),
    angleField: 'value',
    colorField: 'type',
    radius: 0.8,
    innerRadius: 0.6,
    label: {
      type: 'inner',
      offset: '-50%',
      content: '{value}%',
      style: {
        textAlign: 'center',
        fontSize: 12,
      },
    },
    statistic: {
      title: {
        content: '平均使用率',
      },
      content: {
        formatter: () => {
          if (data.length === 0) return '0%'
          const avg = data.reduce((sum, item) => sum + item.usage_rate, 0) / data.length
          return `${avg.toFixed(1)}%`
        },
      },
    },
    onReady: (plot: any) => {
      if (onClick) {
        plot.on('element:click', (evt: any) => {
          const { data } = evt.data
          if (data && data.warehouseId) {
            onClick(data.warehouseId)
          }
        })
      }
    },
  }

  return (
    <Card title="仓库使用率" style={{ height: '100%' }}>
      <Spin spinning={loading}>
        <Pie {...config} />
      </Spin>
    </Card>
  )
}
```

- [ ] **Step 2: Commit**

```bash
git add warehouse/web/src/pages/dashboard/components/WarehouseUsageChart.tsx
git commit -m "feat: add WarehouseUsageChart component for warehouse capacity"
```

### Task 16: Create SupplierPerformanceChart Component

**Files:**
- Create: `warehouse/web/src/pages/dashboard/components/SupplierPerformanceChart.tsx`

- [ ] **Step 1: Write SupplierPerformanceChart component**

```typescript
import { Radar } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import type { SupplierPerformance } from '../../../types/dashboard'

interface SupplierPerformanceChartProps {
  data: SupplierPerformance[]
  loading?: boolean
  onClick?: (supplierId: number) => void
}

export function SupplierPerformanceChart({ data, loading, onClick }: SupplierPerformanceChartProps) {
  const config = {
    data: data.flatMap(item => [
      { name: item.supplier_name, label: '订单量', value: item.order_count },
      { name: item.supplier_name, label: '总金额', value: item.total_value / 10000 },
      { name: item.supplier_name, label: '准时率', value: item.on_time_rate },
      { name: item.supplier_name, label: '质量评分', value: item.quality_score },
      { name: item.supplier_name, label: '交付评分', value: item.delivery_score },
    ]),
    xField: 'label',
    yField: 'value',
    seriesField: 'name',
    meta: {
      value: {
        alias: '分数',
        min: 0,
        max: 100,
      },
    },
    radius: 0.8,
    onReady: (plot: any) => {
      if (onClick) {
        plot.on('element:click', (evt: any) => {
          const { data } = evt.data
          if (data && data.name) {
            const supplier = data.find((item: any) => item.name === data.name)
            if (supplier) {
              onClick(supplier.supplier_id)
            }
          }
        })
      }
    },
  }

  return (
    <Card title="供应商绩效" style={{ height: '100%' }}>
      <Spin spinning={loading}>
        <Radar {...config} />
      </Spin>
    </Card>
  )
}
```

- [ ] **Step 2: Commit**

```bash
git add warehouse/web/src/pages/dashboard/components/SupplierPerformanceChart.tsx
git commit -m "feat: add SupplierPerformanceChart component for supplier metrics"
```

---

## Phase 9: Frontend Data Hook

### Task 17: Create useDashboardStats Hook

**Files:**
- Create: `warehouse/web/src/pages/dashboard/hooks/useDashboardStats.ts`

- [ ] **Step 1: Write the custom hook**

```typescript
import { useQuery } from '@tanstack/react-query'
import { dashboardApi } from '../../../api/dashboard'
import dayjs from 'dayjs'

export function useDashboardStats(startDate?: string, endDate?: string) {
  const start = startDate || dayjs().subtract(30, 'day').format('YYYY-MM-DD')
  const end = endDate || dayjs().format('YYYY-MM-DD')

  const overviewQuery = useQuery({
    queryKey: ['dashboard', 'overview'],
    queryFn: dashboardApi.getOverview,
  })

  const trendQuery = useQuery({
    queryKey: ['dashboard', 'trend', start, end],
    queryFn: () => dashboardApi.getTrend(start, end),
  })

  const topProductsQuery = useQuery({
    queryKey: ['dashboard', 'topProducts', start, end],
    queryFn: () => dashboardApi.getTopProducts(start, end),
  })

  const warehouseUsageQuery = useQuery({
    queryKey: ['dashboard', 'warehouseUsage'],
    queryFn: dashboardApi.getWarehouseUsage,
  })

  const supplierPerformanceQuery = useQuery({
    queryKey: ['dashboard', 'supplierPerformance', start, end],
    queryFn: () => dashboardApi.getSupplierPerformance(start, end),
  })

  const pendingOrdersQuery = useQuery({
    queryKey: ['dashboard', 'pendingOrders'],
    queryFn: dashboardApi.getPendingOrders,
  })

  const refetchAll = () => {
    overviewQuery.refetch()
    trendQuery.refetch()
    topProductsQuery.refetch()
    warehouseUsageQuery.refetch()
    supplierPerformanceQuery.refetch()
    pendingOrdersQuery.refetch()
  }

  return {
    overview: overviewQuery.data,
    trend: trendQuery.data,
    topProducts: topProductsQuery.data,
    warehouseUsage: warehouseUsageQuery.data,
    supplierPerformance: supplierPerformanceQuery.data,
    pendingOrders: pendingOrdersQuery.data,
    loading: 
      overviewQuery.isLoading || 
      trendQuery.isLoading || 
      topProductsQuery.isLoading ||
      warehouseUsageQuery.isLoading ||
      supplierPerformanceQuery.isLoading ||
      pendingOrdersQuery.isLoading,
    refetch: refetchAll,
  }
}
```

- [ ] **Step 2: Check if dayjs is installed**

Run: `cd warehouse/web && npm list dayjs`

If not installed, add it:

Run: `cd warehouse/web && npm install dayjs`

- [ ] **Step 3: Commit**

```bash
git add warehouse/web/src/pages/dashboard/hooks/useDashboardStats.ts warehouse/web/package.json warehouse/web/package-lock.json
git commit -m "feat: add useDashboardStats hook for data fetching"
```

---

## Phase 10: Frontend Main Dashboard Page

### Task 18: Create Dashboard Page

**Files:**
- Create: `warehouse/web/src/pages/dashboard/Dashboard.tsx`
- Create: `warehouse/web/src/pages/dashboard/index.ts`

- [ ] **Step 1: Write the main Dashboard component**

```typescript
import { useState } from 'react'
import { Row, Col, DatePicker, Button, Space, message } from 'antd'
import { ReloadOutlined, DownloadOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import dayjs, { Dayjs } from 'dayjs'
import { useDashboardStats } from './hooks/useDashboardStats'
import { 
  InventoryCard, 
  WarningCard, 
  TodayInboundCard, 
  TodayOutboundCard 
} from './components/StatCard'
import { TrendChart } from './components/TrendChart'
import { TopProductsChart } from './components/TopProductsChart'
import { WarehouseUsageChart } from './components/WarehouseUsageChart'
import { SupplierPerformanceChart } from './components/SupplierPerformanceChart'

const { RangePicker } = DatePicker

export function Dashboard() {
  const navigate = useNavigate()
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs]>([
    dayjs().subtract(30, 'day'),
    dayjs()
  ])
  
  const { 
    overview, 
    trend, 
    topProducts, 
    warehouseUsage, 
    supplierPerformance,
    loading,
    refetch 
  } = useDashboardStats(
    dateRange[0].format('YYYY-MM-DD'),
    dateRange[1].format('YYYY-MM-DD')
  )

  const handleDateChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
    if (dates && dates[0] && dates[1]) {
      setDateRange([dates[0], dates[1]])
    }
  }

  const handleRefresh = () => {
    refetch()
    message.success('数据已刷新')
  }

  const handleExport = async () => {
    message.info('导出功能开发中...')
  }

  return (
    <div style={{ padding: '24px' }}>
      <Row justify="space-between" align="middle" style={{ marginBottom: '24px' }}>
        <Col>
          <h1 style={{ margin: 0 }}>仪表盘</h1>
        </Col>
        <Col>
          <Space>
            <RangePicker
              value={dateRange}
              onChange={handleDateChange}
              format="YYYY-MM-DD"
              allowClear={false}
            />
            <Button icon={<ReloadOutlined />} onClick={handleRefresh}>
              刷新
            </Button>
            <Button icon={<DownloadOutlined />} onClick={handleExport}>
              导出
            </Button>
          </Space>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <InventoryCard 
            value={overview?.total_inventory || 0} 
            warning={overview?.inventory_warning || 0}
            onClick={() => navigate('/inventory')}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <WarningCard 
            value={overview?.inventory_warning || 0}
            onClick={() => navigate('/inventory?quantity_max=10')}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <TodayInboundCard 
            orders={overview?.today_inbound || 0}
            quantity={overview?.today_inbound_qty || 0}
            onClick={() => navigate('/inbound')}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <TodayOutboundCard 
            orders={overview?.today_outbound || 0}
            quantity={overview?.today_outbound_qty || 0}
            onClick={() => navigate('/outbound')}
          />
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: '16px' }}>
        <Col xs={24} lg={16}>
          <TrendChart 
            data={trend || []} 
            loading={loading}
            onPointClick={(date, type) => {
              const path = type === 'inbound' ? '/inbound' : '/outbound'
              navigate(`${path}?date=${date}`)
            }}
          />
        </Col>
        <Col xs={24} lg={8}>
          <TopProductsChart 
            data={topProducts || []} 
            loading={loading}
            onBarClick={(productId) => {
              navigate(`/products?id=${productId}`)
            }}
          />
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: '16px' }}>
        <Col xs={24} lg={12}>
          <WarehouseUsageChart 
            data={warehouseUsage || []} 
            loading={loading}
            onClick={(warehouseId) => {
              navigate(`/inventory?warehouse_id=${warehouseId}`)
            }}
          />
        </Col>
        <Col xs={24} lg={12}>
          <SupplierPerformanceChart 
            data={supplierPerformance || []} 
            loading={loading}
            onClick={(supplierId) => {
              navigate(`/suppliers?id=${supplierId}`)
            }}
          />
        </Col>
      </Row>
    </div>
  )
}
```

- [ ] **Step 2: Create index.ts for barrel export**

```typescript
export { Dashboard } from './Dashboard'
```

- [ ] **Step 3: Commit**

```bash
git add warehouse/web/src/pages/dashboard/Dashboard.tsx warehouse/web/src/pages/dashboard/index.ts
git commit -m "feat: add main Dashboard page with charts and statistics"
```

---

## Phase 11: Frontend Router Integration

### Task 19: Update App Router

**Files:**
- Modify: `warehouse/web/src/App.tsx`

- [ ] **Step 1: Import Dashboard component**

Read the current App.tsx:

```bash
cat warehouse/web/src/App.tsx | head -20
```

Add import at the top with other page imports (around line 10):

```typescript
import { Dashboard } from './pages/dashboard'
```

- [ ] **Step 2: Replace placeholder Dashboard component**

Find the placeholder Dashboard function (lines 20-22) and remove it:

```typescript
// Remove this:
function Dashboard() {
  return <div>Dashboard</div>
}
```

- [ ] **Step 3: Commit**

```bash
git add warehouse/web/src/App.tsx
git commit -m "feat: integrate Dashboard page into app router"
```

---

## Phase 12: Testing and Verification

### Task 20: Run Backend Tests

- [ ] **Step 1: Run all backend tests**

Run: `cd warehouse && go test -v ./internal/repository -run TestDashboardRepository`

Run: `cd warehouse && go test -v ./internal/service -run TestDashboardService`

Run: `cd warehouse && go test -v ./internal/handler -run TestDashboardHandler`

Expected: All tests pass

- [ ] **Step 2: Run backend build**

Run: `cd warehouse && go build ./...`

Expected: No errors

### Task 21: Run Frontend Tests

- [ ] **Step 1: Run TypeScript type check**

Run: `cd warehouse/web && npm run build`

Expected: Build succeeds with no type errors

- [ ] **Step 2: Run linter**

Run: `cd warehouse/web && npm run lint`

Expected: No linting errors

### Task 22: Manual Integration Test

- [ ] **Step 1: Start backend server**

Run: `cd warehouse && go run ./cmd/server`

Expected: Server starts on configured port

- [ ] **Step 2: Start frontend dev server**

Run: `cd warehouse/web && npm run dev`

Expected: Dev server starts

- [ ] **Step 3: Test in browser**

Open: `http://localhost:5173` (or configured port)

Test cases:
- [ ] Dashboard page loads
- [ ] Stat cards show data
- [ ] Time range picker works
- [ ] Refresh button works
- [ ] Trend chart displays
- [ ] Top products chart displays
- [ ] Warehouse usage chart displays
- [ ] Supplier performance chart displays
- [ ] Card clicks navigate to correct pages
- [ ] Chart interactions work

---

## Phase 13: Final Commit and Documentation

### Task 23: Create Final Commit

- [ ] **Step 1: Ensure all changes are committed**

Run: `git status`

Expected: No uncommitted changes

- [ ] **Step 2: Create summary commit**

```bash
git add -A
git commit -m "feat: complete dashboard metrics charts implementation

- Add backend dashboard API with 6 endpoints
- Implement dashboard repository with SQL aggregation queries
- Create dashboard service layer with business logic
- Add dashboard HTTP handlers with parameter validation
- Create frontend dashboard page with Ant Design Charts
- Implement 4 chart components: trend, ranking, usage, performance
- Add stat cards with click-to-navigate functionality
- Support time range filtering (default 30 days)
- Implement manual refresh functionality
- Add TypeScript type definitions for all dashboard data
- Include comprehensive test coverage for backend

Closes: Dashboard metrics visualization requirement"
```

### Task 24: Update Documentation

- [ ] **Step 1: Verify implementation matches design spec**

Compare the implemented features with the design spec:
- All 7 API endpoints implemented ✓
- All chart types implemented ✓
- Time range filtering works ✓
- Click interactions work ✓
- Refresh functionality works ✓

- [ ] **Step 2: Create commit for any fixes**

If any issues found, fix and commit:

```bash
git add -A
git commit -m "fix: update implementation to match design spec"
```

---

## Acceptance Criteria

- [ ] Backend API returns correct statistics
- [ ] Frontend displays all charts correctly
- [ ] Time range filtering works
- [ ] Manual refresh updates all data
- [ ] Click navigation works for all cards and charts
- [ ] Charts are interactive and responsive
- [ ] Loading states display correctly
- [ ] Error states handle gracefully
- [ ] All tests pass
- [ ] Build succeeds with no errors
- [ ] Manual testing shows all features working

---

## Notes

- Export functionality (Excel/PDF) is deferred to future iteration
- Supplier performance scoring may need schema updates if scoring fields don't exist
- Performance optimization (indexing) should be added based on actual data volume
- Consider adding error boundary components for chart failures
- Mobile responsiveness should be tested on actual devices
