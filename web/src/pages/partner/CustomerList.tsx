import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, message, Popconfirm, Tag, Row, Col } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { customerApi } from '../../api/customer'
import type { Customer, CreateCustomerRequest, UpdateCustomerRequest } from '../../types/partner'

export function CustomerList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['customers', page, pageSize],
    queryFn: () => customerApi.list(page, pageSize)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateCustomerRequest) => customerApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateCustomerRequest }) => customerApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => customerApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Customer) => {
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
      createMutation.mutate(values as CreateCustomerRequest)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '客户编码', dataIndex: 'code', width: 120 },
    { title: '客户名称', dataIndex: 'name' },
    { title: '联系人', dataIndex: 'contact', width: 100 },
    { title: '联系电话', dataIndex: 'phone', width: 130 },
    { title: '邮箱', dataIndex: 'email', width: 180 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      )
    },
    {
      title: '操作',
      width: 150,
      render: (_: unknown, record: Customer) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
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
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增客户
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
        title={editingId ? '编辑客户' : '新增客户'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={600}
      >
        <Form form={form} layout="horizontal" labelCol={{ span: 6 }} wrapperCol={{ span: 18 }} style={{ marginTop: 16 }}>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="name" label="客户名称" rules={[{ required: true, message: '请输入客户名称' }]}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="code" label="客户编码">
                <Input />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="contact" label="联系人">
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="phone" label="联系电话">
                <Input />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="email" label="邮箱">
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态" initialValue={1}>
                <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="address" label="地址" labelCol={{ span: 3 }} wrapperCol={{ span: 21 }}>
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
