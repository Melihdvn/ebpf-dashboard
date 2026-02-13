import { Table, Tag, Typography, Card, Space, Row, Col, Statistic } from 'antd';
import { FunctionOutlined } from '@ant-design/icons';
import { Bar } from '@ant-design/charts';
import { api } from '../services/api';
import { useApiData } from '../hooks/useApiData';
import type { SyscallStat } from '../types/types';

const { Title } = Typography;

export default function Syscalls() {
  const { data } = useApiData({ fetchFn: () => api.getSyscalls(100) });

  const sorted = [...data].sort((a, b) => b.count - a.count);
  const totalCalls = data.reduce((s, sc) => s + sc.count, 0);
  const maxCount = sorted.length > 0 ? sorted[0].count : 1;

  const chartData = sorted.slice(0, 10).map((s) => ({
    syscall: s.syscall_name,
    count: s.count,
  }));

  const barConfig = {
    data: chartData,
    xField: 'count',
    yField: 'syscall',
    color: '#cb2b83',
    barStyle: { radius: [0, 4, 4, 0] },
    label: {
      position: 'right' as const,
      style: { fill: '#ffffffaa', fontSize: 11 },
      formatter: (datum: { count?: number }) => (datum.count ?? 0).toLocaleString(),
    },
    xAxis: {
      label: { style: { fill: '#ffffffaa' } },
      grid: { line: { style: { stroke: '#303030' } } },
    },
    yAxis: {
      label: { style: { fill: '#ffffffcc', fontFamily: 'monospace' } },
    },
    theme: 'dark',
  };

  const columns = [
    {
      title: '#',
      key: 'rank',
      width: 50,
      render: (_: unknown, __: unknown, i: number) => (
        <span style={{ color: '#ffffff66' }}>{i + 1}</span>
      ),
    },
    {
      title: 'Syscall',
      dataIndex: 'syscall_name',
      key: 'syscall_name',
      render: (name: string) => (
        <Tag color="purple" style={{ fontFamily: 'monospace' }}>
          {name}
        </Tag>
      ),
    },
    {
      title: 'Count',
      dataIndex: 'count',
      key: 'count',
      sorter: (a: SyscallStat, b: SyscallStat) => a.count - b.count,
      render: (v: number) => (
        <span style={{ fontWeight: 600, color: v > 20000 ? '#d32029' : v > 10000 ? '#d89614' : '#49aa19' }}>
          {v.toLocaleString()}
        </span>
      ),
    },
    {
      title: 'Percentage',
      key: 'percentage',
      render: (_: unknown, record: SyscallStat) => {
        const pct = totalCalls > 0 ? ((record.count / totalCalls) * 100).toFixed(1) : '0';
        return (
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <div
              style={{
                width: `${(record.count / maxCount) * 100}%`,
                maxWidth: 200,
                height: 6,
                borderRadius: 3,
                background: 'linear-gradient(90deg, #cb2b83, #d32029)',
                minWidth: 4,
              }}
            />
            <span style={{ color: '#ffffffaa', fontSize: 12 }}>{pct}%</span>
          </div>
        );
      },
    },
  ];

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
        <Title level={3} style={{ color: '#e6e6e6', margin: 0 }}>
          <FunctionOutlined style={{ marginRight: 10, color: '#cb2b83' }} />
          System Call Statistics
        </Title>
      </div>

      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={8}>
          <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
            <Statistic
              title={<span style={{ color: '#ffffffaa' }}>Total Syscalls</span>}
              value={totalCalls}
              valueStyle={{ color: '#cb2b83' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
            <Statistic
              title={<span style={{ color: '#ffffffaa' }}>Unique Syscalls</span>}
              value={data.length}
              valueStyle={{ color: '#177ddc' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
            <Statistic
              title={<span style={{ color: '#ffffffaa' }}>Most Called</span>}
              value={sorted[0]?.syscall_name ?? '-'}
              valueStyle={{ color: '#d32029', fontSize: 20, fontFamily: 'monospace' }}
            />
          </Card>
        </Col>
      </Row>

      <Card
        title={<span style={{ color: '#e6e6e6' }}>Top 10 Syscalls</span>}
        style={{ background: '#141414', borderColor: '#303030', borderRadius: 12, marginBottom: 16 }}
      >
        <Bar {...barConfig} height={300} />
      </Card>

      <Card
        title={<span style={{ color: '#e6e6e6' }}>All Syscalls</span>}
        style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}
      >
        <Space direction="vertical" style={{ width: '100%' }}>
          <Table
            dataSource={sorted}
            columns={columns}
            rowKey="id"
            pagination={{ pageSize: 10, showSizeChanger: true, showTotal: (t) => `Total ${t} syscalls` }}
            size="middle"
          />
        </Space>
      </Card>
    </div>
  );
}
