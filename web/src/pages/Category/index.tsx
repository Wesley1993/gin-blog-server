import { useEffect, useMemo, useState } from 'react';
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Tag,
  Space,
  Popconfirm,
  Card,
  Upload,
  App as AntApp,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { UploadProps } from 'antd';
import { PlusOutlined, ReloadOutlined, PictureOutlined, DeleteOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import PageHeader from '../../components/PageHeader';
import {
  getCategoryTree,
  createCategory,
  updateCategory,
  deleteCategory,
} from '../../api/category';
import { uploadToOss } from '../../api/upload';
import { usePermission } from '../../hooks/usePermission';
import type { CategoryItem } from '../../api/types';

export default function CategoryPage() {
  const { message } = AntApp.useApp();
  const { hasPerm } = usePermission();
  const [tree, setTree] = useState<CategoryItem[]>([]);
  const [loading, setLoading] = useState(false);

  const [editOpen, setEditOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<CategoryItem | null>(null);
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);
  const [uploadingImage, setUploadingImage] = useState(false);
  const formImage = Form.useWatch('image', form);

  const fetchList = () => {
    setLoading(true);
    getCategoryTree()
      .then((res) => setTree(res.data ?? []))
      .catch(() => undefined)
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchList();
  }, []);

  /** 扁平化分类（用于父分类下拉），编辑时排除自身及其子孙 */
  const flatOptions = useMemo(() => {
    const options: { label: string; value: number }[] = [{ label: '顶级分类', value: 0 }];
    const excludeId = editTarget?.id;
    const walk = (nodes: CategoryItem[], depth: number) => {
      nodes.forEach((n) => {
        if (n.id === excludeId) return; // 排除自身及其子孙
        options.push({ label: `${' '.repeat(depth)}${n.name}`, value: n.id });
        if (n.children) walk(n.children, depth + 1);
      });
    };
    walk(tree, 0);
    return options;
  }, [tree, editTarget]);

  const openCreate = () => {
    setEditTarget(null);
    form.resetFields();
    form.setFieldsValue({ parent_id: 0, sort: 0, status: 1, image: '' });
    setEditOpen(true);
  };

  const openEdit = (record: CategoryItem) => {
    setEditTarget(record);
    form.setFieldsValue({
      name: record.name,
      parent_id: record.parent_id,
      sort: record.sort,
      status: record.status,
      image: record.image ?? '',
    });
    setEditOpen(true);
  };

  const handleSave = async () => {
    const values = await form.validateFields();
    setSaving(true);
    try {
      if (editTarget) {
        await updateCategory({ ...values, id: editTarget.id });
        message.success('分类已更新');
      } else {
        await createCategory(values);
        message.success('分类已创建');
      }
      setEditOpen(false);
      fetchList();
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSaving(false);
    }
  };

  /** 分类图片上传（复用文章封面上传模式：antd Upload + customRequest 调上传 API） */
  const imageUploadProps: UploadProps = {
    accept: 'image/jpeg,image/png,image/gif,image/webp',
    showUploadList: false,
    beforeUpload: (file) => {
      if (file.size > 10 * 1024 * 1024) {
        message.error('图片大小不能超过 10MB');
        return Upload.LIST_IGNORE;
      }
      return true;
    },
    customRequest: async ({ file, onSuccess, onError }) => {
      setUploadingImage(true);
      try {
        const res = await uploadToOss(file as File);
        form.setFieldValue('image', res.data.url);
        message.success('分类图片上传成功');
        onSuccess?.(res.data);
      } catch (err) {
        onError?.(err as Error);
      } finally {
        setUploadingImage(false);
      }
    },
  };

  const handleDelete = async (record: CategoryItem) => {
    try {
      await deleteCategory(record.id);
      message.success('分类已删除');
      fetchList();
    } catch {
      /* 被引用时后端返回错误，拦截器已提示 */
    }
  };

  const canEdit = hasPerm('category:edit');
  const canDelete = hasPerm('category:delete');

  const columns: ColumnsType<CategoryItem> = [
    { title: '分类名称', dataIndex: 'name', width: 240 },
    {
      title: '图片',
      dataIndex: 'image',
      width: 110,
      render: (v: string) =>
        v ? (
          <img
            src={v}
            alt="分类图片"
            className="h-10 w-16 object-cover rounded border border-solid border-[#e7e0d2]"
          />
        ) : (
          '—'
        ),
    },
    { title: '排序', dataIndex: 'sort', width: 100 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (v: number) =>
        v === 1 ? <Tag color="success">启用</Tag> : <Tag color="default">禁用</Tag>,
    },
    {
      title: '创建时间',
      dataIndex: 'create_time',
      width: 200,
      render: (v) => (v ? dayjs(v).format('YYYY-MM-DD HH:mm:ss') : '—'),
    },
  ];

  // 无任一写权限时不渲染操作列，避免空操作栏
  if (canEdit || canDelete) {
    columns.push({
      title: '操作',
      width: 200,
      render: (_, record) => (
        <Space size="small">
          {canEdit && (
            <Button size="small" onClick={() => openEdit(record)}>
              编辑
            </Button>
          )}
          {canDelete && (
            <Popconfirm
              title="确认删除该分类？"
              description="若分类下有文章，将被后端拦截。"
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
        title="分类管理"
        subtitle="文章分类的层级结构维护"
        extra={
          <Space>
            <Button icon={<ReloadOutlined />} onClick={fetchList}>
              刷新
            </Button>
            {hasPerm('category:add') && (
              <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
                新增分类
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
        title={editTarget ? '编辑分类' : '新增分类'}
        open={editOpen}
        onOk={handleSave}
        confirmLoading={saving}
        onCancel={() => setEditOpen(false)}
        destroyOnHidden
      >
        <Form form={form} layout="vertical" className="mt-4">
          <Form.Item
            name="name"
            label="分类名称"
            rules={[{ required: true, message: '请输入分类名称' }]}
          >
            <Input placeholder="如：技术、随笔" />
          </Form.Item>
          <Form.Item name="parent_id" label="父分类" rules={[{ required: true }]}>
            <Select options={flatOptions} placeholder="选择父分类" />
          </Form.Item>
          <Form.Item name="sort" label="排序" rules={[{ required: true }]}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            label="分类图片"
            extra="文章未上传封面时，将使用该图片作为默认封面"
          >
            <div className="flex items-center gap-3 mb-2">
              <Upload {...imageUploadProps}>
                <Button icon={<PictureOutlined />} loading={uploadingImage}>
                  上传图片
                </Button>
              </Upload>
              {formImage && (
                <div className="flex items-center gap-2">
                  <img
                    src={formImage}
                    alt="分类图片预览"
                    className="h-12 w-20 object-cover rounded border border-solid border-[#e7e0d2]"
                  />
                  <Button
                    size="small"
                    danger
                    icon={<DeleteOutlined />}
                    onClick={() => form.setFieldValue('image', '')}
                  >
                    移除
                  </Button>
                </div>
              )}
            </div>
            <Form.Item name="image" noStyle>
              <Input placeholder="图片 URL，也可通过上方按钮上传" />
            </Form.Item>
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
    </div>
  );
}
