import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { MainLayout } from './components/Layout'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Login } from './pages/auth/Login'
import { ChangePassword } from './pages/auth/ChangePassword'

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
            <Route path="warehouses" element={<div>仓库管理</div>} />
            <Route path="products" element={<div>商品管理</div>} />
            <Route path="suppliers" element={<div>供应商管理</div>} />
            <Route path="customers" element={<div>客户管理</div>} />
            <Route path="inbound" element={<div>入库管理</div>} />
            <Route path="outbound" element={<div>出库管理</div>} />
            <Route path="transfers" element={<div>库存调拨</div>} />
            <Route path="audit-logs" element={<div>审计日志</div>} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
