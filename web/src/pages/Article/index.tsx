import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  Table,
  Button,
  Input,
  Select,
  Tag,
  Space,
  Popconfirm,
  Card,
  Image,
  App as AntApp,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PlusOutlined, ThunderboltOutlined, EditOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';
import PageHeader from '../../components/PageHeader';
import {
  getArticlePage,
  searchArticle,
  deleteArticle,
  rebuildIndex,
} from '../../api/article';
import { getCategoryTree } from '../../api/category';
import { usePermission } from '../../hooks/usePermission';
import type { ArticleItem, CategoryItem } from '../../api/types';

export default function ArticlePage() {
  const { message, modal } = AntApp.useApp();
  const navigate = useNavigate();
  const { hasPerm } = usePermission();

  const [list, setList] = useState<ArticleItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [loading, setLoading] = useState(false);

  const [keyword, setKeyword] = useState('');
  const [categoryId, setCategoryId] = useState<number | undefined>();
  const [status, setStatus] = useState<number | undefined>();

  const [categories, setCategories] = useState<CategoryItem[]>([]);

  /** 分类 id → 名称 */
  const categoryNameMap = useMemo(() => {
    const map = new Map<number, string>();
    const walk = (nodes: CategoryItem[]) => {
      nodes.forEach((n) => {
        map.set(n.id, n.name);
        if (n.children) walk(n.children);
      });
    };
    walk(categories);
    return map;
  }, [categories]);

  const fetchList = useCallback(async () => {
    setLoading(true);
    try {
      if (keyword.trim()) {
        // ES 全文搜索
        const res = await searchArticle({ keyword: keyword.trim(), page, page_size: pageSize });
        setList(res.data.list ?? []);
        setTotal(res.data.total ?? 0);
      } else {
        const res = await getArticlePage({
          page,
          page_size: pageSize,
          category_id: categoryId,
          status,
        });
        setList(res.data.list ?? []);
        setTotal(res.data.total ?? 0);
      }
    } catch {
      /* 拦截器已提示 */
    } finally {
      setLoading(false);
    }
  }, [keyword, page, pageSize, categoryId, status]);

  useEffect(() => {
    fetchList();
  }, [fetchList]);

  useEffect(() => {
    getCategoryTree()
      .then((res) => setCategories(res.data ?? []))
      .catch(() => undefined);
  }, []);

  const handleDelete = async (record: ArticleItem) => {
    try {
      await deleteArticle(record.id);
      message.success('文章已删除');
      fetchList();
    } catch {
      /* 拦截器已提示 */
    }
  };

  const handleRebuild = () => {
    modal.confirm({
      title: '重建 ES 索引',
      content: '将删除现有索引并从数据库全量重建，耗时较长，是否继续？',
      okText: '重建',
      cancelText: '取消',
      onOk: async () => {
        await rebuildIndex();
        message.success('索引重建完成');
      },
    });
  };

  const canEdit = hasPerm('article:edit');
  const canDelete = hasPerm('article:delete');

  const columns: ColumnsType<ArticleItem> = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    {
      title: '标题',
      dataIndex: 'title',
      ellipsis: true,
      render: (v: string, record) => (
        <Space>
          {record.cover && (
            <Image src={record.cover} width={40} height={28} style={{ objectFit: 'cover', borderRadius: 4 }} />
          )}
          <span className="font-medium">{v}</span>
          {record.is_repost === 1 && <Tag color="orange">转载</Tag>}
        </Space>
      ),
    },
    {
      title: '分类',
      dataIndex: 'category_id',
      width: 120,
      render: (v: number) => <Tag>{categoryNameMap.get(v) ?? `#${v}`}</Tag>,
    },
    {
      title: '标签',
      dataIndex: 'tags',
      width: 200,
      render: (v?: string) =>
        v
          ? v.split(',').map((t) => (
              <Tag key={t} color="default">
                {t.trim()}
              </Tag>
            ))
          : '—',
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (v: number) =>
        v === 1 ? <Tag color="success">已发布</Tag> : <Tag color="warning">草稿</Tag>,
    },
    {
      title: '发布时间',
      dataIndex: 'published_at',
      key: 'published_at',
      width: 160,
      render: (v: string | undefined, record) =>
        v
          ? dayjs(v).format('YYYY-MM-DD HH:mm')
          : record.create_time
            ? dayjs(record.create_time).format('YYYY-MM-DD HH:mm')
            : '—',
    },
    {
      title: '创建时间',
      dataIndex: 'create_time',
      width: 170,
      render: (v?: string) => (v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '—'),
    },
  ];

  // 无任一行内写权限时不渲染操作列，避免空操作栏；
  // 重建索引按钮由 article:rebuild 权限控制（后端已登记该写路由权限）
  if (canEdit || canDelete) {
    columns.push({
      title: '操作',
      width: 150,
      render: (_, record) => (
        <Space size="small">
          {canEdit && (
            <Button
              size="small"
              icon={<EditOutlined />}
              onClick={() => navigate(`/article/edit/${record.id}`)}
            >
              编辑
            </Button>
          )}
          {canDelete && (
            <Popconfirm title="确认删除该文章？" onConfirm={() => handleDelete(record)}>
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
        title="文章管理"
        subtitle="内容发布、全文检索与索引维护"
        extra={
          <Space>
            {hasPerm('article:rebuild') && (
              <Button icon={<ThunderboltOutlined />} onClick={handleRebuild}>
                重建索引
              </Button>
            )}
            {hasPerm('article:add') && (
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => navigate('/article/edit')}
              >
                写文章
              </Button>
            )}
          </Space>
        }
      />

      <Card variant="borderless">
        <div className="toolbar-row">
          <Input.Search
            placeholder="全文搜索（标题 / 正文 / 标签）"
            allowClear
            style={{ width: 320 }}
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onSearch={() => {
              setPage(1);
            }}
            enterButton="搜索"
          />
          <Select
            placeholder="按分类筛选"
            allowClear
            style={{ width: 160 }}
            value={categoryId}
            disabled={Boolean(keyword.trim())}
            onChange={(v) => {
              setCategoryId(v);
              setPage(1);
            }}
            options={categories.map((c) => ({ label: c.name, value: c.id }))}
          />
          <Select
            placeholder="按状态筛选"
            allowClear
            style={{ width: 140 }}
            value={status}
            disabled={Boolean(keyword.trim())}
            onChange={(v) => {
              setStatus(v);
              setPage(1);
            }}
            options={[
              { label: '已发布', value: 1 },
              { label: '草稿', value: 0 },
            ]}
          />
          {keyword.trim() && (
            <Tag color="processing">搜索模式：{keyword.trim()}</Tag>
          )}
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
    </div>
  );
}
