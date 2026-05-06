import api from './client'
import type { Warehouse, CreateWarehouseRequest, UpdateWarehouseRequest, PaginatedResponse } from '../types/warehouse'

export interface WarehouseFilter {
  name?: string
}

export const warehouseApi = {
  list: async (page = 1, size = 10, filter?: WarehouseFilter): Promise<PaginatedResponse<Warehouse>> => {
    const params = new URLSearchParams()
    params.append('page', String(page))
    params.append('size', String(size))
    
    if (filter?.name) {
      params.append('name', filter.name)
    }
    
    const response = await api.get<PaginatedResponse<Warehouse>>('/warehouses', {
      params: params
    })
    return response.data
  },

  get: async (id: number): Promise<Warehouse> => {
    const response = await api.get<Warehouse>(`/warehouses/${id}`)
    return response.data
  },

  create: async (data: CreateWarehouseRequest): Promise<Warehouse> => {
    const response = await api.post<Warehouse>('/warehouses', data)
    return response.data
  },

  update: async (id: number, data: UpdateWarehouseRequest): Promise<Warehouse> => {
    const response = await api.put<Warehouse>(`/warehouses/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/warehouses/${id}`)
  }
}
