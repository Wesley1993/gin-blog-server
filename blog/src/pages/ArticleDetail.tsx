import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Result, Spin } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import MarkdownPreview from '@uiw/react-markdown-preview';
import { getArticleDetail, getArticles, getCategories } from '../api/blog';
import type { Article, Category } from '../api/blog';

const formatDate = (s: string) => (s ? s.slice(0, 10) : '');

/** 已发布文章列表缓存（不含 content，开销小），避免每次切换重复拉取 */
let articleOrderCache: Article[] | null = null;
let articleOrderPromise: Promise<Article[]> | null = null;

/** 拉取已发布文章全量列表（会话内缓存，请求去重；失败可重试） */
function loadArticleOrder(): Promise<Article[]> {
  if (articleOrderCache) return Promise.resolve(articleOrderCache);
  if (!articleOrderPromise) {
    articleOrderPromise = getArticles({ page: 1, page_size: 1000 })
      .then((res) => {
        articleOrderCache = res.list || [];
        return articleOrderCache;
      })
      .catch((err) => {
        articleOrderPromise = null; // 失败后允许重试，不缓存空结果
        throw err;
      });
  }
  return articleOrderPromise;
}

/** 估算阅读时长：去除 Markdown 语法后统计字数，按约 400 字/分钟计算 */
function readingMinutes(md: string): number {
  const text = (md || '')
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/`[^`\n]*`/g, ' ')
    .replace(/!\[[^\]]*\]\([^)]*\)/g, ' ')
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/^[-*+]\s+/gm, '')
    .replace(/^>\s?/gm, '')
    .replace(/[*_~|]/g, '')
    .replace(/<[^>]+>/g, ' ')
    .replace(/\s+/g, '');
  return Math.max(1, Math.ceil(text.length / 400));
}

export default function ArticleDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [article, setArticle] = useState<Article | null>(null);
  const [cats, setCats] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);
  /** 上一篇 / 下一篇（与首页列表同序：按创建时间排列） */
  const [prevArticle, setPrevArticle] = useState<Article | null>(null);
  const [nextArticle, setNextArticle] = useState<Article | null>(null);

  useEffect(() => {
    getCategories()
      .then((tree) => {
        const flat: Category[] = [];
        const walk = (nodes: Category[]) => {
          for (const n of nodes) {
            flat.push(n);
            if (n.children?.length) walk(n.children);
          }
        };
        walk(tree || []);
        setCats(flat);
      })
      .catch(() => setCats([]));
  }, []);

  useEffect(() => {
    if (!id) return;
    let cancelled = false;
    setLoading(true);
    setNotFound(false);
    getArticleDetail(id)
      .then((a) => {
        if (!cancelled) setArticle(a);
      })
      .catch(() => {
        if (!cancelled) setNotFound(true);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [id]);

  // 相邻文章：在已发布列表中定位当前文章，取前后相邻（路由参数变化时重新计算）
  useEffect(() => {
    if (!id) return;
    let cancelled = false;
    loadArticleOrder()
      .then((list) => {
        if (cancelled) return;
        const idx = list.findIndex((a) => String(a.id) === id);
        setPrevArticle(idx > 0 ? list[idx - 1] : null);
        setNextArticle(idx >= 0 && idx < list.length - 1 ? list[idx + 1] : null);
      })
      .catch(() => {
        if (!cancelled) {
          setPrevArticle(null);
          setNextArticle(null);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [id]);

  const gotoArticle = (target: Article | null) => {
    if (!target) return;
    navigate(`/article/${target.id}`);
  };

  const catName = article
    ? cats.find((c) => c.id === article.category_id)?.name || ''
    : '';
  const tags = article?.tags
    ? article.tags.split(',').map((t) => t.trim()).filter(Boolean)
    : [];

  if (loading) {
    return (
      <div className="container article-page">
        <div className="state-block">
          <Spin size="large" />
        </div>
      </div>
    );
  }

  if (notFound || !article) {
    return (
      <div className="container article-page">
        <Result
          status="404"
          title="文章不存在或尚未发布"
          subTitle="它可能被移走，或仍在编辑之中"
          extra={
            <button className="back-link" onClick={() => navigate('/')}>
              <ArrowLeftOutlined /> 返回首页
            </button>
          }
        />
      </div>
    );
  }

  return (
    <div className="container article-page">
      <button className="back-link" onClick={() => navigate(-1)}>
        <ArrowLeftOutlined /> 返回
      </button>

      <header className="article-head">
        <div className="meta">
          <span>{formatDate(article.create_time)}</span>
          {catName && <span className="cat">◆ {catName}</span>}
          <span>预计阅读 {readingMinutes(article.content)} 分钟</span>
        </div>
        <h1>{article.title}</h1>
        <div className="ornament">◆ ◆ ◆</div>

        {article.is_repost === 1 && (
          <div className="repost-banner">
            <span className="repost-badge">转载</span>
            <span className="repost-text">
              本文为转载作品{article.repost_author ? `，原作者：${article.repost_author}` : ''}
              {article.repost_url && (
                <>
                  ，
                  <a
                    href={article.repost_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="repost-link"
                  >
                    原文链接 ↗
                  </a>
                </>
              )}
            </span>
          </div>
        )}
      </header>

      <div className="article-body" data-color-mode="light">
        <MarkdownPreview source={article.content || ''} />
      </div>

      {tags.length > 0 && (
        <div className="article-tags">
          <span className="label">TAGS —</span>
          {tags.map((t) => (
            <span key={t}># {t}</span>
          ))}
        </div>
      )}

      {/* 上一篇 / 下一篇 */}
      <nav className="article-nav" aria-label="文章切换">
        <div
          className={`article-nav-card ${prevArticle ? '' : 'disabled'}`}
          role={prevArticle ? 'button' : undefined}
          tabIndex={prevArticle ? 0 : undefined}
          onClick={() => gotoArticle(prevArticle)}
          onKeyDown={(e) => e.key === 'Enter' && gotoArticle(prevArticle)}
        >
          <span className="nav-dir">← 上一篇</span>
          <span className="nav-title">{prevArticle ? prevArticle.title : '已是最早一篇'}</span>
          {prevArticle && <span className="nav-date mono">{formatDate(prevArticle.create_time)}</span>}
        </div>
        <div
          className={`article-nav-card right ${nextArticle ? '' : 'disabled'}`}
          role={nextArticle ? 'button' : undefined}
          tabIndex={nextArticle ? 0 : undefined}
          onClick={() => gotoArticle(nextArticle)}
          onKeyDown={(e) => e.key === 'Enter' && gotoArticle(nextArticle)}
        >
          <span className="nav-dir">下一篇 →</span>
          <span className="nav-title">{nextArticle ? nextArticle.title : '已是最新一篇'}</span>
          {nextArticle && <span className="nav-date mono">{formatDate(nextArticle.create_time)}</span>}
        </div>
      </nav>
    </div>
  );
}
