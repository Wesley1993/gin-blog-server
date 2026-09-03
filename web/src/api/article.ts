import request from './request';
import type { ApiResponse, PageData } from './request';
import type { ArticleItem } from './types';

export interface ArticlePageParams {
  page: number;
  page_size: number;
  category_id?: number;
  status?: number;
}

export interface ArticleSearchParams {
  keyword: string;
  page: number;
  page_size: number;
}

export interface ArticleSaveParams {
  id?: number;
  title: string;
  category_id: number;
  cover?: string;
  content: string;
  tags?: string;
  status: number;
  /** 0 原创 1 转载 */
  is_repost: number;
  /** 原文链接（转载时必填） */
  repost_url?: string;
  /** 原作者 */
  repost_author?: string;
}

/** PG 分页列表 */
export function getArticlePage(params: ArticlePageParams) {
  return request.get<unknown, ApiResponse<PageData<ArticleItem>>>('/article/page', { params });
}

/** ES 全文搜索 */
export function searchArticle(params: ArticleSearchParams) {
  return request.get<unknown, ApiResponse<PageData<ArticleItem>>>('/article/es/search', { params });
}

/** 根据 ID 获取文章详情（复用分页接口无法取正文，直接查分页后从列表找） */
export function getArticleById(id: number) {
  return request.get<unknown, ApiResponse<PageData<ArticleItem>>>('/article/page', {
    params: { page: 1, page_size: 100 },
  }).then((res) => {
    const found = res.data.list.find((item) => item.id === id);
    return found ? { data: found } : Promise.reject({ code: 30001, msg: '文章不存在' });
  });
}

/** 创建文章 */
export function createArticle(data: ArticleSaveParams) {
  return request.post<unknown, ApiResponse<null>>('/article/create', data);
}

/** 更新文章 */
export function updateArticle(data: ArticleSaveParams) {
  return request.put<unknown, ApiResponse<null>>('/article/update', data);
}

/** 逻辑删除文章 */
export function deleteArticle(id: number) {
  return request.delete<unknown, ApiResponse<null>>(`/article/${id}`);
}

/** 手动重建 ES 全量索引（超管） */
export function rebuildIndex() {
  return request.post<unknown, ApiResponse<null>>('/article/es/rebuild');
}
