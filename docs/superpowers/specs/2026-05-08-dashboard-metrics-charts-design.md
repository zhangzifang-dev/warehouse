# 仪表盘指标图表设计文档

**日期：** 2026-05-08  
**项目：** 仓库管理系统 (Warehouse Management System)  
**作者：** OpenCode

## 1. 概述

为仓库管理系统添加仪表盘功能，通过可视化图表展示核心业务指标，帮助管理者快速了解仓库运营状况。

### 1.1 目标

- 提供直观的业务指标展示
- 支持数据钻取和详情查看
- 支持数据导出功能
- 提供灵活的时间范围筛选

### 1.2 范围

**包含内容：**
- 统计卡片：总库存量、库存预警、今日入库/出库
- 趋势图表：出入库趋势图、热销产品排行
- 分析图表：仓库使用率、供应商绩效
- 交互功能：数据钻取、跳转详情、导出

**不包含内容：**
- 实时数据推送（WebSocket）
- 自动定时刷新
- 自定义图表配置
- 移动端专属优化（仅响应式适配）

## 2. 架构设计

### 2.1 技术栈

**后端：**
- Go + Gin（现有框架）
- PostgreSQL（现有数据库）
- 分层架构：Handler → Service → Repository

**前端：**
- React 19 + TypeScript
- Ant Design 6（现有UI库）
- Ant Design Charts（新增，用于图表渲染）
- React Router（现有路由）
- TanStack Query（现有数据获取）

### 2.2 整体架构

采用后端专用统计API + 前端可视化的方案：

```
前端层 (React + Ant Design Charts)
    ↓ HTTP API调用
后端API层 (DashboardHandler)
    ↓ 业务逻辑处理
服务层 (DashboardService)
    ↓ 数据查询与聚合
数据访问层 (DashboardRepository)
    ↓ SQL查询
数据库层 (PostgreSQL)
```

**优势：**
- 数据库层面聚合，性能优异
- 职责清晰，易于维护
- 支持灵活的时间范围筛选
- 符合项目现有架构模式

## 3. API设计

### 3.1 接口列表

| 接口 | 方法 | 说明 | 参数 |
|------|------|------|------|
| `/api/v1/dashboard/overview` | GET | 总览统计 | 无 |
| `/api/v1/dashboard/trend` | GET | 出入库趋势 | start_date, end_date |
| `/api/v1/dashboard/top-products` | GET | 热销产品排行 | start_date, end_date, limit |
| `/api/v1/dashboard/warehouse-usage` | GET | 仓库使用率 | 无 |
| `/api/v1/dashboard/supplier-performance` | GET | 供应商绩效 | start_date, end_date, limit |
| `/api/v1/dashboard/pending-orders` | GET | 待处理订单 | 无 |
| `/api/v1/dashboard/export` | GET | 导出数据 | format, start_date, end_date |

### 3.2 请求示例

```bash
# 获取总览统计
GET /api/v1/dashboard/overview

# 获取30天趋势数据
GET /api/v1/dashboard/trend?start_date=2026-04-08&end_date=2026-05-08

# 获取热销产品TOP 10
GET /api/v1/dashboard/top-products?start_date=2026-04-08&end_date=2026-05-08&limit=10

# 导出PDF
GET /api/v1/dashboard/export?format=pdf&start_date=2026-04-08&end_date=2026-05-08
```

### 3.3 响应数据结构

#### 3.3.1 总览统计

```json
{
  "code": 200,
  "data": {
    "total_inventory": 12589.5,
    "inventory_warning": 23,
    "today_inbound": 15,
    "today_inbound_qty": 456.0,
    "today_outbound": 18,
    "today_outbound_qty": 523.0
  }
}
```

#### 3.3.2 趋势数据

```json
{
  "code": 200,
  "data": [
    {
      "date": "2026-04-08",
      "inbound_qty": 120.5,
      "outbound_qty": 95.0
    },
    {
      "date": "2026-04-09",
      "inbound_qty": 150.0,
      "outbound_qty": 130.5
    }
  ]
}
```

#### 3.3.3 热销产品排行

```json
{
  "code": 200,
  "data": [
    {
      "product_id": 123,
      "product_name": "螺丝M8",
      "category": "紧固件",
      "total_qty": 1250.0,
      "order_count": 45
    }
  ]
}
```

#### 3.3.4 仓库使用率

```json
{
  "code": 200,
  "data": [
    {
      "warehouse_id": 1,
      "warehouse_name": "主仓库",
      "capacity": 500,
      "used_capacity": 320,
      "usage_rate": 64.0
    }
  ]
}
```

#### 3.3.5 供应商绩效

```json
{
  "code": 200,
  "data": [
    {
      "supplier_id": 10,
      "supplier_name": "优质供应商A",
      "order_count": 50,
      "total_value": 125000.0,
      "on_time_rate": 95.0,
      "quality_score": 92.0,
      "delivery_score": 88.0
    }
  ]
}
```

#### 3.3.6 待处理订单

```json
{
  "code": 200,
  "data": {
    "inbound_pending": 5,
    "outbound_pending": 8,
    "transfer_pending": 3
  }
}
```

## 4. 数据模型

### 4.1 后端模型定义

```go
// OverviewStats 总览统计
type OverviewStats struct {
    TotalInventory   float64 `json:"total_inventory"`
    InventoryWarning int     `json:"inventory_warning"`
    TodayInbound     int     `json:"today_inbound"`
    TodayInboundQty  float64 `json:"today_inbound_qty"`
    TodayOutbound    int     `json:"today_outbound"`
    TodayOutboundQty float64 `json:"today_outbound_qty"`
}

// TrendData 趋势数据
type TrendData struct {
    Date        string  `json:"date"`
    InboundQty  float64 `json:"inbound_qty"`
    OutboundQty float64 `json:"outbound_qty"`
}

// TopProduct 热销产品
type TopProduct struct {
    ProductID   int64   `json:"product_id"`
    ProductName string  `json:"product_name"`
    Category    string  `json:"category"`
    TotalQty    float64 `json:"total_qty"`
    OrderCount  int     `json:"order_count"`
}

// WarehouseUsage 仓库使用率
type WarehouseUsage struct {
    WarehouseID   int64   `json:"warehouse_id"`
    WarehouseName string  `json:"warehouse_name"`
    Capacity      float64 `json:"capacity"`
    UsedCapacity  float64 `json:"used_capacity"`
    UsageRate     float64 `json:"usage_rate"`
}

// SupplierPerformance 供应商绩效
type SupplierPerformance struct {
    SupplierID    int64   `json:"supplier_id"`
    SupplierName  string  `json:"supplier_name"`
    OrderCount    int     `json:"order_count"`
    TotalValue    float64 `json:"total_value"`
    OnTimeRate    float64 `json:"on_time_rate"`
    QualityScore  float64 `json:"quality_score"`
    DeliveryScore float64 `json:"delivery_score"`
}

// PendingOrders 待处理订单
type PendingOrders struct {
    InboundPending  int `json:"inbound_pending"`
    OutboundPending int `json:"outbound_pending"`
    TransferPending int `json:"transfer_pending"`
}
```

### 4.2 前端类型定义

```typescript
interface OverviewStats {
  total_inventory: number
  inventory_warning: number
  today_inbound: number
  today_inbound_qty: number
  today_outbound: number
  today_outbound_qty: number
}

interface TrendData {
  date: string
  inbound_qty: number
  outbound_qty: number
}

interface TopProduct {
  product_id: number
  product_name: string
  category: string
  total_qty: number
  order_count: number
}

interface WarehouseUsage {
  warehouse_id: number
  warehouse_name: string
  capacity: number
  used_capacity: number
  usage_rate: number
}

interface SupplierPerformance {
  supplier_id: number
  supplier_name: string
  order_count: number
  total_value: number
  on_time_rate: number
  quality_score: number
  delivery_score: number
}

interface PendingOrders {
  inbound_pending: number
  outbound_pending: number
  transfer_pending: number
}
```

## 5. 前端实现

### 5.1 页面布局

采用经典网格布局，响应式设计：

```
┌─────────────────────────────────────────────────────┐
│ 页面头部：标题 + 时间选择器 + 刷新 + 导出          │
├─────────┬─────────┬─────────┬─────────┤
│ 总库存量 │ 库存预警 │ 今日入库 │ 今日出库 │  统计卡片
├─────────┴─────────┴─────────┴─────────┤
│                                           │
│   出入库趋势图 (折线图)     │ 热销产品 │  第一行图表
│   点击钻取详情             │ 排行 (柱状图)│
│                           │            │
├───────────────────────────┼────────────┤
│                           │            │
│  仓库使用率 (饼图)         │ 供应商绩效 │  第二行图表
│                           │ (雷达图)   │
│                           │            │
└───────────────────────────┴────────────┘
```

### 5.2 组件结构

```
src/pages/dashboard/
├── Dashboard.tsx              # 主页面组件
├── components/
│   ├── StatCard.tsx          # 统计卡片组件
│   ├── TrendChart.tsx        # 出入库趋势图
│   ├── TopProductsChart.tsx  # 热销产品排行图
│   ├── WarehouseUsageChart.tsx    # 仓库使用率图
│   ├── SupplierPerformanceChart.tsx # 供应商绩效图
│   └── DrillDownModal.tsx    # 钻取详情弹窗
├── hooks/
│   └── useDashboardStats.ts  # 数据获取hook
└── types/
    └── index.ts              # 类型定义
```

### 5.3 图表类型与交互

#### 5.3.1 出入库趋势图

**图表类型：** 折线图（双折线）

**配置：**
- X轴：日期（格式 YYYY-MM-DD）
- Y轴：数量
- 双折线：入库（蓝色 #1890ff）、出库（橙色 #fa8c16）
- 支持图例切换显示/隐藏
- 支持tooltip悬浮显示详细数据

**交互：**
- 点击数据点：弹出Modal显示该日期的入库/出库单列表
- 支持缩放和平移（数据量大时）

**钻取逻辑：**
```typescript
// 点击某天数据点
const handlePointClick = (date: string, type: 'inbound' | 'outbound') => {
  // 跳转到订单列表页，带上日期筛选参数
  navigate(`/${type}?date=${date}`)
}
```

#### 5.3.2 热销产品排行

**图表类型：** 横向柱状图

**配置：**
- Y轴：产品名称（显示前10名）
- X轴：出库数量
- 按数量降序排列
- 不同类别用不同颜色区分

**交互：**
- 点击柱状：跳转到产品详情页或库存详情页
- 悬浮显示：产品名称、类别、订单数

#### 5.3.3 仓库使用率

**图表类型：** 饼图/环形图

**配置：**
- 显示各仓库的使用率占比
- 颜色区分不同仓库
- 中心显示总使用率

**交互：**
- 点击扇形：跳转到对应仓库的库存列表
- 悬浮显示：仓库名称、已用容量、总容量

#### 5.3.4 供应商绩效

**图表类型：** 雷达图

**配置：**
- 维度：订单量、总金额、准时率、质量评分、交付评分
- 显示前5名供应商
- 不同供应商用不同颜色

**交互：**
- 点击节点：跳转到供应商详情页
- 悬浮显示：供应商名称、各维度分数

### 5.4 数据刷新与导出

**刷新机制：**
- 手动刷新：用户点击刷新按钮
- 刷新时显示loading状态
- 使用TanStack Query的refetch功能

**导出功能：**
- 支持Excel和PDF两种格式
- 导出内容：统计卡片数据 + 趋势数据 + 排行数据
- 导出时显示进度提示
- 导出完成后自动下载文件

```typescript
const handleExport = async (format: 'excel' | 'pdf') => {
  const response = await api.get(`/dashboard/export`, {
    params: { format, start_date, end_date },
    responseType: 'blob'
  })
  // 创建下载链接
  const url = window.URL.createObjectURL(response.data)
  const link = document.createElement('a')
  link.href = url
  link.download = `dashboard-${start_date}-${end_date}.${format === 'excel' ? 'xlsx' : 'pdf'}`
  link.click()
}
```

### 5.5 时间范围筛选

**默认范围：** 最近30天

**选择器：** Ant Design DatePicker.RangePicker

**交互：**
- 用户选择时间范围后，自动重新请求数据
- 显示loading状态
- 时间范围校验：start_date ≤ end_date，时间跨度 ≤ 365天

```typescript
const [dateRange, setDateRange] = useState<[Dayjs, Dayjs]>([
  dayjs().subtract(30, 'day'),
  dayjs()
])

const handleDateChange = (dates: [Dayjs, Dayjs]) => {
  setDateRange(dates)
  // 触发数据重新获取
  refetch()
}
```

## 6. 后端实现

### 6.1 文件结构

```
warehouse/internal/
├── handler/
│   └── dashboard.go           # DashboardHandler
├── service/
│   └── dashboard.go           # DashboardService
├── repository/
│   └── dashboard.go           # DashboardRepository
├── model/
│   └── dashboard.go           # 数据模型定义
└── router/
    └── router.go              # 路由注册（修改）
```

### 6.2 数据查询逻辑

#### 6.2.1 总库存量

```sql
SELECT COALESCE(SUM(quantity), 0) as total_inventory
FROM inventories;
```

#### 6.2.2 库存预警

```sql
SELECT COUNT(*) as inventory_warning
FROM inventories
WHERE quantity < 10;
```

**说明：** 预警阈值暂定为10，后续可配置化

#### 6.2.3 今日入库/出库统计

```sql
-- 今日入库
SELECT 
  COUNT(*) as order_count,
  COALESCE(SUM(total_qty), 0) as total_qty
FROM inbound_orders
WHERE DATE(created_at) = CURRENT_DATE;

-- 今日出库
SELECT 
  COUNT(*) as order_count,
  COALESCE(SUM(total_qty), 0) as total_qty
FROM outbound_orders
WHERE DATE(created_at) = CURRENT_DATE;
```

#### 6.2.4 出入库趋势

```sql
-- 入库趋势
SELECT 
  DATE(created_at) as date,
  COALESCE(SUM(total_qty), 0) as inbound_qty
FROM inbound_orders
WHERE created_at >= ? AND created_at <= ?
GROUP BY DATE(created_at)
ORDER BY date;

-- 出库趋势
SELECT 
  DATE(created_at) as date,
  COALESCE(SUM(total_qty), 0) as outbound_qty
FROM outbound_orders
WHERE created_at >= ? AND created_at <= ?
GROUP BY DATE(created_at)
ORDER BY date;
```

**实现逻辑：**
- 分别查询入库和出库趋势
- 在Service层合并数据（按日期join）
- 对于没有数据的日期，填充0值

#### 6.2.5 热销产品排行

```sql
SELECT 
  p.id as product_id,
  p.name as product_name,
  COALESCE(c.name, '') as category,
  SUM(oi.quantity) as total_qty,
  COUNT(DISTINCT o.id) as order_count
FROM outbound_items oi
JOIN products p ON oi.product_id = p.id
LEFT JOIN categories c ON p.category_id = c.id
JOIN outbound_orders o ON oi.order_id = o.id
WHERE o.created_at >= ? AND o.created_at <= ?
  AND o.status = 'completed'
GROUP BY p.id, p.name, c.name
ORDER BY total_qty DESC
LIMIT ?;
```

#### 6.2.6 仓库使用率

```sql
SELECT 
  w.id as warehouse_id,
  w.name as warehouse_name,
  COUNT(l.id) as capacity,
  COUNT(DISTINCT CASE WHEN i.quantity > 0 THEN i.location_id END) as used_capacity
FROM warehouses w
LEFT JOIN locations l ON w.id = l.warehouse_id
LEFT JOIN inventories i ON l.id = i.location_id
GROUP BY w.id, w.name;
```

**说明：**
- capacity：仓库总库位数
- used_capacity：有库存的库位数
- usage_rate：在Service层计算 (used_capacity / capacity * 100)

#### 6.2.7 供应商绩效

**注意：** 当前数据模型中可能没有供应商评分字段，需要确认：

**方案1：** 如果有评分字段
```sql
SELECT 
  s.id as supplier_id,
  s.name as supplier_name,
  COUNT(DISTINCT io.id) as order_count,
  COALESCE(SUM(io.total_amount), 0) as total_value,
  AVG(s.on_time_rate) as on_time_rate,
  AVG(s.quality_score) as quality_score,
  AVG(s.delivery_score) as delivery_score
FROM suppliers s
LEFT JOIN inbound_orders io ON s.id = io.supplier_id
WHERE io.created_at >= ? AND io.created_at <= ?
GROUP BY s.id, s.name
ORDER BY order_count DESC, total_value DESC
LIMIT ?;
```

**方案2：** 如果没有评分字段，仅展示订单量和金额
```sql
SELECT 
  s.id as supplier_id,
  s.name as supplier_name,
  COUNT(DISTINCT io.id) as order_count,
  COALESCE(SUM(io.total_amount), 0) as total_value,
  0 as on_time_rate,    -- 暂时填充
  0 as quality_score,   -- 暂时填充
  0 as delivery_score   -- 暂时填充
FROM suppliers s
LEFT JOIN inbound_orders io ON s.id = io.supplier_id
WHERE io.created_at >= ? AND io.created_at <= ?
GROUP BY s.id, s.name
ORDER BY order_count DESC, total_value DESC
LIMIT ?;
```

**实现策略：** 先检查数据模型，根据实际情况选择方案

#### 6.2.8 待处理订单

```sql
-- 待入库
SELECT COUNT(*) as inbound_pending
FROM inbound_orders
WHERE status IN ('pending', 'approved');

-- 待出库
SELECT COUNT(*) as outbound_pending
FROM outbound_orders
WHERE status IN ('pending', 'approved');

-- 待调拨
SELECT COUNT(*) as transfer_pending
FROM stock_transfers
WHERE status = 'pending';
```

### 6.3 性能优化

#### 6.3.1 索引优化

建议添加以下索引以提升查询性能：

```sql
-- 库存表
CREATE INDEX idx_inventories_quantity ON inventories(quantity);
CREATE INDEX idx_inventories_location ON inventories(location_id) WHERE quantity > 0;

-- 入库订单表
CREATE INDEX idx_inbound_orders_created_at ON inbound_orders(created_at);
CREATE INDEX idx_inbound_orders_status ON inbound_orders(status);

-- 出库订单表
CREATE INDEX idx_outbound_orders_created_at ON outbound_orders(created_at);
CREATE INDEX idx_outbound_orders_status ON outbound_orders(status);

-- 出库明细表
CREATE INDEX idx_outbound_items_order_product ON outbound_items(order_id, product_id);

-- 调拨单表
CREATE INDEX idx_stock_transfers_status ON stock_transfers(status);
```

#### 6.3.2 查询优化策略

1. **使用COALESCE处理NULL值**：避免前端处理null值
2. **添加LIMIT限制**：排行数据限制返回条数
3. **使用索引字段**：时间范围查询、状态筛选使用索引字段
4. **避免全表扫描**：WHERE条件使用索引列

#### 6.3.3 并发查询

对于多个独立的统计查询，使用goroutine并发执行：

```go
func (s *DashboardService) GetOverview(ctx context.Context) (*OverviewStats, error) {
    var (
        totalInventory   float64
        inventoryWarning int
        todayInbound     int
        todayOutbound    int
        err              error
    )

    var wg sync.WaitGroup
    var mu sync.Mutex
    errors := make([]error, 0)

    wg.Add(4)
    
    go func() {
        defer wg.Done()
        val, e := s.repo.GetTotalInventory(ctx)
        if e != nil {
            mu.Lock()
            errors = append(errors, e)
            mu.Unlock()
            return
        }
        mu.Lock()
        totalInventory = val
        mu.Unlock()
    }()

    // ... 其他并发查询

    wg.Wait()

    if len(errors) > 0 {
        return nil, errors[0]
    }

    return &OverviewStats{
        TotalInventory: totalInventory,
        // ...
    }, nil
}
```

### 6.4 导出功能

**实现方案：**

使用 Go 的 `excelize` 库生成 Excel，使用 `go-pdf/fpdf` 生成 PDF。

```go
func (s *DashboardService) ExportDashboard(ctx context.Context, format string, startDate, endDate time.Time) ([]byte, error) {
    // 获取所有数据
    overview, _ := s.GetOverview(ctx)
    trend, _ := s.GetTrendData(ctx, startDate, endDate)
    topProducts, _ := s.GetTopProducts(ctx, startDate, endDate, 10)
    // ...

    if format == "excel" {
        return s.generateExcel(overview, trend, topProducts)
    } else if format == "pdf" {
        return s.generatePDF(overview, trend, topProducts)
    }

    return nil, errors.New("unsupported format")
}
```

**Excel格式：**
- Sheet1: 总览统计
- Sheet2: 趋势数据
- Sheet3: 热销产品
- Sheet4: 仓库使用率
- Sheet5: 供应商绩效

**PDF格式：**
- 第一页：总览统计 + 图表截图
- 第二页：详细数据表格

## 7. 测试策略

### 7.1 后端测试

#### 7.1.1 单元测试

```go
// repository/dashboard_test.go
func TestDashboardRepository_GetOverviewStats(t *testing.T) {
    // 准备测试数据
    setupTestData(t)
    
    tests := []struct {
        name    string
        want    *model.OverviewStats
        wantErr bool
    }{
        {
            name: "正常情况",
            want: &model.OverviewStats{
                TotalInventory:   12589.5,
                InventoryWarning: 23,
                TodayInbound:     15,
                TodayOutbound:    18,
            },
            wantErr: false,
        },
        {
            name: "空数据库",
            want: &model.OverviewStats{
                TotalInventory:   0,
                InventoryWarning: 0,
                TodayInbound:     0,
                TodayOutbound:    0,
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := repo.GetOverviewStats(context.Background())
            if (err != nil) != tt.wantErr {
                t.Errorf("GetOverviewStats() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetOverviewStats() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestDashboardRepository_GetTrendData(t *testing.T) {
    tests := []struct {
        name      string
        startDate time.Time
        endDate   time.Time
        wantLen   int
        wantErr   bool
    }{
        {
            name:      "30天范围",
            startDate: time.Now().AddDate(0, 0, -30),
            endDate:   time.Now(),
            wantLen:   30,
            wantErr:   false,
        },
        {
            name:      "无效范围（start > end）",
            startDate: time.Now(),
            endDate:   time.Now().AddDate(0, 0, -30),
            wantLen:   0,
            wantErr:   true,
        },
    }

    // ... 测试实现
}
```

#### 7.1.2 集成测试

```go
// handler/dashboard_test.go
func TestDashboardAPI(t *testing.T) {
    router := setupTestRouter(t)
    
    t.Run("获取总览统计", func(t *testing.T) {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/api/v1/dashboard/overview", nil)
        req.Header.Set("Authorization", "Bearer "+testToken)
        
        router.ServeHTTP(w, req)
        
        assert.Equal(t, 200, w.Code)
        
        var response map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &response)
        assert.NotNil(t, response["data"])
    })

    t.Run("获取趋势数据", func(t *testing.T) {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/api/v1/dashboard/trend?start_date=2026-04-08&end_date=2026-05-08", nil)
        req.Header.Set("Authorization", "Bearer "+testToken)
        
        router.ServeHTTP(w, req)
        
        assert.Equal(t, 200, w.Code)
    })

    t.Run("参数校验", func(t *testing.T) {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/api/v1/dashboard/trend?start_date=invalid", nil)
        req.Header.Set("Authorization", "Bearer "+testToken)
        
        router.ServeHTTP(w, req)
        
        assert.Equal(t, 400, w.Code)
    })
}
```

### 7.2 前端测试

#### 7.2.1 组件测试

```typescript
// Dashboard.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Dashboard from './Dashboard'

describe('Dashboard', () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } }
  })

  const wrapper = ({ children }) => (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>{children}</BrowserRouter>
    </QueryClientProvider>
  )

  it('should render stat cards correctly', async () => {
    render(<Dashboard />, { wrapper })
    
    await waitFor(() => {
      expect(screen.getByText('总库存量')).toBeInTheDocument()
      expect(screen.getByText('库存预警')).toBeInTheDocument()
      expect(screen.getByText('今日入库')).toBeInTheDocument()
      expect(screen.getByText('今日出库')).toBeInTheDocument()
    })
  })

  it('should handle date range change', async () => {
    render(<Dashboard />, { wrapper })
    
    const datePickers = screen.getAllByRole('textbox')
    fireEvent.click(datePickers[0])
    
    // 选择日期范围
    // ...
    
    await waitFor(() => {
      // 验证数据重新加载
    })
  })

  it('should handle refresh button click', async () => {
    render(<Dashboard />, { wrapper })
    
    const refreshButton = screen.getByText('刷新')
    fireEvent.click(refreshButton)
    
    await waitFor(() => {
      // 验证数据重新加载
    })
  })

  it('should navigate to detail page when card clicked', async () => {
    const mockNavigate = jest.fn()
    jest.mock('react-router-dom', () => ({
      ...jest.requireActual('react-router-dom'),
      useNavigate: () => mockNavigate
    }))

    render(<Dashboard />, { wrapper })
    
    const inventoryCard = screen.getByText('总库存量').closest('.stat-card')
    fireEvent.click(inventoryCard)
    
    expect(mockNavigate).toHaveBeenCalledWith('/inventory')
  })

  it('should show drill-down modal when chart clicked', async () => {
    render(<Dashboard />, { wrapper })
    
    // 模拟点击图表数据点
    // ...
    
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument()
    })
  })
})
```

#### 7.2.2 Hook测试

```typescript
// useDashboardStats.test.ts
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useDashboardStats } from './useDashboardStats'

describe('useDashboardStats', () => {
  it('should fetch and cache data', async () => {
    const queryClient = new QueryClient()
    const wrapper = ({ children }) => (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    )

    const { result } = renderHook(() => useDashboardStats(), { wrapper })

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.data).toBeDefined()
  })

  it('should handle loading and error states', async () => {
    // 测试loading状态
    // 测试error状态
  })

  it('should refetch on refresh', async () => {
    const queryClient = new QueryClient()
    const wrapper = ({ children }) => (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    )

    const { result } = renderHook(() => useDashboardStats(), { wrapper })

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    const initialData = result.current.data

    result.current.refetch()

    await waitFor(() => {
      expect(result.current.isFetching).toBe(true)
    })

    await waitFor(() => {
      expect(result.current.isFetching).toBe(false)
    })
  })
})
```

### 7.3 测试覆盖率目标

- **后端：** 核心业务逻辑 80%+ 覆盖率
- **前端：** 组件交互逻辑 70%+ 覆盖率

### 7.4 质量检查清单

- [ ] API响应时间 < 500ms（正常数据量）
- [ ] 图表渲染流畅，无卡顿
- [ ] 移动端响应式布局正常
- [ ] 错误提示友好明确
- [ ] 空数据状态显示合理
- [ ] 导出功能正常工作
- [ ] 钻取跳转逻辑正确
- [ ] 时间范围筛选功能正常
- [ ] 刷新功能正常
- [ ] 并发请求无数据错乱

## 8. 实施计划

### 8.1 开发顺序

**阶段1：后端基础（2-3天）**
1. 创建数据模型和DTO
2. 实现Repository层（数据查询）
3. 实现Service层（业务逻辑）
4. 实现Handler层（API接口）
5. 编写单元测试

**阶段2：前端基础（2-3天）**
1. 创建页面组件结构
2. 实现统计卡片组件
3. 实现数据获取hook
4. 集成API调用
5. 编写组件测试

**阶段3：图表实现（3-4天）**
1. 实现出入库趋势图
2. 实现热销产品排行图
3. 实现仓库使用率图
4. 实现供应商绩效图
5. 实现图表钻取功能

**阶段4：高级功能（2-3天）**
1. 实现时间范围筛选
2. 实现刷新功能
3. 实现导出功能（Excel/PDF）
4. 实现点击跳转详情

**阶段5：测试与优化（2天）**
1. 编写集成测试
2. 性能测试与优化
3. 索引优化
4. 代码审查

**总计：约11-15个工作日**

### 8.2 风险与依赖

**风险：**
1. 数据库性能：大量历史数据可能导致查询变慢
   - 缓解措施：添加索引、限制查询范围
2. 供应商评分字段缺失
   - 缓解措施：先实现基础版本，后续补充评分功能
3. 图表库兼容性问题
   - 缓解措施：提前测试Ant Design Charts与现有技术栈的兼容性

**依赖：**
1. Ant Design Charts库安装
2. 可能需要调整现有数据模型（供应商评分）
3. 数据库索引添加（需要DBA配合）

## 9. 验收标准

### 9.1 功能验收

- [ ] 顶部统计卡片正确显示数据
- [ ] 点击统计卡片可跳转到对应列表页
- [ ] 出入库趋势图正确展示30天数据
- [ ] 趋势图支持时间范围筛选
- [ ] 点击趋势图数据点可钻取查看详情
- [ ] 热销产品排行正确显示TOP 10
- [ ] 仓库使用率饼图正确显示各仓库占比
- [ ] 供应商绩效雷达图正确显示各项指标
- [ ] 刷新按钮可手动刷新数据
- [ ] 导出功能支持Excel和PDF格式
- [ ] 导出的文件内容完整准确

### 9.2 性能验收

- [ ] 总览统计API响应时间 < 200ms
- [ ] 趋势数据API响应时间 < 500ms
- [ ] 图表渲染时间 < 1s
- [ ] 导出功能生成时间 < 5s

### 9.3 用户体验验收

- [ ] 页面布局清晰美观
- [ ] 移动端显示正常
- [ ] 错误提示友好
- [ ] 加载状态清晰
- [ ] 空数据状态合理

## 10. 后续优化方向

### 10.1 短期优化（1-2个月）

1. **缓存层：** 添加Redis缓存统计结果，提升性能
2. **配置化：** 库存预警阈值可配置
3. **权限控制：** 不同角色查看不同指标
4. **数据对比：** 支持同比、环比对比

### 10.2 长期优化（3-6个月）

1. **实时推送：** 使用WebSocket推送实时数据
2. **自定义仪表盘：** 用户可自定义显示的指标和图表
3. **移动端优化：** 开发专属移动端仪表盘
4. **AI洞察：** 基于数据的智能分析和建议

## 11. 附录

### 11.1 参考资料

- [Ant Design Charts 文档](https://charts.ant.design/)
- [TanStack Query 文档](https://tanstack.com/query/latest)
- [Excelize 库文档](https://xuri.me/excelize/)
- [Go PDF FPDF 库](https://github.com/go-pdf/fpdf)

### 11.2 变更历史

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-05-08 | 1.0 | 初始版本 | OpenCode |
