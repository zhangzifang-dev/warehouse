export interface Inventory {
  id: number
  warehouse_id: number
  product_id: number
  location_id: number | null
  quantity: number
  batch_no: string
  created_at: string
  updated_at: string
}

export interface CreateInventoryRequest {
  warehouse_id: number
  product_id: number
  location_id?: number
  quantity?: number
  batch_no?: string
}

export interface UpdateInventoryRequest {
  warehouse_id?: number
  product_id?: number
  location_id?: number
  quantity?: number
  batch_no?: string
}

export interface AdjustQuantityRequest {
  inventory_id: number
  quantity: number
}
