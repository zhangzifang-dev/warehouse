# 主题切换功能设计文档

## 概述

为仓库管理系统添加深色/浅色主题切换功能，用户可在用户头像下拉菜单中切换，主题偏好保存到服务器端，下次登录后自动应用上次的主题设置。

## 数据库设计

### users表新增字段

```sql
ALTER TABLE users ADD COLUMN theme VARCHAR(20) NOT NULL DEFAULT 'light';
```

- 字段名：`theme`
- 类型：VARCHAR(20)
- 默认值：'light'
- 可选值：'light' | 'dark'

## 后端API设计

### 1. 获取用户主题偏好

使用已有的 `GET /api/v1/auth/profile` 接口，返回数据中包含 `theme` 字段。

响应示例：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "theme": "dark",
    ...
  }
}
```

### 2. 更新用户主题偏好

**端点**：`PUT /api/v1/auth/theme`

**请求头**：
- Authorization: Bearer {token}

**请求体**：
```json
{
  "theme": "light"
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success"
}
```

**错误处理**：
- 400：无效的主题值
- 401：未登录
- 500：服务器错误

## 前端设计

### 1. CSS变量定义

创建 `web/src/styles/theme.css` 文件：

**浅色主题**：
```css
:root[data-theme='light'] {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f5;
  --bg-tertiary: #fafafa;
  --text-primary: #333333;
  --text-secondary: #666666;
  --text-tertiary: #999999;
  --border-color: #f0f0f0;
  --shadow-color: rgba(0, 0, 0, 0.08);
}
```

**深色主题**：
```css
:root[data-theme='dark'] {
  --bg-primary: #1f1f1f;
  --bg-secondary: #2a2a2a;
  --bg-tertiary: #333333;
  --text-primary: #ffffff;
  --text-secondary: #cccccc;
  --text-tertiary: #888888;
  --border-color: #3a3a3a;
  --shadow-color: rgba(0, 0, 0, 0.3);
}
```

### 2. Zustand Store扩展

修改 `web/src/stores/authStore.ts`：

```typescript
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
```

### 3. 主题切换逻辑

创建 `web/src/hooks/useTheme.ts`：

```typescript
export function useTheme() {
  const { theme, setTheme } = useAuthStore()
  
  const toggleTheme = async (newTheme: 'light' | 'dark') => {
    // 1. 设置html元素的data-theme属性
    document.documentElement.setAttribute('data-theme', newTheme)
    // 2. 更新store
    setTheme(newTheme)
    // 3. 调用API保存到后端
    await authApi.updateTheme(newTheme)
  }
  
  return { theme, toggleTheme }
}
```

### 4. 用户菜单修改

修改 `web/src/components/Layout/MainLayout.tsx`，在用户下拉菜单中添加主题切换选项：

```tsx
const userMenuItems = [
  {
    key: 'theme',
    icon: theme === 'dark' ? <SunOutlined /> : <MoonOutlined />,
    label: theme === 'dark' ? '切换到浅色模式' : '切换到深色模式',
    onClick: () => toggleTheme(theme === 'dark' ? 'light' : 'dark'),
  },
  // ... 其他菜单项
]
```

### 5. 应用初始化

修改 `web/src/App.tsx`，在应用启动时：

1. 从localStorage恢复登录状态（已有）
2. 调用 `/api/v1/auth/profile` 获取用户信息（包含theme）
3. 设置 `document.documentElement.setAttribute('data-theme', theme)`

### 6. Ant Design主题集成

在 `web/src/App.tsx` 或 `main.tsx` 中使用 ConfigProvider：

```tsx
import { ConfigProvider, theme } from 'antd'

<ConfigProvider
  theme={{
    algorithm: theme === 'dark' ? theme.darkAlgorithm : theme.defaultAlgorithm,
  }}
>
  <App />
</ConfigProvider>
```

### 7. 自定义样式调整

修改以下文件以使用CSS变量：

- `web/src/components/Layout/MainLayout.tsx`：侧边栏、顶部栏背景
- `web/src/index.css`：全局背景、文字颜色

## 实现步骤

1. **后端**
   - 修改 User model 添加 Theme 字段
   - 执行数据库迁移
   - 添加 UpdateTheme handler 和 service 方法
   - 修改 Login/GetProfile 返回 theme 字段

2. **前端**
   - 创建 theme.css 文件
   - 扩展 authStore 添加 theme 状态
   - 创建 useTheme hook
   - 修改 MainLayout 添加主题切换菜单项
   - 修改 App.tsx 初始化主题
   - 调整样式使用CSS变量

3. **测试**
   - 测试主题切换功能
   - 测试主题持久化（刷新页面、重新登录）

## 文件变更清单

### 后端新增/修改文件
- `internal/model/user.go` - 添加 Theme 字段
- `internal/handler/auth.go` - 添加 UpdateTheme handler
- `internal/service/auth.go` - 添加 UpdateTheme 方法
- `internal/repository/user.go` - 添加 UpdateTheme 方法

### 前端新增/修改文件
- `web/src/styles/theme.css` - 新增主题CSS变量
- `web/src/stores/authStore.ts` - 添加 theme 状态
- `web/src/hooks/useTheme.ts` - 新增主题hook
- `web/src/components/Layout/MainLayout.tsx` - 添加主题切换菜单
- `web/src/App.tsx` - 初始化主题、ConfigProvider配置
- `web/src/index.css` - 使用CSS变量
- `web/src/api/auth.ts` - 添加 updateTheme API

## 注意事项

1. 主题切换应该是即时生效的，不需要刷新页面
2. 主题偏好保存失败时，前端仍应切换主题，只显示错误提示
3. 深色主题下需要确保所有文字在深色背景上清晰可读
4. 表格斑马纹、hover效果等需要适配深色主题
