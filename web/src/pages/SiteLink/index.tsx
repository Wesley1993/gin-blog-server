import { useEffect, useState } from 'react';
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  Switch,
  Space,
  Popconfirm,
  Card,
  Typography,
  App as AntApp,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import PageHeader from '../../components/PageHeader';
import {
  getSiteLinks,
  createSiteLink,
  updateSiteLink,
  deleteSiteLink,
} from '../../api/siteLink';
import type { SiteLink } from '../../api/types';

/** 常用网站：展示端侧栏「常用网站」卡片的数据维护 */
export default function SiteLinkPage() {
  const { message } = AntApp.useApp();
  const [list, setList] = useState<SiteLink[]>([]);
  const [loading, setLoading] = useState(false);

  const [editOpen, setEditOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<SiteLink | null>(null);
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);

  const fetchList = () => {
    setLoading(true);
    getSiteLinks()
      .then((res) => setList(res.data ?? []))
      .catch(() => undefined)
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchList();
  }, []);

  const openCreate = () => {
    setEditTarget(null);
    form.resetFields();
    form.setFieldsValue({ sort: 0, status: 1 });
    setEditOpen(true);
  };

  const openEdit = (record: SiteLink) => {
    setEditTarget(record);
    form.setFieldsValue({
      name: record.name,
      url: record.url,
      icon: record.icon,
      description: record.description,
      sort: record.sort,
      status: record.status,
    });
    setEditOpen(true);
  };

  const handleSave = async () => {
    const values = await form.validateFields();
    setSaving(true);
    try {
      const payload = {
        name: values.name,
        url: values.url,
        icon: values.icon ?? '',
        description: values.description ?? '',
        sort: values.sort ?? 0,
        status: values.status ?? 1,
      };
      if (editTarget) {
        await updateSiteLink(editTarget.id, payload);
        message.success('常用网站已更新');
      } else {
        await createSiteLink(payload);
        message.success('常用网站已创建');
      }
      setEditOpen(false);
      fetchList();
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (record: SiteLink) => {
    try {
      await deleteSiteLink(record.id);
      message.success('常用网站已删除');
      fetchList();
    } catch {
      /* 拦截器已提示 */
    }
  };

  /** 状态开关：直接调用更新接口切换启用/停用 */
  const handleToggleStatus = async (record: SiteLink, checked: boolean) => {
    try {
      await updateSiteLink(record.id, {
        name: record.name,
        url: record.url,
        icon: record.icon ?? '',
        description: record.description ?? '',
        sort: record.sort,
        status: checked ? 1 : 0,
      });
      message.success(checked ? '已启用' : '已停用');
      fetchList();
    } catch {
      /* 拦截器已提示 */
    }
  };

  const columns: ColumnsType<SiteLink> = [
    {
      title: '图标',
      dataIndex: 'icon',
      width: 60,
      render: (v: string) =>
        v ? (
          <img
            src={v}
            alt="图标"
            style={{
              width: 28,
              height: 28,
              borderRadius: 6,
              objectFit: 'cover',
              background: '#f5f5f5',
            }}
          />
        ) : (
          '—'
        ),
    },
    { title: '名称', dataIndex: 'name', width: 180 },
    {
      title: '链接',
      dataIndex: 'url',
      width: 260,
      ellipsis: true,
      render: (v: string) => (
        <Typography.Link href={v} target="_blank" rel="noreferrer">
          {v}
        </Typography.Link>
      ),
    },
    { title: '描述', dataIndex: 'description', ellipsis: true, render: (v) => v || '—' },
    { title: '排序', dataIndex: 'sort', width: 80 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (v: number, record) => (
        <Switch checked={v === 1} onChange={(checked) => handleToggleStatus(record, checked)} />
      ),
    },
    {
      title: '操作',
      width: 160,
      render: (_, record) => (
        <Space size="small">
          <Button size="small" onClick={() => openEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确认删除该常用网站？" onConfirm={() => handleDelete(record)}>
            <Button size="small" danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="page-container">
      <PageHeader
        title="常用网站"
        subtitle="展示端侧栏「常用网站」卡片的数据维护"
        extra={
          <Space>
            <Button icon={<ReloadOutlined />} onClick={fetchList}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
              新增网站
            </Button>
          </Space>
        }
      />

      <Card variant="borderless">
        <Table
          rowKey="id"
          columns={columns}
          dataSource={list}
          loading={loading}
          pagination={false}
        />
      </Card>

      <Modal
        title={editTarget ? '编辑常用网站' : '新增常用网站'}
        open={editOpen}
        onOk={handleSave}
        confirmLoading={saving}
        onCancel={() => setEditOpen(false)}
        destroyOnHidden
      >
        <Form form={form} layout="vertical" className="mt-4">
          <Form.Item
            name="name"
            label="网站名称"
            rules={[
              { required: true, message: '请输入网站名称' },
              { max: 50, message: '名称不能超过 50 个字符' },
            ]}
          >
            <Input placeholder="如：Go 官网" maxLength={50} />
          </Form.Item>
          <Form.Item
            name="url"
            label="链接地址"
            rules={[
              { required: true, message: '请输入链接地址' },
              { type: 'url', message: '请输入合法的 URL，如 https://go.dev' },
              { max: 255, message: '链接不能超过 255 个字符' },
            ]}
          >
            <Input placeholder="如：https://go.dev" maxLength={255} />
          </Form.Item>
          <Form.Item
            name="icon"
            label="图标"
            rules={[
              { type: 'url', message: '请输入合法的图标 URL，如 https://go.dev/favicon.ico' },
              { max: 500, message: '图标 URL 不能超过 500 个字符' },
            ]}
          >
            <Input placeholder="可选，图标图片 URL，如 https://go.dev/favicon.ico" maxLength={500} />
          </Form.Item>
          <Form.Item
            name="description"
            label="描述"
            rules={[{ max: 200, message: '描述不能超过 200 个字符' }]}
          >
            <Input placeholder="可选，一句话描述该网站" maxLength={200} />
          </Form.Item>
          <Form.Item name="sort" label="排序" extra="数值越小越靠前" rules={[{ required: true }]}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="status" label="状态" rules={[{ required: true }]}>
            <Select
              options={[
                { label: '启用', value: 1 },
                { label: '停用', value: 0 },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
