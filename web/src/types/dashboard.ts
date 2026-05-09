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
