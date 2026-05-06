import React, { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, Select, message, Popconfirm, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, RightOutlined, DownOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { categoryApi, type CategoryFilter } from '../../api/category'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '../../types/product'

export function CategoryList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [filter, setFilter] = useState<CategoryFilter>({})
  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['categories', page, pageSize, filter.name],
    queryFn: () => categoryApi.list(page, pageSize, filter)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateCategoryRequest) => categoryApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateCategoryRequest }) => categoryApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => categoryApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Category) => {
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
      createMutation.mutate(values as CreateCategoryRequest)
    }
  }

  const handleNameFilterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFilter(prev => ({ ...prev, name: e.target.value || undefined }))
    setPage(1)
  }

  const buildTree = (items: Category[]): Category[] => {
    const map = new Map<number, Category>()
    const roots: Category[] = []

    items.forEach(item => {
      map.set(item.id, { ...item, children: [] })
    })

    items.forEach(item => {
      const node = map.get(item.id)!
      if (item.parent_id && map.has(item.parent_id)) {
        map.get(item.parent_id)!.children!.push(node)
      } else {
        roots.push(node)
      }
    })

    return roots
  }

  const treeData = data?.items ? buildTree(data.items) : []

  const columns = [
    {
      title: '',
      width: 24,
      render: (_: unknown, record: Category) => {
        if (record.children && record.children.length > 0) {
          return null
        }
        return null
      }
    },
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '分类名称', dataIndex: 'name' },
    { title: '排序', dataIndex: 'sort_order', width: 60 },
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
      render: (_: unknown, record: Category) => (
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

  const customExpandIcon = (props: { expanded: boolean; onExpand: (record: Category, e: React.MouseEvent<HTMLElement>) => void; record: Category }) => {
    const { expanded, onExpand, record } = props
    if (record.children && record.children.length > 0) {
      return expanded ? (
        <DownOutlined
          onClick={(e) => onExpand(record, e)}
          style={{ cursor: 'pointer', fontSize: 10, color: '#666' }}
        />
      ) : (
        <RightOutlined
          onClick={(e) => onExpand(record, e)}
          style={{ cursor: 'pointer', fontSize: 10, color: '#666' }}
        />
      )
    }
    return <span style={{ display: 'inline-block', width: 14 }} />
  }

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16 }}>
        <Space>
          <Input
            placeholder="分类名称"
            style={{ width: 150 }}
            value={filter.name || ''}
            onChange={handleNameFilterChange}
            allowClear
          />
          <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
            新增分类
          </Button>
        </Space>
      </div>
      <Table
        columns={columns}
        dataSource={treeData}
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
        defaultExpandAllRows
        expandable={{
          expandIcon: customExpandIcon,
          columnWidth: 24,
        }}
        indentSize={0}
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={editingId ? '编辑分类' : '新增分类'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="分类名称" rules={[{ required: true, message: '请输入分类名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="parent_id" label="父级分类">
            <Select
              allowClear
              placeholder="选择父级分类（可选）"
              options={data?.items?.filter((c: Category) => c.id !== editingId).map((c: Category) => ({ value: c.id, label: c.name })) || []}
            />
          </Form.Item>
          <Form.Item name="sort_order" label="排序" initialValue={0}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
