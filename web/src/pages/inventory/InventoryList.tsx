import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, Select, message, Tag } from 'antd'
import { PlusOutlined, EditOutlined, ToolOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { inventoryApi } from '../../api/inventory'
import { warehouseApi } from '../../api/warehouse'
import { productApi } from '../../api/product'
import type { Inventory, CreateInventoryRequest, UpdateInventoryRequest, AdjustQuantityRequest } from '../../types/inventory'

export function InventoryList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [warehouseFilter, setWarehouseFilter] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)
  const [adjustModalOpen, setAdjustModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [adjustForm] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['inventory', page, pageSize, warehouseFilter],
    queryFn: () => inventoryApi.list(page, pageSize, warehouseFilter)
  })

  const { data: warehouses } = useQuery({
    queryKey: ['warehouses-all'],
    queryFn: () => warehouseApi.list(1, 100)
  })

  const { data: products } = useQuery({
    queryKey: ['products-all'],
    queryFn: () => productApi.list(1, 100)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateInventoryRequest) => inventoryApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateInventoryRequest }) => inventoryApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const adjustMutation = useMutation({
    mutationFn: (data: AdjustQuantityRequest) => inventoryApi.adjust(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] })
      messageApi.success('调整成功')
      handleCloseAdjustModal()
    },
    onError: () => messageApi.error('调整失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Inventory) => {
    setEditingId(record.id)
    form.setFieldsValue(record)
    setModalOpen(true)
  }

  const handleCloseModal = () => {
    setModalOpen(false)
    setEditingId(null)
    form.resetFields()
  }

  const handleOpenAdjust = (record: Inventory) => {
    adjustForm.setFieldsValue({ inventory_id: record.id, quantity: 0 })
    setAdjustModalOpen(true)
  }

  const handleCloseAdjustModal = () => {
    setAdjustModalOpen(false)
    adjustForm.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: values })
    } else {
      createMutation.mutate(values as CreateInventoryRequest)
    }
  }

  const handleAdjustSubmit = async () => {
    const values = await adjustForm.validateFields()
    adjustMutation.mutate(values)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    {
      title: '仓库',
      dataIndex: 'warehouse_id',
      width: 150,
      render: (id: number) => warehouses?.items.find(w => w.id === id)?.name || id
    },
    {
      title: '商品',
      dataIndex: 'product_id',
      render: (id: number) => products?.items.find(p => p.id === id)?.name || id
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      width: 120,
      render: (qty: number) => (
        <Tag color={qty > 0 ? 'green' : qty < 0 ? 'red' : 'default'}>
          {qty}
        </Tag>
      )
    },
    { title: '批次号', dataIndex: 'batch_no', width: 120 },
    {
      title: '操作',
      width: 180,
      render: (_: unknown, record: Inventory) => (
        <Space>
          <Button type="link" icon={<ToolOutlined />} onClick={() => handleOpenAdjust(record)}>
            调整
          </Button>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
        </Space>
      )
    }
  ]

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <Select
          allowClear
          placeholder="筛选仓库"
          style={{ width: 200 }}
          value={warehouseFilter}
          onChange={setWarehouseFilter}
          options={warehouses?.items.map(w => ({ value: w.id, label: w.name }))}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增库存
        </Button>
      </div>
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
      <Modal
        title={editingId ? '编辑库存' : '新增库存'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="warehouse_id" label="仓库" rules={[{ required: true, message: '请选择仓库' }]}>
            <Select
              options={warehouses?.items.map(w => ({ value: w.id, label: w.name }))}
            />
          </Form.Item>
          <Form.Item name="product_id" label="商品" rules={[{ required: true, message: '请选择商品' }]}>
            <Select
              showSearch
              optionFilterProp="label"
              options={products?.items.map((p: { id: number; sku: string; name: string }) => ({ value: p.id, label: `${p.sku} - ${p.name}` }))}
            />
          </Form.Item>
          <Form.Item name="quantity" label="数量">
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="batch_no" label="批次号">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="库存调整"
        open={adjustModalOpen}
        onOk={handleAdjustSubmit}
        onCancel={handleCloseAdjustModal}
        confirmLoading={adjustMutation.isPending}
      >
        <Form form={adjustForm} layout="vertical">
          <Form.Item name="inventory_id" label="库存ID" hidden>
            <InputNumber />
          </Form.Item>
          <Form.Item name="quantity" label="调整数量（正数入库，负数出库）" rules={[{ required: true, message: '请输入调整数量' }]}>
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
