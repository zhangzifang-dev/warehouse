import { useState } from 'react'
import { Layout, Menu, Dropdown, Avatar, Button, theme, Tooltip } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  UserOutlined,
  SafetyOutlined,
  HomeOutlined,
  EnvironmentOutlined,
  AppstoreOutlined,
  TagsOutlined,
  ShopOutlined,
  IdcardOutlined,
  DatabaseOutlined,
  LoginOutlined,
  ExportOutlined,
  SwapOutlined,
  AuditOutlined,
  LockOutlined,
  LogoutOutlined,
  InboxOutlined,
  MoonOutlined,
  SunOutlined,
} from '@ant-design/icons'
import { useNavigate, useLocation, Outlet } from 'react-router-dom'
import { useAuthStore } from '../../stores/authStore'
import { useTheme } from '../../hooks/useTheme'

const { Header, Sider, Content } = Layout

const menuItems = [
  { key: '/dashboard', icon: <DashboardOutlined />, label: '仪表盘' },
  { key: '/users', icon: <UserOutlined />, label: '用户管理' },
  { key: '/roles', icon: <SafetyOutlined />, label: '角色管理' },
  { key: '/warehouses', icon: <HomeOutlined />, label: '仓库管理' },
  { key: '/locations', icon: <EnvironmentOutlined />, label: '库位管理' },
  { key: '/products', icon: <AppstoreOutlined />, label: '商品管理' },
  { key: '/categories', icon: <TagsOutlined />, label: '分类管理' },
  { key: '/suppliers', icon: <ShopOutlined />, label: '供应商管理' },
  { key: '/customers', icon: <IdcardOutlined />, label: '客户管理' },
  { key: '/inventory', icon: <DatabaseOutlined />, label: '库存管理' },
  { key: '/inbound', icon: <LoginOutlined />, label: '入库管理' },
  { key: '/outbound', icon: <ExportOutlined />, label: '出库管理' },
  { key: '/transfers', icon: <SwapOutlined />, label: '库存调拨' },
  { key: '/audit-logs', icon: <AuditOutlined />, label: '审计日志' },
]

const pageTitle: Record<string, string> = {
  '/dashboard': '仪表盘',
  '/users': '用户管理',
  '/roles': '角色管理',
  '/permissions': '权限管理',
  '/warehouses': '仓库管理',
  '/locations': '库位管理',
  '/products': '商品管理',
  '/categories': '分类管理',
  '/suppliers': '供应商管理',
  '/customers': '客户管理',
  '/inventory': '库存管理',
  '/inbound': '入库管理',
  '/outbound': '出库管理',
  '/transfers': '库存调拨',
  '/audit-logs': '审计日志',
  '/change-password': '修改密码',
}

export function MainLayout() {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const { token: { colorBgContainer, colorBorder, colorText, colorTextSecondary } } = theme.useToken()
  const { theme: currentTheme, toggleTheme } = useTheme()

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key)
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const userMenuItems = [
    {
      key: 'theme',
      icon: currentTheme === 'light' ? <MoonOutlined /> : <SunOutlined />,
      label: currentTheme === 'light' ? '深色模式' : '浅色模式',
      onClick: toggleTheme,
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

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed} collapsedWidth={60} width={150} theme="light" style={{ borderRight: `1px solid ${colorBorder}` }}>
        <div style={{
          height: 40,
          display: 'flex',
          alignItems: 'center',
          justifyContent: collapsed ? 'center' : 'flex-start',
          padding: collapsed ? 0 : '0 0 0 16px',
          borderBottom: `1px solid ${colorBorder}`,
          flexShrink: 0,
        }}>
          {collapsed ? (
            <Tooltip title="WMS" placement="right">
              <InboxOutlined style={{ fontSize: 18, color: '#1890ff', cursor: 'pointer' }} />
            </Tooltip>
          ) : (
            <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
              <InboxOutlined style={{ fontSize: 18, color: '#1890ff' }} />
              <span style={{ fontSize: 15, fontWeight: 600, color: colorText, letterSpacing: 0.5 }}>WMS</span>
            </div>
          )}
        </div>
        <div className="sider-menu-scroll" style={{
          flex: 1,
          overflowY: 'auto',
          overflowX: 'hidden',
        }}>
          <Menu
            theme="light"
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={handleMenuClick}
          />
        </div>
        <div style={{
          height: 32,
          display: 'flex',
          alignItems: 'center',
          justifyContent: collapsed ? 'center' : 'flex-start',
          padding: collapsed ? 0 : '0 0 0 16px',
          borderTop: `1px solid ${colorBorder}`,
          flexShrink: 0,
        }}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{ fontSize: '14px', color: colorTextSecondary }}
          />
        </div>
      </Sider>
      <Layout>
        <Header style={{
          height: 32,
          lineHeight: '32px',
          padding: '0 16px',
          background: colorBgContainer,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          borderBottom: `1px solid ${colorBorder}`,
          boxShadow: '0 1px 2px rgba(0, 0, 0, 0.03)',
        }}>
          <span style={{ fontSize: 14, fontWeight: 500, color: colorText }}>{pageTitle[location.pathname] || '未知页面'}</span>
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 6 }}>
              <Avatar icon={<UserOutlined />} size="small" />
              <span style={{ fontSize: 12, color: colorText }}>{user?.username}</span>
            </div>
          </Dropdown>
        </Header>
        <Content style={{
          padding: 8,
          background: colorBgContainer,
          minHeight: 280,
        }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
