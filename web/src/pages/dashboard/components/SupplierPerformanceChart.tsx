import { Radar } from '@ant-design/charts'
import { Card, Spin } from 'antd'
import type { SupplierPerformance } from '../../../types/dashboard'

interface SupplierPerformanceChartProps {
  data: SupplierPerformance[]
  loading?: boolean
  onClick?: (supplierId: number) => void
}

export function SupplierPerformanceChart({ data, loading, onClick }: SupplierPerformanceChartProps) {
  const config = {
    data: data.flatMap(item => [
      { name: item.supplier_name, label: '订单量', value: item.order_count },
      { name: item.supplier_name, label: '总金额', value: item.total_value / 10000 },
      { name: item.supplier_name, label: '准时率', value: item.on_time_rate },
      { name: item.supplier_name, label: '质量评分', value: item.quality_score },
      { name: item.supplier_name, label: '交付评分', value: item.delivery_score },
    ]),
    xField: 'label',
    yField: 'value',
    seriesField: 'name',
    meta: {
      value: {
        alias: '分数',
        min: 0,
        max: 100,
      },
    },
    radius: 0.8,
    onReady: (plot: any) => {
      if (onClick) {
        plot.on('element:click', (evt: any) => {
          const { data } = evt.data
          if (data && data.name) {
            const supplier = data.find((item: any) => item.name === data.name)
            if (supplier) {
              onClick(supplier.supplier_id)
            }
          }
        })
      }
    },
  }

  return (
    <Card title="供应商绩效" style={{ height: '100%' }}>
      <Spin spinning={loading}>
        <Radar {...config} />
      </Spin>
    </Card>
  )
}
