import { Pie } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import { useAuthStore } from '../../../stores/authStore'
import type { WarehouseUsage } from '../../../types/dashboard'

interface WarehouseUsageChartProps {
  data: WarehouseUsage[]
  loading?: boolean
  onClick?: (warehouseId: number) => void
}

export function WarehouseUsageChart({ data, loading, onClick }: WarehouseUsageChartProps) {
  const { theme } = useAuthStore()
  
  const config = {
    data: data.map(item => ({
      type: item.warehouse_name,
      value: item.usage_rate,
      warehouseId: item.warehouse_id,
    })),
    angleField: 'value',
    colorField: 'type',
    radius: 0.8,
    innerRadius: 0.6,
    theme: theme === 'dark' ? 'classicDark' : 'classic',
    label: {
      text: 'value',
      position: 'inside',
      style: {
        textAlign: 'center',
        fontSize: 12,
      },
    },
    legend: {
      position: 'bottom' as const,
    },
    tooltip: {
      title: 'type',
      items: [{ channel: 'value' }],
    },
    onReady: (plot: any) => {
      if (onClick && plot) {
        plot.on('element:click', (evt: any) => {
          try {
            const eventData = evt?.data?.data || evt?.data
            if (eventData?.warehouseId) {
              onClick(eventData.warehouseId)
            }
          } catch (error) {
            console.error('Chart click error:', error)
          }
        })
      }
    },
  }

  return (
    <Card title="仓库使用率" style={{ height: '100%' }}>
      <Spin spinning={loading}>
        <Pie {...config} />
      </Spin>
    </Card>
  )
}
