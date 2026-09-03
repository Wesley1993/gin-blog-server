import { create } from 'zustand';
import { getProfile } from '../api/blog';
import type { PublicProfile } from '../api/blog';
import { profile as fallbackProfile } from '../data/profile';
import type { Profile } from '../data/profile';

interface ProfileState {
  /** 展示用个人资料（接口数据或静态兜底） */
  profile: Profile;
  /** 是否已完成一次加载（成功或兜底均算完成） */
  loaded: boolean;
  fetchProfile: () => Promise<void>;
}

/** 起止月份转为展示文案，如 2026-01 → 2026.01；end 为空表示至今 */
function formatPeriod(start: string, end: string): string {
  const fmt = (s: string) => s.replace(/-/g, '.');
  const from = start ? fmt(start) : '';
  const to = end ? fmt(end) : '至今';
  return from ? `${from} - ${to}` : to;
}

/** 将后端契约数据映射为前端 Profile 结构；关键字段缺失则返回 null 走兜底 */
function mapRemote(p: PublicProfile | null | undefined): Profile | null {
  if (!p || !p.name?.trim()) return null;
  const contacts = p.contacts || ({} as PublicProfile['contacts']);
  return {
    name: p.name,
    title: p.title || fallbackProfile.title,
    avatar: p.avatar || '',
    bio: p.bio || fallbackProfile.bio,
    location: contacts.address || fallbackProfile.location,
    email: contacts.email || '',
    github: contacts.github || '',
    wechat: contacts.wechat || '',
    qq: contacts.qq || '',
    skills: (p.skills || [])
      .filter((s) => s && s.name)
      .map((s) => ({ name: s.name, level: s.level || 0, group: s.group || '其他' })),
    projects: (p.projects || [])
      .filter((pr) => pr && pr.name)
      .map((pr) => ({
        name: pr.name,
        desc: pr.desc || '',
        tech: pr.tech || [],
        period: formatPeriod(pr.start, pr.end),
        link: pr.link || '',
      })),
  };
}

/** 个人资料全局状态：关于页与首页侧边栏共用，接口失败回退静态数据 */
export const useProfileStore = create<ProfileState>((set, get) => ({
  profile: fallbackProfile,
  loaded: false,
  fetchProfile: async () => {
    if (get().loaded) return;
    try {
      const remote = await getProfile();
      const mapped = mapRemote(remote);
      set({ profile: mapped || fallbackProfile, loaded: true });
    } catch {
      // 接口不可用/失败：静态兜底，页面不空白
      set({ profile: fallbackProfile, loaded: true });
    }
  },
}));
