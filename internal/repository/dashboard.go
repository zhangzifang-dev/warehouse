package repository

import (
	"context"
	"log"
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

	var totalInv float64
	err := r.db.NewSelect().
		Table("inventories").
		Where("deleted_at IS NULL").
		ColumnExpr("COALESCE(SUM(quantity), 0)").
		Scan(ctx, &totalInv)
	if err != nil {
		log.Printf("[ERROR] GetOverviewStats - total_inventory query failed: %v", err)
		return nil, err
	}
	stats.TotalInventory = totalInv

	var warning int
	err = r.db.NewSelect().
		Table("inventories").
		Where("deleted_at IS NULL").
		Where("quantity < 10").
		ColumnExpr("COUNT(*)").
		Scan(ctx, &warning)
	if err != nil {
		log.Printf("[ERROR] GetOverviewStats - inventory_warning query failed: %v", err)
		return nil, err
	}
	stats.InventoryWarning = warning

	today := time.Now().Format("2006-01-02")
	var todayIn int
	var todayInQty float64
	err = r.db.NewSelect().
		Table("inbound_orders").
		Where("deleted_at IS NULL").
		Where("DATE(created_at) = ?", today).
		ColumnExpr("COUNT(*)").
		ColumnExpr("COALESCE(SUM(total_quantity), 0)").
		Scan(ctx, &todayIn, &todayInQty)
	if err != nil {
		log.Printf("[ERROR] GetOverviewStats - today_inbound query failed: %v", err)
		return nil, err
	}
	stats.TodayInbound = todayIn
	stats.TodayInboundQty = todayInQty

	var todayOut int
	var todayOutQty float64
	err = r.db.NewSelect().
		Table("outbound_orders").
		Where("deleted_at IS NULL").
		Where("DATE(created_at) = ?", today).
		ColumnExpr("COUNT(*)").
		ColumnExpr("COALESCE(SUM(total_quantity), 0)").
		Scan(ctx, &todayOut, &todayOutQty)
	if err != nil {
		log.Printf("[ERROR] GetOverviewStats - today_outbound query failed: %v", err)
		return nil, err
	}
	stats.TodayOutbound = todayOut
	stats.TodayOutboundQty = todayOutQty

	return stats, nil
}

func (r *DashboardRepository) GetTrendData(ctx context.Context, params *model.DashboardQueryParams) ([]model.TrendData, error) {
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
		ColumnExpr("COALESCE(SUM(total_quantity), 0) as qty").
		GroupExpr("DATE(created_at)").
		OrderExpr("date").
		Scan(ctx, &inboundTrend)
	if err != nil {
		return nil, err
	}

	var outboundTrend []DailyData
	err = r.db.NewSelect().
		Table("outbound_orders").
		Where("deleted_at IS NULL").
		Where("created_at >= ?", params.StartDate).
		Where("created_at <= ?", params.EndDate).
		ColumnExpr("DATE(created_at) as date").
		ColumnExpr("COALESCE(SUM(total_quantity), 0) as qty").
		GroupExpr("DATE(created_at)").
		OrderExpr("date").
		Scan(ctx, &outboundTrend)
	if err != nil {
		return nil, err
	}

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

	result := make([]model.TrendData, 0, len(dateMap))
	for _, v := range dateMap {
		result = append(result, *v)
	}

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
		Join("JOIN outbound_orders o ON oi.outbound_order_id = o.id").
		Where("o.deleted_at IS NULL").
		Where("o.created_at >= ?", params.StartDate).
		Where("o.created_at <= ?", params.EndDate).
		Where("o.status >= 1").
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
	log.Printf("[ERROR] GetWarehouseUsage query failed: %v", err)
	if err != nil {
		return nil, err
	}

	for i := range usage {
		if usage[i].Capacity > 0 {
			usage[i].UsageRate = float64(usage[i].UsedCapacity) / float64(usage[i].Capacity) * 100
		}
	}

	return usage, nil
}

func (r *DashboardRepository) GetSupplierPerformance(ctx context.Context, params *model.DashboardQueryParams) ([]model.SupplierPerformance, error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var performance []model.SupplierPerformance

	query := `
		SELECT 
			s.id as supplier_id,
			s.name as supplier_name,
			COUNT(DISTINCT io.id) as order_count,
			COALESCE(SUM(io.total_quantity), 0) as total_value,
			0.0 as on_time_rate,
			0.0 as quality_score,
			0.0 as delivery_score
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
		log.Printf("[ERROR] GetSupplierPerformance query failed: %v", err)
		return nil, err
	}

	return performance, nil
}

func (r *DashboardRepository) GetPendingOrders(ctx context.Context) (*model.PendingOrders, error) {
	pending := &model.PendingOrders{}

	err := r.db.NewSelect().
		Table("inbound_orders").
		Where("deleted_at IS NULL").
		Where("status IN ('pending', 'approved')").
		ColumnExpr("COUNT(*) as count").
		Scan(ctx, &pending.InboundPending)
	if err != nil {
		return nil, err
	}

	err = r.db.NewSelect().
		Table("outbound_orders").
		Where("deleted_at IS NULL").
		Where("status IN ('pending', 'approved')").
		ColumnExpr("COUNT(*) as count").
		Scan(ctx, &pending.OutboundPending)
	if err != nil {
		return nil, err
	}

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
