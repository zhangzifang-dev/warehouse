export interface Category {
  id: number
  name: string
  parent_id: number | null
  sort_order: number
  status: number
  created_at: string
  updated_at: string
  children?: Category[]
}

export interface CreateCategoryRequest {
  name: string
  parent_id?: number
  sort_order?: number
  status?: number
}

export interface UpdateCategoryRequest {
  name?: string
  parent_id?: number
  sort_order?: number
  status?: number
}

export interface Product {
  id: number
  sku: string
  name: string
  category_id: number | null
  specification: string
  unit: string
  barcode: string
  price: number
  description: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateProductRequest {
  sku: string
  name: string
  category_id?: number
  specification?: string
  unit?: string
  barcode?: string
  price?: number
  description?: string
  status?: number
}

export interface UpdateProductRequest {
  sku?: string
  name?: string
  category_id?: number
  specification?: string
  unit?: string
  barcode?: string
  price?: number
  description?: string
  status?: number
}
