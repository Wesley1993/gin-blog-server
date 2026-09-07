import { get } from './request';

/** 文章 */
export interface Article {
  id: number;
  title: string;
  category_id: number;
  cover: string;
  content: string;
  /** 摘要（后端列表可能不返回，前端需降级处理） */
  summary?: string;
  tags: string;
  status: number;
  /** 是否转载：0=原创 1=转载；旧文章可能缺失该字段，缺失时按原创处理 */
  is_repost?: number;
  /** 转载原文链接 */
  repost_url?: string;
  /** 转载原作者 */
  repost_author?: string;
  /** 发布时间（YYYY-MM-DD HH:mm:ss，旧数据可能为空） */
  published_at?: string;
  create_time: string;
  update_time: string;
}

/** 分类 */
export interface Category {
  id: number;
  name: string;
  parent_id: number;
  sort: number;
  status: number;
  children?: Category[];
}

/** 分页数据 */
export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

/** 站点公开配置 */
export interface SiteConfig {
  site_name: string;
  site_desc: string;
  copyright: string;
  founded_at: string;
  icp: string;
  police_icp: string;
}

export interface ArticleQuery {
  page?: number;
  page_size?: number;
  category_id?: number;
}

export interface SearchQuery {
  keyword: string;
  page?: number;
  page_size?: number;
}

/** 个人简介接口：联系方式（后端契约） */
export interface ProfileContacts {
  github: string;
  email: string;
  wechat: string;
  qq: string;
  address: string;
}

/** 个人简介接口：技能项 */
export interface ProfileSkill {
  name: string;
  level: number;
  group: string;
}

/** 个人简介接口：项目经历（起止均为 2026-01 形式，end 为空表示至今） */
export interface ProfileProject {
  name: string;
  desc: string;
  tech: string[];
  start: string;
  end: string;
  link: string;
}

/** 个人公开简介（GET /api/public/profile） */
export interface PublicProfile {
  name: string;
  avatar: string;
  title: string;
  bio: string;
  contacts: ProfileContacts;
  skills: ProfileSkill[];
  projects: ProfileProject[];
}

/** 文章分页（仅已发布） */
export const getArticles = (params: ArticleQuery) =>
  get<PageData<Article>>('/api/public/article/page', params);

/** 文章详情 */
export const getArticleDetail = (id: number | string) =>
  get<Article>(`/api/public/article/${id}`);

/** ES 全文搜索（仅已发布） */
export const searchArticles = (params: SearchQuery) =>
  get<PageData<Article>>('/api/public/article/es/search', params);

/** 分类树 */
export const getCategories = () => get<Category[]>('/api/public/category/tree');

/** 站点公开配置 */
export const getSiteConfig = () => get<SiteConfig>('/api/public/site/config');

/** 个人公开简介（无认证；失败静默，前端有静态兜底） */
export const getProfile = () =>
  get<PublicProfile>('/api/public/profile', undefined, { silent: true });

/** 常用网站（GET /api/public/links） */
export interface SiteLink {
  id: number;
  name: string;
  url: string;
  description: string;
  sort: number;
}

/** 常用网站列表（无认证；失败静默，无数据时侧栏卡片隐藏） */
export const getLinks = () =>
  get<SiteLink[]>('/api/public/links', undefined, { silent: true });
