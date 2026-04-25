import api from './client'
import type { StockTransfer, CreateStockTransferRequest } from '../types/order'
import type { PaginatedResponse } from '../types/warehouse'

export const transferApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<StockTransfer>> => {
    const response = await api.get<PaginatedResponse<StockTransfer>>('/stock-transfers', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<StockTransfer> => {
    const response = await api.get<StockTransfer>(`/stock-transfers/${id}`)
    return response.data
  },

  create: async (data: CreateStockTransferRequest): Promise<StockTransfer> => {
    const response = await api.post<StockTransfer>('/stock-transfers', data)
    return response.data
  },

  confirm: async (id: number): Promise<StockTransfer> => {
    const response = await api.post<StockTransfer>(`/stock-transfers/${id}/confirm`)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/stock-transfers/${id}`)
  }
}
