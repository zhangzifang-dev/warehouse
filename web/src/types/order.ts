export interface InboundOrder {
  id: number
  order_no: string
  supplier_id: number | null
  warehouse_id: number
  total_quantity: number
  status: number
  remark: string
  created_at: string
  updated_at: string
  items?: InboundItem[]
}

export interface InboundItem {
  id: number
  order_id: number
  product_id: number
  location_id: number | null
  quantity: number
  batch_no: string
  created_at: string
  updated_at: string
}

export interface CreateInboundOrderRequest {
  supplier_id?: number
  warehouse_id: number
  remark?: string
  items: CreateInboundItemRequest[]
}

export interface CreateInboundItemRequest {
  product_id: number
  location_id?: number
  quantity: number
  batch_no?: string
}

export interface OutboundOrder {
  id: number
  order_no: string
  customer_id: number | null
  warehouse_id: number
  total_quantity: number
  status: number
  remark: string
  created_at: string
  updated_at: string
  items?: OutboundItem[]
}

export interface OutboundItem {
  id: number
  order_id: number
  product_id: number
  location_id: number | null
  quantity: number
  batch_no: string
  created_at: string
  updated_at: string
}

export interface CreateOutboundOrderRequest {
  customer_id?: number
  warehouse_id: number
  remark?: string
  items: CreateOutboundItemRequest[]
}

export interface CreateOutboundItemRequest {
  product_id: number
  location_id?: number
  quantity: number
  batch_no?: string
}

export interface StockTransfer {
  id: number
  order_no: string
  source_warehouse_id: number
  target_warehouse_id: number
  status: number
  created_at: string
  updated_at: string
  items?: StockTransferItem[]
}

export interface StockTransferItem {
  id: number
  transfer_id: number
  product_id: number
  location_id: number | null
  quantity: number
  batch_no: string
  created_at: string
  updated_at: string
}

export interface CreateStockTransferRequest {
  source_warehouse_id: number
  target_warehouse_id: number
  items: CreateStockTransferItemRequest[]
}

export interface CreateStockTransferItemRequest {
  product_id: number
  location_id?: number
  quantity: number
  batch_no?: string
}
