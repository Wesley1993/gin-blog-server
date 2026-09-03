import { useEffect, useState } from 'react';
import { Spin } from 'antd';
import { getUserInfo } from './api/auth';
import { getRoleList } from './api/role';
import { useAuth } from './store/useAuth';
import AppRouter from './router';

/** 应用根组件：登录后拉取用户信息 / 菜单树 / 权限标识，并判定超管身份 */
export default function App() {
  const token = useAuth((s) => s.token);
  const userInfo = useAuth((s) => s.userInfo);
  const setUser = useAuth((s) => s.setUser);
  const setSuper = useAuth((s) => s.setSuper);
  const logout = useAuth((s) => s.logout);

  // 用户信息加载态：已登录刷新页面时，菜单树未返回前展示全屏 Loading，避免空菜单闪烁
  const [userLoading, setUserLoading] = useState(false);

  useEffect(() => {
    if (!token || userInfo) return;
    let cancelled = false;
    setUserLoading(true);
    getUserInfo()
      .then((res) => {
        if (cancelled) return;
        const { user, menus, permissions } = res.data;
        setUser(user, menus ?? [], permissions ?? []);
        // 通过角色列表判定当前用户是否为超级管理员
        return getRoleList().then((roleRes) => {
          if (cancelled) return;
          const role = roleRes.data?.find((r) => r.id === user.role_id);
          setSuper(Boolean(role?.is_super));
        });
      })
      .catch(() => {
        if (!cancelled) logout();
      })
      .finally(() => {
        if (!cancelled) setUserLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [token, userInfo, setUser, setSuper, logout]);

  if (token && !userInfo && userLoading) {
    return (
      <div className="h-screen flex items-center justify-center bg-[#f4efe6]">
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  return <AppRouter />;
}
