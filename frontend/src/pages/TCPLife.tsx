import { Table, Tag, Typography, Card, Space, Input, Statistic, Row, Col } from 'antd';
import { ApiOutlined, SearchOutlined } from '@ant-design/icons';
import { useState } from 'react';
import { mockTCPLifeEvents } from '../mocks/mockData';
import type { TCPLifeEvent } from '../types/types';

const { Title } = Typography;

const getDurationColor = (ms: number) => {
  if (ms < 1000) return '#49aa19';
  if (ms < 10000) return '#d89614';
  if (ms < 100000) return '#d87a16';
  return '#d32029';
};

const formatDuration = (ms: number) => {
  if (ms < 1000) return `${ms.toFixed(0)} ms`;
  if (ms < 60000) return `${(ms / 1000).toFixed(1)} s`;
  return `${(ms / 60000).toFixed(1)} min`;
};

const formatKB = (kb: number) => {
  if (kb < 1024) return `${kb.toFixed(1)} KB`;
  return `${(kb / 1024).toFixed(1)} MB`;
};

const columns = [
  {
    title: 'PID',
    dataIndex: 'pid',
    key: 'pid',
    width: 80,
    render: (pid: number) => <span style={{ fontFamily: 'monospace', color: '#13a8a8' }}>{pid}</span>,
  },
  {
    title: 'Command',
    dataIndex: 'comm',
    key: 'comm',
    width: 120,
    render: (c: string) => <Tag color="cyan">{c}</Tag>,
    filters: [...new Set(mockTCPLifeEvents.map((t) => t.comm))].map((c) => ({ text: c, value: c })),
    onFilter: (value: unknown, record: TCPLifeEvent) => record.comm === value,
  },
  {
    title: 'Local',
    key: 'local',
    render: (_: unknown, r: TCPLifeEvent) => (
      <span style={{ fontFamily: 'monospace', fontSize: 12 }}>
        {r.local_addr}:<span style={{ color: '#177ddc' }}>{r.local_port}</span>
      </span>
    ),
  },
  {
    title: '',
    key: 'arrow',
    width: 40,
    render: () => <span style={{ color: '#ffffff44' }}>⇌</span>,
  },
  {
    title: 'Remote',
    key: 'remote',
    render: (_: unknown, r: TCPLifeEvent) => (
      <span style={{ fontFamily: 'monospace', fontSize: 12 }}>
        {r.remote_addr}:<span style={{ color: '#d89614' }}>{r.remote_port}</span>
      </span>
    ),
  },
  {
    title: 'TX',
    dataIndex: 'tx_kb',
    key: 'tx_kb',
    width: 100,
    render: (v: number) => <span style={{ color: '#49aa19' }}>↑ {formatKB(v)}</span>,
    sorter: (a: TCPLifeEvent, b: TCPLifeEvent) => a.tx_kb - b.tx_kb,
  },
  {
    title: 'RX',
    dataIndex: 'rx_kb',
    key: 'rx_kb',
    width: 100,
    render: (v: number) => <span style={{ color: '#177ddc' }}>↓ {formatKB(v)}</span>,
    sorter: (a: TCPLifeEvent, b: TCPLifeEvent) => a.rx_kb - b.rx_kb,
  },
  {
    title: 'Duration',
    dataIndex: 'duration_ms',
    key: 'duration_ms',
    width: 120,
    render: (v: number) => (
      <Tag color={getDurationColor(v)} style={{ fontWeight: 600 }}>
        {formatDuration(v)}
      </Tag>
    ),
    sorter: (a: TCPLifeEvent, b: TCPLifeEvent) => a.duration_ms - b.duration_ms,
  },
];

export default function TCPLife() {
  const [search, setSearch] = useState('');

  const filtered = mockTCPLifeEvents.filter(
    (t) =>
      t.comm.toLowerCase().includes(search.toLowerCase()) ||
      t.local_addr.includes(search) ||
      t.remote_addr.includes(search) ||
      t.pid.toString().includes(search)
  );

  const totalTX = mockTCPLifeEvents.reduce((s, t) => s + t.tx_kb, 0);
  const totalRX = mockTCPLifeEvents.reduce((s, t) => s + t.rx_kb, 0);
  const avgDuration = mockTCPLifeEvents.reduce((s, t) => s + t.duration_ms, 0) / mockTCPLifeEvents.length;

  return (
    <div>
      <Title level={3} style={{ color: '#e6e6e6', marginBottom: 24 }}>
        <ApiOutlined style={{ marginRight: 10, color: '#13a8a8' }} />
        TCP Connection Lifecycle
      </Title>

      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={8}>
          <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
            <Statistic
              title={<span style={{ color: '#ffffffaa' }}>Total TX</span>}
              value={formatKB(totalTX)}
              valueStyle={{ color: '#49aa19' }}
              prefix="↑"
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
            <Statistic
              title={<span style={{ color: '#ffffffaa' }}>Total RX</span>}
              value={formatKB(totalRX)}
              valueStyle={{ color: '#177ddc' }}
              prefix="↓"
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
            <Statistic
              title={<span style={{ color: '#ffffffaa' }}>Avg Duration</span>}
              value={formatDuration(avgDuration)}
              valueStyle={{ color: '#d89614' }}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <Input
            placeholder="Search by PID, command or IP address..."
            prefix={<SearchOutlined style={{ color: '#ffffff55' }} />}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            style={{ maxWidth: 400, background: '#1f1f1f', borderColor: '#303030' }}
            allowClear
          />
          <Table
            dataSource={filtered}
            columns={columns}
            rowKey="id"
            pagination={{ pageSize: 10, showSizeChanger: true, showTotal: (t) => `Total ${t} connections` }}
            size="middle"
            scroll={{ x: 900 }}
          />
        </Space>
      </Card>
    </div>
  );
}
