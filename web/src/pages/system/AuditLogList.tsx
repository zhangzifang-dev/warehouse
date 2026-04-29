import { useState } from 'react'
import { Table, Input, DatePicker, Space, Modal, Descriptions, Tag, theme } from 'antd'
import { useQuery } from '@tanstack/react-query'
import { auditLogApi, type AuditLogFilter } from '../../api/auditLog'
import type { AuditLog } from '../../types/system'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker

export function AuditLogList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [filter, setFilter] = useState<AuditLogFilter>({})
  const [detailOpen, setDetailOpen] = useState(false)
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null)
  const { token } = theme.useToken()

  const { data, isLoading } = useQuery({
    queryKey: ['auditLogs', page, pageSize, filter],
    queryFn: () => auditLogApi.list(page, pageSize, filter)
  })

  const handleViewDetail = (record: AuditLog) => {
    setSelectedLog(record)
    setDetailOpen(true)
  }

  const handleTableChange = (tableName: string) => {
    setFilter(prev => ({ ...prev, table_name: tableName || undefined }))
    setPage(1)
  }

  const handleRecordIdChange = (value: string) => {
    const recordId = value ? parseInt(value) : undefined
    setFilter(prev => ({ ...prev, record_id: recordId }))
    setPage(1)
  }

  const handleDateChange = (dates: [dayjs.Dayjs | null, dayjs.Dayjs | null] | null) => {
    if (dates && dates[0] && dates[1]) {
      setFilter(prev => ({
        ...prev,
        start_time: dates[0]!.toISOString(),
        end_time: dates[1]!.toISOString()
      }))
    } else {
      setFilter(prev => {
        const { start_time, end_time, ...rest } = prev
        return rest
      })
    }
    setPage(1)
  }

  const actionColors: Record<string, string> = {
    create: 'green',
    update: 'blue',
    delete: 'red'
  }

  const actionLabels: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除'
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '表名', dataIndex: 'table_name', width: 120 },
    { title: '记录ID', dataIndex: 'record_id', width: 100 },
    {
      title: '操作',
      dataIndex: 'action',
      width: 80,
      render: (action: string) => (
        <Tag color={actionColors[action] || 'default'}>
          {actionLabels[action] || action}
        </Tag>
      )
    },
    { title: '操作人ID', dataIndex: 'operated_by', width: 100 },
    {
      title: '操作时间',
      dataIndex: 'operated_at',
      width: 180,
      render: (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm:ss')
    },
    { title: 'IP地址', dataIndex: 'ip_address', width: 140 },
    {
      title: '操作',
      width: 80,
      render: (_: unknown, record: AuditLog) => (
        <a onClick={() => handleViewDetail(record)}>详情</a>
      )
    }
  ]

  const parseValue = (value: Record<string, unknown> | null | undefined): Record<string, unknown> => {
    if (!value || Object.keys(value).length === 0) {
      return {}
    }
    try {
      const dataStr = value.data as string
      if (dataStr) {
        return JSON.parse(dataStr)
      }
    } catch {
      return {}
    }
    return value
  }

  const renderValueComparison = (oldValue: Record<string, unknown> | null | undefined, newValue: Record<string, unknown> | null | undefined) => {
    const oldData = parseValue(oldValue)
    const newData = parseValue(newValue)
    
    const allKeys = [...new Set([...Object.keys(oldData), ...Object.keys(newData)])]
    
    if (allKeys.length === 0) {
      return <span style={{ color: token.colorTextSecondary }}>无数据</span>
    }

    const dataSource = allKeys.map((key, index) => {
      const oldVal = oldData[key]
      const newVal = newData[key]
      const isDifferent = JSON.stringify(oldVal) !== JSON.stringify(newVal)
      
      return {
        key: index + 1,
        field: key,
        oldValue: oldVal !== undefined ? String(oldVal) : '-',
        newValue: newVal !== undefined ? String(newVal) : '-',
        isDifferent
      }
    })

    return (
      <Table 
        dataSource={dataSource}
        pagination={false}
        size="small"
        rowKey="key"
      >
        <Table.Column title="No" dataIndex="key" width={50} />
        <Table.Column title="字段" dataIndex="field" width={120} />
        <Table.Column 
          title="旧值" 
          dataIndex="oldValue" 
          ellipsis
          render={(val: string) => <span style={{ color: token.colorTextSecondary }}>{val}</span>}
        />
        <Table.Column 
          title="新值" 
          dataIndex="newValue" 
          ellipsis
          render={(val: string, record: { isDifferent: boolean }) => (
            <span style={{ fontWeight: record.isDifferent ? 600 : 400, color: token.colorText }}>
              {val}
            </span>
          )}
        />
      </Table>
    )
  }

  return (
    <>
      <div style={{ marginBottom: 16 }}>
        <Space wrap>
          <Input
            placeholder="表名"
            style={{ width: 150 }}
            value={filter.table_name || ''}
            onChange={e => handleTableChange(e.target.value)}
            allowClear
          />
          <Input
            placeholder="记录ID"
            style={{ width: 120 }}
            value={filter.record_id || ''}
            onChange={e => handleRecordIdChange(e.target.value)}
            allowClear
          />
          <RangePicker
            showTime
            onChange={(dates) => handleDateChange(dates as [dayjs.Dayjs | null, dayjs.Dayjs | null] | null)}
          />
        </Space>
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
        title="审计日志详情"
        open={detailOpen}
        onCancel={() => setDetailOpen(false)}
        footer={null}
        width={700}
      >
        {selectedLog && (
          <>
            <Descriptions column={2} size="small" bordered labelStyle={{ width: 100 }} style={{ marginBottom: 16 }}>
              <Descriptions.Item label="ID">{selectedLog.id}</Descriptions.Item>
              <Descriptions.Item label="表名">{selectedLog.table_name}</Descriptions.Item>
              <Descriptions.Item label="记录ID">{selectedLog.record_id}</Descriptions.Item>
              <Descriptions.Item label="操作">
                <Tag color={actionColors[selectedLog.action] || 'default'}>
                  {actionLabels[selectedLog.action] || selectedLog.action}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="操作人ID">{selectedLog.operated_by}</Descriptions.Item>
              <Descriptions.Item label="IP地址">{selectedLog.ip_address || '-'}</Descriptions.Item>
              <Descriptions.Item label="操作时间" span={2}>
                {dayjs(selectedLog.operated_at).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
            </Descriptions>
            {renderValueComparison(selectedLog.old_value, selectedLog.new_value)}
          </>
        )}
      </Modal>
    </>
  )
}
