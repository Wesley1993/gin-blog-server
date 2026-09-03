import { useCallback } from 'react';
import { useAuth } from '../store/useAuth';

/**
 * 按钮权限判定 hook。
 * 通过 zustand 选择器订阅 permissions / isSuper，权限变化时自动触发组件重渲染
 * （不使用 getState，避免读取快照后不响应更新）。
 * 超管（permissions 含 '*'）对任意权限标识恒真。
 */
export function usePermission() {
  const permissions = useAuth((s) => s.permissions);
  const isSuper = useAuth((s) => s.isSuper);

  const hasPerm = useCallback(
    (perm: string) => isSuper || permissions.includes(perm),
    [permissions, isSuper],
  );

  return { hasPerm, isSuper };
}
