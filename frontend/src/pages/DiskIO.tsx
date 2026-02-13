import { Typography, Card, Row, Col, Statistic } from 'antd';
import { HddOutlined } from '@ant-design/icons';
import { Column } from '@ant-design/charts';
import { mockDiskLatency } from '../mocks/mockData';

const { Title } = Typography;

export default function DiskIO() {
  const chartData = mockDiskLatency.map((d) => ({
    range: d.range_max <= 1023
      ? `${d.range_min}-${d.range_max} µs`
      : `${(d.range_min / 1024).toFixed(0)}-${(d.range_max / 1024).toFixed(0)} ms`,
    count: d.count,
  }));

  const totalIO = mockDiskLatency.reduce((s, d) => s + d.count, 0);
  const peakBucket = mockDiskLatency.reduce((max, d) => (d.count > max.count ? d : max), mockDiskLatency[0]);
  const avgLatency =
    mockDiskLatency.reduce((s, d) => s + ((d.range_min + d.range_max) / 2) * d.count, 0) / totalIO;

  const config = {
    data: chartData,
    xField: 'range',
    yField: 'count',
    color: '#d89614',
    columnStyle: {
      radius: [4, 4, 0, 0],
    },
    label: {
      position: 'top' as const,
      style: {
        fill: '#ffffffaa',
        fontSize: 11,
      },
    },
    xAxis: {
      label: {
        autoRotate: true,
        style: { fill: '#ffffffaa', fontSize: 11 },
      },
    },
    yAxis: {
      label: {
        style: { fill: '#ffffffaa' },
      },
      grid: {
        line: { style: { stroke: '#303030' } },
      },
    },
    theme: 'dark',
  };

  return (
    <div>
      <Title level={3} style={{ color: '#e6e6e6', marginBottom: 24 }}>
        <HddOutlined style={{ marginRight: 10, color: '#d89614' }} />
        Disk I/O Latency
      </Title>

      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={8}>
          <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
            <Statistic
              title={<span style={{ color: '#ffffffaa' }}>Total I/O Events</span>}
              value={totalIO}
              valueStyle={{ color: '#d89614' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
            <Statistic
              title={<span style={{ color: '#ffffffaa' }}>Peak Latency Range</span>}
              value={`${peakBucket.range_min}-${peakBucket.range_max} µs`}
              valueStyle={{ color: '#d32029', fontSize: 20 }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}>
            <Statistic
              title={<span style={{ color: '#ffffffaa' }}>Avg Latency</span>}
              value={avgLatency.toFixed(1)}
              suffix="µs"
              valueStyle={{ color: '#49aa19' }}
            />
          </Card>
        </Col>
      </Row>

      <Card
        title={<span style={{ color: '#e6e6e6' }}>Latency Distribution Histogram</span>}
        style={{ background: '#141414', borderColor: '#303030', borderRadius: 12 }}
      >
        <Column {...config} />
      </Card>
    </div>
  );
}
