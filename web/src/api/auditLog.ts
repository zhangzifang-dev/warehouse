import api from './client'
import type { AuditLog } from '../types/system'
import type { PaginatedResponse } from '../types/warehouse'

export interface AuditLogFilter {
  table_name?: string[]
  record_id?: number
  operated_by?: number
  operated_by_name?: string[]
  action?: string[]
  start_time?: string
  end_time?: string
}

export const auditLogApi = {
  list: async (page = 1, size = 10, filter?: AuditLogFilter): Promise<PaginatedResponse<AuditLog>> => {
    const params = new URLSearchParams()
    params.append('page', String(page))
    params.append('size', String(size))
    
    if (filter) {
      if (filter.table_name) {
        filter.table_name.forEach(v => params.append('table_name', v))
      }
      if (filter.record_id) {
        params.append('record_id', String(filter.record_id))
      }
      if (filter.operated_by) {
        params.append('operated_by', String(filter.operated_by))
      }
      if (filter.operated_by_name) {
        filter.operated_by_name.forEach(v => params.append('operated_by_name', v))
      }
      if (filter.action) {
        filter.action.forEach(v => params.append('action', v))
      }
      if (filter.start_time) {
        params.append('start_time', filter.start_time)
      }
      if (filter.end_time) {
        params.append('end_time', filter.end_time)
      }
    }
    
    const response = await api.get<PaginatedResponse<AuditLog>>('/audit-logs', {
      params: params
    })
    return response.data
  },

  get: async (id: number): Promise<AuditLog> => {
    const response = await api.get<AuditLog>(`/audit-logs/${id}`)
    return response.data
  },

  getTableNames: async (): Promise<string[]> => {
    const response = await api.get<string[]>('/audit-logs/table-names')
    return response.data
  }
}
