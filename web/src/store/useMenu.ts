import { create } from 'zustand';

interface MenuState {
  collapsed: boolean;
  selectedKeys: string[];
  openKeys: string[];
  toggleCollapsed: () => void;
  setSelectedKeys: (keys: string[]) => void;
  setOpenKeys: (keys: string[]) => void;
}

/** 侧边菜单交互状态（折叠、选中、展开） */
export const useMenu = create<MenuState>((set) => ({
  collapsed: false,
  selectedKeys: [],
  openKeys: [],
  toggleCollapsed: () => set((state) => ({ collapsed: !state.collapsed })),
  setSelectedKeys: (keys) => set({ selectedKeys: keys }),
  setOpenKeys: (keys) => set({ openKeys: keys }),
}));
