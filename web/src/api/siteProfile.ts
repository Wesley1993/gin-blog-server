import request from './request';
import type { ApiResponse } from './request';
import type { SiteProfile } from './types';

/** 获取站长个人资料（管理端） */
export function getSiteProfile() {
  return request.get<unknown, ApiResponse<SiteProfile>>('/site/profile');
}

/** 保存站长个人资料（整体覆盖保存） */
export function saveSiteProfile(data: SiteProfile) {
  return request.put<unknown, ApiResponse<null>>('/site/profile', data);
}
