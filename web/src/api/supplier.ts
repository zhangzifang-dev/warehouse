import api from './client'
import type { Supplier, CreateSupplierRequest, UpdateSupplierRequest } from '../types/partner'
import type { PaginatedResponse } from '../types/warehouse'

export interface SupplierFilter {
  code?: string
  name?: string
  contact?: string
  phone?: string
  status?: number
}

export const supplierApi = {
  list: async (page = 1, size = 10, filter?: SupplierFilter): Promise<PaginatedResponse<Supplier>> => {
    const params = new URLSearchParams()
    params.append('page', String(page))
    params.append('size', String(size))
    
    if (filter?.code) {
      params.append('code', filter.code)
    }
    if (filter?.name) {
      params.append('name', filter.name)
    }
    if (filter?.contact) {
      params.append('contact', filter.contact)
    }
    if (filter?.phone) {
      params.append('phone', filter.phone)
    }
    if (filter?.status !== undefined) {
      params.append('status', String(filter.status))
    }
    
    const response = await api.get<PaginatedResponse<Supplier>>('/suppliers', {
      params: params
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
