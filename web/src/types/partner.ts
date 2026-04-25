export interface Supplier {
  id: number
  name: string
  code: string
  contact: string
  phone: string
  email: string
  address: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateSupplierRequest {
  name: string
  code?: string
  contact?: string
  phone?: string
  email?: string
  address?: string
  status?: number
}

export interface UpdateSupplierRequest {
  name?: string
  code?: string
  contact?: string
  phone?: string
  email?: string
  address?: string
  status?: number
}

export interface Customer {
  id: number
  name: string
  code: string
  contact: string
  phone: string
  email: string
  address: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateCustomerRequest {
  name: string
  code?: string
  contact?: string
  phone?: string
  email?: string
  address?: string
  status?: number
}

export interface UpdateCustomerRequest {
  name?: string
  code?: string
  contact?: string
  phone?: string
  email?: string
  address?: string
  status?: number
}
