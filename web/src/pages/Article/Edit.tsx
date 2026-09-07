import { useEffect, useState } from 'react';
import {
  Button,
  Input,
  Select,
  Upload,
  Card,
  Space,
  Segmented,
  Spin,
  Switch,
  DatePicker,
  App as AntApp,
} from 'antd';
import type { UploadProps } from 'antd';
import { ArrowLeftOutlined, SaveOutlined, PictureOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import MDEditor from '@uiw/react-md-editor';
import PageHeader from '../../components/PageHeader';
import { getArticleById, createArticle, updateArticle } from '../../api/article';
import { getCategoryTree } from '../../api/category';
import { uploadToOss } from '../../api/upload';
import type { CategoryItem } from '../../api/types';

export default function ArticleEdit() {
  const { message } = AntApp.useApp();
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const articleId = id ? Number(id) : undefined;

  const [loading, setLoading] = useState(Boolean(articleId));
  const [saving, setSaving] = useState(false);
  const [title, setTitle] = useState('');
  const [categoryId, setCategoryId] = useState<number>();
  const [tags, setTags] = useState<string[]>([]);
  const [cover, setCover] = useState('');
  const [uploadingCover, setUploadingCover] = useState(false);
  const [content, setContent] = useState('');
  const [status, setStatus] = useState<number>(0);
  const [isRepost, setIsRepost] = useState(false);
  const [repostUrl, setRepostUrl] = useState('');
  const [repostAuthor, setRepostAuthor] = useState('');
  const [publishedAt, setPublishedAt] = useState<Dayjs | null>(null);
  const [categories, setCategories] = useState<CategoryItem[]>([]);

  // 编辑器高度随视口自适应，高分屏下充分利用纵向空间（最小 520）
  const [editorHeight, setEditorHeight] = useState(() =>
    Math.max(520, (typeof window !== 'undefined' ? window.innerHeight : 1080) - 480),
  );

  useEffect(() => {
    const onResize = () => setEditorHeight(Math.max(520, window.innerHeight - 480));
    window.addEventListener('resize', onResize);
    return () => window.removeEventListener('resize', onResize);
  }, []);

  useEffect(() => {
    getCategoryTree()
      .then((res) => setCategories(res.data ?? []))
      .catch(() => undefined);
  }, []);

  useEffect(() => {
    if (!articleId) return;
    setLoading(true);
    getArticleById(articleId)
      .then((res) => {
        const article = res.data;
        setTitle(article.title);
        setCategoryId(article.category_id);
        setTags(article.tags ? article.tags.split(',').map((t) => t.trim()).filter(Boolean) : []);
        setCover(article.cover ?? '');
        setContent(article.content ?? '');
        setStatus(article.status);
        setIsRepost(article.is_repost === 1);
        setRepostUrl(article.repost_url ?? '');
        setRepostAuthor(article.repost_author ?? '');
        if (article.published_at) {
          setPublishedAt(dayjs(article.published_at));
        }
      })
      .catch(() => {
        message.error('文章加载失败');
        navigate('/article');
      })
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [articleId]);

  const uploadProps: UploadProps = {
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
      setUploadingCover(true);
      try {
        const res = await uploadToOss(file as File);
        setCover(res.data.url);
        message.success('封面上传成功');
        onSuccess?.(res.data);
      } catch (err) {
        onError?.(err as Error);
      } finally {
        setUploadingCover(false);
      }
    },
  };

  const handleSave = async () => {
    if (!title.trim()) {
      message.warning('请输入文章标题');
      return;
    }
    if (!categoryId) {
      message.warning('请选择分类');
      return;
    }
    if (!content.trim()) {
      message.warning('请输入文章内容');
      return;
    }
    if (isRepost && !repostUrl.trim()) {
      message.warning('转载文章必须填写原文链接');
      return;
    }
    setSaving(true);
    try {
      const payload = {
        title: title.trim(),
        category_id: categoryId,
        cover: cover || undefined,
        content,
        tags: tags.join(','),
        status,
        is_repost: isRepost ? 1 : 0,
        repost_url: isRepost ? repostUrl.trim() : '',
        repost_author: isRepost ? repostAuthor.trim() : '',
        published_at: publishedAt ? publishedAt.format('YYYY-MM-DD HH:mm:ss') : '',
      };
      if (articleId) {
        await updateArticle({ ...payload, id: articleId });
        message.success('文章已更新');
      } else {
        await createArticle(payload);
        message.success('文章已创建');
      }
      navigate('/article');
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="page-container">
      <PageHeader
        title={articleId ? '编辑文章' : '撰写文章'}
        subtitle="Markdown 书写，发布后自动同步搜索索引"
        extra={
          <Space>
            <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/article')}>
              返回列表
            </Button>
            <Button type="primary" icon={<SaveOutlined />} loading={saving} onClick={handleSave}>
              保存
            </Button>
          </Space>
        }
      />

      <Spin spinning={loading}>
        <div className="flex flex-col gap-4">
          {/* 元信息区 */}
          <Card variant="borderless">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="md:col-span-3">
                <div className="text-xs text-[#8a7f6f] mb-1.5 tracking-wide">标题</div>
                <Input
                  size="large"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="文章标题"
                  maxLength={120}
                  showCount
                />
              </div>
              <div>
                <div className="text-xs text-[#8a7f6f] mb-1.5 tracking-wide">分类</div>
                <Select
                  style={{ width: '100%' }}
                  value={categoryId}
                  onChange={setCategoryId}
                  placeholder="选择分类"
                  options={categories.map((c) => ({ label: c.name, value: c.id }))}
                />
              </div>
              <div>
                <div className="text-xs text-[#8a7f6f] mb-1.5 tracking-wide">标签</div>
                <Select
                  mode="tags"
                  style={{ width: '100%' }}
                  value={tags}
                  onChange={setTags}
                  placeholder="输入后回车添加"
                  tokenSeparators={[',']}
                />
              </div>
              <div>
                <div className="text-xs text-[#8a7f6f] mb-1.5 tracking-wide">状态</div>
                <Segmented
                  value={status}
                  onChange={(v) => setStatus(Number(v))}
                  options={[
                    { label: '草稿', value: 0 },
                    { label: '发布', value: 1 },
                  ]}
                />
              </div>
              <div>
                <div className="text-xs text-[#8a7f6f] mb-1.5 tracking-wide">发布时间</div>
                <DatePicker
                  showTime
                  value={publishedAt}
                  onChange={(val) => setPublishedAt(val)}
                  format="YYYY-MM-DD HH:mm:ss"
                  placeholder="默认当前时间"
                  style={{ width: '100%' }}
                />
              </div>
              <div className="md:col-span-3">
                <div className="text-xs text-[#8a7f6f] mb-1.5 tracking-wide">封面图（≤10MB，jpg/png/gif/webp）</div>
                <div className="flex items-center gap-4">
                  <Upload {...uploadProps}>
                    <Button icon={<PictureOutlined />} loading={uploadingCover}>上传封面</Button>
                  </Upload>
                  {cover && (
                    <img
                      src={cover}
                      alt="封面预览"
                      className="h-16 w-28 object-cover rounded-md border border-solid border-[#e7e0d2]"
                    />
                  )}
                </div>
                <div className="text-xs text-[#b0a694] mt-1.5">未上传封面时，将使用所选分类的图片作为默认封面</div>
              </div>
              <div className="md:col-span-3">
                <div className="text-xs text-[#8a7f6f] mb-1.5 tracking-wide">转载信息</div>
                <div className="flex flex-col gap-3">
                  <Space>
                    <Switch
                      checked={isRepost}
                      onChange={(checked) => {
                        setIsRepost(checked);
                        if (!checked) {
                          // 关闭转载时清空原文链接与原作者
                          setRepostUrl('');
                          setRepostAuthor('');
                        }
                      }}
                    />
                    <span className="text-sm text-[#6d6357]">为转载文章</span>
                    {!isRepost && <span className="text-xs text-[#b0a695]">默认为原创</span>}
                  </Space>
                  {isRepost && (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <div className="text-xs text-[#8a7f6f] mb-1.5 tracking-wide">
                          原文链接 <span className="text-[#b45309]">*</span>
                        </div>
                        <Input
                          value={repostUrl}
                          onChange={(e) => setRepostUrl(e.target.value)}
                          placeholder="如：https://example.com/article"
                          maxLength={500}
                          allowClear
                        />
                      </div>
                      <div>
                        <div className="text-xs text-[#8a7f6f] mb-1.5 tracking-wide">原作者</div>
                        <Input
                          value={repostAuthor}
                          onChange={(e) => setRepostAuthor(e.target.value)}
                          placeholder="原文作者名称"
                          maxLength={100}
                          allowClear
                        />
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </Card>

          {/* 正文编辑器 */}
          <Card variant="borderless" className="md-editor-wrap" data-color-mode="light">
            <MDEditor
              value={content}
              onChange={(v) => setContent(v ?? '')}
              height={editorHeight}
              preview="live"
              textareaProps={{ placeholder: '开始书写正文……' }}
            />
          </Card>
        </div>
      </Spin>
    </div>
  );
}
