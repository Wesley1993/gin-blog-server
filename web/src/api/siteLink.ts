import request from './request';
import type { ApiResponse } from './request';
import type { SiteLink } from './types';

/** 新增/更新常用网站的请求体 */
export interface SiteLinkSaveParams {
  name: string;
  url: string;
  /** 图标图片 URL */
  icon: string;
  description: string;
  sort: number;
  /** 1 启用 0 停用 */
  status: number;
}

/** 获取全部常用网站（含停用，按 sort 升序） */
export function getSiteLinks() {
  return request.get<unknown, ApiResponse<SiteLink[]>>('/site/links');
}

/** 新增常用网站 */
export function createSiteLink(data: SiteLinkSaveParams) {
  return request.post<unknown, ApiResponse<null>>('/site/links', data);
}

/** 更新常用网站 */
export function updateSiteLink(id: number, data: SiteLinkSaveParams) {
  return request.put<unknown, ApiResponse<null>>(`/site/links/${id}`, data);
}

/** 删除常用网站 */
export function deleteSiteLink(id: number) {
  return request.delete<unknown, ApiResponse<null>>(`/site/links/${id}`);
}
