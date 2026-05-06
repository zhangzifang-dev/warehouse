import api from './client'
import type { InboundOrder, CreateInboundOrderRequest } from '../types/order'
import type { PaginatedResponse } from '../types/warehouse'

export interface InboundOrderFilter {
  order_no?: string
  supplier_id?: number
  warehouse_id?: number
  quantity_min?: number
  quantity_max?: number
  created_at_start?: string
  created_at_end?: string
}

export const inboundApi = {
  list: async (page = 1, size = 10, filter?: InboundOrderFilter): Promise<PaginatedResponse<InboundOrder>> => {
    const response = await api.get<PaginatedResponse<InboundOrder>>('/inbound-orders', {
      params: { page, size, ...filter }
    })
    return response.data
  },

  get: async (id: number): Promise<InboundOrder> => {
    const response = await api.get<InboundOrder>(`/inbound-orders/${id}`)
    return response.data
  },

  create: async (data: CreateInboundOrderRequest): Promise<InboundOrder> => {
    const response = await api.post<InboundOrder>('/inbound-orders', data)
    return response.data
  },

  confirm: async (id: number): Promise<InboundOrder> => {
    const response = await api.post<InboundOrder>(`/inbound-orders/${id}/confirm`)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/inbound-orders/${id}`)
  }
}
