import { create } from 'zustand';
import { getSiteConfig } from '../api/blog';

interface SiteState {
  siteName: string;
  siteDesc: string;
  copyright: string;
  /** 建站日期，如 2024-01-01 */
  foundedAt: string;
  /** ICP 备案号 */
  icp: string;
  /** 公安备案号 */
  policeIcp: string;
  loaded: boolean;
  fetchConfig: () => Promise<void>;
}

/** 站点配置全局状态（页头/页脚/关于页共用） */
export const useSiteStore = create<SiteState>((set, get) => ({
  siteName: '无之空间',
  siteDesc: '',
  copyright: '',
  foundedAt: '',
  icp: '',
  policeIcp: '',
  loaded: false,
  fetchConfig: async () => {
    if (get().loaded) return;
    try {
      const cfg = await getSiteConfig();
      set({
        siteName: cfg.site_name || '无之空间',
        siteDesc: cfg.site_desc || '',
        copyright: cfg.copyright || '',
        foundedAt: cfg.founded_at || '',
        icp: cfg.icp || '',
        policeIcp: cfg.police_icp || '',
        loaded: true,
      });
      if (cfg.site_name) {
        document.title = `${cfg.site_name} · 手记`;
      }
    } catch {
      set({ loaded: true });
    }
  },
}));
