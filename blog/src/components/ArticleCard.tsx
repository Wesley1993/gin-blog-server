import { useNavigate } from 'react-router-dom';
import type { Article, Category } from '../api/blog';

interface Props {
  article: Article;
  /** 序号，从 1 开始 */
  index: number;
  categories: Category[];
}

const formatDate = (s: string) => (s ? s.slice(0, 10) : '');

/** 简介：优先用 summary 字段；列表未返回时降级为标签，再无则用占位文案 */
function introOf(article: Article): string {
  if (article.summary?.trim()) return article.summary.trim();
  const tags = article.tags
    ? article.tags.split(',').map((t) => t.trim()).filter(Boolean)
    : [];
  if (tags.length > 0) return tags.map((t) => `# ${t}`).join('　');
  return '◆ 一篇尚未贴签的札记，展开正文自见分晓。';
}

/** 目录式文章卡片：左封面 + 右条目 */
export default function ArticleCard({ article, index, categories }: Props) {
  const navigate = useNavigate();

  const catName = categories.find((c) => c.id === article.category_id)?.name || '';

  const glyph = article.title?.trim().charAt(0) || '文';

  return (
    <article
      className="toc-card"
      style={{ animationDelay: `${Math.min(index, 8) * 0.06}s` }}
      onClick={() => navigate(`/article/${article.id}`)}
    >
      <div className="toc-index">{String(index).padStart(2, '0')}</div>

      <div className="toc-cover">
        {article.cover ? (
          <img src={article.cover} alt={article.title} loading="lazy" />
        ) : (
          <span className="toc-cover-glyph" aria-hidden>
            <span className="toc-cover-note">無圖</span>
            {glyph}
          </span>
        )}
      </div>

      <div>
        <div className="toc-meta">
          <span>{formatDate(article.published_at || article.create_time)}</span>
          {catName && <span className="cat">{catName}</span>}
          {article.is_repost === 1 && <span className="repost-mark">转载</span>}
        </div>
        <h2 className="toc-title">{article.title}</h2>
        <div className="toc-intro">{introOf(article)}</div>
        <div className="toc-enter">READ →</div>
      </div>
    </article>
  );
}
