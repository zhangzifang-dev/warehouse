import { Line } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import { useAuthStore } from '../../../stores/authStore'
import type { TrendData } from '../../../types/dashboard'

interface TrendChartProps {
  data: TrendData[]
  loading?: boolean
  onPointClick?: (date: string, type: 'inbound' | 'outbound') => void
}

export function TrendChart({ data, loading, onPointClick }: TrendChartProps) {
  const { theme } = useAuthStore()
  
  const config = {
    data: data.flatMap(item => [
      { date: item.date, value: item.inbound_qty, type: '入库' },
      { date: item.date, value: item.outbound_qty, type: '出库' }
    ]),
    xField: 'date',
    yField: 'value',
    seriesField: 'type',
    color: ['#1890ff', '#fa8c16'],
    smooth: true,
    theme: theme === 'dark' ? 'classicDark' : 'classic',
    animation: {
      appear: {
        animation: 'path-in',
        duration: 1000,
      },
    },
    point: {
      shape: 'circle',
      size: 4,
    },
    interactions: [
      {
        type: 'marker-active',
      },
    ],
    onReady: (plot: any) => {
      if (onPointClick) {
        plot.on('element:click', (evt: any) => {
          const { data } = evt.data
          if (data && data.date) {
            const type = data.type === '入库' ? 'inbound' : 'outbound'
            onPointClick(data.date, type)
          }
        })
      }
    },
  }

  return (
    <Card title="出入库趋势" style={{ height: '100%' }}>
      <Spin spinning={loading}>
        <Line {...config} />
      </Spin>
    </Card>
  )
}
