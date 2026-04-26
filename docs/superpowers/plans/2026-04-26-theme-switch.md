# 主题切换功能实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为仓库管理系统添加深色/浅色主题切换功能，用户可在用户头像下拉菜单中切换，主题偏好保存到服务器端。

**Architecture:** 后端在users表添加theme字段，前端使用CSS变量定义两套主题，通过ConfigProvider和data-theme属性实现主题切换。

**Tech Stack:** Go + Gin + Bun ORM (后端), React + TypeScript + Ant Design 5 + Zustand (前端)

---

## 文件结构

### 后端文件变更
- `internal/model/user.go` - 添加 Theme 字段
- `internal/handler/auth.go` - 添加 UpdateTheme handler
- `internal/service/auth.go` - 添加 UpdateTheme 方法
- `internal/repository/user.go` - 添加 UpdateTheme 方法

### 前端文件变更
- `web/src/styles/theme.css` - 新建主题CSS变量文件
- `web/src/stores/authStore.ts` - 添加 theme 状态和方法
- `web/src/hooks/useTheme.ts` - 新建主题hook
- `web/src/api/auth.ts` - 添加 updateTheme API
- `web/src/types/auth.ts` - 添加 Theme 字段到 User 类型
- `web/src/components/Layout/MainLayout.tsx` - 添加主题切换菜单项
- `web/src/App.tsx` - 添加 ConfigProvider 主题配置
- `web/src/main.tsx` - 导入 theme.css

---

## Task 1: 数据库迁移 - 添加 theme 字段

**Files:**
- 数据库: `warehouse.users` 表

- [ ] **Step 1: 执行数据库迁移**

```bash
mysql -h 192.168.1.13 -u devuser -p'123456' warehouse -e "ALTER TABLE users ADD COLUMN theme VARCHAR(20) NOT NULL DEFAULT 'light';"
```

- [ ] **Step 2: 验证字段已添加**

```bash
mysql -h 192.168.1.13 -u devuser -p'123456' warehouse -e "DESCRIBE users;" | grep theme
```

Expected: 输出包含 `theme` 字段

---

## Task 2: 后端 Model - 添加 Theme 字段

**Files:**
- Modify: `internal/model/user.go:8-13`

- [ ] **Step 1: 修改 User struct 添加 Theme 字段**

```go
type User struct {
	BaseModel
	Username string `bun:"username,unique" json:"username"`
	Password string `bun:"password_hash" json:"-"`
	Status   int    `bun:"status" json:"status"`
	Theme    string `bun:"theme" json:"theme"`
}
```

- [ ] **Step 2: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse && go build ./...
```

Expected: 无错误

- [ ] **Step 3: Commit**

```bash
git add internal/model/user.go
git commit -m "feat(model): add theme field to User"
```

---

## Task 3: 后端 Repository - 添加 UpdateTheme 方法

**Files:**
- Modify: `internal/repository/user.go`

- [ ] **Step 1: 在 UserRepository 添加 UpdateTheme 方法**

在 `internal/repository/user.go` 文件末尾添加:

```go
func (r *UserRepository) UpdateTheme(ctx context.Context, userID int64, theme string) error {
	_, err := r.db.NewUpdate().
		Model((*model.User)(nil)).
		Set("theme = ?", theme).
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
```

- [ ] **Step 2: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse && go build ./...
```

Expected: 无错误

- [ ] **Step 3: Commit**

```bash
git add internal/repository/user.go
git commit -m "feat(repository): add UpdateTheme method"
```

---

## Task 4: 后端 Service - 添加 UpdateTheme 方法

**Files:**
- Modify: `internal/service/auth.go`

- [ ] **Step 1: 在 UserRepository interface 添加方法签名**

修改 `internal/service/auth.go` 中的 `UserRepository` interface:

```go
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdateTheme(ctx context.Context, userID int64, theme string) error
}
```

- [ ] **Step 2: 在 AuthService 添加 UpdateTheme 方法**

在 `internal/service/auth.go` 文件末尾添加:

```go
func (s *AuthService) UpdateTheme(ctx context.Context, userID int64, theme string) error {
	if theme != "light" && theme != "dark" {
		return apperrors.NewAppError(apperrors.CodeBadRequest, "invalid theme value")
	}
	return s.userRepo.UpdateTheme(ctx, userID, theme)
}
```

- [ ] **Step 3: 在 AuthService interface 添加方法签名**

修改 `internal/service/auth.go`，如果 AuthService interface 定义在外部则在对应位置添加；否则在文件顶部注释或直接添加方法即可。

- [ ] **Step 4: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse && go build ./...
```

Expected: 无错误

- [ ] **Step 5: Commit**

```bash
git add internal/service/auth.go
git commit -m "feat(service): add UpdateTheme method to AuthService"
```

---

## Task 5: 后端 Handler - 添加 UpdateTheme 接口

**Files:**
- Modify: `internal/handler/auth.go`

- [ ] **Step 1: 添加请求结构体**

在 `internal/handler/auth.go` 的 `ChangePasswordRequest` 后面添加:

```go
type UpdateThemeRequest struct {
	Theme string `json:"theme" binding:"required"`
}
```

- [ ] **Step 2: 在 AuthService interface 添加方法签名**

修改 `internal/handler/auth.go` 中的 `AuthService` interface:

```go
type AuthService interface {
	Login(ctx context.Context, username, password string) (string, *model.User, error)
	GetProfile(ctx context.Context, userID int64) (*model.User, error)
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	UpdateTheme(ctx context.Context, userID int64, theme string) error
}
```

- [ ] **Step 3: 添加 UpdateTheme handler 方法**

在 `ChangePassword` 方法后添加:

```go
func (h *AuthHandler) UpdateTheme(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, apperrors.CodeUnauthorized, "user not authenticated")
		return
	}

	var req UpdateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	err := h.authService.UpdateTheme(c.Request.Context(), userID, req.Theme)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, nil)
}
```

- [ ] **Step 4: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse && go build ./...
```

Expected: 无错误

- [ ] **Step 5: Commit**

```bash
git add internal/handler/auth.go
git commit -m "feat(handler): add UpdateTheme endpoint"
```

---

## Task 6: 后端路由注册

**Files:**
- Modify: `internal/router/router.go`

- [ ] **Step 1: 查找路由注册位置**

```bash
grep -n "PUT.*password" /home/zzf/projects/goinvent/warehouse/internal/router/router.go
```

- [ ] **Step 2: 添加路由**

在 `auth.PUT("/password", authHandler.ChangePassword)` 后面添加:

```go
auth.PUT("/theme", authHandler.UpdateTheme)
```

- [ ] **Step 3: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse && go build ./...
```

Expected: 无错误

- [ ] **Step 4: Commit**

```bash
git add internal/router/router.go
git commit -m "feat(router): register UpdateTheme route"
```

---

## Task 7: 后端编译测试

**Files:**
- 无文件变更

- [ ] **Step 1: 完整编译**

```bash
cd /home/zzf/projects/goinvent/warehouse && go build -o bin/warehouse ./cmd/server
```

Expected: 编译成功

- [ ] **Step 2: 启动服务器测试**

```bash
cd /home/zzf/projects/goinvent/warehouse && ./bin/warehouse &
sleep 3
TOKEN=$(curl -s -X POST http://127.0.0.1:8080/api/v1/auth/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}' | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token: ${TOKEN:0:20}..."
```

- [ ] **Step 3: 测试 UpdateTheme API**

```bash
curl -s -X PUT http://127.0.0.1:8080/api/v1/auth/theme \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"theme":"dark"}'
```

Expected: `{"code":0,"message":"success","data":null}`

- [ ] **Step 4: 测试 GetProfile 返回 theme**

```bash
curl -s http://127.0.0.1:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN" | grep -o '"theme":"[^"]*"'
```

Expected: `"theme":"dark"`

- [ ] **Step 5: 停止服务器**

```bash
pkill -f "bin/warehouse"
```

---

## Task 8: 前端类型定义

**Files:**
- Modify: `web/src/types/auth.ts`

- [ ] **Step 1: 查看 User 类型定义**

```bash
cat /home/zzf/projects/goinvent/warehouse/web/src/types/auth.ts
```

- [ ] **Step 2: 添加 Theme 字段到 User 类型**

在 User interface 中添加 theme 字段:

```typescript
export interface User {
  id: number
  username: string
  status: number
  theme: 'light' | 'dark'
  created_at: string
  updated_at: string
}
```

- [ ] **Step 3: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse/web && npm run build 2>&1 | tail -10
```

Expected: 无类型错误

- [ ] **Step 4: Commit**

```bash
git add web/src/types/auth.ts
git commit -m "feat(types): add theme field to User type"
```

---

## Task 9: 前端 API - 添加 updateTheme 方法

**Files:**
- Modify: `web/src/api/auth.ts`

- [ ] **Step 1: 添加 updateTheme API 方法**

在 `web/src/api/auth.ts` 的 `authApi` 对象中添加:

```typescript
export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', data)
    return response.data
  },

  getProfile: async (): Promise<User> => {
    const response = await api.get<User>('/auth/profile')
    return response.data
  },

  changePassword: async (data: ChangePasswordRequest): Promise<void> => {
    await api.put('/auth/password', data)
  },

  updateTheme: async (theme: 'light' | 'dark'): Promise<void> => {
    await api.put('/auth/theme', { theme })
  },
}
```

- [ ] **Step 2: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse/web && npm run build 2>&1 | tail -10
```

Expected: 无错误

- [ ] **Step 3: Commit**

```bash
git add web/src/api/auth.ts
git commit -m "feat(api): add updateTheme method"
```

---

## Task 10: 前端 CSS 主题变量

**Files:**
- Create: `web/src/styles/theme.css`

- [ ] **Step 1: 创建样式目录**

```bash
mkdir -p /home/zzf/projects/goinvent/warehouse/web/src/styles
```

- [ ] **Step 2: 创建 theme.css 文件**

```css
:root[data-theme='light'] {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f5;
  --bg-tertiary: #fafafa;
  --bg-container: #ffffff;
  --text-primary: #333333;
  --text-secondary: #666666;
  --text-tertiary: #999999;
  --border-color: #f0f0f0;
  --shadow-color: rgba(0, 0, 0, 0.08);
}

:root[data-theme='dark'] {
  --bg-primary: #1f1f1f;
  --bg-secondary: #141414;
  --bg-tertiary: #262626;
  --bg-container: #1f1f1f;
  --text-primary: #ffffff;
  --text-secondary: rgba(255, 255, 255, 0.85);
  --text-tertiary: rgba(255, 255, 255, 0.45);
  --border-color: #303030;
  --shadow-color: rgba(0, 0, 0, 0.45);
}

:root[data-theme='dark'] .ant-layout-sider {
  background: var(--bg-secondary) !important;
}

:root[data-theme='dark'] .ant-layout-header {
  background: var(--bg-primary) !important;
}

:root[data-theme='dark'] .ant-layout-content {
  background: var(--bg-tertiary) !important;
}

:root[data-theme='dark'] .ant-menu-light {
  background: var(--bg-secondary) !important;
}

:root[data-theme='dark'] .ant-menu-item {
  color: var(--text-primary) !important;
}

:root[data-theme='dark'] .ant-menu-item:hover {
  color: #1890ff !important;
  background: rgba(255, 255, 255, 0.08) !important;
}

:root[data-theme='dark'] .ant-menu-item-selected {
  background: #1890ff !important;
  color: #fff !important;
}

:root[data-theme='dark'] .ant-table {
  background: var(--bg-container) !important;
}

:root[data-theme='dark'] .ant-table-thead > tr > th {
  background: var(--bg-tertiary) !important;
  color: var(--text-primary) !important;
  border-color: var(--border-color) !important;
}

:root[data-theme='dark'] .ant-table-tbody > tr > td {
  background: var(--bg-container) !important;
  color: var(--text-primary) !important;
  border-color: var(--border-color) !important;
}

:root[data-theme='dark'] .ant-table-tbody > tr:hover > td {
  background: var(--bg-tertiary) !important;
}

:root[data-theme='dark'] .ant-modal-content {
  background: var(--bg-container) !important;
}

:root[data-theme='dark'] .ant-modal-header {
  background: var(--bg-container) !important;
  border-color: var(--border-color) !important;
}

:root[data-theme='dark'] .ant-modal-title {
  color: var(--text-primary) !important;
}

:root[data-theme='dark'] .ant-modal-body {
  background: var(--bg-container) !important;
  color: var(--text-primary) !important;
}

:root[data-theme='dark'] .ant-modal-footer {
  background: var(--bg-container) !important;
  border-color: var(--border-color) !important;
}

:root[data-theme='dark'] .ant-input,
:root[data-theme='dark'] .ant-input-affix-wrapper,
:root[data-theme='dark'] .ant-select-selector,
:root[data-theme='dark'] .ant-picker {
  background: var(--bg-tertiary) !important;
  border-color: var(--border-color) !important;
  color: var(--text-primary) !important;
}

:root[data-theme='dark'] .ant-form-item-label > label {
  color: var(--text-primary) !important;
}

:root[data-theme='dark'] .ant-card {
  background: var(--bg-container) !important;
  border-color: var(--border-color) !important;
}

:root[data-theme='dark'] .ant-card-head {
  background: var(--bg-container) !important;
  border-color: var(--border-color) !important;
}

:root[data-theme='dark'] .ant-card-head-title {
  color: var(--text-primary) !important;
}
```

- [ ] **Step 3: Commit**

```bash
git add web/src/styles/theme.css
git commit -m "feat(styles): add theme CSS variables"
```

---

## Task 11: 前端 Store - 添加 theme 状态

**Files:**
- Modify: `web/src/stores/authStore.ts`

- [ ] **Step 1: 修改 authStore 添加 theme 状态**

完整替换 `web/src/stores/authStore.ts`:

```typescript
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types/auth'

interface AuthState {
  user: User | null
  token: string | null
  theme: 'light' | 'dark'
  isAuthenticated: boolean
  login: (token: string, user: User) => void
  logout: () => void
  setUser: (user: User) => void
  setTheme: (theme: 'light' | 'dark') => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      theme: 'light',
      isAuthenticated: false,
      login: (token: string, user: User) => {
        set({ 
          token, 
          user, 
          isAuthenticated: true,
          theme: user.theme || 'light'
        })
      },
      logout: () => {
        set({ token: null, user: null, isAuthenticated: false, theme: 'light' })
      },
      setUser: (user: User) => {
        set({ user })
      },
      setTheme: (theme: 'light' | 'dark') => {
        set({ theme })
      },
    }),
    {
      name: 'auth-storage',
    }
  )
)
```

- [ ] **Step 2: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse/web && npm run build 2>&1 | tail -10
```

Expected: 无错误

- [ ] **Step 3: Commit**

```bash
git add web/src/stores/authStore.ts
git commit -m "feat(store): add theme state to authStore"
```

---

## Task 12: 前端 Hook - 创建 useTheme

**Files:**
- Create: `web/src/hooks/useTheme.ts`

- [ ] **Step 1: 创建 hooks 目录**

```bash
mkdir -p /home/zzf/projects/goinvent/warehouse/web/src/hooks
```

- [ ] **Step 2: 创建 useTheme.ts**

```typescript
import { useAuthStore } from '../stores/authStore'
import { authApi } from '../api/auth'

export function useTheme() {
  const { theme, setTheme } = useAuthStore()

  const applyTheme = (newTheme: 'light' | 'dark') => {
    document.documentElement.setAttribute('data-theme', newTheme)
  }

  const toggleTheme = async (newTheme: 'light' | 'dark') => {
    applyTheme(newTheme)
    setTheme(newTheme)
    try {
      await authApi.updateTheme(newTheme)
    } catch (error) {
      console.error('Failed to save theme preference:', error)
    }
  }

  const initTheme = (savedTheme: 'light' | 'dark') => {
    applyTheme(savedTheme)
    setTheme(savedTheme)
  }

  return { theme, toggleTheme, initTheme, applyTheme }
}
```

- [ ] **Step 3: Commit**

```bash
git add web/src/hooks/useTheme.ts
git commit -m "feat(hooks): add useTheme hook"
```

---

## Task 13: 前端 App - 导入主题样式和配置 ConfigProvider

**Files:**
- Modify: `web/src/App.tsx`
- Modify: `web/src/main.tsx`

- [ ] **Step 1: 在 main.tsx 导入 theme.css**

在 `web/src/main.tsx` 顶部添加导入:

```typescript
import './styles/theme.css'
```

- [ ] **Step 2: 修改 App.tsx 添加主题配置**

完整替换 `web/src/App.tsx`:

```typescript
import { useEffect } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider, theme as antdTheme } from 'antd'
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
import { UserList } from './pages/system/UserList'
import { RoleList } from './pages/system/RoleList'
import { PermissionList } from './pages/system/PermissionList'
import { AuditLogList } from './pages/system/AuditLogList'
import { useAuthStore } from './stores/authStore'
import { useTheme } from './hooks/useTheme'

function Dashboard() {
  return <div>Dashboard</div>
}

function App() {
  const { theme } = useAuthStore()
  const { applyTheme } = useTheme()

  useEffect(() => {
    applyTheme(theme)
  }, [theme])

  return (
    <ConfigProvider 
      locale={zhCN}
      theme={{
        algorithm: theme === 'dark' ? antdTheme.darkAlgorithm : antdTheme.defaultAlgorithm,
      }}
    >
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
            <Route path="users" element={<UserList />} />
            <Route path="roles" element={<RoleList />} />
            <Route path="permissions" element={<PermissionList />} />
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
            <Route path="audit-logs" element={<AuditLogList />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
```

- [ ] **Step 3: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse/web && npm run build 2>&1 | tail -10
```

Expected: 无错误

- [ ] **Step 4: Commit**

```bash
git add web/src/App.tsx web/src/main.tsx
git commit -m "feat(app): add theme config provider and import theme css"
```

---

## Task 14: 前端 MainLayout - 添加主题切换菜单

**Files:**
- Modify: `web/src/components/Layout/MainLayout.tsx`

- [ ] **Step 1: 添加图标导入**

在文件顶部的 import 中添加:

```typescript
import {
  // ... 现有图标
  SunOutlined,
  MoonOutlined,
} from '@ant-design/icons'
```

- [ ] **Step 2: 添加 useTheme hook 导入**

```typescript
import { useTheme } from '../../hooks/useTheme'
```

- [ ] **Step 3: 在组件中使用 useTheme**

在 `MainLayout` 函数内添加:

```typescript
const { theme, toggleTheme } = useTheme()
```

- [ ] **Step 4: 修改 userMenuItems 添加主题切换选项**

修改 `userMenuItems` 数组，在 `change-password` 之前添加:

```typescript
const userMenuItems = [
  {
    key: 'theme',
    icon: theme === 'dark' ? <SunOutlined /> : <MoonOutlined />,
    label: theme === 'dark' ? '切换到浅色模式' : '切换到深色模式',
    onClick: () => toggleTheme(theme === 'dark' ? 'light' : 'dark'),
  },
  {
    key: 'change-password',
    icon: <LockOutlined />,
    label: '修改密码',
    onClick: () => navigate('/change-password'),
  },
  {
    key: 'logout',
    icon: <LogoutOutlined />,
    label: '退出登录',
    onClick: handleLogout,
  },
]
```

- [ ] **Step 5: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse/web && npm run build 2>&1 | tail -10
```

Expected: 无错误

- [ ] **Step 6: Commit**

```bash
git add web/src/components/Layout/MainLayout.tsx
git commit -m "feat(layout): add theme toggle menu item"
```

---

## Task 15: ProtectedRoute - 初始化主题

**Files:**
- Modify: `web/src/components/ProtectedRoute.tsx`

- [ ] **Step 1: 查看 ProtectedRoute 当前实现**

```bash
cat /home/zzf/projects/goinvent/warehouse/web/src/components/ProtectedRoute.tsx
```

- [ ] **Step 2: 添加主题初始化逻辑**

在 ProtectedRoute 组件中，当获取到用户信息后初始化主题:

```typescript
import { Navigate, Outlet } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import { useTheme } from '../hooks/useTheme'
import { authApi } from '../api/auth'
import { useEffect, useState } from 'react'

export function ProtectedRoute() {
  const { token, isAuthenticated, setUser, user } = useAuthStore()
  const { initTheme, theme } = useTheme()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchProfile = async () => {
      if (token && !user) {
        try {
          const profile = await authApi.getProfile()
          setUser(profile)
          initTheme(profile.theme || 'light')
        } catch (error) {
          console.error('Failed to fetch profile:', error)
        }
      }
      setLoading(false)
    }
    fetchProfile()
  }, [token, user, setUser, initTheme])

  if (!token || !isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (loading) {
    return <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>Loading...</div>
  }

  return <Outlet />
}
```

- [ ] **Step 3: 编译验证**

```bash
cd /home/zzf/projects/goinvent/warehouse/web && npm run build 2>&1 | tail -10
```

Expected: 无错误

- [ ] **Step 4: Commit**

```bash
git add web/src/components/ProtectedRoute.tsx
git commit -m "feat(protected-route): initialize theme from user profile"
```

---

## Task 16: 完整编译和测试

**Files:**
- 无文件变更

- [ ] **Step 1: 编译后端**

```bash
cd /home/zzf/projects/goinvent/warehouse && go build -o bin/warehouse ./cmd/server
```

Expected: 编译成功

- [ ] **Step 2: 编译前端**

```bash
cd /home/zzf/projects/goinvent/warehouse/web && npm run build
```

Expected: 编译成功

- [ ] **Step 3: 启动服务器**

```bash
cd /home/zzf/projects/goinvent/warehouse && ./bin/warehouse &
sleep 3
ss -tlnp | grep 8080
```

Expected: 服务器启动成功

- [ ] **Step 4: 手动测试**

1. 访问 http://192.168.1.13:8080
2. 登录 admin/admin123
3. 点击右上角用户头像
4. 点击"切换到深色模式"
5. 验证主题切换成功
6. 刷新页面，验证主题保持
7. 退出登录后重新登录，验证主题保持

- [ ] **Step 5: 停止服务器**

```bash
pkill -f "bin/warehouse"
```

---

## Task 17: 最终 Commit 和合并

**Files:**
- 无文件变更

- [ ] **Step 1: 推送分支到远程**

```bash
cd /home/zzf/projects/goinvent/warehouse
git push -u origin feature/theme-switch
```

- [ ] **Step 2: 合并到 master**

```bash
git checkout master
git merge feature/theme-switch
git push
```

---

## 验收标准

- [ ] 用户可以在头像下拉菜单中切换主题
- [ ] 主题切换即时生效，无需刷新页面
- [ ] 主题偏好保存到服务器端
- [ ] 刷新页面后主题保持
- [ ] 退出登录后重新登录，主题保持
- [ ] 深色主题下所有组件显示正常
- [ ] 所有 API 测试通过
