import request from './request';
import type { ApiResponse } from './request';

/** 系统运行概览（对应后端 service.SystemOverview） */
export interface SystemOverview {
  app_name: string;
  app_version: string;
  app_mode: string;
  go_version: string;
  goroutines: number;
  mem_alloc_mb: number;
  uptime: string;
  /** healthy / unhealthy / not_configured */
  db_status: string;
  /** healthy / unhealthy / not_configured */
  redis_status: string;
  /** healthy / unhealthy / not_configured */
  es_status: string;
}

/** 获取系统运行概览（应用信息、运行时指标、依赖健康状态） */
export function getSystemOverview() {
  return request.get<unknown, ApiResponse<SystemOverview>>('/system/overview');
}
