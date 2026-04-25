import api from './client'
import type { Product, CreateProductRequest, UpdateProductRequest } from '../types/product'
import type { PaginatedResponse } from '../types/warehouse'

export const productApi = {
  list: async (page = 1, size = 10, categoryId?: number, keyword?: string): Promise<PaginatedResponse<Product>> => {
    const response = await api.get<PaginatedResponse<Product>>('/products', {
      params: { page, size, category_id: categoryId, keyword }
    })
    return response.data
  },

  get: async (id: number): Promise<Product> => {
    const response = await api.get<Product>(`/products/${id}`)
    return response.data
  },

  create: async (data: CreateProductRequest): Promise<Product> => {
    const response = await api.post<Product>('/products', data)
    return response.data
  },

  update: async (id: number, data: UpdateProductRequest): Promise<Product> => {
    const response = await api.put<Product>(`/products/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/products/${id}`)
  }
}
