import { useCallback, useEffect, useState } from 'react';
import {
  Table,
  Button,
  Input,
  Modal,
  Form,
  Select,
  Switch,
  Space,
  Popconfirm,
  Tag,
  Card,
  App as AntApp,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PlusOutlined, KeyOutlined, ReloadOutlined } from '@ant-design/icons';
import PageHeader from '../../components/PageHeader';
import {
  getUserPage,
  createUser,
  updateUser,
  resetPassword,
  deleteUser,
  updateUserStatus,
} from '../../api/user';
import { getRoleList } from '../../api/role';
import { usePermission } from '../../hooks/usePermission';
import type { UserItem, Role } from '../../api/types';

export default function UserPage() {
  const { message } = AntApp.useApp();
  const { hasPerm } = usePermission();
  const [list, setList] = useState<UserItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState('');
  const [loading, setLoading] = useState(false);
  const [roles, setRoles] = useState<Role[]>([]);

  // 新增/编辑弹窗
  const [editOpen, setEditOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<UserItem | null>(null);
  const [editForm] = Form.useForm();
  const [saving, setSaving] = useState(false);

  // 重置密码弹窗
  const [pwdOpen, setPwdOpen] = useState(false);
  const [pwdTarget, setPwdTarget] = useState<UserItem | null>(null);
  const [pwdForm] = Form.useForm();

  const fetchList = useCallback(async () => {
    setLoading(true);
    try {
      const res = await getUserPage({ page, page_size: pageSize, username: keyword || undefined });
      setList(res.data.list ?? []);
      setTotal(res.data.total ?? 0);
    } catch {
      /* 拦截器已提示 */
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, keyword]);

  useEffect(() => {
    fetchList();
  }, [fetchList]);

  useEffect(() => {
    getRoleList()
      .then((res) => setRoles(res.data ?? []))
      .catch(() => undefined);
  }, []);

  const roleName = (roleId: number) =>
    roles.find((r) => r.id === roleId)?.role_name ?? `#${roleId}`;

  const openCreate = () => {
    setEditTarget(null);
    editForm.resetFields();
    editForm.setFieldsValue({ status: 1 });
    setEditOpen(true);
  };

  const openEdit = (record: UserItem) => {
    setEditTarget(record);
    editForm.setFieldsValue({
      username: record.username,
      nickname: record.nickname,
      role_id: record.role_id,
      status: record.status,
    });
    setEditOpen(true);
  };

  const handleSave = async () => {
    const values = await editForm.validateFields();
    setSaving(true);
    try {
      if (editTarget) {
        await updateUser({
          id: editTarget.id,
          nickname: values.nickname,
          role_id: values.role_id,
          status: values.status,
        });
        message.success('用户已更新');
      } else {
        await createUser(values);
        message.success('用户已创建');
      }
      setEditOpen(false);
      fetchList();
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSaving(false);
    }
  };

  const handleResetPwd = async () => {
    const values = await pwdForm.validateFields();
    if (!pwdTarget) return;
    try {
      await resetPassword(pwdTarget.id, values.password);
      message.success('密码已重置');
      setPwdOpen(false);
    } catch {
      /* 拦截器已提示 */
    }
  };

  const handleDelete = async (record: UserItem) => {
    try {
      await deleteUser(record.id);
      message.success('用户已删除');
      fetchList();
    } catch {
      /* 拦截器已提示 */
    }
  };

  const handleToggleStatus = async (record: UserItem, checked: boolean) => {
    try {
      await updateUserStatus(record.id, checked ? 1 : 0);
      message.success(checked ? '已启用' : '已禁用');
      fetchList();
    } catch {
      /* 拦截器已提示 */
    }
  };

  const canEdit = hasPerm('user:edit');
  const canDelete = hasPerm('user:delete');
  const canResetPwd = hasPerm('user:resetPwd');
  const canToggleStatus = hasPerm('user:status');

  const columns: ColumnsType<UserItem> = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '用户名', dataIndex: 'username', width: 160 },
    { title: '昵称', dataIndex: 'nickname', width: 160, render: (v) => v || '—' },
    {
      title: '角色',
      dataIndex: 'role_id',
      width: 160,
      render: (v: number) => <Tag color="processing">{roleName(v)}</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: number, record) => (
        <Switch
          checked={status === 1}
          checkedChildren="启用"
          unCheckedChildren="禁用"
          disabled={!canToggleStatus}
          onChange={(checked) => handleToggleStatus(record, checked)}
        />
      ),
    },
    { title: '创建时间', dataIndex: 'create_time', width: 180, render: (v) => v || '—' },
  ];

  // 无任一行内写权限时不渲染操作列，避免空操作栏
  if (canEdit || canDelete || canResetPwd) {
    columns.push({
      title: '操作',
      width: 220,
      render: (_, record) => (
        <Space size="small">
          {canEdit && (
            <Button size="small" onClick={() => openEdit(record)}>
              编辑
            </Button>
          )}
          {canResetPwd && (
            <Button
              size="small"
              icon={<KeyOutlined />}
              onClick={() => {
                setPwdTarget(record);
                pwdForm.resetFields();
                setPwdOpen(true);
              }}
            >
              重置密码
            </Button>
          )}
          {canDelete && (
            <Popconfirm
              title="确认删除该用户？"
              description="超管账号不可删除，删除后不可恢复。"
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
        title="用户管理"
        subtitle="账号的创建、角色绑定与状态控制"
        extra={
          hasPerm('user:add') ? (
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
              新增用户
            </Button>
          ) : null
        }
      />

      <Card variant="borderless">
        <div className="toolbar-row">
          <Input.Search
            placeholder="搜索用户名"
            allowClear
            style={{ width: 240 }}
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onSearch={() => {
              setPage(1);
              fetchList();
            }}
          />
          <Button
            icon={<ReloadOutlined />}
            onClick={() => {
              setKeyword('');
              setPage(1);
            }}
          >
            重置
          </Button>
        </div>

        <Table
          rowKey="id"
          columns={columns}
          dataSource={list}
          loading={loading}
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: true,
            showTotal: (t) => `共 ${t} 条`,
            onChange: (p, ps) => {
              setPage(p);
              setPageSize(ps);
            },
          }}
        />
      </Card>

      {/* 新增 / 编辑弹窗 */}
      <Modal
        title={editTarget ? '编辑用户' : '新增用户'}
        open={editOpen}
        onOk={handleSave}
        confirmLoading={saving}
        onCancel={() => setEditOpen(false)}
        destroyOnHidden
      >
        <Form form={editForm} layout="vertical" className="mt-4">
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input disabled={Boolean(editTarget)} placeholder="登录账号" />
          </Form.Item>
          {!editTarget && (
            <Form.Item
              name="password"
              label="密码"
              rules={[
                { required: true, message: '请输入密码' },
                { min: 6, message: '密码至少 6 位' },
              ]}
            >
              <Input.Password placeholder="初始密码" />
            </Form.Item>
          )}
          <Form.Item name="nickname" label="昵称">
            <Input placeholder="显示昵称" />
          </Form.Item>
          <Form.Item
            name="role_id"
            label="角色"
            rules={[{ required: true, message: '请选择角色' }]}
          >
            <Select
              placeholder="选择角色"
              options={roles.map((r) => ({ label: r.role_name, value: r.id }))}
            />
          </Form.Item>
          <Form.Item name="status" label="状态" rules={[{ required: true }]}>
            <Select
              options={[
                { label: '启用', value: 1 },
                { label: '禁用', value: 0 },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>

      {/* 重置密码弹窗 */}
      <Modal
        title={`重置密码 — ${pwdTarget?.username ?? ''}`}
        open={pwdOpen}
        onOk={handleResetPwd}
        onCancel={() => setPwdOpen(false)}
        destroyOnHidden
      >
        <Form form={pwdForm} layout="vertical" className="mt-4">
          <Form.Item
            name="password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码至少 6 位' },
            ]}
          >
            <Input.Password placeholder="新密码" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
