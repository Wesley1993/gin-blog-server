import { useEffect, useState } from 'react';
import { VerticalAlignTopOutlined } from '@ant-design/icons';

/** 长文滚过该距离后浮现 */
const SHOW_AFTER_PX = 320;

/** 悬浮「回到顶部」按钮：内容超出一屏且向下滚动后浮现，点击平滑回到页面顶部 */
export default function BackToTop() {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    let ticking = false;
    const update = () => {
      ticking = false;
      // 可滚动距离：内容不足一屏时为 0，按钮始终不出现
      const maxScroll = Math.max(
        0,
        document.documentElement.scrollHeight - window.innerHeight,
      );
      // 只略超一屏的短文章按可滚动距离折半，避免按钮无处可显
      setVisible(window.scrollY > Math.min(SHOW_AFTER_PX, maxScroll / 2));
    };
    const schedule = () => {
      if (ticking) return;
      ticking = true;
      requestAnimationFrame(update);
    };
    update();
    window.addEventListener('scroll', schedule, { passive: true });
    window.addEventListener('resize', schedule, { passive: true });
    return () => {
      window.removeEventListener('scroll', schedule);
      window.removeEventListener('resize', schedule);
    };
  }, []);

  return (
    <button
      type="button"
      className={`back-to-top ${visible ? 'visible' : ''}`}
      onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
      aria-label="回到顶部"
      title="回到顶部"
    >
      <VerticalAlignTopOutlined />
    </button>
  );
}
