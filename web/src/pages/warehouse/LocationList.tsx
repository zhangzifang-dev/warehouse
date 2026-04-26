import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, message, Popconfirm, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { locationApi } from '../../api/location'
import { warehouseApi } from '../../api/warehouse'
import type { Location, CreateLocationRequest, UpdateLocationRequest } from '../../types/warehouse'

export function LocationList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [warehouseFilter, setWarehouseFilter] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['locations', page, pageSize, warehouseFilter],
    queryFn: () => locationApi.list(page, pageSize, warehouseFilter)
  })

  const { data: warehouses } = useQuery({
    queryKey: ['warehouses-all'],
    queryFn: () => warehouseApi.list(1, 100)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateLocationRequest) => locationApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['locations'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateLocationRequest }) => locationApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['locations'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => locationApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['locations'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Location) => {
    setEditingId(record.id)
    form.setFieldsValue(record)
    setModalOpen(true)
  }

  const handleCloseModal = () => {
    setModalOpen(false)
    setEditingId(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: values })
    } else {
      createMutation.mutate(values as CreateLocationRequest)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '库位编码', dataIndex: 'code', width: 100 },
    {
      title: '所属仓库',
      dataIndex: 'warehouse_id',
      width: 100,
      ellipsis: true,
      render: (id: number) => warehouses?.items?.find(w => w.id === id)?.name || id
    },
    { title: '区域', dataIndex: 'zone', width: 60 },
    { title: '货架', dataIndex: 'shelf', width: 60 },
    { title: '层', dataIndex: 'level', width: 60 },
    { title: '位', dataIndex: 'position', width: 60 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 80,
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      )
    },
    {
      title: '操作',
      width: 140,
      render: (_: unknown, record: Location) => (
        <Space>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
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
          options={warehouses?.items?.map(w => ({ value: w.id, label: w.name })) || []}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增库位
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
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={editingId ? '编辑库位' : '新增库位'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="warehouse_id" label="所属仓库" rules={[{ required: true, message: '请选择仓库' }]}>
            <Select
              options={warehouses?.items?.map(w => ({ value: w.id, label: w.name })) || []}
            />
          </Form.Item>
          <Form.Item name="zone" label="区域" rules={[{ required: true, message: '请输入区域' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="shelf" label="货架" rules={[{ required: true, message: '请输入货架' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="level" label="层" rules={[{ required: true, message: '请输入层' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="position" label="位" rules={[{ required: true, message: '请输入位' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
