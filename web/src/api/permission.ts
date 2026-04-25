import api from './client'
import type { Permission } from '../types/system'
import type { PaginatedResponse } from '../types/warehouse'

export const permissionApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<Permission>> => {
    const response = await api.get<PaginatedResponse<Permission>>('/permissions', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<Permission> => {
    const response = await api.get<Permission>(`/permissions/${id}`)
    return response.data
  }
}
