import { useEffect, useState } from 'react';
import {
  Tabs,
  Form,
  Input,
  Button,
  Card,
  Space,
  Spin,
  DatePicker,
  Statistic,
  Row,
  Col,
  Select,
  Checkbox,
  App as AntApp,
} from 'antd';
import {
  SaveOutlined,
  ApiOutlined,
  CalendarOutlined,
  FieldTimeOutlined,
  FileTextOutlined,
  AppstoreOutlined,
  TeamOutlined,
  CloudServerOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import PageHeader from '../../components/PageHeader';
import { getSiteConfig, saveSiteConfig, testOss, getSiteStats } from '../../api/site';
import type { SiteConfig, SiteStats } from '../../api/types';

export default function SitePage() {
  const { message } = AntApp.useApp();
  const [baseForm] = Form.useForm();
  const [ossForm] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [savingBase, setSavingBase] = useState(false);
  const [savingOss, setSavingOss] = useState(false);
  const [testing, setTesting] = useState(false);
  const [config, setConfig] = useState<SiteConfig>({});
  const [stats, setStats] = useState<SiteStats | null>(null);
  const [statsLoading, setStatsLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    getSiteConfig()
      .then((res) => {
        const data = res.data ?? {};
        setConfig(data);
        baseForm.setFieldsValue({
          site_name: data.site_name,
          site_desc: data.site_desc,
          copyright: data.copyright,
          founded_at: data.founded_at ? dayjs(data.founded_at, 'YYYY-MM-DD') : null,
          icp: data.icp,
          police_icp: data.police_icp,
        });
        ossForm.setFieldsValue({
          oss_access_key: data.oss_access_key,
          oss_secret_key: data.oss_secret_key,
          oss_bucket: data.oss_bucket,
          oss_endpoint: data.oss_endpoint,
          oss_domain: data.oss_domain,
          oss_provider: data.oss_provider || 'aliyun',
          oss_region: data.oss_region,
          oss_insecure: data.oss_insecure === 1,
        });
      })
      .catch(() => undefined)
      .finally(() => setLoading(false));

    setStatsLoading(true);
    getSiteStats()
      .then((res) => setStats(res.data))
      .catch(() => undefined)
      .finally(() => setStatsLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const saveBase = async () => {
    const values = await baseForm.validateFields();
    setSavingBase(true);
    try {
      await saveSiteConfig({
        ...config,
        ...values,
        founded_at: (values.founded_at as Dayjs | null)?.format('YYYY-MM-DD') ?? '',
      });
      message.success('基础设置已保存');
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSavingBase(false);
    }
  };

  const saveOss = async () => {
    const values = await ossForm.validateFields();
    // 表单中 Checkbox 为布尔值，提交时转为 0/1 与后端契约对齐
    const payload = { ...values, oss_insecure: values.oss_insecure ? 1 : 0 };
    setSavingOss(true);
    try {
      await saveSiteConfig({ ...config, ...payload });
      setConfig((prev) => ({ ...prev, ...payload }));
      message.success('OSS 设置已保存');
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSavingOss(false);
    }
  };

  const handleTest = async () => {
    setTesting(true);
    try {
      await testOss();
      message.success('OSS 连通正常');
    } catch {
      /* 拦截器已提示 */
    } finally {
      setTesting(false);
    }
  };

  const ossProvider = Form.useWatch('oss_provider', ossForm);
  // RustFS 的 region 为可选（通常不严格校验），仅展示提示不强制必填；endpoint 必填由 rules 控制
  const needRegion = ossProvider === 'tencent' || ossProvider === 's3' || ossProvider === 'rustfs';

  return (
    <div className="page-container">
      <PageHeader title="站点设置" subtitle="站点基础信息与对象存储（OSS）配置" />

      <Card variant="borderless">
        <Spin spinning={loading}>
          <Tabs
            defaultActiveKey="base"
            items={[
              {
                key: 'base',
                label: '基础设置',
                children: (
                  <div className="max-w-xl pt-4">
                    <Form form={baseForm} layout="vertical">
                      <Form.Item name="site_name" label="站点名称">
                        <Input placeholder="如：我的博客" />
                      </Form.Item>
                      <Form.Item name="site_desc" label="站点描述">
                        <Input.TextArea rows={3} placeholder="一句话描述你的站点" />
                      </Form.Item>
                      <Form.Item name="copyright" label="版权信息">
                        <Input placeholder="如：©2026 墨阁" />
                      </Form.Item>
                      <Form.Item
                        name="founded_at"
                        label="建站日期"
                        extra="用于计算站点运行天数，格式：年-月-日"
                      >
                        <DatePicker
                          style={{ width: '100%' }}
                          format="YYYY-MM-DD"
                          placeholder="选择建站日期"
                          allowClear
                        />
                      </Form.Item>
                      <Form.Item
                        name="icp"
                        label="ICP 备案号"
                        extra="展示于站点底部，如：京ICP备XXXXXXXX号"
                      >
                        <Input placeholder="如：京ICP备XXXXXXXX号" maxLength={100} allowClear />
                      </Form.Item>
                      <Form.Item
                        name="police_icp"
                        label="公安备案号"
                        extra="展示于站点底部，如：京公网安备 11010502000000号"
                      >
                        <Input
                          placeholder="如：京公网安备 11010502000000号"
                          maxLength={100}
                          allowClear
                        />
                      </Form.Item>
                      <Form.Item>
                        <Button
                          type="primary"
                          icon={<SaveOutlined />}
                          loading={savingBase}
                          onClick={saveBase}
                        >
                          保存基础设置
                        </Button>
                      </Form.Item>
                    </Form>
                  </div>
                ),
              },
              {
                key: 'stats',
                label: '站点运行信息',
                children: (
                  <Spin spinning={statsLoading}>
                    <div className="pt-4">
                      <Row gutter={[16, 16]}>
                        <Col xs={24} sm={12} lg={8}>
                          <Card variant="borderless">
                            <Statistic
                              title="运行天数"
                              value={stats?.running_days ?? 0}
                              suffix="天"
                              prefix={<FieldTimeOutlined />}
                              valueStyle={{ color: '#b45309' }}
                            />
                          </Card>
                        </Col>
                        <Col xs={24} sm={12} lg={8}>
                          <Card variant="borderless">
                            <Statistic
                              title="建站日期"
                              value={stats?.founded_at || '未设置'}
                              prefix={<CalendarOutlined />}
                            />
                          </Card>
                        </Col>
                        <Col xs={24} sm={12} lg={8}>
                          <Card variant="borderless">
                            <Statistic
                              title="文章总数"
                              value={stats?.article_count ?? 0}
                              prefix={<FileTextOutlined />}
                            />
                          </Card>
                        </Col>
                        <Col xs={24} sm={12} lg={8}>
                          <Card variant="borderless">
                            <Statistic
                              title="分类数"
                              value={stats?.category_count ?? 0}
                              prefix={<AppstoreOutlined />}
                            />
                          </Card>
                        </Col>
                        <Col xs={24} sm={12} lg={8}>
                          <Card variant="borderless">
                            <Statistic
                              title="用户数"
                              value={stats?.user_count ?? 0}
                              prefix={<TeamOutlined />}
                            />
                          </Card>
                        </Col>
                      </Row>
                    </div>
                  </Spin>
                ),
              },
              {
                key: 'oss',
                label: 'OSS 设置',
                children: (
                  <div className="max-w-xl pt-4">
                    <Form
                      form={ossForm}
                      layout="vertical"
                      initialValues={{ oss_provider: 'aliyun', oss_insecure: false }}
                    >
                      <Form.Item name="oss_provider" label="存储厂商">
                        <Select
                          options={[
                            { label: '阿里云', value: 'aliyun' },
                            { label: '腾讯云', value: 'tencent' },
                            { label: '七牛云', value: 'qiniu' },
                            { label: 'Amazon S3', value: 's3' },
                            { label: 'RustFS', value: 'rustfs' },
                          ]}
                        />
                      </Form.Item>
                      <Form.Item
                        name="oss_insecure"
                        valuePropName="checked"
                        extra="自签名证书或私有端点时启用"
                      >
                        <Checkbox>跳过 HTTPS 证书校验</Checkbox>
                      </Form.Item>
                      {needRegion && (
                        <Form.Item
                          name="oss_region"
                          label="Region"
                          extra={
                            ossProvider === 's3'
                              ? '如：us-east-1、ap-northeast-1'
                              : ossProvider === 'rustfs'
                                ? 'RustFS 通常不严格校验，可留空或填 us-east-1'
                                : '如：ap-guangzhou、ap-shanghai'
                          }
                        >
                          <Input
                            placeholder={
                              ossProvider === 's3'
                                ? '如：us-east-1'
                                : ossProvider === 'rustfs'
                                  ? '可留空或填 us-east-1'
                                  : '如：ap-guangzhou'
                            }
                          />
                        </Form.Item>
                      )}
                      <Form.Item name="oss_access_key" label="AccessKey ID">
                        <Input placeholder="AccessKey" autoComplete="off" />
                      </Form.Item>
                      <Form.Item name="oss_secret_key" label="AccessKey Secret">
                        <Input.Password placeholder="Secret" autoComplete="new-password" />
                      </Form.Item>
                      <Form.Item name="oss_bucket" label="Bucket">
                        <Input placeholder="存储空间名称" />
                      </Form.Item>
                      <Form.Item
                        name="oss_endpoint"
                        label="Endpoint"
                        rules={
                          ossProvider === 'rustfs'
                            ? [{ required: true, message: '请填写 RustFS 服务地址' }]
                            : undefined
                        }
                        extra={
                          ossProvider === 'rustfs'
                            ? '自建服务地址，如：http://192.168.1.100:9000'
                            : undefined
                        }
                      >
                        <Input
                          placeholder={
                            ossProvider === 'rustfs'
                              ? '如：http://192.168.1.100:9000'
                              : '如：oss-cn-hangzhou.aliyuncs.com'
                          }
                        />
                      </Form.Item>
                      <Form.Item name="oss_domain" label="自定义域名（可选）">
                        <Input placeholder="如：cdn.example.com" />
                      </Form.Item>
                      <Form.Item>
                        <Space>
                          <Button
                            type="primary"
                            icon={<CloudServerOutlined />}
                            loading={savingOss}
                            onClick={saveOss}
                          >
                            保存 OSS 设置
                          </Button>
                          <Button icon={<ApiOutlined />} loading={testing} onClick={handleTest}>
                            连通测试
                          </Button>
                        </Space>
                      </Form.Item>
                    </Form>
                  </div>
                ),
              },
            ]}
          />
        </Spin>
      </Card>
    </div>
  );
}
