import { useEffect, useMemo, useState } from 'react';
import {
  Table,
  Tag,
  Card,
  Button,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Switch,
  TreeSelect,
  Space,
  Popconfirm,
  App as AntApp,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import PageHeader from '../../components/PageHeader';
import { getMenuList, createMenu, updateMenu, deleteMenu } from '../../api/menu';
import { usePermission } from '../../hooks/usePermission';
import type { MenuItem } from '../../api/types';

const menuTypeMap: Record<number, { label: string; color: string }> = {
  1: { label: '目录', color: 'geekblue' },
  2: { label: '页面', color: 'green' },
  3: { label: '按钮', color: 'orange' },
};

export default function MenuPage() {
  const { message } = AntApp.useApp();
  const { hasPerm } = usePermission();
  const [tree, setTree] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState(false);

  const [editOpen, setEditOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<MenuItem | null>(null);
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);
  const menuType = Form.useWatch('menu_type', form);

  const fetchList = () => {
    setLoading(true);
    getMenuList()
      .then((res) => setTree(res.data ?? []))
      .catch(() => undefined)
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchList();
  }, []);

  /** 父级菜单树选项（编辑时排除自身及其子孙，避免成环） */
  const parentTree = useMemo(() => {
    const excludeId = editTarget?.id;
    const walk = (nodes: MenuItem[]): MenuItem[] =>
      nodes
        .filter((n) => n.id !== excludeId)
        .map((n) => ({
          ...n,
          children: n.children ? walk(n.children) : undefined,
        }));
    return walk(tree);
  }, [tree, editTarget]);

  const openCreate = () => {
    setEditTarget(null);
    form.resetFields();
    form.setFieldsValue({ menu_type: 1, sort: 0, status: true });
    setEditOpen(true);
  };

  const openEdit = (record: MenuItem) => {
    setEditTarget(record);
    form.setFieldsValue({
      menu_name: record.menu_name,
      parent_id: record.parent_id > 0 ? record.parent_id : undefined,
      menu_type: record.menu_type,
      path: record.path ?? '',
      perms: record.perms ?? '',
      sort: record.sort,
      status: record.status === 1,
    });
    setEditOpen(true);
  };

  const handleSave = async () => {
    const values = await form.validateFields();
    const payload = {
      menu_name: values.menu_name,
      parent_id: values.parent_id ?? 0,
      menu_type: values.menu_type,
      path: values.menu_type === 2 ? (values.path ?? '') : '',
      perms: values.menu_type === 3 ? (values.perms ?? '') : '',
      sort: values.sort ?? 0,
      status: values.status ? 1 : 0,
    };
    setSaving(true);
    try {
      if (editTarget) {
        await updateMenu({ ...payload, id: editTarget.id });
        message.success('菜单已更新');
      } else {
        await createMenu(payload);
        message.success('菜单已创建');
      }
      setEditOpen(false);
      fetchList();
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (record: MenuItem) => {
    try {
      await deleteMenu(record.id);
      message.success('菜单已删除');
      fetchList();
    } catch {
      /* 拦截器已提示 */
    }
  };

  const canEdit = hasPerm('menu:edit');
  const canDelete = hasPerm('menu:delete');

  const columns: ColumnsType<MenuItem> = [
    { title: '菜单名称', dataIndex: 'menu_name', width: 260 },
    {
      title: '类型',
      dataIndex: 'menu_type',
      width: 100,
      render: (v: number) => {
        const conf = menuTypeMap[v] ?? { label: String(v), color: 'default' };
        return <Tag color={conf.color}>{conf.label}</Tag>;
      },
    },
    { title: '路由路径', dataIndex: 'path', width: 200, render: (v) => v || '—' },
    { title: '权限标识', dataIndex: 'perms', width: 200, render: (v) => v || '—' },
    { title: '排序', dataIndex: 'sort', width: 80 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (v: number) =>
        v === 1 ? <Tag color="success">启用</Tag> : <Tag color="default">禁用</Tag>,
    },
  ];

  // 无任一写权限时不渲染操作列，避免空操作栏
  if (canEdit || canDelete) {
    columns.push({
      title: '操作',
      width: 160,
      render: (_, record) => (
        <Space size="small">
          {canEdit && (
            <Button size="small" onClick={() => openEdit(record)}>
              编辑
            </Button>
          )}
          {canDelete && (
            <Popconfirm
              title="确认删除该菜单？"
              description="删除后，角色中关联该菜单的权限将失效。"
              onConfirm={() => handleDelete(record)}
            >
              <Button size="small" danger>
                删除
              </Button>
            </Popconfirm>
          )}
        </Space>
      ),
    });
  }

  return (
    <div className="page-container">
      <PageHeader
        title="菜单管理"
        subtitle="系统菜单结构与按钮权限标识的维护"
        extra={
          <Space>
            <Button icon={<ReloadOutlined />} onClick={fetchList}>
              刷新
            </Button>
            {hasPerm('menu:add') && (
              <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
                新增菜单
              </Button>
            )}
          </Space>
        }
      />

      <Card variant="borderless">
        <Table
          rowKey="id"
          columns={columns}
          dataSource={tree}
          loading={loading}
          pagination={false}
          expandable={{ defaultExpandAllRows: true, childrenColumnName: 'children' }}
        />
      </Card>

      <Modal
        title={editTarget ? '编辑菜单' : '新增菜单'}
        width={640}
        open={editOpen}
        onOk={handleSave}
        confirmLoading={saving}
        onCancel={() => setEditOpen(false)}
        destroyOnHidden
      >
        <Form form={form} layout="vertical" className="mt-4">
          <Form.Item
            name="menu_name"
            label="菜单名称"
            rules={[{ required: true, message: '请输入菜单名称' }]}
          >
            <Input placeholder="如：文章管理" />
          </Form.Item>

          <Form.Item name="parent_id" label="父级菜单">
            <TreeSelect
              treeData={parentTree}
              fieldNames={{ label: 'menu_name', value: 'id', children: 'children' }}
              treeDefaultExpandAll
              allowClear
              placeholder="不选则为顶级菜单"
            />
          </Form.Item>

          <Form.Item
            name="menu_type"
            label="菜单类型"
            rules={[{ required: true, message: '请选择菜单类型' }]}
          >
            <Select
              options={[
                { label: '目录', value: 1 },
                { label: '页面', value: 2 },
                { label: '按钮', value: 3 },
              ]}
            />
          </Form.Item>

          {menuType === 2 && (
            <Form.Item name="path" label="路由路径">
              <Input placeholder="如：/articles" />
            </Form.Item>
          )}

          {menuType === 3 && (
            <Form.Item name="perms" label="权限标识">
              <Input placeholder="如：article:add" />
            </Form.Item>
          )}

          <Form.Item name="sort" label="排序">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item name="status" label="状态" valuePropName="checked">
            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
