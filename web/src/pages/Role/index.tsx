import { useEffect, useMemo, useState } from 'react';
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Tree,
  Tag,
  Space,
  Popconfirm,
  Card,
  Checkbox,
  App as AntApp,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { DataNode } from 'antd/es/tree';
import { PlusOutlined } from '@ant-design/icons';
import PageHeader from '../../components/PageHeader';
import { getRoleList, createRole, updateRole, deleteRole } from '../../api/role';
import { getMenuList } from '../../api/menu';
import { usePermission } from '../../hooks/usePermission';
import type { Role, MenuItem } from '../../api/types';

/** 菜单树 → antd Tree 数据（按钮类型作为叶子节点，供勾选权限） */
function toTreeData(nodes: MenuItem[]): DataNode[] {
  return nodes.map((n) => ({
    key: n.id,
    title: n.menu_type === 3 ? `${n.menu_name}（按钮）` : n.menu_name,
    children: n.children && n.children.length > 0 ? toTreeData(n.children) : undefined,
  }));
}

/** 收集树中全部节点 key（用于"全选"） */
function collectKeys(nodes: MenuItem[]): number[] {
  const keys: number[] = [];
  const walk = (list: MenuItem[]) => {
    list.forEach((n) => {
      keys.push(n.id);
      if (n.children) walk(n.children);
    });
  };
  walk(nodes);
  return keys;
}

/** 扫平菜单树，取全部按钮节点（menu_type=3 且 perms 非空）的权限标识选项 */
function collectButtonPerms(nodes: MenuItem[]): { label: string; value: string }[] {
  const options: { label: string; value: string }[] = [];
  const walk = (list: MenuItem[]) => {
    list.forEach((n) => {
      if (n.menu_type === 3 && n.perms) {
        options.push({ label: `${n.menu_name}（${n.perms}）`, value: n.perms });
      }
      if (n.children) walk(n.children);
    });
  };
  walk(nodes);
  return options;
}

export default function RolePage() {
  const { message } = AntApp.useApp();
  const { hasPerm } = usePermission();
  const [list, setList] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [menus, setMenus] = useState<MenuItem[]>([]);
  const [allMenuKeys, setAllMenuKeys] = useState<number[]>([]);

  const [editOpen, setEditOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Role | null>(null);
  const [form] = Form.useForm();
  const [checkedMenus, setCheckedMenus] = useState<number[]>([]);
  const [perms, setPerms] = useState<string[]>([]);
  const [saving, setSaving] = useState(false);

  /** 全部按钮权限选项（来自菜单树 menu_type=3 节点） */
  const buttonPermOptions = useMemo(() => collectButtonPerms(menus), [menus]);
  const knownPermSet = useMemo(() => new Set(buttonPermOptions.map((o) => o.value)), [buttonPermOptions]);

  const fetchList = () => {
    setLoading(true);
    getRoleList()
      .then((res) => setList(res.data ?? []))
      .catch(() => undefined)
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchList();
    getMenuList()
      .then((res) => {
        setMenus(res.data ?? []);
        setAllMenuKeys(collectKeys(res.data ?? []));
      })
      .catch(() => undefined);
  }, []);

  const openCreate = () => {
    setEditTarget(null);
    form.resetFields();
    setCheckedMenus([]);
    setPerms([]);
    setEditOpen(true);
  };

  const openEdit = (record: Role) => {
    setEditTarget(record);
    form.setFieldsValue({ role_name: record.role_name });
    setCheckedMenus(record.menu_ids ?? []);
    setPerms(record.button_perms ?? []);
    setEditOpen(true);
  };

  /** Checkbox 变更：更新已勾选标识，同时保留不在选项列表中的历史自定义标识，避免数据丢失 */
  const handlePermCheck = (checked: string[]) => {
    const legacy = perms.filter((p) => !knownPermSet.has(p));
    setPerms([...checked, ...legacy]);
  };

  const handleSave = async () => {
    const values = await form.validateFields();
    setSaving(true);
    try {
      const payload = {
        role_name: values.role_name,
        menu_ids: checkedMenus,
        button_perms: perms,
      };
      if (editTarget) {
        await updateRole({ ...payload, id: editTarget.id });
        message.success('角色已更新');
      } else {
        await createRole(payload);
        message.success('角色已创建');
      }
      setEditOpen(false);
      fetchList();
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (record: Role) => {
    try {
      await deleteRole(record.id);
      message.success('角色已删除');
      fetchList();
    } catch {
      /* 拦截器已提示 */
    }
  };

  const canEdit = hasPerm('role:edit');
  const canDelete = hasPerm('role:delete');

  const columns: ColumnsType<Role> = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    {
      title: '角色名称',
      dataIndex: 'role_name',
      width: 180,
      render: (v: string, record) => (
        <Space>
          <span className="font-medium">{v}</span>
          {record.is_super === 1 && <Tag color="gold">超级管理员</Tag>}
        </Space>
      ),
    },
    {
      title: '菜单权限',
      dataIndex: 'menu_ids',
      render: (v: number[] | null) => (v && v.length > 0 ? `已选 ${v.length} 项` : '—'),
    },
    {
      title: '按钮权限',
      dataIndex: 'button_perms',
      render: (v: string[] | null) =>
        v && v.length > 0 ? (
          <Space size={[0, 4]} wrap>
            {v.map((p) => (
              <Tag key={p}>{p}</Tag>
            ))}
          </Space>
        ) : (
          '—'
        ),
    },
    { title: '创建时间', dataIndex: 'create_time', width: 180, render: (v) => v || '—' },
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
              title="确认删除该角色？"
              description="删除后绑定该角色的用户将失去权限。"
              disabled={record.is_super === 1}
              onConfirm={() => handleDelete(record)}
            >
              <Button size="small" danger disabled={record.is_super === 1}>
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
        title="角色管理"
        subtitle="角色、菜单权限与按钮权限标识的分配"
        extra={
          hasPerm('role:add') ? (
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
              新增角色
            </Button>
          ) : null
        }
      />

      <Card variant="borderless">
        <Table rowKey="id" columns={columns} dataSource={list} loading={loading} pagination={false} />
      </Card>

      <Modal
        title={editTarget ? '编辑角色' : '新增角色'}
        open={editOpen}
        width={640}
        onOk={handleSave}
        confirmLoading={saving}
        onCancel={() => setEditOpen(false)}
        destroyOnHidden
      >
        <Form form={form} layout="vertical" className="mt-4">
          <Form.Item
            name="role_name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          >
            <Input placeholder="如：编辑、运营" />
          </Form.Item>

          <Form.Item label="菜单权限">
            <div className="flex justify-end mb-2">
              <Space>
                <Button size="small" onClick={() => setCheckedMenus(allMenuKeys)}>
                  全选
                </Button>
                <Button size="small" onClick={() => setCheckedMenus([])}>
                  清空
                </Button>
              </Space>
            </div>
            <div className="border border-solid border-[#e7e0d2] rounded-lg p-3 max-h-64 overflow-auto">
              {menus.length > 0 ? (
                <Tree
                  checkable
                  treeData={toTreeData(menus)}
                  checkedKeys={checkedMenus}
                  onCheck={(keys) => setCheckedMenus(keys as number[])}
                  defaultExpandAll
                />
              ) : (
                <span className="text-sm text-[#8a7f6f]">暂无菜单数据</span>
              )}
            </div>
          </Form.Item>

          <Form.Item label="按钮权限标识" extra="勾选菜单树中已登记的按钮权限；提交仍为字符串数组契约">
            {buttonPermOptions.length > 0 ? (
              <div className="border border-solid border-[#e7e0d2] rounded-lg p-3 max-h-64 overflow-auto">
                <Checkbox.Group
                  options={buttonPermOptions}
                  value={perms.filter((p) => knownPermSet.has(p))}
                  onChange={(values) => handlePermCheck(values as string[])}
                />
              </div>
            ) : (
              <span className="text-sm text-[#8a7f6f]">暂无可选按钮权限（请先在菜单管理中登记）</span>
            )}
            {perms.some((p) => !knownPermSet.has(p)) && (
              <div className="mt-2">
                <span className="text-sm text-[#8a7f6f] mr-2">未登记的历史标识（保留提交）：</span>
                <Space size={[4, 8]} wrap>
                  {perms
                    .filter((p) => !knownPermSet.has(p))
                    .map((p) => (
                      <Tag key={p}>{p}</Tag>
                    ))}
                </Space>
              </div>
            )}
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
