import { create } from 'zustand';
import type { CurrentUser, MenuItem } from '../api/types';

interface AuthState {
  token: string | null;
  userInfo: CurrentUser | null;
  menus: MenuItem[];
  permissions: string[];
  isSuper: boolean;
  setToken: (token: string) => void;
  setUser: (user: CurrentUser, menus: MenuItem[], perms: string[]) => void;
  patchUser: (patch: Partial<CurrentUser>) => void;
  setSuper: (isSuper: boolean) => void;
  logout: () => void;
}

/** 根据菜单树推断是否超级管理员（拥有全部菜单视为超管标识的兜底判断） */
export const useAuth = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  userInfo: null,
  menus: [],
  permissions: [],
  isSuper: false,
  setToken: (token) => {
    localStorage.setItem('token', token);
    set({ token });
  },
  setUser: (user, menus, perms) =>
    set({
      userInfo: user,
      menus,
      permissions: perms,
      isSuper: perms.includes('*'),
    }),
  patchUser: (patch) =>
    set((state) => ({
      userInfo: state.userInfo ? { ...state.userInfo, ...patch } : state.userInfo,
    })),
  setSuper: (isSuper) => set({ isSuper }),
  logout: () => {
    localStorage.removeItem('token');
    set({ token: null, userInfo: null, menus: [], permissions: [], isSuper: false });
  },
}));

/** 判断当前用户是否拥有某个按钮权限标识 */
export function hasPermission(perm: string): boolean {
  const { isSuper, permissions } = useAuth.getState();
  return isSuper || permissions.includes(perm);
}
