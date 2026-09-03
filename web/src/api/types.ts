// 后端实体类型定义（对应 REQUIREMENTS §5 数据模型）

export interface MenuItem {
  id: number;
  parent_id: number;
  menu_name: string;
  menu_type: number; // 1 目录 2 页面 3 按钮
  path?: string;
  perms?: string;
  sort: number;
  status: number;
  children?: MenuItem[];
}

export interface Role {
  id: number;
  role_name: string;
  menu_ids: number[] | null;
  button_perms: string[] | null;
  is_super: number;
  create_time?: string;
}

export interface UserItem {
  id: number;
  username: string;
  nickname?: string;
  role_id: number;
  status: number;
  create_time?: string;
}

export interface CategoryItem {
  id: number;
  name: string;
  parent_id: number;
  sort: number;
  status: number;
  /** 分类图片（文章未上传封面时作为默认封面） */
  image?: string;
  children?: CategoryItem[];
  create_time?: string;
}

export interface ArticleItem {
  id: number;
  title: string;
  category_id: number;
  cover?: string;
  content?: string;
  tags?: string;
  status: number; // 0 草稿 1 已发布
  /** 0 原创 1 转载 */
  is_repost?: number;
  /** 原文链接（转载时必填） */
  repost_url?: string;
  /** 原作者 */
  repost_author?: string;
  create_time?: string;
  update_time?: string;
}

export interface SiteConfig {
  site_name?: string;
  site_desc?: string;
  copyright?: string;
  founded_at?: string;
  /** ICP 备案号 */
  icp?: string;
  /** 公安备案号 */
  police_icp?: string;
  oss_access_key?: string;
  oss_secret_key?: string;
  oss_bucket?: string;
  oss_endpoint?: string;
  oss_domain?: string;
  /** 存储厂商：aliyun | tencent | qiniu | s3 | rustfs */
  oss_provider?: string;
  /** 存储区域（腾讯云/Amazon S3 使用，RustFS 可选） */
  oss_region?: string;
  /** 跳过 HTTPS 证书校验（自签名证书/私有端点）：0 否 1 是 */
  oss_insecure?: number;
}

export interface CurrentUser {
  id: number;
  username: string;
  nickname?: string;
  avatar?: string;
  role_id: number;
}

export interface Profile {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  bio: string;
  role_id: number;
}

export interface SiteStats {
  founded_at: string;
  running_days: number;
  article_count: number;
  category_count: number;
  user_count: number;
}

/** 站长联系方式 */
export interface SiteProfileContact {
  github: string;
  email: string;
  wechat: string;
  qq: string;
  address: string;
}

/** 技能项 */
export interface SiteProfileSkill {
  name: string;
  /** 熟练度 0-100 */
  level: number;
  group: string;
}

/** 项目经历项 */
export interface SiteProfileProject {
  name: string;
  desc: string;
  tech: string[];
  /** 开始时间 如 2026-01 */
  start: string;
  /** 结束时间，空表示至今 */
  end: string;
  link: string;
}

/** 站长个人资料（展示端「关于页/联系站长」数据源） */
export interface SiteProfile {
  name: string;
  avatar: string;
  title: string;
  bio: string;
  contacts: SiteProfileContact;
  skills: SiteProfileSkill[];
  projects: SiteProfileProject[];
}

/** 常用网站（友情链接） */
export interface SiteLink {
  id: number;
  name: string;
  url: string;
  /** 图标图片 URL */
  icon?: string;
  description?: string;
  sort: number;
  /** 1 启用 0 停用 */
  status: number;
  create_time?: string;
  update_time?: string;
}

export interface UserInfoResp {
  user: CurrentUser;
  menus: MenuItem[];
  permissions: string[];
}
