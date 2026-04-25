import api from './client'
import type { Location, CreateLocationRequest, UpdateLocationRequest, PaginatedResponse } from '../types/warehouse'

export const locationApi = {
  list: async (page = 1, size = 10, warehouseId?: number): Promise<PaginatedResponse<Location>> => {
    const response = await api.get<PaginatedResponse<Location>>('/locations', {
      params: { page, size, warehouse_id: warehouseId }
    })
    return response.data
  },

  get: async (id: number): Promise<Location> => {
    const response = await api.get<Location>(`/locations/${id}`)
    return response.data
  },

  create: async (data: CreateLocationRequest): Promise<Location> => {
    const response = await api.post<Location>('/locations', data)
    return response.data
  },

  update: async (id: number, data: UpdateLocationRequest): Promise<Location> => {
    const response = await api.put<Location>(`/locations/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/locations/${id}`)
  }
}
