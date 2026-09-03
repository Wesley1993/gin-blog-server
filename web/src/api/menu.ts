import request from './request';
import type { ApiResponse } from './request';
import type { MenuItem } from './types';

/** 获取全部菜单树 */
export function getMenuList() {
  return request.get<unknown, ApiResponse<MenuItem[]>>('/menu/list');
}

/** 创建菜单 */
export function createMenu(data: Partial<MenuItem>) {
  return request.post<unknown, ApiResponse<null>>('/menu/create', data);
}

/** 更新菜单 */
export function updateMenu(data: Partial<MenuItem>) {
  return request.put<unknown, ApiResponse<null>>('/menu/update', data);
}

/** 删除菜单 */
export function deleteMenu(id: number) {
  return request.delete<unknown, ApiResponse<null>>(`/menu/${id}`);
}
