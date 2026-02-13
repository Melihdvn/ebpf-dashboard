import { Table, Tag, Input, Typography, Card, Space } from 'antd';
import { CodeOutlined, SearchOutlined } from '@ant-design/icons';
import { useState } from 'react';
import { api } from '../services/api';
import { useApiData } from '../hooks/useApiData';
import type { ProcessEvent } from '../types/types';

const { Title } = Typography;

export default function Processes() {
  const { data } = useApiData({ fetchFn: () => api.getProcesses(100) });
  const [search, setSearch] = useState('');

  const filtered = data.filter(
    (p) =>
      p.comm.toLowerCase().includes(search.toLowerCase()) ||
      p.args.toLowerCase().includes(search.toLowerCase()) ||
      p.pid.toString().includes(search)
  );

  const commFilters = [...new Set(data.map((p) => p.comm))].map((c) => ({ text: c, value: c }));

  const columns = [
    {
      title: 'PID',
      dataIndex: 'pid',
      key: 'pid',
      width: 90,
      sorter: (a: ProcessEvent, b: ProcessEvent) => parseInt(String(a.pid)) - parseInt(String(b.pid)),
      render: (pid: string) => <span style={{ fontFamily: 'monospace', color: '#177ddc' }}>{pid}</span>,
    },
    {
      title: 'Command',
      dataIndex: 'comm',
      key: 'comm',
      width: 150,
      render: (comm: string) => <Tag color="blue">{comm}</Tag>,
      filters: commFilters,
      onFilter: (value: unknown, record: ProcessEvent) => record.comm === value,
    },
    {
      title: 'Arguments',
      dataIndex: 'args',
      key: 'args',
      ellipsis: true,
      render: (args: string) => <span style={{ fontFamily: 'monospace', fontSize: 12, color: '#ffffffcc' }}>{args}</span>,
    },
    {
      title: 'Time',
      dataIndex: 'time',
      key: 'time',
      width: 110,
      render: (t: string) => <span style={{ color: '#ffffff88' }}>{t}</span>,
    },
    {
      title: 'Timestamp',
      dataIndex: 'timestamp',
      key: 'timestamp',
      width: 200,
      render: (t: string) => <span style={{ color: '#ffffff66', fontSize: 12 }}>{new Date(t).toLocaleString()}</span>,
      sorter: (a: ProcessEvent, b: ProcessEvent) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
    },
  ];

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
        <Title level={3} style={{ color: '#e6e6e6', margin: 0 }}>
          <CodeOutlined style={{ marginRight: 10, color: '#177ddc' }} />
          Process Monitoring
        </Title>
      </div>

      <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <Input
            placeholder="Search by PID, command or arguments..."
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
            pagination={{ pageSize: 10, showSizeChanger: true, showTotal: (t) => `Total ${t} processes` }}
            size="middle"
            scroll={{ x: 700 }}
          />
        </Space>
      </Card>
    </div>
  );
}
