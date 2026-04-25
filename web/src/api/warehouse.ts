import api from './client'
import type { Warehouse, CreateWarehouseRequest, UpdateWarehouseRequest, PaginatedResponse } from '../types/warehouse'

export const warehouseApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<Warehouse>> => {
    const response = await api.get<PaginatedResponse<Warehouse>>('/warehouses', {
      params: { page, size }
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
