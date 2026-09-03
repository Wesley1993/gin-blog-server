import { lazy, type LazyExoticComponent, type ComponentType } from 'react';

type PageComponent = LazyExoticComponent<ComponentType<any>>;

// 组件注册表：path → 懒加载组件
export const componentRegistry: Record<string, PageComponent> = {
  '/dashboard': lazy(() => import('../pages/Dashboard')),
  '/article': lazy(() => import('../pages/Article')),
  '/article/edit': lazy(() => import('../pages/Article/Edit')),
  '/category': lazy(() => import('../pages/Category')),
  '/user': lazy(() => import('../pages/User')),
  '/role': lazy(() => import('../pages/Role')),
  '/menu': lazy(() => import('../pages/Menu')),
  '/site': lazy(() => import('../pages/Site')),
  '/profile-config': lazy(() => import('../pages/SiteProfile')),
  '/links': lazy(() => import('../pages/SiteLink')),
  '/profile': lazy(() => import('../pages/Profile')),
};

// 已静态注册的路径（不走动态路由）
export const staticPaths = [
  '/dashboard',
  '/article',
  '/article/edit',
  '/category',
  '/user',
  '/role',
  '/menu',
  '/site',
  '/profile-config',
  '/links',
  '/profile',
];

// 获取动态路由路径（菜单中有但不在静态路由中的）
export function getDynamicPaths(menuPaths: string[]): string[] {
  return menuPaths.filter((p) => !staticPaths.includes(p));
}

// 获取注册表中的组件
export function getComponent(path: string): PageComponent | undefined {
  return componentRegistry[path];
}
