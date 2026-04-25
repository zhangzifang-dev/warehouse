import api from './client'
import type { Supplier, CreateSupplierRequest, UpdateSupplierRequest } from '../types/partner'
import type { PaginatedResponse } from '../types/warehouse'

export const supplierApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<Supplier>> => {
    const response = await api.get<PaginatedResponse<Supplier>>('/suppliers', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<Supplier> => {
    const response = await api.get<Supplier>(`/suppliers/${id}`)
    return response.data
  },

  create: async (data: CreateSupplierRequest): Promise<Supplier> => {
    const response = await api.post<Supplier>('/suppliers', data)
    return response.data
  },

  update: async (id: number, data: UpdateSupplierRequest): Promise<Supplier> => {
    const response = await api.put<Supplier>(`/suppliers/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/suppliers/${id}`)
  }
}
