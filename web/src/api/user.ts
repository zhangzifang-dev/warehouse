import api from './client'
import type { User, CreateUserRequest, UpdateUserRequest, Role } from '../types/system'
import type { PaginatedResponse } from '../types/warehouse'

export const userApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<User>> => {
    const response = await api.get<{ users: User[]; total: number; page: number; size: number }>('/users', {
      params: { page, size }
    })
    return { items: response.data.users, total: response.data.total, page: response.data.page, size: response.data.size }
  },

  get: async (id: number): Promise<User> => {
    const response = await api.get<User>(`/users/${id}`)
    return response.data
  },

  create: async (data: CreateUserRequest): Promise<User> => {
    const response = await api.post<User>('/users', data)
    return response.data
  },

  update: async (id: number, data: UpdateUserRequest): Promise<User> => {
    const response = await api.put<User>(`/users/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/users/${id}`)
  },

  getRoles: async (id: number): Promise<Role[]> => {
    const response = await api.get<Role[]>(`/users/${id}/roles`)
    return response.data
  },

  assignRoles: async (id: number, roleIds: number[]): Promise<void> => {
    await api.post(`/users/${id}/roles`, { role_ids: roleIds })
  }
}
