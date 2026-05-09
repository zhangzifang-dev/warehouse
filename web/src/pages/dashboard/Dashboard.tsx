import { useState } from 'react'
import { Row, Col, DatePicker, Button, Space, message } from 'antd'
import { ReloadOutlined, DownloadOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import dayjs, { Dayjs } from 'dayjs'
import { useDashboardStats } from './hooks/useDashboardStats'
import { ErrorBoundary } from '../../components/ErrorBoundary'
import { 
  InventoryCard, 
  WarningCard, 
  TodayInboundCard, 
  TodayOutboundCard 
} from './components/StatCard'
import { TrendChart } from './components/TrendChart'
import { TopProductsChart } from './components/TopProductsChart'
import { WarehouseUsageChart } from './components/WarehouseUsageChart'
import { SupplierPerformanceChart } from './components/SupplierPerformanceChart'

const { RangePicker } = DatePicker

export function Dashboard() {
  const navigate = useNavigate()
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs]>([
    dayjs().subtract(30, 'day'),
    dayjs()
  ])
  
  const { 
    overview, 
    trend, 
    topProducts, 
    warehouseUsage, 
    supplierPerformance,
    loading,
    refetch 
  } = useDashboardStats(
    dateRange[0].format('YYYY-MM-DD'),
    dateRange[1].format('YYYY-MM-DD')
  )

  const handleDateChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
    if (dates && dates[0] && dates[1]) {
      setDateRange([dates[0], dates[1]])
    }
  }

  const handleRefresh = () => {
    refetch()
    message.success('数据已刷新')
  }

  const handleExport = async () => {
    message.info('导出功能开发中...')
  }

  return (
    <div style={{ padding: '24px' }}>
      <Row justify="end" style={{ marginBottom: '24px' }}>
        <Col>
          <Space>
            <RangePicker
              value={dateRange}
              onChange={handleDateChange}
              format="YYYY-MM-DD"
              allowClear={false}
            />
            <Button icon={<ReloadOutlined />} onClick={handleRefresh}>
              刷新
            </Button>
            <Button icon={<DownloadOutlined />} onClick={handleExport}>
              导出
            </Button>
          </Space>
        </Col>
      </Row>

        <ErrorBoundary>
          <Row gutter={[16, 16]}>
            <Col xs={24} sm={12} lg={6}>
              <InventoryCard 
                value={overview?.total_inventory || 0} 
                warning={overview?.inventory_warning || 0}
                onClick={() => navigate('/inventory')}
              />
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <WarningCard 
                value={overview?.inventory_warning || 0}
                onClick={() => navigate('/inventory?quantity_max=10')}
              />
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <TodayInboundCard 
                orders={overview?.today_inbound || 0}
                quantity={overview?.today_inbound_qty || 0}
                onClick={() => navigate('/inbound')}
              />
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <TodayOutboundCard 
                orders={overview?.today_outbound || 0}
                quantity={overview?.today_outbound_qty || 0}
                onClick={() => navigate('/outbound')}
              />
            </Col>
          </Row>

        <Row gutter={[16, 16]} style={{ marginTop: '16px' }}>
          <Col xs={24} lg={16}>
            <ErrorBoundary>
              <TrendChart 
                data={trend || []} 
                loading={loading}
                onPointClick={(date, type) => {
                  const path = type === 'inbound' ? '/inbound' : '/outbound'
                  navigate(`${path}?date=${date}`)
                }}
              />
            </ErrorBoundary>
          </Col>
          <Col xs={24} lg={8}>
            <ErrorBoundary>
              <TopProductsChart 
                data={topProducts || []} 
                loading={loading}
                onBarClick={(productId) => {
                  navigate(`/products?id=${productId}`)
                }}
              />
            </ErrorBoundary>
          </Col>
        </Row>

        <Row gutter={[16, 16]} style={{ marginTop: '16px' }}>
          <Col xs={24} lg={12}>
            <ErrorBoundary>
              <WarehouseUsageChart 
                data={warehouseUsage || []} 
                loading={loading}
                onClick={(warehouseId) => {
                  navigate(`/inventory?warehouse_id=${warehouseId}`)
                }}
              />
            </ErrorBoundary>
          </Col>
          <Col xs={24} lg={12}>
            <ErrorBoundary>
              <SupplierPerformanceChart 
                data={supplierPerformance || []} 
                loading={loading}
                onClick={(supplierId) => {
                  navigate(`/suppliers?id=${supplierId}`)
                }}
              />
            </ErrorBoundary>
          </Col>
        </Row>
      </ErrorBoundary>
    </div>
  )
}
