import api from './client'
import type { Inventory, CreateInventoryRequest, UpdateInventoryRequest, AdjustQuantityRequest } from '../types/inventory'
import type { PaginatedResponse } from '../types/warehouse'

export const inventoryApi = {
  list: async (page = 1, size = 10, warehouseId?: number, productId?: number): Promise<PaginatedResponse<Inventory>> => {
    const response = await api.get<PaginatedResponse<Inventory>>('/inventory', {
      params: { page, size, warehouse_id: warehouseId, product_id: productId }
    })
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
