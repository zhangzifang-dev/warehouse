import api from './client'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '../types/product'
import type { PaginatedResponse } from '../types/warehouse'

export const categoryApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<Category>> => {
    const response = await api.get<PaginatedResponse<Category>>('/categories', {
      params: { page, size }
    })
    return response.data
  },

  tree: async (): Promise<Category[]> => {
    const response = await api.get<Category[]>('/categories/tree')
    return response.data
  },

  get: async (id: number): Promise<Category> => {
    const response = await api.get<Category>(`/categories/${id}`)
    return response.data
  },

  create: async (data: CreateCategoryRequest): Promise<Category> => {
    const response = await api.post<Category>('/categories', data)
    return response.data
  },

  update: async (id: number, data: UpdateCategoryRequest): Promise<Category> => {
    const response = await api.put<Category>(`/categories/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/categories/${id}`)
  }
}
