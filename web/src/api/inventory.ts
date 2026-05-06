import api from './client'
import type { Inventory, CreateInventoryRequest, UpdateInventoryRequest, AdjustQuantityRequest } from '../types/inventory'
import type { PaginatedResponse } from '../types/warehouse'

export interface InventoryFilter {
  product_name?: string
  product_id?: number
  warehouse_id?: number
  quantity_min?: number
  quantity_max?: number
  batch_no?: string
}

export const inventoryApi = {
  list: async (page = 1, size = 10, filter?: InventoryFilter): Promise<PaginatedResponse<Inventory>> => {
    const params = new URLSearchParams()
    params.append('page', String(page))
    params.append('size', String(size))
    
    if (filter) {
      if (filter.product_name) params.append('product_name', filter.product_name)
      if (filter.product_id !== undefined) params.append('product_id', String(filter.product_id))
      if (filter.warehouse_id !== undefined) params.append('warehouse_id', String(filter.warehouse_id))
      if (filter.batch_no) params.append('batch_no', filter.batch_no)
      if (filter.quantity_min !== undefined) params.append('quantity_min', String(filter.quantity_min))
      if (filter.quantity_max !== undefined) params.append('quantity_max', String(filter.quantity_max))
    }
    
    const response = await api.get<PaginatedResponse<Inventory>>('/inventory', { params })
    return response.data
  },

  get: async (id: number): Promise<Inventory> => {
    const response = await api.get<Inventory>(`/inventory/${id}`)
    return response.data
  },

  create: async (data: CreateInventoryRequest): Promise<Inventory> => {
    const response = await api.post<Inventory>('/inventory', data)
    return response.data
  },

  update: async (id: number, data: UpdateInventoryRequest): Promise<Inventory> => {
    const response = await api.put<Inventory>(`/inventory/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/inventory/${id}`)
  },

  adjust: async (data: AdjustQuantityRequest): Promise<Inventory> => {
    const response = await api.post<Inventory>('/inventory/adjust', data)
    return response.data
  }
}
