# 出库订单过滤器功能实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为出库订单管理添加7个过滤器功能（订单编号、客户、仓库、数量范围、创建时间范围）

**Architecture:** 参考入库订单过滤器实现模式，将supplier_id替换为customer_id，其他过滤器实现完全相同。采用分层架构：Model → Repository → Service → Handler → Frontend

**Tech Stack:** Go (backend), React + TypeScript + Ant Design (frontend)

---

## 任务分解

### Task 1: 后端Model层 - 添加过滤器结构体

**Files:**
- Modify: `internal/model/order.go:44-56`

- [ ] **Step 1: 添加OutboundOrderQueryFilter结构体**

在`internal/model/order.go`中，在`OutboundOrder`结构体定义之前添加过滤器结构体：

```go
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
```

- [ ] **Step 2: 验证代码编译**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go build ./internal/model`
Expected: 编译成功，无错误

- [ ] **Step 3: Commit**

```bash
git add internal/model/order.go
git commit -m "feat: add OutboundOrderQueryFilter struct for filter support"
```

---

### Task 2: 后端Repository层 - 实现过滤器查询

**Files:**
- Modify: `internal/repository/outbound_order.go:74-92`
- Create: `internal/repository/outbound_order_test.go`

- [ ] **Step 1: 编写失败的Repository测试**

创建`internal/repository/outbound_order_test.go`：

```go
package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func setupOutboundOrderTest(t *testing.T) (*OutboundOrderRepository, *bun.DB, context.Context) {
	t.Helper()
	sqlDB, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test_db?parseTime=true")
	if err != nil {
		t.Fatal(err)
	}
	db := bun.NewDB(sqlDB, mysqldialect.New())
	repo := NewOutboundOrderRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}

func TestOutboundOrderRepository_ListWithFilter(t *testing.T) {
	repo, _, ctx := setupOutboundOrderTest(t)
	customerID := int64(1)
	warehouseID := int64(1)
	quantityMin := 10.0
	quantityMax := 100.0
	startTime := time.Now()
	endTime := time.Now().Add(24 * time.Hour)

	filter := &model.OutboundOrderQueryFilter{
		OrderNo:        "SO-2024",
		CustomerID:     &customerID,
		WarehouseID:    &warehouseID,
		QuantityMin:    &quantityMin,
		QuantityMax:    &quantityMax,
		CreatedAtStart: &startTime,
		CreatedAtEnd:   &endTime,
		Page:           1,
		PageSize:       10,
	}

	_, _, err := repo.ListWithFilter(ctx, filter)
	if err == nil {
		t.Error("ListWithFilter() should return error with mock DB")
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go test ./internal/repository -run TestOutboundOrderRepository_ListWithFilter -v`
Expected: FAIL - `repo.ListWithFilter undefined`

- [ ] **Step 3: 实现ListWithFilter方法**

在`internal/repository/outbound_order.go`的`List`方法之后添加：

```go
func (r *OutboundOrderRepository) ListWithFilter(ctx context.Context, filter *model.OutboundOrderQueryFilter) ([]model.OutboundOrder, int, error) {
	var orders []model.OutboundOrder
	q := r.db.NewSelect().
		Model(&orders).
		Where("deleted_at IS NULL")

	if filter.OrderNo != "" {
		q = q.Where("order_no LIKE ?", "%"+filter.OrderNo+"%")
	}

	if filter.CustomerID != nil {
		q = q.Where("customer_id = ?", *filter.CustomerID)
	}

	if filter.WarehouseID != nil {
		q = q.Where("warehouse_id = ?", *filter.WarehouseID)
	}

	if filter.QuantityMin != nil {
		q = q.Where("total_quantity >= ?", *filter.QuantityMin)
	}

	if filter.QuantityMax != nil {
		q = q.Where("total_quantity <= ?", *filter.QuantityMax)
	}

	if filter.CreatedAtStart != nil {
		q = q.Where("created_at >= ?", filter.CreatedAtStart)
	}

	if filter.CreatedAtEnd != nil {
		q = q.Where("created_at <= ?", filter.CreatedAtEnd)
	}

	total, err := q.
		Order("id DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}
```

- [ ] **Step 4: 运行测试确认通过**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go test ./internal/repository -run TestOutboundOrderRepository_ListWithFilter -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/repository/outbound_order.go internal/repository/outbound_order_test.go
git commit -m "feat: implement ListWithFilter in OutboundOrderRepository with all filter types"
```

---

### Task 3: 后端Service层 - 添加过滤器服务方法

**Files:**
- Modify: `internal/service/outbound_order.go:11-18`
- Modify: `internal/service/outbound_order_test.go`

- [ ] **Step 1: 编写失败的Service测试**

在`internal/service/outbound_order_test.go`中添加测试：

```go
func TestOutboundOrderService_ListWithFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockOutboundOrderRepository(ctrl)
	mockItemRepo := NewMockOutboundItemRepository(ctrl)
	mockInventorySvc := NewMockInventoryServiceForOutbound(ctrl)
	svc := NewOutboundOrderService(mockRepo, mockItemRepo, mockInventorySvc, nil)

	customerID := int64(1)
	warehouseID := int64(1)
	quantityMin := 10.0
	quantityMax := 100.0

	filter := &model.OutboundOrderQueryFilter{
		CustomerID:  &customerID,
		WarehouseID: &warehouseID,
		QuantityMin: &quantityMin,
		QuantityMax: &quantityMax,
		Page:        1,
		PageSize:    10,
	}

	mockRepo.EXPECT().
		ListWithFilter(gomock.Any(), filter).
		Return([]model.OutboundOrder{}, 0, nil)

	result, err := svc.ListWithFilter(context.Background(), filter)
	if err != nil {
		t.Errorf("ListWithFilter() error = %v", err)
	}
	if result == nil {
		t.Error("ListWithFilter() should return result")
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go test ./internal/service -run TestOutboundOrderService_ListWithFilter -v`
Expected: FAIL - `svc.ListWithFilter undefined`

- [ ] **Step 3: 更新Repository接口定义**

在`internal/service/outbound_order.go`的`OutboundOrderRepository`接口中添加方法：

```go
type OutboundOrderRepository interface {
	Create(ctx context.Context, order *model.OutboundOrder) error
	GetByID(ctx context.Context, id int64) (*model.OutboundOrder, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*model.OutboundOrder, error)
	List(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.OutboundOrder, int, error)
	ListWithFilter(ctx context.Context, filter *model.OutboundOrderQueryFilter) ([]model.OutboundOrder, int, error)
	Update(ctx context.Context, order *model.OutboundOrder) error
	Delete(ctx context.Context, id int64) error
}
```

- [ ] **Step 4: 实现ListWithFilter服务方法**

在`internal/service/outbound_order.go`的`List`方法之后添加：

```go
func (s *OutboundOrderService) ListWithFilter(ctx context.Context, filter *model.OutboundOrderQueryFilter) (*ListOutboundOrdersResult, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	orders, total, err := s.orderRepo.ListWithFilter(ctx, filter)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list outbound orders")
	}

	return &ListOutboundOrdersResult{
		Orders: orders,
		Total:  total,
	}, nil
}
```

- [ ] **Step 5: 运行测试确认通过**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go test ./internal/service -run TestOutboundOrderService_ListWithFilter -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/service/outbound_order.go internal/service/outbound_order_test.go
git commit -m "feat: add ListWithFilter service method with validation"
```

---

### Task 4: 后端Handler层 - 实现过滤器参数解析

**Files:**
- Modify: `internal/handler/outbound_order.go:39-46`
- Modify: `internal/handler/outbound_order.go:98-116`
- Modify: `internal/handler/outbound_order_test.go`

- [ ] **Step 1: 编写失败的Handler测试**

在`internal/handler/outbound_order_test.go`中添加测试：

```go
func TestOutboundOrderHandler_List_WithFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := service.NewMockOutboundOrderService(ctrl)
	handler := NewOutboundOrderHandler(mockSvc)

	router := gin.New()
	router.GET("/outbound-orders", handler.List)

	customerID := int64(1)
	warehouseID := int64(2)
	quantityMin := 10.0
	quantityMax := 100.0
	startTime := time.Now()
	endTime := startTime.Add(24 * time.Hour)

	expectedFilter := &model.OutboundOrderQueryFilter{
		OrderNo:        "SO-2024",
		CustomerID:     &customerID,
		WarehouseID:    &warehouseID,
		QuantityMin:    &quantityMin,
		QuantityMax:    &quantityMax,
		CreatedAtStart: &startTime,
		CreatedAtEnd:   &endTime,
		Page:           1,
		PageSize:       10,
	}

	mockSvc.EXPECT().
		ListWithFilter(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, filter *model.OutboundOrderQueryFilter) (*service.ListOutboundOrdersResult, error) {
			if filter.OrderNo != expectedFilter.OrderNo {
				t.Errorf("expected order_no %s, got %s", expectedFilter.OrderNo, filter.OrderNo)
			}
			if filter.CustomerID == nil || *filter.CustomerID != *expectedFilter.CustomerID {
				t.Errorf("customer_id mismatch")
			}
			return &service.ListOutboundOrdersResult{Orders: []model.OutboundOrder{}, Total: 0}, nil
		})

	startTimeStr := startTime.Format(time.RFC3339)
	endTimeStr := endTime.Format(time.RFC3339)

	req := httptest.NewRequest("GET", fmt.Sprintf("/outbound-orders?order_no=SO-2024&customer_id=1&warehouse_id=2&quantity_min=10&quantity_max=100&created_at_start=%s&created_at_end=%s", startTimeStr, endTimeStr), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go test ./internal/handler -run TestOutboundOrderHandler_List_WithFilter -v`
Expected: FAIL - mock expects ListWithFilter but got List

- [ ] **Step 3: 更新Handler接口**

在`internal/handler/outbound_order.go`的`OutboundOrderService`接口中添加方法：

```go
type OutboundOrderService interface {
	Create(ctx context.Context, input *service.CreateOutboundOrderInput) (*model.OutboundOrder, error)
	GetByID(ctx context.Context, id int64) (*model.OutboundOrder, error)
	List(ctx context.Context, page, pageSize int, warehouseID, status int) (*service.ListOutboundOrdersResult, error)
	ListWithFilter(ctx context.Context, filter *model.OutboundOrderQueryFilter) (*service.ListOutboundOrdersResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateOutboundOrderInput) (*model.OutboundOrder, error)
	Delete(ctx context.Context, id int64) error
	Confirm(ctx context.Context, id int64) (*model.OutboundOrder, error)
}
```

- [ ] **Step 4: 实现Handler过滤器参数解析**

修改`internal/handler/outbound_order.go`的`List`方法，导入`time`包并替换整个方法：

```go
func (h *OutboundOrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	filter := &model.OutboundOrderQueryFilter{
		Page:     page,
		PageSize: pageSize,
	}

	if orderNo := c.Query("order_no"); orderNo != "" {
		filter.OrderNo = orderNo
	}

	if customerIDStr := c.Query("customer_id"); customerIDStr != "" {
		if customerID, err := strconv.ParseInt(customerIDStr, 10, 64); err == nil {
			filter.CustomerID = &customerID
		}
	}

	if warehouseIDStr := c.Query("warehouse_id"); warehouseIDStr != "" {
		if warehouseID, err := strconv.ParseInt(warehouseIDStr, 10, 64); err == nil {
			filter.WarehouseID = &warehouseID
		}
	}

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

	if createdAtStartStr := c.Query("created_at_start"); createdAtStartStr != "" {
		if createdAtStart, err := time.Parse(time.RFC3339, createdAtStartStr); err == nil {
			filter.CreatedAtStart = &createdAtStart
		}
	}

	if createdAtEndStr := c.Query("created_at_end"); createdAtEndStr != "" {
		if createdAtEnd, err := time.Parse(time.RFC3339, createdAtEndStr); err == nil {
			filter.CreatedAtEnd = &createdAtEnd
		}
	}

	result, err := h.outboundOrderService.ListWithFilter(c.Request.Context(), filter)
	if err != nil {
		handleOutboundOrderError(c, err)
		return
	}

	response.Success(c, OutboundOrderListResponse{
		Items: result.Orders,
		Total: result.Total,
		Page:  page,
		Size:  pageSize,
	})
}
```

- [ ] **Step 5: 运行测试确认通过**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go test ./internal/handler -run TestOutboundOrderHandler_List_WithFilter -v`
Expected: PASS

- [ ] **Step 6: 运行所有后端测试**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go test ./internal/... -v`
Expected: 所有测试通过

- [ ] **Step 7: Commit**

```bash
git add internal/handler/outbound_order.go internal/handler/outbound_order_test.go
git commit -m "feat: implement filter parameter parsing in OutboundOrderHandler"
```

---

### Task 5: 前端API层 - 添加过滤器接口

**Files:**
- Modify: `web/src/api/outbound.ts`

- [ ] **Step 1: 添加OutboundOrderFilter接口**

在`web/src/api/outbound.ts`中，在`outboundApi`定义之前添加接口定义：

```typescript
export interface OutboundOrderFilter {
  order_no?: string
  customer_id?: number
  warehouse_id?: number
  quantity_min?: number
  quantity_max?: number
  created_at_start?: string
  created_at_end?: string
}
```

- [ ] **Step 2: 更新list方法支持过滤器**

修改`outboundApi.list`方法：

```typescript
  list: async (page = 1, size = 10, filter?: OutboundOrderFilter): Promise<PaginatedResponse<OutboundOrder>> => {
    const response = await api.get<PaginatedResponse<OutboundOrder>>('/outbound-orders', {
      params: { page, size, ...filter }
    })
    return response.data
  },
```

- [ ] **Step 3: 验证TypeScript编译**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality/web && npm run typecheck`
Expected: 编译成功，无类型错误

- [ ] **Step 4: Commit**

```bash
git add web/src/api/outbound.ts
git commit -m "feat: add OutboundOrderFilter interface and update list API"
```

---

### Task 6: 前端UI层 - 添加过滤器组件

**Files:**
- Modify: `web/src/pages/order/OutboundOrderList.tsx`

- [ ] **Step 1: 导入必要的组件和图标**

在文件顶部添加导入：

```typescript
import { useState } from 'react'
import { Table, Button, Space, Drawer, Descriptions, Tag, message, Popconfirm, Input, Select, DatePicker, Form, Row, Col, Card, InputNumber } from 'antd'
import { EyeOutlined, CheckOutlined, DeleteOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { outboundApi, type OutboundOrderFilter } from '../../api/outbound'
import { warehouseApi } from '../../api/warehouse'
import { customerApi } from '../../api/customer'
import { productApi } from '../../api/product'
import type { OutboundOrder } from '../../types/order'

const { RangePicker } = DatePicker
```

- [ ] **Step 2: 添加filter state和form**

在`OutboundOrderList`组件中，在现有state之后添加：

```typescript
  const [filter, setFilter] = useState<OutboundOrderFilter>({})
  const [form] = Form.useForm()
```

- [ ] **Step 3: 更新useQuery包含filter**

修改`useQuery`调用：

```typescript
  const { data, isLoading } = useQuery({
    queryKey: ['outbound-orders', page, pageSize, filter],
    queryFn: () => outboundApi.list(page, pageSize, filter)
  })
```

- [ ] **Step 4: 添加过滤器处理函数**

在`handleViewDetail`函数之后添加：

```typescript
  const handleSearch = () => {
    const values = form.getFieldsValue()
    const newFilter: OutboundOrderFilter = {}
    
    if (values.order_no) newFilter.order_no = values.order_no
    if (values.customer_id) newFilter.customer_id = values.customer_id
    if (values.warehouse_id) newFilter.warehouse_id = values.warehouse_id
    if (values.quantity_min !== undefined && values.quantity_min !== null) {
      newFilter.quantity_min = values.quantity_min
    }
    if (values.quantity_max !== undefined && values.quantity_max !== null) {
      newFilter.quantity_max = values.quantity_max
    }
    if (values.created_at_range && values.created_at_range[0] && values.created_at_range[1]) {
      newFilter.created_at_start = values.created_at_range[0].format('YYYY-MM-DDTHH:mm:ssZ')
      newFilter.created_at_end = values.created_at_range[1].format('YYYY-MM-DDTHH:mm:ssZ')
    }
    
    setFilter(newFilter)
    setPage(1)
  }

  const handleReset = () => {
    form.resetFields()
    setFilter({})
    setPage(1)
  }
```

- [ ] **Step 5: 添加过滤器UI组件**

在return语句中，在`{contextHolder}`之后、`<Table>`之前添加过滤器卡片：

```typescript
      <Card style={{ marginBottom: 16 }}>
        <Form form={form} layout="inline">
          <Row gutter={16} style={{ width: '100%' }}>
            <Col>
              <Form.Item name="order_no" label="订单编号">
                <Input placeholder="订单编号" style={{ width: 150 }} allowClear />
              </Form.Item>
            </Col>
            <Col>
              <Form.Item name="customer_id" label="客户">
                <Select
                  placeholder="选择客户"
                  style={{ width: 150 }}
                  allowClear
                  showSearch
                  filterOption={(input, option) =>
                    (option?.label ?? '').toString().toLowerCase().includes(input.toLowerCase())
                  }
                  options={customers?.items?.map((c: { id: number; name: string }) => ({
                    label: c.name,
                    value: c.id
                  }))}
                />
              </Form.Item>
            </Col>
            <Col>
              <Form.Item name="warehouse_id" label="仓库">
                <Select
                  placeholder="选择仓库"
                  style={{ width: 150 }}
                  allowClear
                  showSearch
                  filterOption={(input, option) =>
                    (option?.label ?? '').toString().toLowerCase().includes(input.toLowerCase())
                  }
                  options={warehouses?.items?.map((w: { id: number; name: string }) => ({
                    label: w.name,
                    value: w.id
                  }))}
                />
              </Form.Item>
            </Col>
            <Col>
              <Form.Item label="数量范围">
                <Space>
                  <Form.Item name="quantity_min" noStyle>
                    <InputNumber placeholder="最小" style={{ width: 120 }} min={0} />
                  </Form.Item>
                  <span>-</span>
                  <Form.Item name="quantity_max" noStyle>
                    <InputNumber placeholder="最大" style={{ width: 120 }} min={0} />
                  </Form.Item>
                </Space>
              </Form.Item>
            </Col>
            <Col>
              <Form.Item name="created_at_range" label="创建时间">
                <RangePicker showTime style={{ width: 360 }} />
              </Form.Item>
            </Col>
            <Col>
              <Space>
                <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
                  查询
                </Button>
                <Button icon={<ReloadOutlined />} onClick={handleReset}>
                  重置
                </Button>
              </Space>
            </Col>
          </Row>
        </Form>
      </Card>
```

- [ ] **Step 6: 验证前端编译**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality/web && npm run typecheck`
Expected: 编译成功，无类型错误

- [ ] **Step 7: Commit**

```bash
git add web/src/pages/order/OutboundOrderList.tsx
git commit -m "feat: add filter UI for outbound order list with all 7 filters"
```

---

### Task 7: 集成测试和最终验证

**Files:**
- 无新文件

- [ ] **Step 1: 运行所有后端测试**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go test ./internal/... -v`
Expected: 所有测试通过

- [ ] **Step 2: 运行前端类型检查**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality/web && npm run typecheck`
Expected: 无类型错误

- [ ] **Step 3: 运行前端lint**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality/web && npm run lint`
Expected: 无lint错误

- [ ] **Step 4: 启动后端服务验证API**

Run: `cd /home/zzf/projects/goinvent/warehouse/.worktrees/filter-functionality && go run cmd/server/main.go`
Expected: 服务正常启动（手动验证）

- [ ] **Step 5: 测试API端点**

使用curl或Postman测试：

```bash
# 测试无过滤器
curl "http://localhost:8080/api/v1/outbound-orders?page=1&size=10"

# 测试订单编号过滤
curl "http://localhost:8080/api/v1/outbound-orders?page=1&size=10&order_no=SO-2024"

# 测试客户过滤
curl "http://localhost:8080/api/v1/outbound-orders?page=1&size=10&customer_id=1"

# 测试数量范围
curl "http://localhost:8080/api/v1/outbound-orders?page=1&size=10&quantity_min=10&quantity_max=100"
```

Expected: API响应正常，返回过滤后的数据

- [ ] **Step 6: 最终commit（如果前面有遗漏）**

```bash
git status
# 如有未提交文件，补充提交
```

---

## 完成标准

- [ ] 后端所有测试通过
- [ ] 前端TypeScript编译无错误
- [ ] 前端lint检查通过
- [ ] API端点可通过所有过滤器参数查询
- [ ] UI显示过滤器组件并正常工作
- [ ] 所有代码已提交到git
