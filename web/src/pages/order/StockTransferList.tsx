import { useState } from 'react'
import { Table, Button, Space, Drawer, Descriptions, Tag, message, Popconfirm, Input, Select, DatePicker, Form, Row, Col, Card } from 'antd'
import { EyeOutlined, CheckOutlined, DeleteOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { transferApi, type StockTransferFilter } from '../../api/transfer'
import { warehouseApi } from '../../api/warehouse'
import { productApi } from '../../api/product'
import type { StockTransfer } from '../../types/order'

const { RangePicker } = DatePicker

const statusMap: Record<number, { text: string; color: string }> = {
  0: { text: '待确认', color: 'orange' },
  1: { text: '已完成', color: 'green' },
  2: { text: '已取消', color: 'red' }
}

export function StockTransferList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [selectedOrder, setSelectedOrder] = useState<StockTransfer | null>(null)
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()
  const [filter, setFilter] = useState<StockTransferFilter>({})
  const [form] = Form.useForm()

  const { data, isLoading } = useQuery({
    queryKey: ['stock-transfers', page, pageSize, filter],
    queryFn: () => transferApi.list(page, pageSize, filter)
  })

  const { data: warehouses } = useQuery({
    queryKey: ['warehouses-all'],
    queryFn: () => warehouseApi.list(1, 100)
  })

  const { data: products } = useQuery({
    queryKey: ['products-all'],
    queryFn: () => productApi.list(1, 100)
  })

  const confirmMutation = useMutation({
    mutationFn: (id: number) => transferApi.confirm(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stock-transfers'] })
      messageApi.success('确认成功')
    },
    onError: () => messageApi.error('确认失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => transferApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stock-transfers'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleViewDetail = async (id: number) => {
    const order = await transferApi.get(id)
    setSelectedOrder(order)
    setDrawerOpen(true)
  }

  const handleSearch = () => {
    const values = form.getFieldsValue()
    const newFilter: StockTransferFilter = {}
    
    if (values.order_no) newFilter.order_no = values.order_no
    if (values.source_warehouse_id) newFilter.source_warehouse_id = values.source_warehouse_id
    if (values.target_warehouse_id) newFilter.target_warehouse_id = values.target_warehouse_id
    if (values.created_at_range && values.created_at_range[0] && values.created_at_range[1]) {
      newFilter.created_at_start = values.created_at_range[0].format('YYYY-MM-DDTHH:mm:ssZ')
      newFilter.created_at_end = values.created_at_range[1].format('YYYY-MM-DDTHH:mm:ssZ')
    }
    
    setFilter(newFilter)
    setPage(1)
  }

  const handleReset = () => {
    form.resetFields()
    setFilter({})
    setPage(1)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '调拨单号', dataIndex: 'order_no', width: 140 },
    {
      title: '调出仓库',
      dataIndex: 'source_warehouse_id',
      width: 100,
      ellipsis: true,
      render: (id: number) => warehouses?.items?.find(w => w.id === id)?.name || id
    },
    {
      title: '调入仓库',
      dataIndex: 'target_warehouse_id',
      width: 100,
      ellipsis: true,
      render: (id: number) => warehouses?.items?.find(w => w.id === id)?.name || id
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 80,
      render: (status: number) => {
        const s = statusMap[status] || { text: '未知', color: 'default' }
        return <Tag color={s.color}>{s.text}</Tag>
      }
    },
    { title: '创建时间', dataIndex: 'created_at', width: 150 },
    {
      title: '操作',
      width: 160,
      render: (_: unknown, record: StockTransfer) => (
        <Space>
          <Button type="link" size="small" icon={<EyeOutlined />} onClick={() => handleViewDetail(record.id)}>
            详情
          </Button>
          {record.status === 0 && (
            <Popconfirm title="确认调拨?" onConfirm={() => confirmMutation.mutate(record.id)}>
              <Button type="link" size="small" icon={<CheckOutlined />}>
                确认
              </Button>
            </Popconfirm>
          )}
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const itemColumns = [
    { title: '商品', dataIndex: 'product_id', render: (id: number) => products?.items?.find((p: { id: number; name: string }) => p.id === id)?.name || id },
    { title: '数量', dataIndex: 'quantity', width: 100 },
    { title: '批次号', dataIndex: 'batch_no', width: 120 }
  ]

  return (
    <>
      {contextHolder}
      <Card style={{ marginBottom: 16 }}>
        <Form form={form} layout="inline">
          <Row gutter={16}>
            <Col>
              <Form.Item name="order_no" label="订单编号">
                <Input placeholder="输入订单编号" style={{ width: 150 }} allowClear />
              </Form.Item>
            </Col>
            <Col>
              <Form.Item name="source_warehouse_id" label="调出仓库">
                <Select
                  placeholder="选择调出仓库"
                  style={{ width: 150 }}
                  allowClear
                  showSearch
                  filterOption={(input, option) =>
                    (option?.label ?? '').toString().toLowerCase().includes(input.toLowerCase())
                  }
                  options={warehouses?.items.map((w: { id: number; name: string }) => ({
                    label: w.name,
                    value: w.id
                  }))}
                />
              </Form.Item>
            </Col>
            <Col>
              <Form.Item name="target_warehouse_id" label="调入仓库">
                <Select
                  placeholder="选择调入仓库"
                  style={{ width: 150 }}
                  allowClear
                  showSearch
                  filterOption={(input, option) =>
                    (option?.label ?? '').toString().toLowerCase().includes(input.toLowerCase())
                  }
                  options={warehouses?.items.map((w: { id: number; name: string }) => ({
                    label: w.name,
                    value: w.id
                  }))}
                />
              </Form.Item>
            </Col>
            <Col>
              <Form.Item name="created_at_range" label="创建时间">
                <RangePicker showTime style={{ width: 350 }} />
              </Form.Item>
            </Col>
            <Col>
              <Space>
                <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
                  搜索
                </Button>
                <Button icon={<ReloadOutlined />} onClick={handleReset}>
                  重置
                </Button>
              </Space>
            </Col>
          </Row>
        </Form>
      </Card>
      <Table
        columns={columns}
        dataSource={data?.items}
        rowKey="id"
        loading={isLoading}
        pagination={{
          current: page,
          pageSize,
          total: data?.total,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
          onChange: (p, ps) => {
            setPage(p)
            setPageSize(ps)
          }
        }}
        scroll={{ x: 'max-content' }}
      />
      <Drawer
        title="调拨单详情"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        width={600}
      >
        {selectedOrder && (
          <>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="调拨单号">{selectedOrder.order_no}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={statusMap[selectedOrder.status]?.color}>
                  {statusMap[selectedOrder.status]?.text}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="调出仓库">
                {warehouses?.items?.find(w => w.id === selectedOrder.source_warehouse_id)?.name}
              </Descriptions.Item>
              <Descriptions.Item label="调入仓库">
                {warehouses?.items?.find(w => w.id === selectedOrder.target_warehouse_id)?.name}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">{selectedOrder.created_at}</Descriptions.Item>
            </Descriptions>
            <h4 style={{ marginTop: 16 }}>商品明细</h4>
            <Table
              columns={itemColumns}
              dataSource={selectedOrder.items || []}
              rowKey="id"
              size="small"
              pagination={false}
            />
          </>
        )}
      </Drawer>
    </>
  )
}
