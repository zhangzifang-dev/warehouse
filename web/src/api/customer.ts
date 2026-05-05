import api from './client'
import type { Customer, CreateCustomerRequest, UpdateCustomerRequest } from '../types/partner'
import type { PaginatedResponse } from '../types/warehouse'

export interface CustomerFilter {
  code?: string
  name?: string
  phone?: string
  status?: number
}

export const customerApi = {
  list: async (page = 1, size = 10, filter?: CustomerFilter): Promise<PaginatedResponse<Customer>> => {
    const response = await api.get<PaginatedResponse<Customer>>('/customers', {
      params: { 
        page, 
        size,
        ...filter
      }
    })
    return response.data
  },

  get: async (id: number): Promise<Customer> => {
    const response = await api.get<Customer>(`/customers/${id}`)
    return response.data
  },

  create: async (data: CreateCustomerRequest): Promise<Customer> => {
    const response = await api.post<Customer>('/customers', data)
    return response.data
  },

  update: async (id: number, data: UpdateCustomerRequest): Promise<Customer> => {
    const response = await api.put<Customer>(`/customers/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/customers/${id}`)
  }
}
