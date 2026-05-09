import { useQuery } from '@tanstack/react-query'
import { dashboardApi } from '../../../api/dashboard'
import dayjs from 'dayjs'

export function useDashboardStats(startDate?: string, endDate?: string) {
  const start = startDate || dayjs().subtract(30, 'day').format('YYYY-MM-DD')
  const end = endDate || dayjs().format('YYYY-MM-DD')

  const overviewQuery = useQuery({
    queryKey: ['dashboard', 'overview'],
    queryFn: dashboardApi.getOverview,
  })

  const trendQuery = useQuery({
    queryKey: ['dashboard', 'trend', start, end],
    queryFn: () => dashboardApi.getTrend(start, end),
  })

  const topProductsQuery = useQuery({
    queryKey: ['dashboard', 'topProducts', start, end],
    queryFn: () => dashboardApi.getTopProducts(start, end),
  })

  const warehouseUsageQuery = useQuery({
    queryKey: ['dashboard', 'warehouseUsage'],
    queryFn: dashboardApi.getWarehouseUsage,
  })

  const supplierPerformanceQuery = useQuery({
    queryKey: ['dashboard', 'supplierPerformance', start, end],
    queryFn: () => dashboardApi.getSupplierPerformance(start, end),
  })

  const pendingOrdersQuery = useQuery({
    queryKey: ['dashboard', 'pendingOrders'],
    queryFn: dashboardApi.getPendingOrders,
  })

  const refetchAll = () => {
    overviewQuery.refetch()
    trendQuery.refetch()
    topProductsQuery.refetch()
    warehouseUsageQuery.refetch()
    supplierPerformanceQuery.refetch()
    pendingOrdersQuery.refetch()
  }

  return {
    overview: overviewQuery.data,
    trend: trendQuery.data,
    topProducts: topProductsQuery.data,
    warehouseUsage: warehouseUsageQuery.data,
    supplierPerformance: supplierPerformanceQuery.data,
    pendingOrders: pendingOrdersQuery.data,
    loading: 
      overviewQuery.isLoading || 
      trendQuery.isLoading || 
      topProductsQuery.isLoading ||
      warehouseUsageQuery.isLoading ||
      supplierPerformanceQuery.isLoading ||
      pendingOrdersQuery.isLoading,
    refetch: refetchAll,
  }
}
