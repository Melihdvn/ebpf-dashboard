import { useState } from 'react';
import { BrowserRouter, Routes, Route, useNavigate, useLocation } from 'react-router-dom';
import { ConfigProvider, Layout, Menu, theme, Typography } from 'antd';
import {
  DashboardOutlined,
  CodeOutlined,
  GlobalOutlined,
  HddOutlined,
  DashboardFilled,
  ApiOutlined,
  FunctionOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons';

import Dashboard from './pages/Dashboard';
import Processes from './pages/Processes';
import Network from './pages/Network';
import DiskIO from './pages/DiskIO';
import CPUProfilePage from './pages/CPUProfile';
import TCPLife from './pages/TCPLife';
import Syscalls from './pages/Syscalls';

const { Header, Sider, Content } = Layout;
const { Title } = Typography;

const menuItems = [
  { key: '/', icon: <DashboardOutlined />, label: 'Dashboard' },
  { key: '/processes', icon: <CodeOutlined />, label: 'Processes' },
  { key: '/network', icon: <GlobalOutlined />, label: 'Network' },
  { key: '/disk', icon: <HddOutlined />, label: 'Disk I/O' },
  { key: '/cpu', icon: <DashboardFilled />, label: 'CPU Profile' },
  { key: '/tcplife', icon: <ApiOutlined />, label: 'TCP Life' },
  { key: '/syscalls', icon: <FunctionOutlined />, label: 'Syscalls' },
];

function AppLayout() {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={setCollapsed}
        trigger={null}
        width={240}
        style={{
          background: '#0a0a0a',
          borderRight: '1px solid #1f1f1f',
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
          zIndex: 10,
        }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: collapsed ? 'center' : 'flex-start',
            padding: collapsed ? '0' : '0 20px',
            borderBottom: '1px solid #1f1f1f',
            gap: 10,
          }}
        >
          <div
            style={{
              width: 32,
              height: 32,
              borderRadius: 8,
              background: 'linear-gradient(135deg, #177ddc, #13a8a8)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: 16,
              fontWeight: 700,
              color: '#fff',
              flexShrink: 0,
            }}
          >
            e
          </div>
          {!collapsed && (
            <Title level={5} style={{ color: '#e6e6e6', margin: 0, whiteSpace: 'nowrap' }}>
              eBPF Dashboard
            </Title>
          )}
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{
            background: 'transparent',
            border: 'none',
            marginTop: 8,
          }}
        />
      </Sider>

      <Layout style={{ marginLeft: collapsed ? 80 : 240, transition: 'margin-left 0.2s' }}>
        <Header
          style={{
            background: '#0d0d0d',
            borderBottom: '1px solid #1f1f1f',
            padding: '0 24px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            position: 'sticky',
            top: 0,
            zIndex: 5,
            height: 64,
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <div
              onClick={() => setCollapsed(!collapsed)}
              style={{ cursor: 'pointer', color: '#ffffffaa', fontSize: 18 }}
            >
              {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            </div>
            <span style={{ color: '#ffffff88', fontSize: 13 }}>
              Real-time eBPF System Monitoring
            </span>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <div
              style={{
                width: 8,
                height: 8,
                borderRadius: '50%',
                background: '#49aa19',
                boxShadow: '0 0 8px #49aa19',
                animation: 'pulse 2s infinite',
              }}
            />
            <span style={{ color: '#49aa19', fontSize: 12 }}>Mock Data</span>
          </div>
        </Header>
        <Content
          style={{
            margin: 24,
            minHeight: 'calc(100vh - 112px)',
          }}
        >
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/processes" element={<Processes />} />
            <Route path="/network" element={<Network />} />
            <Route path="/disk" element={<DiskIO />} />
            <Route path="/cpu" element={<CPUProfilePage />} />
            <Route path="/tcplife" element={<TCPLife />} />
            <Route path="/syscalls" element={<Syscalls />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  );
}

export default function App() {
  return (
    <ConfigProvider
      theme={{
        algorithm: theme.darkAlgorithm,
        token: {
          colorPrimary: '#177ddc',
          colorBgContainer: '#141414',
          colorBgElevated: '#1f1f1f',
          borderRadius: 8,
          colorBorderSecondary: '#303030',
          fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
        },
        components: {
          Menu: {
            darkItemBg: 'transparent',
            darkItemSelectedBg: '#177ddc22',
            darkItemHoverBg: '#ffffff0a',
            itemBorderRadius: 8,
            itemMarginInline: 8,
          },
          Table: {
            headerBg: '#1a1a1a',
            rowHoverBg: '#ffffff08',
          },
          Card: {
            paddingLG: 20,
          },
        },
      }}
    >
      <BrowserRouter>
        <AppLayout />
      </BrowserRouter>
    </ConfigProvider>
  );
}
