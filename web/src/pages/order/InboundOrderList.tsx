import { useState } from 'react'
import { Table, Button, Space, Drawer, Descriptions, Tag, message, Popconfirm } from 'antd'
import { EyeOutlined, CheckOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { inboundApi } from '../../api/inbound'
import { warehouseApi } from '../../api/warehouse'
import { supplierApi } from '../../api/supplier'
import { productApi } from '../../api/product'
import type { InboundOrder } from '../../types/order'

const statusMap: Record<number, { text: string; color: string }> = {
  0: { text: '待确认', color: 'orange' },
  1: { text: '已完成', color: 'green' },
  2: { text: '已取消', color: 'red' }
}

export function InboundOrderList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [selectedOrder, setSelectedOrder] = useState<InboundOrder | null>(null)
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['inbound-orders', page, pageSize],
    queryFn: () => inboundApi.list(page, pageSize)
  })

  const { data: warehouses } = useQuery({
    queryKey: ['warehouses-all'],
    queryFn: () => warehouseApi.list(1, 100)
  })

  const { data: suppliers } = useQuery({
    queryKey: ['suppliers-all'],
    queryFn: () => supplierApi.list(1, 100)
  })

  const { data: products } = useQuery({
    queryKey: ['products-all'],
    queryFn: () => productApi.list(1, 100)
  })

  const confirmMutation = useMutation({
    mutationFn: (id: number) => inboundApi.confirm(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inbound-orders'] })
      messageApi.success('确认成功')
    },
    onError: () => messageApi.error('确认失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => inboundApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inbound-orders'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleViewDetail = async (id: number) => {
    const order = await inboundApi.get(id)
    setSelectedOrder(order)
    setDrawerOpen(true)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '订单编号', dataIndex: 'order_no', width: 180 },
    {
      title: '供应商',
      dataIndex: 'supplier_id',
      width: 150,
      render: (id: number | null) => id ? suppliers?.items.find((s: { id: number; name: string }) => s.id === id)?.name || id : '-'
    },
    {
      title: '仓库',
      dataIndex: 'warehouse_id',
      width: 150,
      render: (id: number) => warehouses?.items.find((w: { id: number; name: string }) => w.id === id)?.name || id
    },
    { title: '总数量', dataIndex: 'total_quantity', width: 100 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => {
        const s = statusMap[status] || { text: '未知', color: 'default' }
        return <Tag color={s.color}>{s.text}</Tag>
      }
    },
    { title: '备注', dataIndex: 'remark', ellipsis: true },
    { title: '创建时间', dataIndex: 'created_at', width: 180 },
    {
      title: '操作',
      width: 180,
      render: (_: unknown, record: InboundOrder) => (
        <Space>
          <Button type="link" icon={<EyeOutlined />} onClick={() => handleViewDetail(record.id)}>
            详情
          </Button>
          {record.status === 0 && (
            <Popconfirm title="确认入库?" onConfirm={() => confirmMutation.mutate(record.id)}>
              <Button type="link" icon={<CheckOutlined />}>
                确认
              </Button>
            </Popconfirm>
          )}
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const itemColumns = [
    { title: '商品', dataIndex: 'product_id', render: (id: number) => products?.items.find((p: { id: number; name: string }) => p.id === id)?.name || id },
    { title: '数量', dataIndex: 'quantity', width: 100 },
    { title: '批次号', dataIndex: 'batch_no', width: 120 }
  ]

  return (
    <>
      {contextHolder}
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
      />
      <Drawer
        title="入库单详情"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        width={600}
      >
        {selectedOrder && (
          <>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="订单编号">{selectedOrder.order_no}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={statusMap[selectedOrder.status]?.color}>
                  {statusMap[selectedOrder.status]?.text}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="供应商">
                {selectedOrder.supplier_id ? suppliers?.items.find((s: { id: number; name: string }) => s.id === selectedOrder.supplier_id)?.name : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="仓库">
                {warehouses?.items.find((w: { id: number; name: string }) => w.id === selectedOrder.warehouse_id)?.name}
              </Descriptions.Item>
              <Descriptions.Item label="总数量">{selectedOrder.total_quantity}</Descriptions.Item>
              <Descriptions.Item label="创建时间">{selectedOrder.created_at}</Descriptions.Item>
              <Descriptions.Item label="备注" span={2}>{selectedOrder.remark || '-'}</Descriptions.Item>
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
