import { Pie } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import type { WarehouseUsage } from '../../../types/dashboard'

interface WarehouseUsageChartProps {
  data: WarehouseUsage[]
  loading?: boolean
  onClick?: (warehouseId: number) => void
}

export function WarehouseUsageChart({ data, loading, onClick }: WarehouseUsageChartProps) {
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
    label: {
      type: 'inner',
      offset: '-50%',
      content: '{value}%',
      style: {
        textAlign: 'center',
        fontSize: 12,
      },
    },
    statistic: {
      title: {
        content: '平均使用率',
      },
      content: {
        formatter: () => {
          if (data.length === 0) return '0%'
          const avg = data.reduce((sum, item) => sum + item.usage_rate, 0) / data.length
          return `${avg.toFixed(1)}%`
        },
      },
    },
    onReady: (plot: any) => {
      if (onClick) {
        plot.on('element:click', (evt: any) => {
          const { data } = evt.data
          if (data && data.warehouseId) {
            onClick(data.warehouseId)
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
