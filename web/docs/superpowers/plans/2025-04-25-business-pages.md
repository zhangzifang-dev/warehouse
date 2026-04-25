# Business Pages Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement 10 business page modules with CRUD operations, tables, modals, and drawers for a warehouse management system.

**Architecture:** React frontend using Ant Design components, TanStack Query for data fetching, React Router for navigation. Each module follows a consistent pattern: TypeScript types, API functions, list page with table, create/edit modals, and detail drawers for orders.

**Tech Stack:** React 19, Ant Design 6, TanStack Query 5, React Router 7, Axios, Zustand

---

## File Structure

```
src/
├── types/
│   ├── warehouse.ts      # Warehouse, Location interfaces
│   ├── product.ts        # Category, Product interfaces
│   ├── inventory.ts      # Inventory interfaces
│   ├── partner.ts        # Supplier, Customer interfaces
│   └── order.ts          # InboundOrder, OutboundOrder, StockTransfer interfaces
├── api/
│   ├── warehouse.ts      # Warehouse API
│   ├── location.ts       # Location API
│   ├── category.ts       # Category API
│   ├── product.ts        # Product API
│   ├── inventory.ts      # Inventory API
│   ├── supplier.ts       # Supplier API
│   ├── customer.ts       # Customer API
│   ├── inbound.ts        # Inbound Order API
│   ├── outbound.ts       # Outbound Order API
│   └── transfer.ts       # Stock Transfer API
├── pages/
│   ├── warehouse/
│   │   ├── WarehouseList.tsx
│   │   └── LocationList.tsx
│   ├── product/
│   │   ├── CategoryList.tsx
│   │   └── ProductList.tsx
│   ├── inventory/
│   │   └── InventoryList.tsx
│   ├── partner/
│   │   ├── SupplierList.tsx
│   │   └── CustomerList.tsx
│   └── order/
│       ├── InboundOrderList.tsx
│       ├── OutboundOrderList.tsx
│       └── StockTransferList.tsx
└── App.tsx               # Update routes
```

---

### Task 1: Create TypeScript Types

**Files:**
- Create: `src/types/warehouse.ts`
- Create: `src/types/product.ts`
- Create: `src/types/inventory.ts`
- Create: `src/types/partner.ts`
- Create: `src/types/order.ts`

- [ ] **Step 1: Create warehouse types**

```typescript
// src/types/warehouse.ts

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
```

- [ ] **Step 2: Create product types**

```typescript
// src/types/product.ts

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
```

- [ ] **Step 3: Create inventory types**

```typescript
// src/types/inventory.ts

export interface Inventory {
  id: number
  warehouse_id: number
  product_id: number
  location_id: number | null
  quantity: number
  batch_no: string
  created_at: string
  updated_at: string
}

export interface CreateInventoryRequest {
  warehouse_id: number
  product_id: number
  location_id?: number
  quantity?: number
  batch_no?: string
}

export interface UpdateInventoryRequest {
  warehouse_id?: number
  product_id?: number
  location_id?: number
  quantity?: number
  batch_no?: string
}

export interface AdjustQuantityRequest {
  inventory_id: number
  quantity: number
}
```

- [ ] **Step 4: Create partner types**

```typescript
// src/types/partner.ts

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
```

- [ ] **Step 5: Create order types**

```typescript
// src/types/order.ts

export interface InboundOrder {
  id: number
  order_no: string
  supplier_id: number | null
  warehouse_id: number
  total_quantity: number
  status: number
  remark: string
  created_at: string
  updated_at: string
  items?: InboundItem[]
}

export interface InboundItem {
  id: number
  order_id: number
  product_id: number
  location_id: number | null
  quantity: number
  batch_no: string
  created_at: string
  updated_at: string
}

export interface CreateInboundOrderRequest {
  supplier_id?: number
  warehouse_id: number
  remark?: string
  items: CreateInboundItemRequest[]
}

export interface CreateInboundItemRequest {
  product_id: number
  location_id?: number
  quantity: number
  batch_no?: string
}

export interface OutboundOrder {
  id: number
  order_no: string
  customer_id: number | null
  warehouse_id: number
  total_quantity: number
  status: number
  remark: string
  created_at: string
  updated_at: string
  items?: OutboundItem[]
}

export interface OutboundItem {
  id: number
  order_id: number
  product_id: number
  location_id: number | null
  quantity: number
  batch_no: string
  created_at: string
  updated_at: string
}

export interface CreateOutboundOrderRequest {
  customer_id?: number
  warehouse_id: number
  remark?: string
  items: CreateOutboundItemRequest[]
}

export interface CreateOutboundItemRequest {
  product_id: number
  location_id?: number
  quantity: number
  batch_no?: string
}

export interface StockTransfer {
  id: number
  order_no: string
  from_warehouse_id: number
  to_warehouse_id: number
  total_quantity: number
  status: number
  remark: string
  created_at: string
  updated_at: string
  items?: StockTransferItem[]
}

export interface StockTransferItem {
  id: number
  transfer_id: number
  product_id: number
  location_id: number | null
  quantity: number
  batch_no: string
  created_at: string
  updated_at: string
}

export interface CreateStockTransferRequest {
  from_warehouse_id: number
  to_warehouse_id: number
  remark?: string
  items: CreateStockTransferItemRequest[]
}

export interface CreateStockTransferItemRequest {
  product_id: number
  location_id?: number
  quantity: number
  batch_no?: string
}
```

- [ ] **Step 6: Commit types**

```bash
git add src/types/warehouse.ts src/types/product.ts src/types/inventory.ts src/types/partner.ts src/types/order.ts
git commit -m "feat: add TypeScript types for business modules"
```

---

### Task 2: Create API Functions

**Files:**
- Create: `src/api/warehouse.ts`
- Create: `src/api/location.ts`
- Create: `src/api/category.ts`
- Create: `src/api/product.ts`
- Create: `src/api/inventory.ts`
- Create: `src/api/supplier.ts`
- Create: `src/api/customer.ts`
- Create: `src/api/inbound.ts`
- Create: `src/api/outbound.ts`
- Create: `src/api/transfer.ts`

- [ ] **Step 1: Create warehouse API**

```typescript
// src/api/warehouse.ts

import api from './client'
import type { Warehouse, CreateWarehouseRequest, UpdateWarehouseRequest, PaginatedResponse } from '../types/warehouse'

export const warehouseApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<Warehouse>> => {
    const response = await api.get<PaginatedResponse<Warehouse>>('/warehouses', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<Warehouse> => {
    const response = await api.get<Warehouse>(`/warehouses/${id}`)
    return response.data
  },

  create: async (data: CreateWarehouseRequest): Promise<Warehouse> => {
    const response = await api.post<Warehouse>('/warehouses', data)
    return response.data
  },

  update: async (id: number, data: UpdateWarehouseRequest): Promise<Warehouse> => {
    const response = await api.put<Warehouse>(`/warehouses/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/warehouses/${id}`)
  }
}
```

- [ ] **Step 2: Create location API**

```typescript
// src/api/location.ts

import api from './client'
import type { Location, CreateLocationRequest, UpdateLocationRequest, PaginatedResponse } from '../types/warehouse'

export const locationApi = {
  list: async (page = 1, size = 10, warehouseId?: number): Promise<PaginatedResponse<Location>> => {
    const response = await api.get<PaginatedResponse<Location>>('/locations', {
      params: { page, size, warehouse_id: warehouseId }
    })
    return response.data
  },

  get: async (id: number): Promise<Location> => {
    const response = await api.get<Location>(`/locations/${id}`)
    return response.data
  },

  create: async (data: CreateLocationRequest): Promise<Location> => {
    const response = await api.post<Location>('/locations', data)
    return response.data
  },

  update: async (id: number, data: UpdateLocationRequest): Promise<Location> => {
    const response = await api.put<Location>(`/locations/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/locations/${id}`)
  }
}
```

- [ ] **Step 3: Create category API**

```typescript
// src/api/category.ts

import api from './client'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest, PaginatedResponse } from '../types/product'

export const categoryApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<Category>> => {
    const response = await api.get<PaginatedResponse<Category>>('/categories', {
      params: { page, size }
    })
    return response.data
  },

  tree: async (): Promise<Category[]> => {
    const response = await api.get<Category[]>('/categories/tree')
    return response.data
  },

  get: async (id: number): Promise<Category> => {
    const response = await api.get<Category>(`/categories/${id}`)
    return response.data
  },

  create: async (data: CreateCategoryRequest): Promise<Category> => {
    const response = await api.post<Category>('/categories', data)
    return response.data
  },

  update: async (id: number, data: UpdateCategoryRequest): Promise<Category> => {
    const response = await api.put<Category>(`/categories/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/categories/${id}`)
  }
}
```

- [ ] **Step 4: Create product API**

```typescript
// src/api/product.ts

import api from './client'
import type { Product, CreateProductRequest, UpdateProductRequest, PaginatedResponse } from '../types/product'

export const productApi = {
  list: async (page = 1, size = 10, categoryId?: number, keyword?: string): Promise<PaginatedResponse<Product>> => {
    const response = await api.get<PaginatedResponse<Product>>('/products', {
      params: { page, size, category_id: categoryId, keyword }
    })
    return response.data
  },

  get: async (id: number): Promise<Product> => {
    const response = await api.get<Product>(`/products/${id}`)
    return response.data
  },

  create: async (data: CreateProductRequest): Promise<Product> => {
    const response = await api.post<Product>('/products', data)
    return response.data
  },

  update: async (id: number, data: UpdateProductRequest): Promise<Product> => {
    const response = await api.put<Product>(`/products/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/products/${id}`)
  }
}
```

- [ ] **Step 5: Create inventory API**

```typescript
// src/api/inventory.ts

import api from './client'
import type { Inventory, CreateInventoryRequest, UpdateInventoryRequest, AdjustQuantityRequest, PaginatedResponse } from '../types/inventory'

export const inventoryApi = {
  list: async (page = 1, size = 10, warehouseId?: number, productId?: number): Promise<PaginatedResponse<Inventory>> => {
    const response = await api.get<PaginatedResponse<Inventory>>('/inventory', {
      params: { page, size, warehouse_id: warehouseId, product_id: productId }
    })
    return response.data
  },

  get: async (id: number): Promise<Inventory> => {
    const response = await api.get<Inventory>(`/inventory/${id}`)
    return response.data
  },

  create: async (data: CreateInventoryRequest): Promise<Inventory> => {
    const response = await api.post<Inventory>('/inventory', data)
    return response.data
  },

  update: async (id: number, data: UpdateInventoryRequest): Promise<Inventory> => {
    const response = await api.put<Inventory>(`/inventory/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/inventory/${id}`)
  },

  adjust: async (data: AdjustQuantityRequest): Promise<Inventory> => {
    const response = await api.post<Inventory>('/inventory/adjust', data)
    return response.data
  }
}
```

- [ ] **Step 6: Create supplier API**

```typescript
// src/api/supplier.ts

import api from './client'
import type { Supplier, CreateSupplierRequest, UpdateSupplierRequest, PaginatedResponse } from '../types/partner'

export const supplierApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<Supplier>> => {
    const response = await api.get<PaginatedResponse<Supplier>>('/suppliers', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<Supplier> => {
    const response = await api.get<Supplier>(`/suppliers/${id}`)
    return response.data
  },

  create: async (data: CreateSupplierRequest): Promise<Supplier> => {
    const response = await api.post<Supplier>('/suppliers', data)
    return response.data
  },

  update: async (id: number, data: UpdateSupplierRequest): Promise<Supplier> => {
    const response = await api.put<Supplier>(`/suppliers/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/suppliers/${id}`)
  }
}
```

- [ ] **Step 7: Create customer API**

```typescript
// src/api/customer.ts

import api from './client'
import type { Customer, CreateCustomerRequest, UpdateCustomerRequest, PaginatedResponse } from '../types/partner'

export const customerApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<Customer>> => {
    const response = await api.get<PaginatedResponse<Customer>>('/customers', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<Customer> => {
    const response = await api.get<Customer>(`/customers/${id}`)
    return response.data
  },

  create: async (data: CreateCustomerRequest): Promise<Customer> => {
    const response = await api.post<Customer>('/customers', data)
    return response.data
  },

  update: async (id: number, data: UpdateCustomerRequest): Promise<Customer> => {
    const response = await api.put<Customer>(`/customers/${id}`, data)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/customers/${id}`)
  }
}
```

- [ ] **Step 8: Create inbound order API**

```typescript
// src/api/inbound.ts

import api from './client'
import type { InboundOrder, CreateInboundOrderRequest, PaginatedResponse } from '../types/order'

export const inboundApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<InboundOrder>> => {
    const response = await api.get<PaginatedResponse<InboundOrder>>('/inbound-orders', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<InboundOrder> => {
    const response = await api.get<InboundOrder>(`/inbound-orders/${id}`)
    return response.data
  },

  create: async (data: CreateInboundOrderRequest): Promise<InboundOrder> => {
    const response = await api.post<InboundOrder>('/inbound-orders', data)
    return response.data
  },

  confirm: async (id: number): Promise<InboundOrder> => {
    const response = await api.post<InboundOrder>(`/inbound-orders/${id}/confirm`)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/inbound-orders/${id}`)
  }
}
```

- [ ] **Step 9: Create outbound order API**

```typescript
// src/api/outbound.ts

import api from './client'
import type { OutboundOrder, CreateOutboundOrderRequest, PaginatedResponse } from '../types/order'

export const outboundApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<OutboundOrder>> => {
    const response = await api.get<PaginatedResponse<OutboundOrder>>('/outbound-orders', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<OutboundOrder> => {
    const response = await api.get<OutboundOrder>(`/outbound-orders/${id}`)
    return response.data
  },

  create: async (data: CreateOutboundOrderRequest): Promise<OutboundOrder> => {
    const response = await api.post<OutboundOrder>('/outbound-orders', data)
    return response.data
  },

  confirm: async (id: number): Promise<OutboundOrder> => {
    const response = await api.post<OutboundOrder>(`/outbound-orders/${id}/confirm`)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/outbound-orders/${id}`)
  }
}
```

- [ ] **Step 10: Create stock transfer API**

```typescript
// src/api/transfer.ts

import api from './client'
import type { StockTransfer, CreateStockTransferRequest, PaginatedResponse } from '../types/order'

export const transferApi = {
  list: async (page = 1, size = 10): Promise<PaginatedResponse<StockTransfer>> => {
    const response = await api.get<PaginatedResponse<StockTransfer>>('/stock-transfers', {
      params: { page, size }
    })
    return response.data
  },

  get: async (id: number): Promise<StockTransfer> => {
    const response = await api.get<StockTransfer>(`/stock-transfers/${id}`)
    return response.data
  },

  create: async (data: CreateStockTransferRequest): Promise<StockTransfer> => {
    const response = await api.post<StockTransfer>('/stock-transfers', data)
    return response.data
  },

  confirm: async (id: number): Promise<StockTransfer> => {
    const response = await api.post<StockTransfer>(`/stock-transfers/${id}/confirm`)
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/stock-transfers/${id}`)
  }
}
```

- [ ] **Step 11: Commit API functions**

```bash
git add src/api/
git commit -m "feat: add API functions for business modules"
```

---

### Task 3: Create WarehouseList Page

**Files:**
- Create: `src/pages/warehouse/WarehouseList.tsx`
- Create: `src/pages/warehouse/index.ts`

- [ ] **Step 1: Create WarehouseList component**

```typescript
// src/pages/warehouse/WarehouseList.tsx

import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, message, Popconfirm, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { warehouseApi } from '../../api/warehouse'
import type { Warehouse, CreateWarehouseRequest, UpdateWarehouseRequest } from '../../types/warehouse'

export function WarehouseList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['warehouses', page, pageSize],
    queryFn: () => warehouseApi.list(page, pageSize)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateWarehouseRequest) => warehouseApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['warehouses'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateWarehouseRequest }) => warehouseApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['warehouses'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => warehouseApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['warehouses'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Warehouse) => {
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
      createMutation.mutate(values as CreateWarehouseRequest)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '仓库编码', dataIndex: 'code', width: 120 },
    { title: '仓库名称', dataIndex: 'name' },
    { title: '地址', dataIndex: 'address', ellipsis: true },
    { title: '联系人', dataIndex: 'contact', width: 100 },
    { title: '联系电话', dataIndex: 'phone', width: 130 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      )
    },
    {
      title: '操作',
      width: 150,
      render: (_: unknown, record: Warehouse) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增仓库
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
      />
      <Modal
        title={editingId ? '编辑仓库' : '新增仓库'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="仓库名称" rules={[{ required: true, message: '请输入仓库名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="code" label="仓库编码" rules={[{ required: true, message: '请输入仓库编码' }]}>
            <Input disabled={!!editingId} />
          </Form.Item>
          <Form.Item name="address" label="地址">
            <Input />
          </Form.Item>
          <Form.Item name="contact" label="联系人">
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="联系电话">
            <Input />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
```

- [ ] **Step 2: Create index file**

```typescript
// src/pages/warehouse/index.ts

export { WarehouseList } from './WarehouseList'
export { LocationList } from './LocationList'
```

- [ ] **Step 3: Commit warehouse list**

```bash
git add src/pages/warehouse/
git commit -m "feat: add WarehouseList page with CRUD"
```

---

### Task 4: Create LocationList Page

**Files:**
- Create: `src/pages/warehouse/LocationList.tsx`

- [ ] **Step 1: Create LocationList component**

```typescript
// src/pages/warehouse/LocationList.tsx

import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, message, Popconfirm, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { locationApi } from '../../api/location'
import { warehouseApi } from '../../api/warehouse'
import type { Location, CreateLocationRequest, UpdateLocationRequest } from '../../types/warehouse'

export function LocationList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [warehouseFilter, setWarehouseFilter] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['locations', page, pageSize, warehouseFilter],
    queryFn: () => locationApi.list(page, pageSize, warehouseFilter)
  })

  const { data: warehouses } = useQuery({
    queryKey: ['warehouses-all'],
    queryFn: () => warehouseApi.list(1, 100)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateLocationRequest) => locationApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['locations'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateLocationRequest }) => locationApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['locations'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => locationApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['locations'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Location) => {
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
      createMutation.mutate(values as CreateLocationRequest)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '库位编码', dataIndex: 'code', width: 150 },
    {
      title: '所属仓库',
      dataIndex: 'warehouse_id',
      width: 150,
      render: (id: number) => warehouses?.items.find(w => w.id === id)?.name || id
    },
    { title: '区域', dataIndex: 'zone', width: 80 },
    { title: '货架', dataIndex: 'shelf', width: 80 },
    { title: '层', dataIndex: 'level', width: 80 },
    { title: '位', dataIndex: 'position', width: 80 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      )
    },
    {
      title: '操作',
      width: 150,
      render: (_: unknown, record: Location) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <Select
          allowClear
          placeholder="筛选仓库"
          style={{ width: 200 }}
          value={warehouseFilter}
          onChange={setWarehouseFilter}
          options={warehouses?.items.map(w => ({ value: w.id, label: w.name }))}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增库位
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
      />
      <Modal
        title={editingId ? '编辑库位' : '新增库位'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="warehouse_id" label="所属仓库" rules={[{ required: true, message: '请选择仓库' }]}>
            <Select
              options={warehouses?.items.map(w => ({ value: w.id, label: w.name }))}
            />
          </Form.Item>
          <Form.Item name="zone" label="区域" rules={[{ required: true, message: '请输入区域' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="shelf" label="货架" rules={[{ required: true, message: '请输入货架' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="level" label="层" rules={[{ required: true, message: '请输入层' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="position" label="位" rules={[{ required: true, message: '请输入位' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
```

- [ ] **Step 2: Update index file**

```typescript
// src/pages/warehouse/index.ts

export { WarehouseList } from './WarehouseList'
export { LocationList } from './LocationList'
```

- [ ] **Step 3: Commit location list**

```bash
git add src/pages/warehouse/
git commit -m "feat: add LocationList page with CRUD and warehouse filter"
```

---

### Task 5: Create CategoryList Page

**Files:**
- Create: `src/pages/product/CategoryList.tsx`
- Create: `src/pages/product/index.ts`

- [ ] **Step 1: Create CategoryList component**

```typescript
// src/pages/product/CategoryList.tsx

import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, Select, message, Popconfirm, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { categoryApi } from '../../api/category'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '../../types/product'

export function CategoryList() {
  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list(1, 100)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateCategoryRequest) => categoryApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateCategoryRequest }) => categoryApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => categoryApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Category) => {
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
      createMutation.mutate(values as CreateCategoryRequest)
    }
  }

  const buildTree = (items: Category[]): Category[] => {
    const map = new Map<number, Category>()
    const roots: Category[] = []

    items.forEach(item => {
      map.set(item.id, { ...item, children: [] })
    })

    items.forEach(item => {
      const node = map.get(item.id)!
      if (item.parent_id && map.has(item.parent_id)) {
        map.get(item.parent_id)!.children!.push(node)
      } else {
        roots.push(node)
      }
    })

    return roots
  }

  const treeData = data?.items ? buildTree(data.items) : []

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '分类名称', dataIndex: 'name' },
    { title: '排序', dataIndex: 'sort_order', width: 80 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      )
    },
    {
      title: '操作',
      width: 150,
      render: (_: unknown, record: Category) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增分类
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={treeData}
        rowKey="id"
        loading={isLoading}
        pagination={false}
        defaultExpandAllRows
      />
      <Modal
        title={editingId ? '编辑分类' : '新增分类'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="分类名称" rules={[{ required: true, message: '请输入分类名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="parent_id" label="父级分类">
            <Select
              allowClear
              placeholder="选择父级分类（可选）"
              options={data?.items
                .filter(c => c.id !== editingId)
                .map(c => ({ value: c.id, label: c.name }))}
            />
          </Form.Item>
          <Form.Item name="sort_order" label="排序" initialValue={0}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
```

- [ ] **Step 2: Create product index file**

```typescript
// src/pages/product/index.ts

export { CategoryList } from './CategoryList'
export { ProductList } from './ProductList'
```

- [ ] **Step 3: Commit category list**

```bash
git add src/pages/product/
git commit -m "feat: add CategoryList page with tree table"
```

---

### Task 6: Create ProductList Page

**Files:**
- Create: `src/pages/product/ProductList.tsx`

- [ ] **Step 1: Create ProductList component**

```typescript
// src/pages/product/ProductList.tsx

import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, Select, message, Popconfirm, Tag, TreeSelect } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { productApi } from '../../api/product'
import { categoryApi } from '../../api/category'
import type { Product, CreateProductRequest, UpdateProductRequest, Category } from '../../types/product'

export function ProductList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [categoryFilter, setCategoryFilter] = useState<number | undefined>()
  const [keyword, setKeyword] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['products', page, pageSize, categoryFilter, keyword],
    queryFn: () => productApi.list(page, pageSize, categoryFilter, keyword || undefined)
  })

  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list(1, 100)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateProductRequest) => productApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateProductRequest }) => productApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => productApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Product) => {
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
      createMutation.mutate(values as CreateProductRequest)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  const buildTreeData = (items: Category[]): { value: number; title: string; children?: { value: number; title: string }[] }[] => {
    const map = new Map<number, Category & { children: Category[] }>()
    const roots: Category[] = []

    items.forEach(item => {
      map.set(item.id, { ...item, children: [] })
    })

    items.forEach(item => {
      const node = map.get(item.id)!
      if (item.parent_id && map.has(item.parent_id)) {
        map.get(item.parent_id)!.children.push(node)
      } else {
        roots.push(node)
      }
    })

    return roots.map(r => ({
      value: r.id,
      title: r.name,
      children: r.children?.map(c => ({ value: c.id, title: c.name }))
    }))
  }

  const categoryTreeData = categories?.items ? buildTreeData(categories.items) : []

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: 'SKU', dataIndex: 'sku', width: 120 },
    { title: '商品名称', dataIndex: 'name' },
    { title: '规格', dataIndex: 'specification', width: 120, ellipsis: true },
    { title: '单位', dataIndex: 'unit', width: 80 },
    { title: '价格', dataIndex: 'price', width: 100 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      )
    },
    {
      title: '操作',
      width: 150,
      render: (_: unknown, record: Product) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <TreeSelect
          allowClear
          placeholder="筛选分类"
          style={{ width: 200 }}
          value={categoryFilter}
          onChange={(v) => { setCategoryFilter(v); setPage(1) }}
          treeData={categoryTreeData}
        />
        <Input.Search
          placeholder="搜索商品名称/SKU"
          style={{ width: 250 }}
          onSearch={handleSearch}
          allowClear
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增商品
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
      />
      <Modal
        title={editingId ? '编辑商品' : '新增商品'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="sku" label="SKU" rules={[{ required: true, message: '请输入SKU' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="name" label="商品名称" rules={[{ required: true, message: '请输入商品名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="category_id" label="分类">
            <TreeSelect
              placeholder="选择分类"
              treeData={categoryTreeData}
              allowClear
            />
          </Form.Item>
          <Form.Item name="specification" label="规格">
            <Input />
          </Form.Item>
          <Form.Item name="unit" label="单位">
            <Input />
          </Form.Item>
          <Form.Item name="barcode" label="条码">
            <Input />
          </Form.Item>
          <Form.Item name="price" label="价格">
            <InputNumber min={0} precision={2} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
```

- [ ] **Step 2: Update product index**

```typescript
// src/pages/product/index.ts

export { CategoryList } from './CategoryList'
export { ProductList } from './ProductList'
```

- [ ] **Step 3: Commit product list**

```bash
git add src/pages/product/
git commit -m "feat: add ProductList page with search and category filter"
```

---

### Task 7: Create InventoryList Page

**Files:**
- Create: `src/pages/inventory/InventoryList.tsx`
- Create: `src/pages/inventory/index.ts`

- [ ] **Step 1: Create InventoryList component**

```typescript
// src/pages/inventory/InventoryList.tsx

import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, Select, message, Tag } from 'antd'
import { PlusOutlined, EditOutlined, AdjustOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { inventoryApi } from '../../api/inventory'
import { warehouseApi } from '../../api/warehouse'
import { productApi } from '../../api/product'
import type { Inventory, CreateInventoryRequest, UpdateInventoryRequest, AdjustQuantityRequest } from '../../types/inventory'

export function InventoryList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [warehouseFilter, setWarehouseFilter] = useState<number | undefined>()
  const [modalOpen, setModalOpen] = useState(false)
  const [adjustModalOpen, setAdjustModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [adjustingId, setAdjustingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [adjustForm] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['inventory', page, pageSize, warehouseFilter],
    queryFn: () => inventoryApi.list(page, pageSize, warehouseFilter)
  })

  const { data: warehouses } = useQuery({
    queryKey: ['warehouses-all'],
    queryFn: () => warehouseApi.list(1, 100)
  })

  const { data: products } = useQuery({
    queryKey: ['products-all'],
    queryFn: () => productApi.list(1, 100)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateInventoryRequest) => inventoryApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateInventoryRequest }) => inventoryApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const adjustMutation = useMutation({
    mutationFn: (data: AdjustQuantityRequest) => inventoryApi.adjust(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] })
      messageApi.success('调整成功')
      handleCloseAdjustModal()
    },
    onError: () => messageApi.error('调整失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Inventory) => {
    setEditingId(record.id)
    form.setFieldsValue(record)
    setModalOpen(true)
  }

  const handleCloseModal = () => {
    setModalOpen(false)
    setEditingId(null)
    form.resetFields()
  }

  const handleOpenAdjust = (record: Inventory) => {
    setAdjustingId(record.id)
    adjustForm.setFieldsValue({ inventory_id: record.id, quantity: 0 })
    setAdjustModalOpen(true)
  }

  const handleCloseAdjustModal = () => {
    setAdjustModalOpen(false)
    setAdjustingId(null)
    adjustForm.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: values })
    } else {
      createMutation.mutate(values as CreateInventoryRequest)
    }
  }

  const handleAdjustSubmit = async () => {
    const values = await adjustForm.validateFields()
    adjustMutation.mutate(values)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    {
      title: '仓库',
      dataIndex: 'warehouse_id',
      width: 150,
      render: (id: number) => warehouses?.items.find(w => w.id === id)?.name || id
    },
    {
      title: '商品',
      dataIndex: 'product_id',
      render: (id: number) => products?.items.find(p => p.id === id)?.name || id
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      width: 120,
      render: (qty: number) => (
        <Tag color={qty > 0 ? 'green' : qty < 0 ? 'red' : 'default'}>
          {qty}
        </Tag>
      )
    },
    { title: '批次号', dataIndex: 'batch_no', width: 120 },
    {
      title: '操作',
      width: 180,
      render: (_: unknown, record: Inventory) => (
        <Space>
          <Button type="link" icon={<AdjustOutlined />} onClick={() => handleOpenAdjust(record)}>
            调整
          </Button>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
        </Space>
      )
    }
  ]

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <Select
          allowClear
          placeholder="筛选仓库"
          style={{ width: 200 }}
          value={warehouseFilter}
          onChange={setWarehouseFilter}
          options={warehouses?.items.map(w => ({ value: w.id, label: w.name }))}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增库存
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
      />
      <Modal
        title={editingId ? '编辑库存' : '新增库存'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="warehouse_id" label="仓库" rules={[{ required: true, message: '请选择仓库' }]}>
            <Select
              options={warehouses?.items.map(w => ({ value: w.id, label: w.name }))}
            />
          </Form.Item>
          <Form.Item name="product_id" label="商品" rules={[{ required: true, message: '请选择商品' }]}>
            <Select
              showSearch
              optionFilterProp="label"
              options={products?.items.map(p => ({ value: p.id, label: `${p.sku} - ${p.name}` }))}
            />
          </Form.Item>
          <Form.Item name="quantity" label="数量">
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="batch_no" label="批次号">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="库存调整"
        open={adjustModalOpen}
        onOk={handleAdjustSubmit}
        onCancel={handleCloseAdjustModal}
        confirmLoading={adjustMutation.isPending}
      >
        <Form form={adjustForm} layout="vertical">
          <Form.Item name="inventory_id" label="库存ID" hidden>
            <InputNumber />
          </Form.Item>
          <Form.Item name="quantity" label="调整数量（正数入库，负数出库）" rules={[{ required: true, message: '请输入调整数量' }]}>
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
```

- [ ] **Step 2: Create inventory index**

```typescript
// src/pages/inventory/index.ts

export { InventoryList } from './InventoryList'
```

- [ ] **Step 3: Commit inventory list**

```bash
git add src/pages/inventory/
git commit -m "feat: add InventoryList page with stock adjustment"
```

---

### Task 8: Create SupplierList Page

**Files:**
- Create: `src/pages/partner/SupplierList.tsx`
- Create: `src/pages/partner/index.ts`

- [ ] **Step 1: Create SupplierList component**

```typescript
// src/pages/partner/SupplierList.tsx

import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, message, Popconfirm, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { supplierApi } from '../../api/supplier'
import type { Supplier, CreateSupplierRequest, UpdateSupplierRequest } from '../types/partner'

export function SupplierList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['suppliers', page, pageSize],
    queryFn: () => supplierApi.list(page, pageSize)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateSupplierRequest) => supplierApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['suppliers'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateSupplierRequest }) => supplierApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['suppliers'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => supplierApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['suppliers'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Supplier) => {
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
      createMutation.mutate(values as CreateSupplierRequest)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '供应商编码', dataIndex: 'code', width: 120 },
    { title: '供应商名称', dataIndex: 'name' },
    { title: '联系人', dataIndex: 'contact', width: 100 },
    { title: '联系电话', dataIndex: 'phone', width: 130 },
    { title: '邮箱', dataIndex: 'email', width: 180 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      )
    },
    {
      title: '操作',
      width: 150,
      render: (_: unknown, record: Supplier) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增供应商
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
      />
      <Modal
        title={editingId ? '编辑供应商' : '新增供应商'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="供应商名称" rules={[{ required: true, message: '请输入供应商名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="code" label="供应商编码">
            <Input />
          </Form.Item>
          <Form.Item name="contact" label="联系人">
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="联系电话">
            <Input />
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input />
          </Form.Item>
          <Form.Item name="address" label="地址">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
```

- [ ] **Step 2: Create partner index**

```typescript
// src/pages/partner/index.ts

export { SupplierList } from './SupplierList'
export { CustomerList } from './CustomerList'
```

- [ ] **Step 3: Commit supplier list**

```bash
git add src/pages/partner/
git commit -m "feat: add SupplierList page with CRUD"
```

---

### Task 9: Create CustomerList Page

**Files:**
- Create: `src/pages/partner/CustomerList.tsx`

- [ ] **Step 1: Create CustomerList component**

```typescript
// src/pages/partner/CustomerList.tsx

import { useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, message, Popconfirm, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { customerApi } from '../../api/customer'
import type { Customer, CreateCustomerRequest, UpdateCustomerRequest } from '../types/partner'

export function CustomerList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['customers', page, pageSize],
    queryFn: () => customerApi.list(page, pageSize)
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateCustomerRequest) => customerApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      messageApi.success('创建成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('创建失败')
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateCustomerRequest }) => customerApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      messageApi.success('更新成功')
      handleCloseModal()
    },
    onError: () => messageApi.error('更新失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => customerApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleOpenCreate = () => {
    setEditingId(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleOpenEdit = (record: Customer) => {
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
      createMutation.mutate(values as CreateCustomerRequest)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '客户编码', dataIndex: 'code', width: 120 },
    { title: '客户名称', dataIndex: 'name' },
    { title: '联系人', dataIndex: 'contact', width: 100 },
    { title: '联系电话', dataIndex: 'phone', width: 130 },
    { title: '邮箱', dataIndex: 'email', width: 180 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      )
    },
    {
      title: '操作',
      width: 150,
      render: (_: unknown, record: Customer) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <>
      {contextHolder}
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增客户
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
      />
      <Modal
        title={editingId ? '编辑客户' : '新增客户'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="客户名称" rules={[{ required: true, message: '请输入客户名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="code" label="客户编码">
            <Input />
          </Form.Item>
          <Form.Item name="contact" label="联系人">
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="联系电话">
            <Input />
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input />
          </Form.Item>
          <Form.Item name="address" label="地址">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '禁用' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
```

- [ ] **Step 2: Update partner index**

```typescript
// src/pages/partner/index.ts

export { SupplierList } from './SupplierList'
export { CustomerList } from './CustomerList'
```

- [ ] **Step 3: Commit customer list**

```bash
git add src/pages/partner/
git commit -m "feat: add CustomerList page with CRUD"
```

---

### Task 10: Create InboundOrderList Page

**Files:**
- Create: `src/pages/order/InboundOrderList.tsx`
- Create: `src/pages/order/index.ts`

- [ ] **Step 1: Create InboundOrderList component**

```typescript
// src/pages/order/InboundOrderList.tsx

import { useState } from 'react'
import { Table, Button, Space, Drawer, Descriptions, Tag, message, Popconfirm, Badge } from 'antd'
import { EyeOutlined, CheckOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { inboundApi } from '../../api/inbound'
import { warehouseApi } from '../../api/warehouse'
import { supplierApi } from '../../api/supplier'
import { productApi } from '../../api/product'
import type { InboundOrder, InboundItem } from '../types/order'

const statusMap: Record<number, { text: string; color: string }> = {
  0: { text: '待确认', color: 'orange' },
  1: { text: '已完成', color: 'green' },
  2: { text: '已取消', color: 'red' }
}

export function InboundOrderList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [selectedOrder, setSelectedOrder] = useState<InboundOrder | null>(null)
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['inbound-orders', page, pageSize],
    queryFn: () => inboundApi.list(page, pageSize)
  })

  const { data: warehouses } = useQuery({
    queryKey: ['warehouses-all'],
    queryFn: () => warehouseApi.list(1, 100)
  })

  const { data: suppliers } = useQuery({
    queryKey: ['suppliers-all'],
    queryFn: () => supplierApi.list(1, 100)
  })

  const { data: products } = useQuery({
    queryKey: ['products-all'],
    queryFn: () => productApi.list(1, 100)
  })

  const confirmMutation = useMutation({
    mutationFn: (id: number) => inboundApi.confirm(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inbound-orders'] })
      messageApi.success('确认成功')
    },
    onError: () => messageApi.error('确认失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => inboundApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inbound-orders'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleViewDetail = async (id: number) => {
    const order = await inboundApi.get(id)
    setSelectedOrder(order)
    setDrawerOpen(true)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '订单编号', dataIndex: 'order_no', width: 180 },
    {
      title: '供应商',
      dataIndex: 'supplier_id',
      width: 150,
      render: (id: number | null) => id ? suppliers?.items.find(s => s.id === id)?.name || id : '-'
    },
    {
      title: '仓库',
      dataIndex: 'warehouse_id',
      width: 150,
      render: (id: number) => warehouses?.items.find(w => w.id === id)?.name || id
    },
    { title: '总数量', dataIndex: 'total_quantity', width: 100 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => {
        const s = statusMap[status] || { text: '未知', color: 'default' }
        return <Tag color={s.color}>{s.text}</Tag>
      }
    },
    { title: '备注', dataIndex: 'remark', ellipsis: true },
    { title: '创建时间', dataIndex: 'created_at', width: 180 },
    {
      title: '操作',
      width: 180,
      render: (_: unknown, record: InboundOrder) => (
        <Space>
          <Button type="link" icon={<EyeOutlined />} onClick={() => handleViewDetail(record.id)}>
            详情
          </Button>
          {record.status === 0 && (
            <Popconfirm title="确认入库?" onConfirm={() => confirmMutation.mutate(record.id)}>
              <Button type="link" icon={<CheckOutlined />}>
                确认
              </Button>
            </Popconfirm>
          )}
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const itemColumns = [
    { title: '商品', dataIndex: 'product_id', render: (id: number) => products?.items.find(p => p.id === id)?.name || id },
    { title: '数量', dataIndex: 'quantity', width: 100 },
    { title: '批次号', dataIndex: 'batch_no', width: 120 }
  ]

  return (
    <>
      {contextHolder}
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
      <Drawer
        title="入库单详情"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        width={600}
      >
        {selectedOrder && (
          <>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="订单编号">{selectedOrder.order_no}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={statusMap[selectedOrder.status]?.color}>
                  {statusMap[selectedOrder.status]?.text}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="供应商">
                {selectedOrder.supplier_id ? suppliers?.items.find(s => s.id === selectedOrder.supplier_id)?.name : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="仓库">
                {warehouses?.items.find(w => w.id === selectedOrder.warehouse_id)?.name}
              </Descriptions.Item>
              <Descriptions.Item label="总数量">{selectedOrder.total_quantity}</Descriptions.Item>
              <Descriptions.Item label="创建时间">{selectedOrder.created_at}</Descriptions.Item>
              <Descriptions.Item label="备注" span={2}>{selectedOrder.remark || '-'}</Descriptions.Item>
            </Descriptions>
            <h4 style={{ marginTop: 16 }}>商品明细</h4>
            <Table
              columns={itemColumns}
              dataSource={selectedOrder.items || []}
              rowKey="id"
              size="small"
              pagination={false}
            />
          </>
        )}
      </Drawer>
    </>
  )
}
```

- [ ] **Step 2: Create order index**

```typescript
// src/pages/order/index.ts

export { InboundOrderList } from './InboundOrderList'
export { OutboundOrderList } from './OutboundOrderList'
export { StockTransferList } from './StockTransferList'
```

- [ ] **Step 3: Commit inbound order list**

```bash
git add src/pages/order/
git commit -m "feat: add InboundOrderList page with detail drawer"
```

---

### Task 11: Create OutboundOrderList Page

**Files:**
- Create: `src/pages/order/OutboundOrderList.tsx`

- [ ] **Step 1: Create OutboundOrderList component**

```typescript
// src/pages/order/OutboundOrderList.tsx

import { useState } from 'react'
import { Table, Button, Space, Drawer, Descriptions, Tag, message, Popconfirm } from 'antd'
import { EyeOutlined, CheckOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { outboundApi } from '../../api/outbound'
import { warehouseApi } from '../../api/warehouse'
import { customerApi } from '../../api/customer'
import { productApi } from '../../api/product'
import type { OutboundOrder } from '../types/order'

const statusMap: Record<number, { text: string; color: string }> = {
  0: { text: '待确认', color: 'orange' },
  1: { text: '已完成', color: 'green' },
  2: { text: '已取消', color: 'red' }
}

export function OutboundOrderList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [selectedOrder, setSelectedOrder] = useState<OutboundOrder | null>(null)
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['outbound-orders', page, pageSize],
    queryFn: () => outboundApi.list(page, pageSize)
  })

  const { data: warehouses } = useQuery({
    queryKey: ['warehouses-all'],
    queryFn: () => warehouseApi.list(1, 100)
  })

  const { data: customers } = useQuery({
    queryKey: ['customers-all'],
    queryFn: () => customerApi.list(1, 100)
  })

  const { data: products } = useQuery({
    queryKey: ['products-all'],
    queryFn: () => productApi.list(1, 100)
  })

  const confirmMutation = useMutation({
    mutationFn: (id: number) => outboundApi.confirm(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['outbound-orders'] })
      messageApi.success('确认成功')
    },
    onError: () => messageApi.error('确认失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => outboundApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['outbound-orders'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleViewDetail = async (id: number) => {
    const order = await outboundApi.get(id)
    setSelectedOrder(order)
    setDrawerOpen(true)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '订单编号', dataIndex: 'order_no', width: 180 },
    {
      title: '客户',
      dataIndex: 'customer_id',
      width: 150,
      render: (id: number | null) => id ? customers?.items.find(c => c.id === id)?.name || id : '-'
    },
    {
      title: '仓库',
      dataIndex: 'warehouse_id',
      width: 150,
      render: (id: number) => warehouses?.items.find(w => w.id === id)?.name || id
    },
    { title: '总数量', dataIndex: 'total_quantity', width: 100 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => {
        const s = statusMap[status] || { text: '未知', color: 'default' }
        return <Tag color={s.color}>{s.text}</Tag>
      }
    },
    { title: '备注', dataIndex: 'remark', ellipsis: true },
    { title: '创建时间', dataIndex: 'created_at', width: 180 },
    {
      title: '操作',
      width: 180,
      render: (_: unknown, record: OutboundOrder) => (
        <Space>
          <Button type="link" icon={<EyeOutlined />} onClick={() => handleViewDetail(record.id)}>
            详情
          </Button>
          {record.status === 0 && (
            <Popconfirm title="确认出库?" onConfirm={() => confirmMutation.mutate(record.id)}>
              <Button type="link" icon={<CheckOutlined />}>
                确认
              </Button>
            </Popconfirm>
          )}
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const itemColumns = [
    { title: '商品', dataIndex: 'product_id', render: (id: number) => products?.items.find(p => p.id === id)?.name || id },
    { title: '数量', dataIndex: 'quantity', width: 100 },
    { title: '批次号', dataIndex: 'batch_no', width: 120 }
  ]

  return (
    <>
      {contextHolder}
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
      <Drawer
        title="出库单详情"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        width={600}
      >
        {selectedOrder && (
          <>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="订单编号">{selectedOrder.order_no}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={statusMap[selectedOrder.status]?.color}>
                  {statusMap[selectedOrder.status]?.text}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="客户">
                {selectedOrder.customer_id ? customers?.items.find(c => c.id === selectedOrder.customer_id)?.name : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="仓库">
                {warehouses?.items.find(w => w.id === selectedOrder.warehouse_id)?.name}
              </Descriptions.Item>
              <Descriptions.Item label="总数量">{selectedOrder.total_quantity}</Descriptions.Item>
              <Descriptions.Item label="创建时间">{selectedOrder.created_at}</Descriptions.Item>
              <Descriptions.Item label="备注" span={2}>{selectedOrder.remark || '-'}</Descriptions.Item>
            </Descriptions>
            <h4 style={{ marginTop: 16 }}>商品明细</h4>
            <Table
              columns={itemColumns}
              dataSource={selectedOrder.items || []}
              rowKey="id"
              size="small"
              pagination={false}
            />
          </>
        )}
      </Drawer>
    </>
  )
}
```

- [ ] **Step 2: Update order index**

```typescript
// src/pages/order/index.ts

export { InboundOrderList } from './InboundOrderList'
export { OutboundOrderList } from './OutboundOrderList'
export { StockTransferList } from './StockTransferList'
```

- [ ] **Step 3: Commit outbound order list**

```bash
git add src/pages/order/
git commit -m "feat: add OutboundOrderList page with detail drawer"
```

---

### Task 12: Create StockTransferList Page

**Files:**
- Create: `src/pages/order/StockTransferList.tsx`

- [ ] **Step 1: Create StockTransferList component**

```typescript
// src/pages/order/StockTransferList.tsx

import { useState } from 'react'
import { Table, Button, Space, Drawer, Descriptions, Tag, message, Popconfirm } from 'antd'
import { EyeOutlined, CheckOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { transferApi } from '../../api/transfer'
import { warehouseApi } from '../../api/warehouse'
import { productApi } from '../../api/product'
import type { StockTransfer } from '../types/order'

const statusMap: Record<number, { text: string; color: string }> = {
  0: { text: '待确认', color: 'orange' },
  1: { text: '已完成', color: 'green' },
  2: { text: '已取消', color: 'red' }
}

export function StockTransferList() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [selectedOrder, setSelectedOrder] = useState<StockTransfer | null>(null)
  const queryClient = useQueryClient()
  const [messageApi, contextHolder] = message.useMessage()

  const { data, isLoading } = useQuery({
    queryKey: ['stock-transfers', page, pageSize],
    queryFn: () => transferApi.list(page, pageSize)
  })

  const { data: warehouses } = useQuery({
    queryKey: ['warehouses-all'],
    queryFn: () => warehouseApi.list(1, 100)
  })

  const { data: products } = useQuery({
    queryKey: ['products-all'],
    queryFn: () => productApi.list(1, 100)
  })

  const confirmMutation = useMutation({
    mutationFn: (id: number) => transferApi.confirm(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stock-transfers'] })
      messageApi.success('确认成功')
    },
    onError: () => messageApi.error('确认失败')
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => transferApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stock-transfers'] })
      messageApi.success('删除成功')
    },
    onError: () => messageApi.error('删除失败')
  })

  const handleViewDetail = async (id: number) => {
    const order = await transferApi.get(id)
    setSelectedOrder(order)
    setDrawerOpen(true)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '调拨单号', dataIndex: 'order_no', width: 180 },
    {
      title: '调出仓库',
      dataIndex: 'from_warehouse_id',
      width: 150,
      render: (id: number) => warehouses?.items.find(w => w.id === id)?.name || id
    },
    {
      title: '调入仓库',
      dataIndex: 'to_warehouse_id',
      width: 150,
      render: (id: number) => warehouses?.items.find(w => w.id === id)?.name || id
    },
    { title: '总数量', dataIndex: 'total_quantity', width: 100 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number) => {
        const s = statusMap[status] || { text: '未知', color: 'default' }
        return <Tag color={s.color}>{s.text}</Tag>
      }
    },
    { title: '备注', dataIndex: 'remark', ellipsis: true },
    { title: '创建时间', dataIndex: 'created_at', width: 180 },
    {
      title: '操作',
      width: 180,
      render: (_: unknown, record: StockTransfer) => (
        <Space>
          <Button type="link" icon={<EyeOutlined />} onClick={() => handleViewDetail(record.id)}>
            详情
          </Button>
          {record.status === 0 && (
            <Popconfirm title="确认调拨?" onConfirm={() => confirmMutation.mutate(record.id)}>
              <Button type="link" icon={<CheckOutlined />}>
                确认
              </Button>
            </Popconfirm>
          )}
          <Popconfirm title="确定删除?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const itemColumns = [
    { title: '商品', dataIndex: 'product_id', render: (id: number) => products?.items.find(p => p.id === id)?.name || id },
    { title: '数量', dataIndex: 'quantity', width: 100 },
    { title: '批次号', dataIndex: 'batch_no', width: 120 }
  ]

  return (
    <>
      {contextHolder}
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
      <Drawer
        title="调拨单详情"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        width={600}
      >
        {selectedOrder && (
          <>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="调拨单号">{selectedOrder.order_no}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={statusMap[selectedOrder.status]?.color}>
                  {statusMap[selectedOrder.status]?.text}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="调出仓库">
                {warehouses?.items.find(w => w.id === selectedOrder.from_warehouse_id)?.name}
              </Descriptions.Item>
              <Descriptions.Item label="调入仓库">
                {warehouses?.items.find(w => w.id === selectedOrder.to_warehouse_id)?.name}
              </Descriptions.Item>
              <Descriptions.Item label="总数量">{selectedOrder.total_quantity}</Descriptions.Item>
              <Descriptions.Item label="创建时间">{selectedOrder.created_at}</Descriptions.Item>
              <Descriptions.Item label="备注" span={2}>{selectedOrder.remark || '-'}</Descriptions.Item>
            </Descriptions>
            <h4 style={{ marginTop: 16 }}>商品明细</h4>
            <Table
              columns={itemColumns}
              dataSource={selectedOrder.items || []}
              rowKey="id"
              size="small"
              pagination={false}
            />
          </>
        )}
      </Drawer>
    </>
  )
}
```

- [ ] **Step 2: Update order index**

```typescript
// src/pages/order/index.ts

export { InboundOrderList } from './InboundOrderList'
export { OutboundOrderList } from './OutboundOrderList'
export { StockTransferList } from './StockTransferList'
```

- [ ] **Step 3: Commit stock transfer list**

```bash
git add src/pages/order/
git commit -m "feat: add StockTransferList page with detail drawer"
```

---

### Task 13: Update App.tsx Routes

**Files:**
- Modify: `src/App.tsx`

- [ ] **Step 1: Update App.tsx with new routes**

```typescript
// src/App.tsx

import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { MainLayout } from './components/Layout'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Login } from './pages/auth/Login'
import { ChangePassword } from './pages/auth/ChangePassword'
import { WarehouseList, LocationList } from './pages/warehouse'
import { CategoryList, ProductList } from './pages/product'
import { InventoryList } from './pages/inventory'
import { SupplierList, CustomerList } from './pages/partner'
import { InboundOrderList, OutboundOrderList, StockTransferList } from './pages/order'

function Dashboard() {
  return <div>Dashboard</div>
}

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            path="/change-password"
            element={
              <ProtectedRoute>
                <ChangePassword />
              </ProtectedRoute>
            }
          />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="users" element={<div>用户管理</div>} />
            <Route path="roles" element={<div>角色管理</div>} />
            <Route path="warehouses" element={<WarehouseList />} />
            <Route path="locations" element={<LocationList />} />
            <Route path="categories" element={<CategoryList />} />
            <Route path="products" element={<ProductList />} />
            <Route path="inventory" element={<InventoryList />} />
            <Route path="suppliers" element={<SupplierList />} />
            <Route path="customers" element={<CustomerList />} />
            <Route path="inbound" element={<InboundOrderList />} />
            <Route path="outbound" element={<OutboundOrderList />} />
            <Route path="transfers" element={<StockTransferList />} />
            <Route path="audit-logs" element={<div>审计日志</div>} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
```

- [ ] **Step 2: Commit route updates**

```bash
git add src/App.tsx
git commit -m "feat: add routes for all business pages"
```

---

### Task 14: Run Build and Verify

**Files:**
- None (verification only)

- [ ] **Step 1: Run TypeScript build**

Run: `npm run build`
Expected: Build succeeds without errors

- [ ] **Step 2: Run lint**

Run: `npm run lint`
Expected: No lint errors

- [ ] **Step 3: Final commit if needed**

```bash
git status
# If any uncommitted changes:
git add -A
git commit -m "fix: resolve build issues"
```
