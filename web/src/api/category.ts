import request from './request';
import type { ApiResponse } from './request';
import type { CategoryItem } from './types';

export interface CategorySaveParams {
  id?: number;
  name: string;
  parent_id: number;
  sort: number;
  status: number;
  image?: string;
}

/** 获取分类树 */
export function getCategoryTree() {
  return request.get<unknown, ApiResponse<CategoryItem[]>>('/category/tree');
}

/** 新增分类 */
export function createCategory(data: CategorySaveParams) {
  return request.post<unknown, ApiResponse<null>>('/category/create', data);
}

/** 更新分类 */
export function updateCategory(data: CategorySaveParams) {
  return request.put<unknown, ApiResponse<null>>('/category/update', data);
}

/** 删除分类（被文章引用时后端返回错误） */
export function deleteCategory(id: number) {
  return request.delete<unknown, ApiResponse<null>>(`/category/${id}`);
}
