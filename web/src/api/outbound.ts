import api from './client'
import type { OutboundOrder, CreateOutboundOrderRequest } from '../types/order'
import type { PaginatedResponse } from '../types/warehouse'

export interface OutboundOrderFilter {
  order_no?: string
  customer_id?: number
  warehouse_id?: number
  quantity_min?: number
  quantity_max?: number
  created_at_start?: string
  created_at_end?: string
}

export const outboundApi = {
  list: async (page = 1, size = 10, filter?: OutboundOrderFilter): Promise<PaginatedResponse<OutboundOrder>> => {
    const response = await api.get<PaginatedResponse<OutboundOrder>>('/outbound-orders', {
      params: { page, size, ...filter }
    })
    return response.data
  },

  get: async (id: number): Promise<OutboundOrder> => {
    const response = await api.get<OutboundOrder>(`/outbound-orders/${id}`)
    return response.data
  },

  create: async (data: CreateOutboundOrderRequest): Promise<OutboundOrder> => {
    const response = await api.post<OutboundOrder>('/outbound-orders', data)
    return response.data
  },

  confirm: async (id: number): Promise<OutboundOrder> => {
    const response = await api.post<OutboundOrder>(`/outbound-orders/${id}/confirm`)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/outbound-orders/${id}`)
  }
}
