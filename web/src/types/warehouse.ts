export interface Warehouse {
  id: number
  name: string
  code: string
  address: string
  contact: string
  phone: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateWarehouseRequest {
  name: string
  code: string
  address?: string
  contact?: string
  phone?: string
  status?: number
}

export interface UpdateWarehouseRequest {
  name?: string
  address?: string
  contact?: string
  phone?: string
  status?: number
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  size: number
}

export interface Location {
  id: number
  warehouse_id: number
  zone: string
  shelf: string
  level: string
  position: string
  code: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateLocationRequest {
  warehouse_id: number
  zone: string
  shelf: string
  level: string
  position: string
  status?: number
}

export interface UpdateLocationRequest {
  warehouse_id?: number
  zone?: string
  shelf?: string
  level?: string
  position?: string
  status?: number
}
