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
import { UserList } from './pages/system/UserList'
import { RoleList } from './pages/system/RoleList'
import { PermissionList } from './pages/system/PermissionList'
import { AuditLogList } from './pages/system/AuditLogList'

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
