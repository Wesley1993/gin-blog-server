import { Suspense } from 'react';
import { Navigate, Route, Routes, useLocation } from 'react-router-dom';
import type { ReactNode } from 'react';
import { Spin } from 'antd';
import { useAuth } from '../store/useAuth';
import type { MenuItem } from '../api/types';
import { getComponent, staticPaths as STATIC_PATHS } from './componentRegistry';
import MainLayout from '../components/Layout';
import Login from '../pages/Login';
import Dashboard from '../pages/Dashboard';
import UserPage from '../pages/User';
import RolePage from '../pages/Role';
import MenuPage from '../pages/Menu';
import CategoryPage from '../pages/Category';
import ArticlePage from '../pages/Article';
import ArticleEdit from '../pages/Article/Edit';
import SitePage from '../pages/Site';
import SiteProfilePage from '../pages/SiteProfile';
import SiteLinkPage from '../pages/SiteLink';
import ProfilePage from '../pages/Profile';

/** 递归收集菜单树中所有页面类型（menu_type === 2）的路径 */
function collectPagePaths(menus: MenuItem[]): string[] {
  const paths: string[] = [];
  const walk = (items: MenuItem[]) => {
    for (const item of items) {
      if (item.menu_type === 2 && item.path) {
        paths.push(item.path);
      }
      if (item.children?.length) {
        walk(item.children);
      }
    }
  };
  walk(menus);
  return paths;
}

/** 路由权限守卫：基于菜单树拦截无权限路由（需在 RequireAuth 内部使用） */
function RequirePermission({ path, children }: { path: string; children: ReactNode }) {
  const menus = useAuth((s) => s.menus);
  const isSuper = useAuth((s) => s.isSuper);
  // 超管直接放行；菜单未加载完成（空数组）时放行，避免闪烁
  if (isSuper || menus.length === 0) {
    return <>{children}</>;
  }
  const allowed = collectPagePaths(menus).includes(path);
  if (!allowed) {
    return <Navigate to="/article" replace />;
  }
  return <>{children}</>;
}

/** 路由守卫：未登录跳转登录页 */
function RequireAuth({ children }: { children: ReactNode }) {
  const token = useAuth((s) => s.token);
  const location = useLocation();
  if (!token) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  return <>{children}</>;
}

/** 动态路由页面：按路径从组件注册表中懒加载对应组件 */
function DynamicPage({ path }: { path: string }) {
  const Component = getComponent(path);
  if (!Component) {
    return <div className="page-container text-center py-20 text-[#8a7f6f]">页面建设中...</div>;
  }
  return <Component />;
}

export default function AppRouter() {
  const menus = useAuth((s) => s.menus);
  const dynamicPaths = collectPagePaths(menus).filter((p) => !STATIC_PATHS.includes(p));

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route
        path="/"
        element={
          <RequireAuth>
            <MainLayout />
          </RequireAuth>
        }
      >
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route
          path="user"
          element={
            <RequirePermission path="/user">
              <UserPage />
            </RequirePermission>
          }
        />
        <Route
          path="role"
          element={
            <RequirePermission path="/role">
              <RolePage />
            </RequirePermission>
          }
        />
        <Route
          path="menu"
          element={
            <RequirePermission path="/menu">
              <MenuPage />
            </RequirePermission>
          }
        />
        <Route
          path="category"
          element={
            <RequirePermission path="/category">
              <CategoryPage />
            </RequirePermission>
          }
        />
        <Route path="article" element={<ArticlePage />} />
        <Route path="article/edit" element={<ArticleEdit />} />
        <Route path="article/edit/:id" element={<ArticleEdit />} />
        <Route
          path="site"
          element={
            <RequirePermission path="/site">
              <SitePage />
            </RequirePermission>
          }
        />
        <Route
          path="profile-config"
          element={
            <RequirePermission path="/profile-config">
              <SiteProfilePage />
            </RequirePermission>
          }
        />
        <Route
          path="links"
          element={
            <RequirePermission path="/links">
              <SiteLinkPage />
            </RequirePermission>
          }
        />
        <Route path="profile" element={<ProfilePage />} />
        {/* 动态路由：直接内联渲染为 <Route>，不使用自定义组件包装（Routes 仅接受 Route/Fragment 子节点） */}
        {dynamicPaths.map((path) => (
          <Route
            key={path}
            path={path}
            element={
              <RequirePermission path={path}>
                <Suspense fallback={<div className="flex items-center justify-center h-64"><Spin size="large" /></div>}>
                  <DynamicPage path={path} />
                </Suspense>
              </RequirePermission>
            }
          />
        ))}
        <Route path="*" element={<Navigate to="/article" replace />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
