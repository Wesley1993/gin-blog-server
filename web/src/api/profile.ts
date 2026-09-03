import request from './request';
import type { ApiResponse } from './request';
import type { Profile } from './types';

export interface UpdateProfileData {
  nickname: string;
  avatar: string;
  bio: string;
}

/** 获取当前用户个人信息 */
export function getProfile() {
  return request.get<unknown, ApiResponse<Profile>>('/profile/info');
}

/** 更新个人信息（头像、昵称、简介） */
export function updateProfile(data: UpdateProfileData) {
  return request.put<unknown, ApiResponse<Profile>>('/profile/update', data);
}
