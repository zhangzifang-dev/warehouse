import { useState } from 'react'
import { Table } from 'antd'
import { useQuery } from '@tanstack/react-query'
import { permissionApi } from '../../api/permission'

export function PermissionList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  const { data, isLoading } = useQuery({
    queryKey: ['permissions', page, pageSize],
    queryFn: () => permissionApi.list(page, pageSize)
  })

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '权限名称', dataIndex: 'name', width: 150 },
    { title: '权限编码', dataIndex: 'code', width: 200 },
    { title: '资源', dataIndex: 'resource', width: 150 },
    { title: '操作', dataIndex: 'action', width: 100 },
    { title: '描述', dataIndex: 'description', ellipsis: true }
  ]

  return (
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
  )
}
