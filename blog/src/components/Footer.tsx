import { SafetyCertificateOutlined } from '@ant-design/icons';
import { useSiteStore } from '../store/site';

/** 页脚：版权与备案信息 */
export default function Footer() {
  const siteName = useSiteStore((s) => s.siteName);
  const copyright = useSiteStore((s) => s.copyright);
  const icp = useSiteStore((s) => s.icp);
  const policeIcp = useSiteStore((s) => s.policeIcp);

  return (
    <footer className="footer">
      <div className="container footer-inner">
        <div className="brand">
          {siteName.slice(0, Math.max(0, siteName.length - 1))}
          <span className="accent">{siteName.slice(-1) || '·'}</span>
          <span className="masthead-sub">技术笔记</span>
        </div>
        <div className="footer-meta">
          <div className="copy">
            {copyright || `© ${new Date().getFullYear()} ${siteName} · 保留所有权利`}
            {icp && (
              <>
                <span className="footer-sep">◆</span>
                <a
                  className="footer-icp"
                  href="https://beian.miit.gov.cn/"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {icp}
                </a>
              </>
            )}
          </div>
          {policeIcp && (
            <div className="copy">
              <a
                className="footer-icp footer-police"
                href="https://beian.mps.gov.cn/"
                target="_blank"
                rel="noopener noreferrer"
              >
                <SafetyCertificateOutlined /> {policeIcp}
              </a>
            </div>
          )}
          <div className="copy">SET IN SERIF · PRINTED ON PAPER &amp; PIXEL</div>
        </div>
      </div>
    </footer>
  );
}
