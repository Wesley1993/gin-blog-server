import request from './request';
import type { ApiResponse } from './request';
import type { UserInfoResp } from './types';

export interface LoginParams {
  username: string;
  password: string;
}

/** 账号密码登录 */
export function login(data: LoginParams) {
  return request.post<unknown, ApiResponse<{ token: string }>>('/auth/login', data);
}

/** 登出 */
export function logout() {
  return request.post<unknown, ApiResponse<null>>('/auth/logout');
}

/** 获取当前用户信息 + 菜单树 + 权限标识集合 */
export function getUserInfo() {
  return request.get<unknown, ApiResponse<UserInfoResp>>('/auth/userinfo');
}

export interface ChangePwdParams {
  old_password: string;
  new_password: string;
}

/** 修改当前登录用户密码（成功后后端会删除 token 强制下线） */
export function changePassword(data: ChangePwdParams) {
  return request.put<unknown, ApiResponse<null>>('/user/password', data);
}
