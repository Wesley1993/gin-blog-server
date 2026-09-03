import request from './request';
import type { ApiResponse } from './request';
import type { SiteConfig, SiteStats } from './types';

/** 获取站点配置 */
export function getSiteConfig() {
  return request.get<unknown, ApiResponse<SiteConfig>>('/site/config');
}

/** 保存站点配置（刷新 Redis 缓存） */
export function saveSiteConfig(data: SiteConfig) {
  return request.put<unknown, ApiResponse<null>>('/site/config', data);
}

/** OSS 连通测试 */
export function testOss() {
  return request.post<unknown, ApiResponse<null>>('/site/testOss');
}

/** 获取站点运行统计（建站日期、运行天数、文章/分类/用户数） */
export function getSiteStats() {
  return request.get<unknown, ApiResponse<SiteStats>>('/site/stats');
}
