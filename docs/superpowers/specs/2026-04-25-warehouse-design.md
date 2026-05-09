# 仓库管理系统设计文档

## 概述

面向中小企业的实体仓库管理系统，支持商品库存管理、出入库操作、货位管理、供应商客户管理，以及完整的操作审计功能。

**技术栈**：
- 后端：Go + Gin + Bun ORM + MySQL/MariaDB
- 前端：React + TypeScript + Ant Design + Vite
- 架构：嵌入式SPA，前后端融合部署

---

## 1. 系统架构

### 1.1 架构图

```
┌─────────────────────────────────────────────────┐
│                  Go Binary                       │
│  ┌─────────────────┐   ┌──────────────────────┐ │
│  │  Gin Router     │──▶│  Static File Server  │ │
│  └────────┬────────┘   │  (embedded React)    │ │
│           │            └──────────────────────┘ │
│           ▼                                     │
│  ┌─────────────────┐                           │
│  │  API Handlers   │                           │
│  └────────┬────────┘                           │
│           ▼                                     │
│  ┌─────────────────┐   ┌──────────────────────┐│
│  │  Service Layer  │──▶│  Audit Logger        ││
│  └────────┬────────┘   └──────────────────────┘│
│           ▼                                     │
│  ┌─────────────────┐                           │
│  │  Bun ORM        │                           │
│  └────────┬────────┘                           │
│           ▼                                     │
│  ┌─────────────────┐                           │
│  │  MySQL/MariaDB  │                           │
│  └─────────────────┘                           │
└─────────────────────────────────────────────────┘
```

### 1.2 分层架构

| 层级 | 职责 |
|------|------|
| Router层 | Gin路由，处理静态资源和API请求分发 |
| Handler层 | HTTP请求/响应处理，参数验证 |
| Service层 | 业务逻辑，事务管理，审计记录 |
| Repository层 | Bun ORM封装，数据库操作 |
| Model层 | 数据模型定义 |

### 1.3 请求处理流程

```
请求 → CORS中间件 → JWT认证中间件 → RBAC权限中间件 → 审计中间件 → Handler → Service → Repository → 数据库
```

---

## 2. 数据模型设计

### 2.1 基础模型

所有业务表继承此基础模型：

```go
type BaseModel struct {
    ID        int64     `bun:"id,pk,autoincrement"`
    CreatedAt time.Time `bun:"created_at,notnull"`
    CreatedBy int64     `bun:"created_by,notnull"`
    UpdatedAt time.Time `bun:"updated_at,notnull"`
    UpdatedBy int64     `bun:"updated_by,notnull"`
    DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero"`
}
```

### 2.2 用户权限模块

| 表名 | 说明 |
|------|------|
| users | 用户表（用户名、密码hash、状态） |
| roles | 角色表（角色名、描述、状态） |
| permissions | 权限表（权限码、名称、资源、动作） |
| user_roles | 用户-角色关联表 |
| role_permissions | 角色-权限关联表 |

### 2.3 核心业务模块

| 表名 | 说明 |
|------|------|
| warehouses | 仓库表（名称、地址、状态） |
| locations | 货位表（仓库ID、区域、货架、层级、位置） |
| categories | 商品分类表（名称、父分类ID） |
| products | 商品表（SKU、名称、分类ID、规格、单位） |
| inventory | 库存表（仓库ID、商品ID、货位ID、数量、批次号） |
| suppliers | 供应商表（名称、联系人、电话、地址） |
| customers | 客户表（名称、联系人、电话、地址） |

### 2.4 出入库模块

| 表名 | 说明 |
|------|------|
| inbound_orders | 入库单（单号、供应商ID、仓库ID、状态、总数量） |
| inbound_items | 入库明细（入库单ID、商品ID、货位ID、数量、批次号） |
| outbound_orders | 出库单（单号、客户ID、仓库ID、状态、总数量） |
| outbound_items | 出库明细（出库单ID、商品ID、货位ID、数量） |
| stock_transfers | 调拨单（单号、源仓库、目标仓库、状态） |

### 2.5 审计模块

| 表名 | 说明 |
|------|------|
| audit_logs | 审计日志表 |

**audit_logs 表结构**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 主键 |
| table_name | varchar(64) | 表名 |
| record_id | int64 | 记录ID |
| action | varchar(16) | 操作类型（create/update/delete） |
| old_value | json | 修改前值 |
| new_value | json | 修改后值 |
| operated_by | int64 | 操作人ID |
| operated_at | datetime | 操作时间 |
| ip_address | varchar(45) | IP地址 |

---

## 3. 审计系统设计

### 3.1 审计触发机制

在Service层实现审计拦截：

```go
type AuditService struct {
    service   interface{}
    auditRepo *AuditRepository
}
```

### 3.2 审计触发时机

| 操作 | 触发点 | 记录内容 |
|------|--------|----------|
| 创建 | Create后 | new_value = 新记录 |
| 更新 | Update前查询旧值 | old_value, new_value |
| 删除 | 软删除前 | old_value = 删除前记录 |
| 批量操作 | 逐条记录 | 每条记录独立审计 |

### 3.3 审计查询

- 按 table_name + record_id 查询某记录的所有变更历史
- 按 operated_by 查询某用户的所有操作
- 按 operated_at 时间范围查询

---

## 4. API设计

### 4.1 通用规范

```
基础路径: /api/v1
认证方式: JWT Token (Header: Authorization: Bearer <token>)
响应格式: JSON

统一响应结构:
{
  "code": 0,
  "message": "success",
  "data": {}
}

分页响应:
{
  "code": 0,
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 4.2 API列表

**认证模块**
```
POST   /auth/login          - 登录
POST   /auth/logout         - 登出
GET    /auth/profile        - 当前用户信息
PUT    /auth/password       - 修改密码
```

**用户权限模块**
```
GET    /users               - 用户列表
POST   /users               - 创建用户
PUT    /users/:id           - 更新用户
DELETE /users/:id           - 删除用户

GET    /roles               - 角色列表
POST   /roles               - 创建角色
PUT    /roles/:id           - 更新角色

GET    /permissions         - 权限列表
```

**仓库货位模块**
```
GET    /warehouses          - 仓库列表
POST   /warehouses          - 创建仓库
PUT    /warehouses/:id      - 更新仓库
DELETE /warehouses/:id      - 删除仓库

GET    /locations           - 货位列表
POST   /locations           - 创建货位
PUT    /locations/:id       - 更新货位
DELETE /locations/:id       - 删除货位
```

**商品库存模块**
```
GET    /categories          - 分类列表
POST   /categories          - 创建分类
PUT    /categories/:id      - 更新分类

GET    /products            - 商品列表
POST   /products            - 创建商品
PUT    /products/:id        - 更新商品
DELETE /products/:id        - 删除商品

GET    /inventory           - 库存查询
POST   /inventory/check     - 库存盘点
```

**出入库模块**
```
GET    /inbound-orders      - 入库单列表
POST   /inbound-orders      - 创建入库单
PUT    /inbound-orders/:id  - 更新入库单
POST   /inbound-orders/:id/confirm - 确认入库

GET    /outbound-orders     - 出库单列表
POST   /outbound-orders     - 创建出库单
PUT    /outbound-orders/:id - 更新出库单
POST   /outbound-orders/:id/confirm - 确认出库

GET    /stock-transfers     - 调拨单列表
POST   /stock-transfers     - 创建调拨单
```

**供应商客户模块**
```
GET    /suppliers           - 供应商列表
POST   /suppliers           - 创建供应商
PUT    /suppliers/:id       - 更新供应商

GET    /customers           - 客户列表
POST   /customers           - 创建客户
PUT    /customers/:id       - 更新客户
```

**审计模块**
```
GET    /audit-logs          - 审计日志列表
GET    /audit-logs/:id      - 审计日志详情
```

---

## 5. 前端架构

### 5.1 技术栈

| 技术 | 版本/说明 |
|------|-----------|
| React | 18 |
| TypeScript | - |
| Vite | 构建工具 |
| React Router | v6 |
| Ant Design | 5.x |
| Axios | HTTP请求 |
| Zustand | 状态管理 |
| React Query | 服务端状态缓存 |

### 5.2 项目结构

```
web/
├── src/
│   ├── api/              # API请求封装
│   ├── components/       # 通用组件
│   │   ├── Layout/       # 布局组件
│   │   ├── Table/        # 表格封装
│   │   └── Form/         # 表单封装
│   ├── pages/            # 页面组件
│   │   ├── auth/         # 登录、修改密码
│   │   ├── dashboard/    # 仪表盘
│   │   ├── warehouse/    # 仓库、货位管理
│   │   ├── product/      # 商品、分类、库存
│   │   ├── order/        # 出入库单
│   │   ├── partner/      # 供应商、客户
│   │   ├── system/       # 用户、角色、权限
│   │   └── audit/        # 审计日志
│   ├── hooks/            # 自定义Hooks
│   ├── stores/           # Zustand状态
│   ├── types/            # TypeScript类型定义
│   ├── utils/            # 工具函数
│   └── App.tsx
├── public/
├── vite.config.ts
└── package.json
```

### 5.3 权限控制

- 路由守卫：根据权限动态注册路由
- 按钮级权限：`hasPermission()`函数判断

### 5.4 审计日志展示

- 根据table_name映射到业务名称
- 点击record_id跳转到对应业务详情页
- old_value/new_value以表格对比形式展示字段变更

### 5.5 静态资源嵌入

```go
//go:embed all:web/dist
var staticFS embed.FS
```

---

## 6. RBAC权限系统

### 6.1 权限模型

```
用户 ──M:N──▶ 角色 ──M:N──▶ 权限
```

### 6.2 权限粒度

资源-动作模式：`资源:动作`

**资源**：商品、仓库、入库单、出库单、供应商、客户、用户、角色、审计日志

**动作**：list, create, update, delete, export

**权限示例**：
- product:list - 查看商品列表
- product:create - 创建商品
- product:update - 更新商品
- product:delete - 删除商品
- inbound:confirm - 确认入库
- audit:export - 导出审计日志

### 6.3 预置角色

| 角色 | 权限范围 |
|------|----------|
| 超级管理员 | 所有权限 |
| 仓库管理员 | 仓库、货位、库存、出入库管理 |
| 采购员 | 供应商、入库单 |
| 销售员 | 客户、出库单 |
| 库管员 | 库存查询、盘点 |
| 审计员 | 审计日志查看 |

### 6.4 权限校验流程

```
请求 → JWT解析用户ID → 查询用户角色 → 查询角色权限 → 匹配API权限码 → 放行/拒绝
```

---

## 7. 部署与配置

### 7.1 配置管理

```yaml
# config.yaml
server:
  port: 8080
  mode: release

database:
  driver: mysql
  host: localhost
  port: 3306
  name: warehouse
  user: root
  password: ****
  max_open_conns: 100
  max_idle_conns: 10

jwt:
  secret: your-secret-key
  expire: 24h

log:
  level: info
  output: stdout
  file: ./logs/app.log
```

### 7.2 部署方式

**单机部署（初期）**：
```bash
go build -o warehouse
./warehouse --config config.yaml
```

**前后端分离部署（后期）**：
```
Nginx → React静态资源
     → 反向代理 /api/* → Go服务
```

### 7.3 项目目录结构

```
warehouse/
├── cmd/
│   └── server/main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── model/
│   ├── middleware/
│   └── pkg/
├── pkg/
├── migrations/
├── config/
│   └── config.yaml
├── web/
│   ├── src/
│   ├── dist/
│   └── package.json
├── embed.go
├── go.mod
└── Makefile
```

---

## 8. 错误处理与日志

### 8.1 错误码设计

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 内部错误 |
| 1001+ | 业务错误码 |

### 8.2 日志规范

```
[时间] [级别] [trace_id] [user_id] 消息 字段...
```

- 请求入口生成trace_id，全链路传递
- 敏感信息脱敏（密码、token等）

---

## 9. 测试策略

### 9.1 后端测试

- 单元测试：Service层、工具函数
- 集成测试：API端到端测试
- 覆盖率目标：核心业务 > 80%

### 9.2 前端测试

- 组件测试：React Testing Library
- E2E测试：Playwright（关键流程）

### 9.3 测试数据

- 使用fixtures或factory模式生成
- 测试数据库独立，每次测试前清理
