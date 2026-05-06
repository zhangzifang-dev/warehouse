# Filter Functionality for Management Pages - Design Specification

**Date:** 2026-05-01  
**Status:** Approved  
**Scope:** Add filter functionality to 8 management pages

## Overview

Add comprehensive filter functionality to 8 management pages in the warehouse system, enabling users to search and filter data through text inputs, dropdowns, and range selectors.

## Requirements

### Pages and Filters

#### 1. Warehouse Management
- **仓库名称** (Warehouse Name) - Text input, fuzzy search

#### 2. Category Management
- **分类名称** (Category Name) - Text input, fuzzy search

#### 3. Supplier Management
- **供应商编码** (Supplier Code) - Text input, fuzzy search
- **供应商名称** (Supplier Name) - Text input, fuzzy search
- **联系人** (Contact) - Text input, fuzzy search
- **联系电话** (Phone) - Text input, fuzzy search
- **状态** (Status) - Dropdown, exact match

#### 4. Customer Management
- **客户编码** (Customer Code) - Text input, fuzzy search
- **客户名称** (Customer Name) - Text input, fuzzy search
- **联系电话** (Phone) - Text input, fuzzy search
- **状态** (Status) - Dropdown, exact match

#### 5. Inventory Management
- **商品** (Product) - Text input, fuzzy search on product name
- **数量** (Quantity) - Range filter with min/max inputs
- **批次号** (Batch No) - Text input, fuzzy search

#### 6. Inbound Order Management
- **订单编号** (Order No) - Text input, fuzzy search
- **供应商** (Supplier) - Dropdown, exact match
- **仓库** (Warehouse) - Dropdown, exact match
- **数量** (Quantity) - Range filter with min/max inputs
- **创建时间** (Created At) - Date range with start/end pickers

#### 7. Outbound Order Management
- **订单编号** (Order No) - Text input, fuzzy search
- **客户** (Customer) - Dropdown, exact match
- **仓库** (Warehouse) - Dropdown, exact match
- **数量** (Quantity) - Range filter with min/max inputs
- **创建时间** (Created At) - Date range with start/end pickers

#### 8. Stock Transfer Management
- **订单编号** (Order No) - Text input, fuzzy search
- **调出仓库** (Source Warehouse) - Dropdown, exact match
- **调入仓库** (Target Warehouse) - Dropdown, exact match
- **创建时间** (Created At) - Date range with start/end pickers

## Architecture

### Pattern
Follow the existing audit log filter implementation across three layers:
- **Handler**: Parse query parameters → filter struct
- **Service**: Pass filter struct to repository
- **Repository**: Build dynamic WHERE clauses
- **Frontend**: Filter state + UI controls

### Backend Implementation

#### Layer 1: Filter Structures (Service Layer)

```go
// Warehouse
type WarehouseQueryFilter struct {
    Name     string
    Page     int
    PageSize int
}

// Category
type CategoryQueryFilter struct {
    Name     string
    Page     int
    PageSize int
}

// Supplier
type SupplierQueryFilter struct {
    Code     string
    Name     string
    Contact  string
    Phone    string
    Status   *int
    Page     int
    PageSize int
}

// Customer
type CustomerQueryFilter struct {
    Code     string
    Name     string
    Phone    string
    Status   *int
    Page     int
    PageSize int
}

// Inventory
type InventoryQueryFilter struct {
    ProductName string
    QuantityMin *float64
    QuantityMax *float64
    BatchNo     string
    Page        int
    PageSize    int
}

// InboundOrder
type InboundOrderQueryFilter struct {
    OrderNo        string
    SupplierID     *int64
    WarehouseID    *int64
    QuantityMin    *float64
    QuantityMax    *float64
    CreatedAtStart *time.Time
    CreatedAtEnd   *time.Time
    Page           int
    PageSize       int
}

// OutboundOrder
type OutboundOrderQueryFilter struct {
    OrderNo        string
    CustomerID     *int64
    WarehouseID    *int64
    QuantityMin    *float64
    QuantityMax    *float64
    CreatedAtStart *time.Time
    CreatedAtEnd   *time.Time
    Page           int
    PageSize       int
}

// StockTransfer
type StockTransferQueryFilter struct {
    OrderNo            string
    SourceWarehouseID  *int64
    TargetWarehouseID  *int64
    CreatedAtStart     *time.Time
    CreatedAtEnd       *time.Time
    Page           int
    PageSize       int
}
```

#### Layer 2: Repository Implementation

**Pattern for building filter queries:**

```go
func (r *InventoryRepository) List(ctx context.Context, filter *InventoryQueryFilter) ([]model.Inventory, int, error) {
    var inventories []model.Inventory
    
    q := r.db.NewSelect().Model(&inventories).
        Relation("Product").
        Relation("Warehouse").
        Where("inventory.deleted_at IS NULL")
    
    // Text search (fuzzy match)
    if filter.ProductName != "" {
        q = q.Where("product.name LIKE ?", "%"+filter.ProductName+"%")
    }
    
    // Range filters (independent min/max)
    if filter.QuantityMin != nil {
        q = q.Where("inventory.quantity >= ?", *filter.QuantityMin)
    }
    if filter.QuantityMax != nil {
        q = q.Where("inventory.quantity <= ?", *filter.QuantityMax)
    }
    
    // Text search
    if filter.BatchNo != "" {
        q = q.Where("inventory.batch_no LIKE ?", "%"+filter.BatchNo+"%")
    }
    
    total, err := q.
        Order("inventory.id DESC").
        Offset((filter.Page - 1) * filter.PageSize).
        Limit(filter.PageSize).
        ScanAndCount(ctx)
    
    return inventories, total, err
}
```

**Text search pattern:** `LIKE '%value%'` for fuzzy matching  
**Range filter pattern:** Separate min/max with `>=` and `<=`, both optional  
**Foreign key filters:** Exact match using ID, pointer to distinguish "not set"  
**Status dropdown:** Exact match, pointer to distinguish "not set" from 0

#### Layer 3: Handler Implementation

**Pattern for parsing query parameters:**

```go
func (h *InventoryHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
    
    filter := &service.InventoryQueryFilter{
        Page:     page,
        PageSize: pageSize,
    }
    
    // Text filters
    if productName := c.Query("product_name"); productName != "" {
        filter.ProductName = productName
    }
    if batchNo := c.Query("batch_no"); batchNo != "" {
        filter.BatchNo = batchNo
    }
    
    // Range filters
    if quantityMinStr := c.Query("quantity_min"); quantityMinStr != "" {
        if quantityMin, err := strconv.ParseFloat(quantityMinStr, 64); err == nil {
            filter.QuantityMin = &quantityMin
        }
    }
    if quantityMaxStr := c.Query("quantity_max"); quantityMaxStr != "" {
        if quantityMax, err := strconv.ParseFloat(quantityMaxStr, 64); err == nil {
            filter.QuantityMax = &quantityMax
        }
    }
    
    result, err := h.inventoryService.List(c.Request.Context(), filter)
    // ... error handling and response
}
```

### Frontend Implementation

#### API Layer

**Filter interfaces:**

```typescript
export interface WarehouseFilter {
  name?: string
}

export interface CategoryFilter {
  name?: string
}

export interface SupplierFilter {
  code?: string
  name?: string
  contact?: string
  phone?: string
  status?: number
}

export interface CustomerFilter {
  code?: string
  name?: string
  phone?: string
  status?: number
}

export interface InventoryFilter {
  product_name?: string
  quantity_min?: number
  quantity_max?: number
  batch_no?: string
}

export interface InboundOrderFilter {
  order_no?: string
  supplier_id?: number
  warehouse_id?: number
  quantity_min?: number
  quantity_max?: number
  created_at_start?: string
  created_at_end?: string
}

export interface OutboundOrderFilter {
  order_no?: string
  customer_id?: number
  warehouse_id?: number
  quantity_min?: number
  quantity_max?: number
  created_at_start?: string
  created_at_end?: string
}

export interface StockTransferFilter {
  order_no?: string
  source_warehouse_id?: number
  target_warehouse_id?: number
  created_at_start?: string
  created_at_end?: string
}
```

**API method pattern:**

```typescript
export const inventoryApi = {
  list: async (page = 1, size = 10, filter?: InventoryFilter): Promise<PaginatedResponse<Inventory>> => {
    const params = new URLSearchParams()
    params.append('page', String(page))
    params.append('size', String(size))
    
    if (filter) {
      if (filter.product_name) params.append('product_name', filter.product_name)
      if (filter.batch_no) params.append('batch_no', filter.batch_no)
      if (filter.quantity_min !== undefined) params.append('quantity_min', String(filter.quantity_min))
      if (filter.quantity_max !== undefined) params.append('quantity_max', String(filter.quantity_max))
    }
    
    const response = await api.get<PaginatedResponse<Inventory>>('/inventory', { params })
    return response.data
  },
  // ... other methods
}
```

#### UI Component Pattern

**Layout:** Single row horizontal layout with wrapping (Space wrap)

**Component usage:**

| Filter Type | Component | Width | Props |
|------------|-----------|-------|-------|
| Text search | Input | 150px | placeholder, allowClear |
| Status dropdown | Select | 120px | options, allowClear |
| Foreign key dropdown | Select | 150px | options (from API), allowClear, showSearch |
| Range min/max | Input (2) | 120px each | placeholder="最小/最大", allowClear |
| Date range | DatePicker (2) | 180px each | placeholder="开始/结束日期", showTime |

**State management:**

```tsx
const [filter, setFilter] = useState<InventoryFilter>({})

const handleFilterChange = (key: string, value: any) => {
  setFilter(prev => ({ ...prev, [key]: value || undefined }))
  setPage(1) // Reset to first page on filter change
}
```

**Query key pattern:**

```tsx
const { data, isLoading } = useQuery({
  queryKey: ['inventory', page, pageSize, filter],
  queryFn: () => inventoryApi.list(page, pageSize, filter)
})
```

## Testing Strategy

### Backend Tests

1. **Unit tests - Repository layer:**
   - Test WHERE clause building with various filter combinations
   - Test empty filters (should return all records)
   - Test partial range filters (only min or only max)
   - Test special characters in text search

2. **Integration tests - Handler layer:**
   - Test API endpoints with various query parameter combinations
   - Test multiple filters applied simultaneously
   - Test pagination combined with filters

### Frontend Tests

1. **Component tests:**
   - Filter state updates correctly
   - Page resets to 1 when filter changes
   - Clear filter functionality works
   - Query key includes filter state

## Implementation Plan

### Phase 1: Simple Modules (1-2 filters)
1. Warehouse (1 filter: name)
2. Category (1 filter: name)

**Rationale:** Validate the pattern with minimal complexity

### Phase 2: Medium Modules (3-5 filters)
3. Supplier (5 filters)
4. Customer (4 filters)

**Rationale:** Build confidence with dropdown filters

### Phase 3: Complex Modules (range filters)
5. Inventory (4 filters, quantity range)
6. Inbound Order (7 filters, quantity + date ranges)
7. Outbound Order (7 filters, quantity + date ranges)
8. Stock Transfer (5 filters, date range)

**Rationale:** Range filter logic is reusable; date range is identical for order modules

## Success Criteria

- All 8 management pages have functional filter controls
- Text search uses fuzzy matching (LIKE '%value%')
- Dropdowns show exact matches
- Range filters work independently (can set only min, only max, or both)
- Filter state is maintained in URL query parameters
- Page resets to 1 when filter changes
- Clear filter functionality works
- All filters properly escape SQL to prevent injection
- Backend tests cover filter logic
- No performance degradation with filters

## Dependencies

- Existing audit log filter implementation (for reference)
- Gin query parameter parsing (c.Query, c.QueryArray)
- Bun query builder (Where, Like clauses)
- Ant Design components (Input, Select, DatePicker)
- React Query for state management

## Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| SQL injection in text search | Use parameterized queries (bun handles this) |
| Performance with large datasets | Add database indexes on filter columns |
| Filter combination complexity | Test thoroughly; start with simple combinations |
| UI clutter with many filters | Use single row layout with wrapping |

## Out of Scope

- Advanced search (AND/OR logic between filters)
- Filter presets/saved searches
- Export filtered results
- Filter by audit log history
