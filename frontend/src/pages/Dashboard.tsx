import { Card, Col, Row, Statistic, Table, Tag, Typography, Spin } from 'antd';
import {
  CodeOutlined,
  GlobalOutlined,
  HddOutlined,
  DashboardOutlined,
  ApiOutlined,
  FunctionOutlined,
} from '@ant-design/icons';
import { api } from '../services/api';
import { useApiData } from '../hooks/useApiData';

const { Title } = Typography;

const recentProcessColumns = [
  { title: 'PID', dataIndex: 'pid', key: 'pid', width: 80 },
  { title: 'Command', dataIndex: 'comm', key: 'comm', render: (t: string) => <Tag color="blue">{t}</Tag> },
  { title: 'Arguments', dataIndex: 'args', key: 'args', ellipsis: true },
  { title: 'Time', dataIndex: 'time', key: 'time', width: 100 },
];

const topSyscallColumns = [
  { title: 'Syscall', dataIndex: 'syscall_name', key: 'syscall_name', render: (t: string) => <Tag color="purple">{t}</Tag> },
  {
    title: 'Count',
    dataIndex: 'count',
    key: 'count',
    render: (v: number) => v.toLocaleString(),
    sorter: (a: { count: number }, b: { count: number }) => a.count - b.count,
  },
];

export default function Dashboard() {
  const processes = useApiData({ fetchFn: () => api.getProcesses(20) });
  const network = useApiData({ fetchFn: () => api.getNetwork(20) });
  const disk = useApiData({ fetchFn: () => api.getDisk() });
  const cpu = useApiData({ fetchFn: () => api.getCPUProfile(20) });
  const tcp = useApiData({ fetchFn: () => api.getTCPLife(20) });
  const syscalls = useApiData({ fetchFn: () => api.getSyscalls(20) });

  const isAnyLoading = processes.loading;

  const summaryCards = [
    {
      title: 'Active Processes',
      value: processes.data.length,
      icon: <CodeOutlined style={{ fontSize: 28, color: '#177ddc' }} />,
      color: 'linear-gradient(135deg, #141e30 0%, #1a2a4a 100%)',
      accent: '#177ddc',
    },
    {
      title: 'Network Connections',
      value: network.data.length,
      icon: <GlobalOutlined style={{ fontSize: 28, color: '#49aa19' }} />,
      color: 'linear-gradient(135deg, #1a2e1a 0%, #1a3a2a 100%)',
      accent: '#49aa19',
    },
    {
      title: 'Disk I/O Events',
      value: disk.data.reduce((s, d) => s + d.count, 0).toLocaleString(),
      icon: <HddOutlined style={{ fontSize: 28, color: '#d89614' }} />,
      color: 'linear-gradient(135deg, #2a2010 0%, #3a2a10 100%)',
      accent: '#d89614',
    },
    {
      title: 'CPU Samples',
      value: cpu.data.reduce((s, c) => s + c.sample_count, 0).toLocaleString(),
      icon: <DashboardOutlined style={{ fontSize: 28, color: '#d32029' }} />,
      color: 'linear-gradient(135deg, #2a1215 0%, #3a1520 100%)',
      accent: '#d32029',
    },
    {
      title: 'TCP Connections',
      value: tcp.data.length,
      icon: <ApiOutlined style={{ fontSize: 28, color: '#13a8a8' }} />,
      color: 'linear-gradient(135deg, #112a2a 0%, #153a3a 100%)',
      accent: '#13a8a8',
    },
    {
      title: 'Unique Syscalls',
      value: syscalls.data.length,
      icon: <FunctionOutlined style={{ fontSize: 28, color: '#cb2b83' }} />,
      color: 'linear-gradient(135deg, #2a1225 0%, #3a1530 100%)',
      accent: '#cb2b83',
    },
  ];

  if (isAnyLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '50vh' }}>
        <Spin size="large" tip="Loading metrics..." />
      </div>
    );
  }

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
        <Title level={3} style={{ color: '#e6e6e6', margin: 0 }}>
          <DashboardOutlined style={{ marginRight: 10, color: '#177ddc' }} />
          System Overview
        </Title>
      </div>

      <Row gutter={[16, 16]}>
        {summaryCards.map((card, i) => (
          <Col xs={24} sm={12} lg={8} key={i}>
            <Card
              style={{
                background: card.color,
                border: `1px solid ${card.accent}33`,
                borderRadius: 12,
              }}
              hoverable
            >
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <Statistic
                  title={<span style={{ color: '#ffffffaa', fontSize: 13 }}>{card.title}</span>}
                  value={card.value}
                  valueStyle={{ color: card.accent, fontWeight: 700, fontSize: 28 }}
                />
                <div
                  style={{
                    width: 56,
                    height: 56,
                    borderRadius: 16,
                    background: `${card.accent}18`,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                >
                  {card.icon}
                </div>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} lg={14}>
          <Card
            title={<span style={{ color: '#e6e6e6' }}>Recent Processes</span>}
            style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}
          >
            <Table
              dataSource={processes.data.slice(0, 8)}
              columns={recentProcessColumns}
              rowKey="id"
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
        <Col xs={24} lg={10}>
          <Card
            title={<span style={{ color: '#e6e6e6' }}>Top Syscalls</span>}
            style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}
          >
            <Table
              dataSource={[...syscalls.data].sort((a, b) => b.count - a.count).slice(0, 8)}
              columns={topSyscallColumns}
              rowKey="id"
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
}
