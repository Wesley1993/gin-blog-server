import axios, { AxiosResponse } from 'axios';
import { message } from 'antd';

/** 后端统一响应结构 */
export interface ApiResult<T> {
  code: number;
  msg: string;
  data: T;
}

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/',
  timeout: 15000,
});

// 响应拦截：拆包统一结构，业务码非 200 视为失败；配置带 silent 时不弹提示（供有兜底的请求使用）
request.interceptors.response.use(
  (res: AxiosResponse<ApiResult<unknown>>) => {
    const body = res.data;
    if (body && body.code === 200) {
      return res;
    }
    if (!(res.config as { silent?: boolean })?.silent) {
      message.error(body?.msg || '请求失败');
    }
    return Promise.reject(new Error(body?.msg || '请求失败'));
  },
  (err) => {
    if (!err?.config?.silent) {
      message.error('网络异常，请稍后重试');
    }
    return Promise.reject(err);
  },
);

/** GET 请求，直接返回 data 字段；silent=true 时失败不弹全局提示 */
export async function get<T>(url: string, params?: object, opts?: { silent?: boolean }): Promise<T> {
  const res = await request.get<ApiResult<T>>(url, { params, ...opts });
  return res.data.data;
}

export default request;
