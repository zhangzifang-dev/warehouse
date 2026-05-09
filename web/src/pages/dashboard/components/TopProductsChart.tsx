import { Bar } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import type { TopProduct } from '../../../types/dashboard'

interface TopProductsChartProps {
  data: TopProduct[]
  loading?: boolean
  onBarClick?: (productId: number) => void
}

export function TopProductsChart({ data, loading, onBarClick }: TopProductsChartProps) {
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
    barStyle: {
      radius: [0, 4, 4, 0],
    },
    onReady: (plot: any) => {
      if (onBarClick) {
        plot.on('element:click', (evt: any) => {
          const { data } = evt.data
          if (data && data.productId) {
            onBarClick(data.productId)
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
