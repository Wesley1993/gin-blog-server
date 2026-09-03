import request from './request';
import type { ApiResponse } from './request';
import type { Role } from './types';

export interface RoleSaveParams {
  id?: number;
  role_name: string;
  menu_ids: number[];
  button_perms: string[];
}

/** 角色列表 */
export function getRoleList() {
  return request.get<unknown, ApiResponse<Role[]>>('/role/list');
}

/** 创建角色 */
export function createRole(data: RoleSaveParams) {
  return request.post<unknown, ApiResponse<null>>('/role/create', data);
}

/** 更新角色 + 权限 */
export function updateRole(data: RoleSaveParams) {
  return request.put<unknown, ApiResponse<null>>('/role/update', data);
}

/** 删除角色 */
export function deleteRole(id: number) {
  return request.delete<unknown, ApiResponse<null>>(`/role/${id}`);
}
