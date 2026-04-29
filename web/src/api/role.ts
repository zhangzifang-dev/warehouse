import api from './client'
import type { Role, CreateRoleRequest, UpdateRoleRequest, Permission } from '../types/system'
import type { PaginatedResponse } from '../types/warehouse'

export const roleApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<Role>> => {
    const response = await api.get<PaginatedResponse<Role>>('/roles', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<Role> => {
    const response = await api.get<Role>(`/roles/${id}`)
    return response.data
  },

  create: async (data: CreateRoleRequest): Promise<Role> => {
    const response = await api.post<Role>('/roles', data)
    return response.data
  },

  update: async (id: number, data: UpdateRoleRequest): Promise<Role> => {
    const response = await api.put<Role>(`/roles/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/roles/${id}`)
  },

  getPermissions: async (id: number): Promise<Permission[]> => {
    const response = await api.get<Permission[]>(`/roles/${id}/permissions`)
    return response.data
  },

  assignPermissions: async (id: number, permissionIds: number[]): Promise<void> => {
    await api.post(`/roles/${id}/permissions`, { permission_ids: permissionIds })
  }
}
