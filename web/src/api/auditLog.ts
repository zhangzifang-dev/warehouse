import api from './client'
import type { AuditLog } from '../types/system'
import type { PaginatedResponse } from '../types/warehouse'

export interface AuditLogFilter {
  table_name?: string
  record_id?: number
  operated_by?: number
  operated_by_name?: string
  start_time?: string
  end_time?: string
}

export const auditLogApi = {
  list: async (page = 1, size = 10, filter?: AuditLogFilter): Promise<PaginatedResponse<AuditLog>> => {
    const response = await api.get<PaginatedResponse<AuditLog>>('/audit-logs', {
      params: { page, size, ...filter }
    })
    return response.data
  },

  get: async (id: number): Promise<AuditLog> => {
    const response = await api.get<AuditLog>(`/audit-logs/${id}`)
    return response.data
  }
}
