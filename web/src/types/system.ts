export interface User {
  id: number
  username: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  username: string
  password: string
  status?: number
}

export interface UpdateUserRequest {
  status?: number
}

export interface Role {
  id: number
  name: string
  code: string
  description: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateRoleRequest {
  name: string
  code: string
  description?: string
  status?: number
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
  status?: number
}

export interface Permission {
  id: number
  name: string
  code: string
  resource: string
  action: string
  description: string
  created_at: string
  updated_at: string
}

export interface AuditLog {
  id: number
  table_name: string
  record_id: number
  action: string
  old_value: Record<string, unknown>
  new_value: Record<string, unknown>
  operated_by: number
  operated_at: string
  ip_address: string
}
