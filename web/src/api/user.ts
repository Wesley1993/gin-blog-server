import request from './request';
import type { ApiResponse, PageData } from './request';
import type { UserItem } from './types';

export interface UserPageParams {
  page: number;
  page_size: number;
  username?: string;
}

export interface UserCreateParams {
  username: string;
  nickname?: string;
  password: string;
  role_id: number;
  status: number;
}

export interface UserUpdateParams {
  id: number;
  nickname?: string;
  role_id: number;
  status: number;
}

/** 用户分页列表 */
export function getUserPage(params: UserPageParams) {
  return request.get<unknown, ApiResponse<PageData<UserItem>>>('/user/page', { params });
}

/** 创建用户 */
export function createUser(data: UserCreateParams) {
  return request.post<unknown, ApiResponse<null>>('/user/create', data);
}

/** 编辑用户 */
export function updateUser(data: UserUpdateParams) {
  return request.put<unknown, ApiResponse<null>>('/user/update', data);
}

/** 重置密码 */
export function resetPassword(id: number, password: string) {
  return request.put<unknown, ApiResponse<null>>(`/user/resetPwd/${id}`, { password });
}

/** 启用/禁用用户（对应 user:status 权限） */
export function updateUserStatus(id: number, status: number) {
  return request.put<unknown, ApiResponse<null>>(`/user/status/${id}`, { status });
}

/** 删除用户 */
export function deleteUser(id: number) {
  return request.delete<unknown, ApiResponse<null>>(`/user/${id}`);
}
