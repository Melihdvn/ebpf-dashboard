import { Table, Tag, Input, Typography, Card, Space } from 'antd';
import { GlobalOutlined, SearchOutlined } from '@ant-design/icons';
import { useState } from 'react';
import { api } from '../services/api';
import { useApiData } from '../hooks/useApiData';
import type { NetworkConnection } from '../types/types';

const { Title } = Typography;

export default function Network() {
  const { data } = useApiData({ fetchFn: () => api.getNetwork(100) });
  const [search, setSearch] = useState('');

  const filtered = data.filter(
    (n) =>
      n.comm.toLowerCase().includes(search.toLowerCase()) ||
      n.source_addr.includes(search) ||
      n.dest_addr.includes(search) ||
      n.pid.toString().includes(search)
  );

  const commFilters = [...new Set(data.map((n) => n.comm))].map((c) => ({ text: c, value: c }));

  const columns = [
    {
      title: 'PID',
      dataIndex: 'pid',
      key: 'pid',
      width: 80,
      render: (pid: string) => <span style={{ fontFamily: 'monospace', color: '#49aa19' }}>{pid}</span>,
    },
    {
      title: 'Command',
      dataIndex: 'comm',
      key: 'comm',
      width: 120,
      render: (c: string) => <Tag color="green">{c}</Tag>,
      filters: commFilters,
      onFilter: (value: unknown, record: NetworkConnection) => record.comm === value,
    },
    {
      title: 'IP',
      dataIndex: 'ip_version',
      key: 'ip_version',
      width: 80,
      render: (v: string) => <Tag color={v === '4' ? 'cyan' : 'magenta'}>IPv{v}</Tag>,
      filters: [
        { text: 'IPv4', value: '4' },
        { text: 'IPv6', value: '6' },
      ],
      onFilter: (value: unknown, record: NetworkConnection) => record.ip_version === value,
    },
    {
      title: 'Source',
      key: 'source',
      render: (_: unknown, r: NetworkConnection) => (
        <span style={{ fontFamily: 'monospace', fontSize: 12 }}>
          {r.source_addr}:<span style={{ color: '#177ddc' }}>{r.source_port}</span>
        </span>
      ),
    },
    {
      title: '',
      key: 'arrow',
      width: 40,
      render: () => <span style={{ color: '#ffffff44' }}>→</span>,
    },
    {
      title: 'Destination',
      key: 'dest',
      render: (_: unknown, r: NetworkConnection) => (
        <span style={{ fontFamily: 'monospace', fontSize: 12 }}>
          {r.dest_addr}:<span style={{ color: '#d89614' }}>{r.dest_port}</span>
        </span>
      ),
    },
    {
      title: 'Timestamp',
      dataIndex: 'timestamp',
      key: 'timestamp',
      width: 180,
      render: (t: string) => <span style={{ color: '#ffffff66', fontSize: 12 }}>{new Date(t).toLocaleString()}</span>,
      sorter: (a: NetworkConnection, b: NetworkConnection) =>
        new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
    },
  ];

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
        <Title level={3} style={{ color: '#e6e6e6', margin: 0 }}>
          <GlobalOutlined style={{ marginRight: 10, color: '#49aa19' }} />
          Network Connections
        </Title>
      </div>

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
            scroll={{ x: 800 }}
          />
        </Space>
      </Card>
    </div>
  );
}
