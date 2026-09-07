import axios from 'axios';
import { message } from 'antd';

export interface ApiResponse<T = unknown> {
  code: number;
  msg: string;
  data: T;
}

export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
});

// 请求拦截器：添加 Authorization header
request.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器：统一错误处理 + 401 跳转登录
request.interceptors.response.use(
  (res) => {
    const body = res.data as ApiResponse;
    if (body.code !== 200) {
      message.error(body.msg || '请求失败');
      return Promise.reject(body);
    }
    return res.data;
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
      return Promise.reject(error);
    }
    // 403：账号无对应操作权限，统一提示「账号权限不足」
    if (error.response?.status === 403) {
      message.error('账号权限不足');
      return Promise.reject(error);
    }
    const msg = error.response?.data?.msg || error.message || '网络异常';
    message.error(msg);
    return Promise.reject(error);
  }
);

export default request;
