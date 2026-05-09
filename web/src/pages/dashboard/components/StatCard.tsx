import { Card, Statistic } from 'antd'
import { ArrowUpOutlined, ArrowDownOutlined, WarningOutlined } from '@ant-design/icons'
import type { ReactNode } from 'react'

interface StatCardProps {
  title: string
  value: number | string
  prefix?: ReactNode
  suffix?: string
  valueStyle?: React.CSSProperties
  onClick?: () => void
  loading?: boolean
}

export function StatCard({ 
  title, 
  value, 
  prefix, 
  suffix, 
  valueStyle, 
  onClick,
  loading 
}: StatCardProps) {
  return (
    <Card 
      hoverable={!!onClick} 
      onClick={onClick}
      loading={loading}
      styles={{
        body: { padding: '20px' }
      }}
    >
      <Statistic
        title={title}
        value={value}
        prefix={prefix}
        suffix={suffix}
        valueStyle={{ fontSize: '24px', ...valueStyle }}
      />
    </Card>
  )
}

export function InventoryCard({ value, warning, onClick }: { 
  value: number
  warning: number
  onClick?: () => void 
}) {
  return (
    <StatCard
      title="总库存量"
      value={value}
      suffix="件"
      onClick={onClick}
      prefix={warning > 0 ? <WarningOutlined style={{ color: '#faad14' }} /> : undefined}
    />
  )
}

export function WarningCard({ value, onClick }: { value: number; onClick?: () => void }) {
  return (
    <StatCard
      title="库存预警"
      value={value}
      suffix="项"
      onClick={onClick}
      valueStyle={value > 0 ? { color: '#ff4d4f' } : undefined}
      prefix={<WarningOutlined />}
    />
  )
}

export function TodayInboundCard({ orders, quantity, onClick }: { 
  orders: number
  quantity: number
  onClick?: () => void 
}) {
  return (
    <StatCard
      title="今日入库"
      value={orders}
      suffix={`单 / ${quantity} 件`}
      onClick={onClick}
      prefix={<ArrowDownOutlined style={{ color: '#52c41a' }} />}
    />
  )
}

export function TodayOutboundCard({ orders, quantity, onClick }: { 
  orders: number
  quantity: number
  onClick?: () => void 
}) {
  return (
    <StatCard
      title="今日出库"
      value={orders}
      suffix={`单 / ${quantity} 件`}
      onClick={onClick}
      prefix={<ArrowUpOutlined style={{ color: '#1890ff' }} />}
    />
  )
}
