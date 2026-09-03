import { useEffect, useMemo, useState } from 'react';
import { Empty, Spin, Tooltip } from 'antd';
import { GithubOutlined, LinkOutlined, MailOutlined, QqOutlined, WechatOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { getArticles, getCategories, getLinks } from '../api/blog';
import type { Article, Category, SiteLink } from '../api/blog';
import ArticleCard from '../components/ArticleCard';
import FeaturedArticles from '../components/FeaturedArticles';
import { useSiteStore } from '../store/site';
import { useProfileStore } from '../store/profile';

/** 一次性全量拉取已发布文章（展示端不分页） */
const FETCH_SIZE = 1000;

/** 展平后的分类（保留原名，另带层级用于缩进展示） */
interface FlatCategory extends Category {
  depth: number;
}

/** 展平分类树（一级 + 二级），不改动 name，保证 id→name 映射干净 */
function flattenCats(tree: Category[], depth = 0): FlatCategory[] {
  const out: FlatCategory[] = [];
  for (const c of tree) {
    out.push({ ...c, depth });
    if (c.children?.length) out.push(...flattenCats(c.children, depth + 1));
  }
  return out;
}

/** 解析 tags 字段：兼容逗号分隔字符串与数组两种结构 */
function parseTags(tags: string | string[] | null | undefined): string[] {
  if (Array.isArray(tags)) return tags.map((t) => String(t).trim()).filter(Boolean);
  return tags ? tags.split(',').map((t) => t.trim()).filter(Boolean) : [];
}

/** 标签名确定性哈希（种子），保证同一标签颜色稳定不闪烁 */
function hashStr(s: string): number {
  let h = 7;
  for (let i = 0; i < s.length; i++) {
    h = (h * 31 + s.charCodeAt(i)) >>> 0;
  }
  return h;
}

/** 以标签名为种子的 HSL 随机色：色相随机、饱和度 55%-75%、亮度 38%-52%，暖纸底上清晰可读 */
function tagColor(name: string): string {
  const h = hashStr(name);
  const hue = h % 360;
  const sat = 55 + ((h >> 8) % 21);
  const lig = 38 + ((h >> 16) % 15);
  return `hsl(${hue}, ${sat}%, ${lig}%)`;
}

export default function Home() {
  const foundedAt = useSiteStore((s) => s.foundedAt);

  // 联系站长：接口优先，失败由 store 内部回退静态资料
  const contact = useProfileStore((s) => s.profile);
  const fetchProfile = useProfileStore((s) => s.fetchProfile);

  useEffect(() => {
    fetchProfile();
  }, [fetchProfile]);

  // 与后端 /api/site/stats 对齐：建站日算第 0 天，即 (now - founded_at) 的整天数
  const runDays = foundedAt && dayjs(foundedAt).isValid() ? dayjs().diff(dayjs(foundedAt), 'day') : 0;

  const [cats, setCats] = useState<FlatCategory[]>([]);
  const [catId, setCatId] = useState<number>(0);
  const [tagFilter, setTagFilter] = useState<string>('');
  /** 全量文章（不带筛选），供标签云汇总与前端筛选共用 */
  const [allArticles, setAllArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(false);
  const [featured, setFeatured] = useState<Article[]>([]);
  /** 侧栏「常用网站」：接口失败或空时整卡隐藏 */
  const [links, setLinks] = useState<SiteLink[]>([]);

  // 分类树
  useEffect(() => {
    getCategories()
      .then((tree) => setCats(flattenCats(tree || [])))
      .catch(() => setCats([]));
  }, []);

  // 最近精选：取最新 10 篇供轮播展示
  useEffect(() => {
    getArticles({ page: 1, page_size: 10 })
      .then((res) => setFeatured(res.list || []))
      .catch(() => setFeatured([]));
  }, []);

  // 常用网站：静默拉取，失败/空则隐藏卡片（不弹全局提示）
  useEffect(() => {
    getLinks()
      .then((ls) => setLinks(Array.isArray(ls) ? ls : []))
      .catch(() => setLinks([]));
  }, []);

  // 全量拉取一次已发布文章；分类 / 标签筛选均在前端叠加应用
  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    getArticles({ page: 1, page_size: FETCH_SIZE })
      .then((res) => {
        if (!cancelled) setAllArticles(res.list || []);
      })
      .catch(() => {
        if (!cancelled) setAllArticles([]);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  // 标签云：汇总全量文章 tags，去重计数，按文章数倒序
  const tagStats = useMemo(() => {
    const counts = new Map<string, number>();
    for (const a of allArticles) {
      for (const t of new Set(parseTags(a.tags))) {
        counts.set(t, (counts.get(t) || 0) + 1);
      }
    }
    return Array.from(counts.entries())
      .map(([name, count]) => ({ name, count }))
      .sort((x, y) => y.count - x.count || x.name.localeCompare(y.name));
  }, [allArticles]);

  // 列表：分类 + 标签双重筛选叠加（均命中才展示）
  const list = useMemo(
    () =>
      allArticles.filter((a) => {
        if (catId && a.category_id !== catId) return false;
        if (tagFilter && !parseTags(a.tags).includes(tagFilter)) return false;
        return true;
      }),
    [allArticles, catId, tagFilter],
  );

  const total = list.length;

  const onPickCat = (id: number) => {
    setCatId(id);
  };

  // 点击标签：选中筛选，再次点击取消
  const onPickTag = (name: string) => {
    setTagFilter((prev) => (prev === name ? '' : name));
  };

  // 分类 / 标签筛选激活时隐藏精选区，聚焦筛选结果
  const hasFilter = catId !== 0 || Boolean(tagFilter);

  return (
    <>
      <div className="container home-container">
        {/* 最近精选（无筛选时展示） */}
        {!hasFilter && !loading && featured.length > 0 && (
          <FeaturedArticles articles={featured} categories={cats} />
        )}

        {/* 分类筛选 */}
        <div className="section-head">
          <span className="zh">目录</span>
          <span className="rule" />
          <span className="meta">CONTENTS · 共 {total} 篇</span>
        </div>
        {cats.length > 0 && (
          <div className="cat-row">
            <button className={`cat-btn ${catId === 0 ? 'on' : ''}`} onClick={() => onPickCat(0)}>
              全部
            </button>
            {cats.map((c) => (
              <button
                key={c.id}
                className={`cat-btn ${catId === c.id ? 'on' : ''}`}
                onClick={() => onPickCat(c.id)}
              >
                {'— '.repeat(c.depth)}
                {c.name}
              </button>
            ))}
          </div>
        )}
        {catId !== 0 && (
          <div className="date-chip">
            筛选：分类「{cats.find((c) => c.id === catId)?.name || ''}」 · {total} 篇
            <button onClick={() => onPickCat(0)} aria-label="取消分类筛选">
              ✕
            </button>
          </div>
        )}
        {tagFilter && (
          <div className="date-chip">
            筛选：标签「{tagFilter}」 · {total} 篇
            <button onClick={() => setTagFilter('')} aria-label="取消标签筛选">
              ✕
            </button>
          </div>
        )}

        {/* 主体：列表 + 侧栏 */}
        <div className="home-body">
          <div>
            <Spin spinning={loading}>
              {list.length === 0 && !loading ? (
                <div className="state-block">
                  <Empty description={hasFilter ? '无符合当前筛选条件的文章' : '暂无文章，静候落笔'} />
                </div>
              ) : (
                list.map((a, i) => (
                  <ArticleCard key={a.id} article={a} index={i + 1} categories={cats} />
                ))
              )}
            </Spin>
          </div>

          <aside>
            <div className="side-block">
              <div className="side-title">
                <span className="zh">标签</span>
                <span className="en">Tags</span>
              </div>
              {tagStats.length > 0 ? (
                <>
                  <div className="tag-cloud">
                    {tagStats.map((t) => {
                      const color = tagColor(t.name);
                      const on = tagFilter === t.name;
                      return (
                        <button
                          key={t.name}
                          className={`tag-cloud-item ${on ? 'on' : ''}`}
                          style={
                            on
                              ? { background: color, borderColor: color, color: 'var(--paper-card)' }
                              : { color, borderColor: color }
                          }
                          onClick={() => onPickTag(t.name)}
                          aria-pressed={on}
                        >
                          # {t.name}
                          <span className="tag-count">{t.count}</span>
                        </button>
                      );
                    })}
                  </div>
                  <div className="side-note">◆ 点选标签可按其筛选文章，再点一次取消</div>
                </>
              ) : (
                <div className="side-note">◆ 尚无标签，静候落笔</div>
              )}
            </div>

            <div className="side-block">
              <div className="side-title">
                <span className="zh">站点信息</span>
                <span className="en">Colophon</span>
              </div>
              {runDays > 0 ? (
                <>
                  <div className="side-stats-num">
                    {runDays}
                    <span>天</span>
                  </div>
                  <div className="side-note">◆ 本站已运行 {runDays} 天</div>
                  <div className="side-note side-note-sub">◆ 建站于 {foundedAt}</div>
                </>
              ) : (
                <div className="side-note">◆ 甫一落笔，来日方长</div>
              )}
            </div>

            <div className="side-block">
              <div className="side-title">
                <span className="zh">联系站长</span>
                <span className="en">Contact</span>
              </div>
              <div className="side-contact">
                {contact.github && (
                  <a href={contact.github} target="_blank" rel="noopener noreferrer">
                    <GithubOutlined /> GitHub
                  </a>
                )}
                {contact.email && (
                  <a href={`mailto:${contact.email}`}>
                    <MailOutlined /> {contact.email}
                  </a>
                )}
                {contact.wechat && (
                  <span>
                    <WechatOutlined /> 微信 · {contact.wechat}
                  </span>
                )}
                {contact.qq && (
                  <span>
                    <QqOutlined /> QQ · {contact.qq}
                  </span>
                )}
              </div>
            </div>

            {links.length > 0 && (
              <div className="side-block">
                <div className="side-title">
                  <span className="zh">常用网站</span>
                  <span className="en">Links</span>
                </div>
                <div className="side-links">
                  {links.map((l) => (
                    <a
                      key={l.id}
                      href={l.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="side-links-item"
                    >
                      <LinkOutlined />
                      <span className="side-links-name">
                        {l.name}
                        {l.description && (
                          <Tooltip title={l.description} placement="right">
                            <span className="side-links-desc">{l.description}</span>
                          </Tooltip>
                        )}
                      </span>
                    </a>
                  ))}
                </div>
                <div className="side-note">◆ 点击可在新窗口打开</div>
              </div>
            )}

            <div className="side-block">
              <div className="side-title">
                <span className="zh">编者的话</span>
                <span className="en">Note</span>
              </div>
              <p style={{ fontSize: 13.5, color: 'var(--ink-soft)', textAlign: 'justify' }}>
                这里记录技术学习中的所得，也收录生活里值得记下的片段。每篇文章都欢迎评论，也期待评论之外的交流——若有所感，欢迎通过关于页与我联系。
              </p>
            </div>
          </aside>
        </div>
      </div>
    </>
  );
}
