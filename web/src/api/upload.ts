import request from './request';
import type { ApiResponse } from './request';

/** OSS 文件上传（multipart/form-data，字段 file；限 10MB 图片） */
export function uploadToOss(file: File) {
  const formData = new FormData();
  formData.append('file', file);
  return request.post<unknown, ApiResponse<{ url: string }>>('/upload/oss', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    timeout: 30000,
  });
}
