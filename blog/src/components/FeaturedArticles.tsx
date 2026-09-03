import { useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { Carousel } from 'antd';
import type { CarouselRef } from 'antd/es/carousel';
import type { Article, Category } from '../api/blog';

interface Props {
  articles: Article[];
  categories: Category[];
}

const formatDate = (s: string) => (s ? s.slice(0, 10) : '');

/** 简介：优先用 summary 字段；列表未返回时降级为标签，再无则用编者按占位 */
function introOf(article: Article): string {
  if (article.summary?.trim()) return article.summary.trim();
  const tags = article.tags
    ? article.tags.split(',').map((t) => t.trim()).filter(Boolean)
    : [];
  if (tags.length > 0) return tags.map((t) => `# ${t}`).join('　');
  return '◆ 一篇尚未贴签的札记，展开正文自见分晓。';
}

/** 封面：有图用图，无图则以墨底题字代之，保持纸墨气质 */
function Cover({ article, index }: { article: Article; index: number }) {
  if (article.cover) {
    return (
      <div className="featured-cover">
        <img src={article.cover} alt={article.title} loading={index === 1 ? 'eager' : 'lazy'} />
      </div>
    );
  }
  const glyph = article.title?.trim().charAt(0) || '文';
  return (
    <div className="featured-cover featured-cover--ink" aria-hidden>
      <span className="ink-glyph">{glyph}</span>
      <span className="ink-seal">NO.{String(index).padStart(2, '0')}</span>
    </div>
  );
}

/** 首页「最近精选」区块：画框式轮播，一屏 3 篇（窄屏降为 2/1 篇），自动播放、悬停暂停 */
export default function FeaturedArticles({ articles, categories }: Props) {
  const navigate = useNavigate();
  const carouselRef = useRef<CarouselRef>(null);

  if (articles.length === 0) return null;

  const catName = (id: number) => categories.find((c) => c.id === id)?.name || '';

  return (
    <section className="featured">
      <div className="section-head">
        <span className="zh">精选</span>
        <span className="rule" />
        <span className="meta">FEATURED · 最近的落笔</span>
      </div>
      <div className="featured-carousel">
        <Carousel
          ref={carouselRef}
          autoplay
          autoplaySpeed={5000}
          pauseOnHover
          dots={{ className: 'featured-dots' }}
          speed={550}
          slidesToShow={Math.min(3, articles.length)}
          slidesToScroll={Math.min(3, articles.length)}
          responsive={[
            { breakpoint: 960, settings: { slidesToShow: Math.min(2, articles.length), slidesToScroll: Math.min(2, articles.length) } },
            { breakpoint: 640, settings: { slidesToShow: 1, slidesToScroll: 1 } },
          ]}
        >
          {articles.map((a, i) => (
            <div key={a.id}>
              <article className="featured-slide" onClick={() => navigate(`/article/${a.id}`)}>
                <Cover article={a} index={i + 1} />
                <div className="featured-body">
                  <div className="featured-kicker">
                    <span className="no">FEATURED · {String(i + 1).padStart(2, '0')} / {String(articles.length).padStart(2, '0')}</span>
                    <span className="date">{formatDate(a.create_time)}</span>
                  </div>
                  <h3 className="featured-title">
                    {a.title}
                    {a.is_repost === 1 && <span className="repost-mark">转载</span>}
                  </h3>
                  <p className="featured-excerpt">{introOf(a)}</p>
                  <div className="featured-foot">
                    {catName(a.category_id) ? (
                      <span className="featured-cat">{catName(a.category_id)}</span>
                    ) : (
                      <span />
                    )}
                    <span className="featured-enter">READ →</span>
                  </div>
                </div>
              </article>
            </div>
          ))}
        </Carousel>

        {/* 纸墨风左右箭头 */}
        {articles.length > 1 && (
          <>
            <button
              className="featured-arrow prev"
              aria-label="上一篇精选"
              onClick={() => carouselRef.current?.prev()}
            >
              ←
            </button>
            <button
              className="featured-arrow next"
              aria-label="下一篇精选"
              onClick={() => carouselRef.current?.next()}
            >
              →
            </button>
          </>
        )}
      </div>
    </section>
  );
}
