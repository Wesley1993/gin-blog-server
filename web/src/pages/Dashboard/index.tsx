import { useCallback, useEffect, useState } from 'react';
import { Card, Col, Descriptions, Row, Spin, Statistic, Tag, Button } from 'antd';
import {
  FileTextOutlined,
  AppstoreOutlined,
  TeamOutlined,
  FieldTimeOutlined,
  ReloadOutlined,
  DatabaseOutlined,
  ThunderboltOutlined,
  SearchOutlined,
  CodeOutlined,
  ClockCircleOutlined,
  DashboardOutlined,
  RocketOutlined,
  FundOutlined,
} from '@ant-design/icons';
import PageHeader from '../../components/PageHeader';
import { getSiteStats } from '../../api/site';
import { getSystemOverview } from '../../api/system';
import type { SiteStats } from '../../api/types';
import type { SystemOverview } from '../../api/system';

/** 依赖健康状态 → Tag 配置 */
const STATUS_MAP: Record<string, { color: string; text: string; dot: string }> = {
  healthy: { color: 'success', text: '正常', dot: '#52c41a' },
  unhealthy: { color: 'error', text: '异常', dot: '#ff4d4f' },
  not_configured: { color: 'default', text: '未配置', dot: '#d6cdbd' },
};

/** 运行模式 → 展示文案 */
const MODE_TEXT: Record<string, string> = {
  dev: '开发',
  test: '测试',
  release: '生产',
};

function StatusBadge({ status }: { status?: string }) {
  const cfg = STATUS_MAP[status ?? ''] ?? STATUS_MAP.not_configured;
  return (
    <span className="inline-flex items-center gap-2">
      <span
        className={`inline-block w-2 h-2 rounded-full ${status === 'healthy' ? 'dash-pulse-dot' : ''}`}
        style={{ background: cfg.dot }}
      />
      <Tag color={cfg.color} className="m-0">
        {cfg.text}
      </Tag>
    </span>
  );
}

export default function DashboardPage() {
  const [stats, setStats] = useState<SiteStats | null>(null);
  const [overview, setOverview] = useState<SystemOverview | null>(null);
  const [loading, setLoading] = useState(false);

  const fetchAll = useCallback(() => {
    setLoading(true);
    Promise.all([
      getSiteStats().catch(() => undefined),
      getSystemOverview().catch(() => undefined),
    ])
      .then(([statsRes, overviewRes]) => {
        setStats(statsRes?.data ?? null);
        setOverview(overviewRes?.data ?? null);
      })
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    fetchAll();
  }, [fetchAll]);

  const statCards = [
    {
      title: '文章总数',
      value: stats?.article_count,
      icon: <FileTextOutlined />,
      accent: '#b45309',
    },
    {
      title: '分类数',
      value: stats?.category_count,
      icon: <AppstoreOutlined />,
      accent: '#92400e',
    },
    {
      title: '用户数',
      value: stats?.user_count,
      icon: <TeamOutlined />,
      accent: '#7c2d12',
    },
    {
      title: '运行天数',
      value: stats?.running_days,
      suffix: '天',
      icon: <FieldTimeOutlined />,
      accent: '#d97706',
    },
  ];

  const dependencies = [
    { key: 'db_status', name: 'PostgreSQL', desc: '主数据库', icon: <DatabaseOutlined /> },
    { key: 'redis_status', name: 'Redis', desc: '缓存与会话', icon: <ThunderboltOutlined /> },
    { key: 'es_status', name: 'Elasticsearch', desc: '全文检索', icon: <SearchOutlined /> },
  ] as const;

  return (
    <div className="page-container">
      <PageHeader
        title="仪表盘"
        subtitle="站点数据与系统运行状态的总览"
        extra={
          <Button icon={<ReloadOutlined spin={loading} />} onClick={fetchAll}>
            刷新
          </Button>
        }
      />

      <Spin spinning={loading}>
        {/* 统计卡片行 */}
        <Row gutter={[16, 16]}>
          {statCards.map((item, idx) => (
            <Col key={item.title} xs={24} sm={12} xl={6}>
              <div className="dash-reveal" style={{ animationDelay: `${idx * 90}ms` }}>
                <Card variant="borderless" className="dash-stat-card">
                  <div className="flex items-start justify-between">
                    <Statistic
                      title={item.title}
                      value={item.value ?? 0}
                      suffix={item.suffix}
                      valueStyle={{ color: item.accent, fontWeight: 700 }}
                    />
                    <span
                      className="w-11 h-11 shrink-0 flex items-center justify-center rounded-full text-xl"
                      style={{ background: 'rgba(180, 83, 9, 0.08)', color: item.accent }}
                    >
                      {item.icon}
                    </span>
                  </div>
                  <div className="mt-2 h-[2px] w-10 rounded" style={{ background: item.accent }} />
                </Card>
              </div>
            </Col>
          ))}
        </Row>

        {/* 环境信息 + 依赖状态 */}
        <Row gutter={[16, 16]} className="mt-4">
          <Col xs={24} xl={14}>
            <div className="dash-reveal" style={{ animationDelay: '360ms' }}>
              <Card
                variant="borderless"
                title={
                  <span className="flex items-center gap-2">
                    <CodeOutlined style={{ color: '#b45309' }} />
                    运行环境
                  </span>
                }
              >
                <Descriptions
                  column={{ xs: 1, sm: 2 }}
                  items={[
                    {
                      key: 'app_name',
                      label: (
                        <span className="inline-flex items-center gap-1.5">
                          <RocketOutlined /> 应用名称
                        </span>
                      ),
                      children: overview?.app_name || '—',
                    },
                    {
                      key: 'app_version',
                      label: (
                        <span className="inline-flex items-center gap-1.5">
                          <FundOutlined /> 应用版本
                        </span>
                      ),
                      children: overview ? (
                        <Tag color="amber" className="m-0">
                          v{overview.app_version}
                        </Tag>
                      ) : (
                        '—'
                      ),
                    },
                    {
                      key: 'app_mode',
                      label: (
                        <span className="inline-flex items-center gap-1.5">
                          <DashboardOutlined /> 运行模式
                        </span>
                      ),
                      children: overview?.app_mode
                        ? `${MODE_TEXT[overview.app_mode] ?? overview.app_mode}（${overview.app_mode}）`
                        : '—',
                    },
                    {
                      key: 'go_version',
                      label: (
                        <span className="inline-flex items-center gap-1.5">
                          <CodeOutlined /> Go 版本
                        </span>
                      ),
                      children: overview?.go_version || '—',
                    },
                    {
                      key: 'goroutines',
                      label: (
                        <span className="inline-flex items-center gap-1.5">
                          <ThunderboltOutlined /> Goroutine 数
                        </span>
                      ),
                      children: overview?.goroutines ?? '—',
                    },
                    {
                      key: 'mem_alloc',
                      label: (
                        <span className="inline-flex items-center gap-1.5">
                          <FundOutlined /> 内存使用
                        </span>
                      ),
                      children: overview ? `${overview.mem_alloc_mb} MB` : '—',
                    },
                    {
                      key: 'uptime',
                      label: (
                        <span className="inline-flex items-center gap-1.5">
                          <ClockCircleOutlined /> 运行时长
                        </span>
                      ),
                      children: overview?.uptime || '—',
                      span: 2,
                    },
                  ]}
                />
              </Card>
            </div>
          </Col>

          <Col xs={24} xl={10}>
            <div className="dash-reveal" style={{ animationDelay: '450ms' }}>
              <Card
                variant="borderless"
                title={
                  <span className="flex items-center gap-2">
                    <DatabaseOutlined style={{ color: '#b45309' }} />
                    依赖状态
                  </span>
                }
                className="h-full"
              >
                <div className="flex flex-col gap-3">
                  {dependencies.map((dep) => {
                    const status = overview?.[dep.key];
                    return (
                      <div
                        key={dep.key}
                        className="flex items-center justify-between px-4 py-3 rounded-lg border border-solid border-[#eee6d6]"
                        style={{ background: 'rgba(244, 239, 230, 0.5)' }}
                      >
                        <span className="flex items-center gap-3">
                          <span
                            className="w-9 h-9 flex items-center justify-center rounded-md text-lg"
                            style={{ background: 'rgba(180, 83, 9, 0.08)', color: '#b45309' }}
                          >
                            {dep.icon}
                          </span>
                          <span className="leading-tight">
                            <div className="font-medium text-[#292521]">{dep.name}</div>
                            <div className="text-xs text-[#8a7f6f]">{dep.desc}</div>
                          </span>
                        </span>
                        <StatusBadge status={status} />
                      </div>
                    );
                  })}
                </div>
                <p className="mt-4 mb-0 text-xs text-[#a69880]">
                  健康状态实时探测：绿色表示连接正常，红色表示连接异常。
                </p>
              </Card>
            </div>
          </Col>
        </Row>
      </Spin>
    </div>
  );
}
