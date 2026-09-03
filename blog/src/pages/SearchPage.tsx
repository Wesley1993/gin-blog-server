import { useCallback, useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Empty, Spin } from 'antd';
import { getCategories, searchArticles } from '../api/blog';
import type { Article, Category } from '../api/blog';
import ArticleCard from '../components/ArticleCard';

/** 搜索页一次性拉取全量结果（展示端不分页） */
const FETCH_SIZE = 1000;

/** 独立搜索页：从 /search?keyword=xx 读取关键词并展示检索结果 */
export default function SearchPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const keyword = (searchParams.get('keyword') || '').trim();

  const [cats, setCats] = useState<Category[]>([]);
  const [list, setList] = useState<Article[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  // 分类名映射（卡片展示用）
  useEffect(() => {
    getCategories()
      .then((tree) => setCats(flattenCats(tree || [])))
      .catch(() => setCats([]));
  }, []);

  // 关键词变化即重新检索
  useEffect(() => {
    if (!keyword) {
      setList([]);
      setTotal(0);
      return;
    }
    let cancelled = false;
    setLoading(true);
    searchArticles({ keyword, page: 1, page_size: FETCH_SIZE })
      .then((res) => {
        if (cancelled) return;
        setList(res.list || []);
        setTotal(res.total || 0);
      })
      .catch(() => {
        if (cancelled) return;
        setList([]);
        setTotal(0);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [keyword]);

  const reSearch = useCallback(
    (v: string) => {
      const kw = v.trim();
      if (!kw) return;
      navigate(`/search?keyword=${encodeURIComponent(kw)}`);
    },
    [navigate],
  );

  return (
    <div className="container search-page">
      <div className="section-head">
        <span className="zh">检索</span>
        <span className="rule" />
        <span className="meta">SEARCH &amp; INDEX</span>
      </div>

      {/* 页内可再次输入的检索框 */}
      <div className="search-input-row">
        <SearchBox defaultValue={keyword} onSearch={reSearch} />
      </div>

      {!keyword ? (
        <div className="state-block">
          <Empty description="请输入关键词，检索文章标题、正文与标签">
            <button className="search-back-home" onClick={() => navigate('/')}>
              ← 返回首页
            </button>
          </Empty>
        </div>
      ) : (
        <>
          <div className="date-chip">
            检索「{keyword}」的结果 · {loading ? '…' : total} 条
          </div>
          <Spin spinning={loading}>
            {list.length === 0 && !loading ? (
              <div className="state-block">
                <Empty description={`未见与「${keyword}」相关的文字，换个词再试`}>
                  <button className="search-back-home" onClick={() => navigate('/')}>
                    ← 返回首页
                  </button>
                </Empty>
              </div>
            ) : (
              list.map((a, i) => (
                <ArticleCard key={a.id} article={a} index={i + 1} categories={cats} />
              ))
            )}
          </Spin>
        </>
      )}
    </div>
  );
}

/** 展平分类树（与首页一致，供卡片取分类名） */
function flattenCats(tree: Category[], depth = 0): Category[] {
  const out: Category[] = [];
  for (const c of tree) {
    out.push({ ...c });
    if (c.children?.length) out.push(...flattenCats(c.children, depth + 1));
  }
  return out;
}

/** 纸墨风格检索框：受控输入，回车或点「检索」提交 */
function SearchBox({
  defaultValue,
  onSearch,
}: {
  defaultValue: string;
  onSearch: (v: string) => void;
}) {
  const [v, setV] = useState(defaultValue);

  // URL 关键词变化时回显
  useEffect(() => {
    setV(defaultValue);
  }, [defaultValue]);

  return (
    <div className="search-box">
      <input
        value={v}
        onChange={(e) => setV(e.target.value)}
        onKeyDown={(e) => e.key === 'Enter' && onSearch(v)}
        placeholder="全文检索文章标题、正文、标签…"
        aria-label="全文检索"
      />
      <button onClick={() => onSearch(v)}>检索</button>
    </div>
  );
}
