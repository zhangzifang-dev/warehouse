import api from './client'
import type { 
  OverviewStats, 
  TrendData, 
  TopProduct, 
  WarehouseUsage, 
  SupplierPerformance,
  PendingOrders 
} from '../types/dashboard'

export const dashboardApi = {
  getOverview: async (): Promise<OverviewStats> => {
    const response = await api.get<OverviewStats>('/dashboard/overview')
    return response.data
  },

  getTrend: async (startDate: string, endDate: string): Promise<TrendData[]> => {
    const response = await api.get<TrendData[]>('/dashboard/trend', {
      params: { start_date: startDate, end_date: endDate }
    })
    return response.data
  },

  getTopProducts: async (startDate: string, endDate: string, limit = 10): Promise<TopProduct[]> => {
    const response = await api.get<TopProduct[]>('/dashboard/top-products', {
      params: { start_date: startDate, end_date: endDate, limit }
    })
    return response.data
  },

  getWarehouseUsage: async (): Promise<WarehouseUsage[]> => {
    const response = await api.get<WarehouseUsage[]>('/dashboard/warehouse-usage')
    return response.data
  },

  getSupplierPerformance: async (startDate: string, endDate: string, limit = 10): Promise<SupplierPerformance[]> => {
    const response = await api.get<SupplierPerformance[]>('/dashboard/supplier-performance', {
      params: { start_date: startDate, end_date: endDate, limit }
    })
    return response.data
  },

  getPendingOrders: async (): Promise<PendingOrders> => {
    const response = await api.get<PendingOrders>('/dashboard/pending-orders')
    return response.data
  }
}
