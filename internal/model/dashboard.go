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
