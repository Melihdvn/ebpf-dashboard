import { Table, Tag, Typography, Card, Space, Input } from 'antd';
import { DashboardOutlined, SearchOutlined } from '@ant-design/icons';
import { useState } from 'react';
import { mockCPUProfiles } from '../mocks/mockData';
import type { CPUProfile } from '../types/types';
import { Bar } from '@ant-design/charts';

const { Title } = Typography;

const columns = [
  {
    title: 'Process',
    dataIndex: 'process_name',
    key: 'process_name',
    width: 130,
    render: (name: string) => <Tag color="red">{name}</Tag>,
    filters: [...new Set(mockCPUProfiles.map((c) => c.process_name))].map((n) => ({ text: n, value: n })),
    onFilter: (value: unknown, record: CPUProfile) => record.process_name === value,
  },
  {
    title: 'Sample Count',
    dataIndex: 'sample_count',
    key: 'sample_count',
    width: 130,
    sorter: (a: CPUProfile, b: CPUProfile) => a.sample_count - b.sample_count,
    render: (v: number) => (
      <span style={{ color: v > 200 ? '#d32029' : v > 100 ? '#d89614' : '#49aa19', fontWeight: 600 }}>
        {v.toLocaleString()}
      </span>
    ),
  },
  {
    title: 'Stack Trace',
    dataIndex: 'stack_trace',
    key: 'stack_trace',
    render: (trace: string) => (
      <div style={{ fontFamily: 'monospace', fontSize: 11, color: '#ffffffcc', wordBreak: 'break-all' as const }}>
        {trace.split(';').map((frame, i) => (
          <span key={i}>
            {i > 0 && <span style={{ color: '#ffffff44' }}> → </span>}
            <span style={{ color: i === 0 ? '#177ddc' : '#ffffffaa' }}>{frame}</span>
          </span>
        ))}
      </div>
    ),
  },
  {
    title: 'Timestamp',
    dataIndex: 'timestamp',
    key: 'timestamp',
    width: 180,
    render: (t: string) => <span style={{ color: '#ffffff66', fontSize: 12 }}>{new Date(t).toLocaleString()}</span>,
  },
];

export default function CPUProfilePage() {
  const [search, setSearch] = useState('');

  const filtered = mockCPUProfiles.filter(
    (c) =>
      c.process_name.toLowerCase().includes(search.toLowerCase()) ||
      c.stack_trace.toLowerCase().includes(search.toLowerCase())
  );

  // Aggregate by process for chart
  const processMap = new Map<string, number>();
  mockCPUProfiles.forEach((c) => {
    processMap.set(c.process_name, (processMap.get(c.process_name) || 0) + c.sample_count);
  });
  const chartData = Array.from(processMap, ([name, count]) => ({ process: name, samples: count }))
    .sort((a, b) => b.samples - a.samples);

  const barConfig = {
    data: chartData,
    xField: 'samples',
    yField: 'process',
    color: '#d32029',
    barStyle: { radius: [0, 4, 4, 0] },
    label: {
      position: 'right' as const,
      style: { fill: '#ffffffaa', fontSize: 11 },
    },
    xAxis: {
      label: { style: { fill: '#ffffffaa' } },
      grid: { line: { style: { stroke: '#303030' } } },
    },
    yAxis: {
      label: { style: { fill: '#ffffffcc' } },
    },
    theme: 'dark',
  };

  return (
    <div>
      <Title level={3} style={{ color: '#e6e6e6', marginBottom: 24 }}>
        <DashboardOutlined style={{ marginRight: 10, color: '#d32029' }} />
        CPU Profiling
      </Title>

      <Card
        title={<span style={{ color: '#e6e6e6' }}>CPU Samples by Process</span>}
        style={{ background: '#141414', borderColor: '#303030', borderRadius: 12, marginBottom: 16 }}
      >
        <Bar {...barConfig} height={280} />
      </Card>

      <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <Input
            placeholder="Search by process name or stack trace..."
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
            pagination={{ pageSize: 10, showSizeChanger: true, showTotal: (t) => `Total ${t} samples` }}
            size="middle"
            scroll={{ x: 800 }}
          />
        </Space>
      </Card>
    </div>
  );
}
