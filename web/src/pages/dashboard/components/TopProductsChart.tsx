import { Bar } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import { useAuthStore } from '../../../stores/authStore'
import type { TopProduct } from '../../../types/dashboard'

interface TopProductsChartProps {
  data: TopProduct[]
  loading?: boolean
  onBarClick?: (productId: number) => void
}

export function TopProductsChart({ data, loading, onBarClick }: TopProductsChartProps) {
  const { theme } = useAuthStore()
  
  const config = {
    data: data.map(item => ({
      name: item.product_name,
      value: item.total_qty,
      productId: item.product_id,
    })),
    xField: 'value',
    yField: 'name',
    seriesField: 'name',
    legend: false,
    color: '#1890ff',
    theme: theme === 'dark' ? 'classicDark' : 'classic',
    barStyle: {
      radius: [0, 4, 4, 0],
    },
    onReady: (plot: any) => {
      if (onBarClick && plot) {
        plot.on('element:click', (evt: any) => {
          try {
            const eventData = evt?.data?.data || evt?.data
            if (eventData?.productId) {
              onBarClick(eventData.productId)
            }
          } catch (error) {
            console.error('Chart click error:', error)
          }
        })
      }
    },
  }

  return (
    <Card title="热销产品排行" style={{ height: '100%' }}>
      <Spin spinning={loading}>
        <Bar {...config} />
      </Spin>
    </Card>
  )
}
