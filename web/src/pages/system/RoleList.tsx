import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, message, Popconfirm, Tag, Transfer } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, SafetyOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { roleApi } from '../../api/role'
import { permissionApi } from '../../api/permission'
import type { Role, CreateRoleRequest, UpdateRoleRequest } from '../../types/system'

export function RoleList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [permModalOpen, setPermModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null)
  const [targetKeys, setTargetKeys] = useState<number[]>([])
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['roles', page, pageSize],
    queryFn: () => roleApi.list(page, pageSize)
  })

  const { data: permissionsData } = useQuery({
    queryKey: ['permissions'],
    queryFn: () => permissionApi.list(1, 100)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateRoleRequest) => roleApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateRoleRequest }) => roleApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => roleApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const assignPermsMutation = useMutation({
    mutationFn: ({ roleId, permIds }: { roleId: number; permIds: number[] }) => roleApi.assignPermissions(roleId, permIds),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rolePermissions'] })
      messageApi.success('分配成功')
      setPermModalOpen(false)
    },
    onError: () => messageApi.error('分配失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Role) => {
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
      createMutation.mutate(values as CreateRoleRequest)
    }
  }

  const handleOpenPermModal = async (roleId: number) => {
    try {
      setSelectedRoleId(roleId)
      const perms = await roleApi.getPermissions(roleId)
      setTargetKeys(perms ? perms.map(p => p.id) : [])
      setPermModalOpen(true)
    } catch (error) {
      console.error('Failed to get role permissions:', error)
      messageApi.error('获取权限失败')
    }
  }

  const handleAssignPermissions = () => {
    if (selectedRoleId) {
      assignPermsMutation.mutate({ roleId: selectedRoleId, permIds: targetKeys })
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '角色名称', dataIndex: 'name', width: 120 },
    { title: '角色编码', dataIndex: 'code', width: 120 },
    { title: '描述', dataIndex: 'description', ellipsis: true },
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
      width: 200,
      render: (_: unknown, record: Role) => (
        <Space>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Button type="link" size="small" icon={<SafetyOutlined />} onClick={() => handleOpenPermModal(record.id)}>
            权限
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

  const transferDataSource = permissionsData?.items?.map(p => ({
    key: p.id,
    title: p.name,
    description: `${p.resource}:${p.action}`
  })) || []

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增角色
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
        title={editingId ? '编辑角色' : '新增角色'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="角色名称" rules={[{ required: true, message: '请输入角色名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="code" label="角色编码" rules={[{ required: true, message: '请输入角色编码' }]}>
            <Input disabled={!!editingId} />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="分配权限"
        open={permModalOpen}
        onOk={handleAssignPermissions}
        onCancel={() => setPermModalOpen(false)}
        confirmLoading={assignPermsMutation.isPending}
        width={700}
      >
        <Transfer
          dataSource={transferDataSource}
          titles={['可分配权限', '已分配权限']}
          targetKeys={targetKeys}
          onChange={(keys) => setTargetKeys(keys as number[])}
          render={item => `${item.title} (${item.description})`}
          listStyle={{ width: 300, height: 400 }}
        />
      </Modal>
    </>
  )
}
