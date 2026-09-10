import { useEffect, useState } from 'react';
import type { Heading } from '../utils/markdown';

interface TocSidebarProps {
  headings: Heading[];
}

/** 文章目录侧栏：显示标题层级，点击滚动定位，滚动时高亮当前标题 */
export default function TocSidebar({ headings }: TocSidebarProps) {
  const [activeId, setActiveId] = useState('');

  useEffect(() => {
    if (headings.length === 0) return;
    const visibleIds = new Set<string>();
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((e) => {
          if (e.isIntersecting) visibleIds.add(e.target.id);
          else visibleIds.delete(e.target.id);
        });
        // 从 headings 顺序中找第一个可见的标题
        const firstVisible = headings.find((h) => visibleIds.has(h.id));
        if (firstVisible) setActiveId(firstVisible.id);
      },
      { rootMargin: '-80px 0px -70% 0px' }, // 只观察视口顶部 30% 区域
    );
    headings.forEach((h) => {
      const el = document.getElementById(h.id);
      if (el) observer.observe(el);
    });
    return () => observer.disconnect();
  }, [headings]);

  const handleClick = (id: string) => {
    const el = document.getElementById(id);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' });
      setActiveId(id);
    }
  };

  if (headings.length === 0) return null;

  return (
    <aside className="toc-sidebar">
      <div className="toc-title">目录</div>
      <nav>
        {headings.map((h) => (
          <a
            key={h.id}
            href={`#${h.id}`}
            className={`toc-item level-${h.level} ${activeId === h.id ? 'active' : ''}`}
            onClick={(e) => {
              e.preventDefault();
              handleClick(h.id);
            }}
          >
            {h.text}
          </a>
        ))}
      </nav>
    </aside>
  );
}
