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

  const renderJsonView = (value: Record<string, unknown> | null | undefined) => {
    if (!value || Object.keys(value).length === 0) {
      return <span style={{ color: token.colorTextSecondary }}>无</span>
    }
    
    try {
      const dataStr = value.data as string
      if (dataStr) {
        const parsed = JSON.parse(dataStr)
        return (
          <div style={{ 
            background: token.colorBgContainer, 
            border: `1px solid ${token.colorBorder}`,
            padding: 8, 
            borderRadius: 4, 
            maxHeight: 200, 
            overflow: 'auto',
            fontSize: 12,
            wordBreak: 'break-all'
          }}>
            {Object.entries(parsed).map(([key, val]) => (
              <div key={key} style={{ marginBottom: 4 }}>
                <span style={{ fontWeight: 500, color: token.colorPrimary }}>{key}: </span>
                <span style={{ color: token.colorText }}>{String(val)}</span>
              </div>
            ))}
          </div>
        )
      }
    } catch {
      return <span style={{ color: token.colorTextSecondary }}>解析失败</span>
    }
    
    return (
      <div style={{ 
        background: token.colorBgContainer, 
        border: `1px solid ${token.colorBorder}`,
        padding: 8, 
        borderRadius: 4, 
        maxHeight: 200, 
        overflow: 'auto',
        fontSize: 12,
        wordBreak: 'break-all',
        color: token.colorText
      }}>
        {JSON.stringify(value)}
      </div>
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
        width={600}
      >
        {selectedLog && (
          <Descriptions column={2} size="small" bordered labelStyle={{ width: 100 }}>
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
            <Descriptions.Item label="旧值" span={2}>
              {renderJsonView(selectedLog.old_value)}
            </Descriptions.Item>
            <Descriptions.Item label="新值" span={2}>
              {renderJsonView(selectedLog.new_value)}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </>
  )
}
