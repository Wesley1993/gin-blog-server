import { useEffect, useState } from 'react';
import { NavLink, useLocation, useNavigate, useSearchParams } from 'react-router-dom';
import { SearchOutlined } from '@ant-design/icons';
import { useSiteStore } from '../store/site';

/** 报头式顶部导航：导航链接与检索框同一行 */
export default function Header() {
  const siteName = useSiteStore((s) => s.siteName);
  const navigate = useNavigate();
  const location = useLocation();
  const [searchParams] = useSearchParams();

  const [kw, setKw] = useState('');

  // 搜索页上回显当前关键词
  useEffect(() => {
    if (location.pathname === '/search') {
      setKw(searchParams.get('keyword') || '');
    }
  }, [location.pathname, searchParams]);

  const doSearch = () => {
    const keyword = kw.trim();
    if (!keyword) return;
    navigate(`/search?keyword=${encodeURIComponent(keyword)}`);
  };

  return (
    <header className="masthead">
      <div className="container">
        <div className="masthead-top">
          <span className="est-mark">EST. 2026 · Coding之路</span>
          <span>WORDS &amp; CODE · A PERSONAL JOURNAL</span>
        </div>
        <div className="masthead-main">
          <h1 className="masthead-title" onClick={() => navigate('/')}>
            {siteName.slice(0, Math.max(0, siteName.length - 1))}
            <span className="accent">{siteName.slice(-1) || '·'}</span>
            {/* <span className="masthead-sub">手记</span> */}
          </h1>
          <div className="masthead-navrow">
            <nav className="masthead-nav">
              <NavLink to="/" end>
                首页
              </NavLink>
              <NavLink to="/about">关于</NavLink>
            </nav>
            <div className="masthead-search">
              <input
                value={kw}
                onChange={(e) => setKw(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && doSearch()}
                placeholder="检索文章…"
                aria-label="检索文章"
              />
              <button onClick={doSearch} aria-label="搜索">
                <SearchOutlined />
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
