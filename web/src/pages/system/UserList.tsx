import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, message, Popconfirm, Tag, Transfer } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, UserOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { userApi } from '../../api/user'
import { roleApi } from '../../api/role'
import type { User, CreateUserRequest, UpdateUserRequest } from '../../types/system'

export function UserList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [roleModalOpen, setRoleModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [selectedUserId, setSelectedUserId] = useState<number | null>(null)
  const [targetKeys, setTargetKeys] = useState<number[]>([])
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['users', page, pageSize],
    queryFn: () => userApi.list(page, pageSize)
  })

  const { data: rolesData } = useQuery({
    queryKey: ['roles'],
    queryFn: () => roleApi.list(1, 100)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateUserRequest) => userApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserRequest }) => userApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => userApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const assignRolesMutation = useMutation({
    mutationFn: ({ userId, roleIds }: { userId: number; roleIds: number[] }) => userApi.assignRoles(userId, roleIds),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['userRoles'] })
      messageApi.success('分配成功')
      setRoleModalOpen(false)
    },
    onError: () => messageApi.error('分配失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: User) => {
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
      createMutation.mutate(values as CreateUserRequest)
    }
  }

  const handleOpenRoleModal = async (userId: number) => {
    setSelectedUserId(userId)
    const roles = await userApi.getRoles(userId)
    setTargetKeys(roles.map(r => r.id))
    setRoleModalOpen(true)
  }

  const handleAssignRoles = () => {
    if (selectedUserId) {
      assignRolesMutation.mutate({ userId: selectedUserId, roleIds: targetKeys })
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '用户名', dataIndex: 'username', width: 150 },
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
      width: 220,
      render: (_: unknown, record: User) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Button type="link" icon={<UserOutlined />} onClick={() => handleOpenRoleModal(record.id)}>
            分配角色
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

  const transferDataSource = rolesData?.items?.map(r => ({
    key: r.id,
    title: r.name,
    description: r.code
  })) || []

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增用户
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
        title={editingId ? '编辑用户' : '新增用户'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="username" label="用户名" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input disabled={!!editingId} />
          </Form.Item>
          {!editingId && (
            <Form.Item name="password" label="密码" rules={[{ required: true, message: '请输入密码' }]}>
              <Input.Password />
            </Form.Item>
          )}
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="分配角色"
        open={roleModalOpen}
        onOk={handleAssignRoles}
        onCancel={() => setRoleModalOpen(false)}
        confirmLoading={assignRolesMutation.isPending}
        width={600}
      >
        <Transfer
          dataSource={transferDataSource}
          titles={['可选角色', '已分配角色']}
          targetKeys={targetKeys}
          onChange={(keys) => setTargetKeys(keys as number[])}
          render={item => item.title}
          listStyle={{ width: 250, height: 300 }}
          showSearch
          filterOption={(input, item) => 
            item.title.toLowerCase().includes(input.toLowerCase()) ||
            item.description?.toLowerCase().includes(input.toLowerCase())
          }
        />
      </Modal>
    </>
  )
}
