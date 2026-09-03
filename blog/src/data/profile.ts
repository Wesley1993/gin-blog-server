/**
 * 个人信息数据结构与静态占位数据
 * 接口 /api/public/profile 不可用或返回空时，作为兜底展示数据源。
 */

export interface Project {
  name: string;
  desc: string;
  tech: string[];
  period: string;
  /** 项目链接（可选） */
  link?: string;
}

export interface Skill {
  name: string;
  /** 熟练度 0-100 */
  level: number;
  group: string;
}

export interface Profile {
  name: string;
  title: string;
  avatar: string;
  bio: string;
  location: string;
  email: string;
  github: string;
  wechat: string;
  qq: string;
  skills: Skill[];
  projects: Project[];
}

export const profile: Profile = {
  name: '吴之',
  title: '后端开发工程师 / 文字爱好者',
  avatar: '',
  bio: '相信代码与文字是同一件事的两面：都是把混沌的想法整理成清晰的秩序。日常写 Go 与 TypeScript，偶尔写随笔。这个站点既是我沉淀技术笔记的地方，也是记录生活片段的角落。愿读到这里的你，也有所收获。',
  location: '杭州',
  email: 'hello@wuzhispace.com',
  github: 'https://github.com/wuzhi',
  wechat: 'wuzhispace',
  qq: '85263741',
  skills: [
    { name: 'Go', level: 90, group: '后端' },
    { name: 'Gin / GORM', level: 85, group: '后端' },
    { name: 'PostgreSQL', level: 80, group: '后端' },
    { name: 'Redis', level: 78, group: '后端' },
    { name: 'Elasticsearch', level: 70, group: '后端' },
    { name: 'TypeScript / React', level: 75, group: '前端' },
    { name: 'Docker / CI', level: 65, group: '工程化' },
    { name: 'Linux', level: 72, group: '工程化' },
  ],
  projects: [
    {
      name: 'gin-blog-server',
      desc: '基于 Gin + GORM + PostgreSQL 的博客系统，包含管理后台与展示前台，集成 JWT 鉴权、RBAC 权限、Redis 缓存、Elasticsearch 全文搜索与阿里云 OSS 上传。',
      tech: ['Go', 'Gin', 'GORM', 'PostgreSQL', 'Redis', 'Elasticsearch'],
      period: '2026.05 - 至今',
    },
    {
      name: '分布式短链服务',
      desc: '高并发短链接生成与跳转服务，基于发号器 + 布隆过滤器去重，支持统计埋点与灰度发布。',
      tech: ['Go', 'Redis', 'Kafka', 'ClickHouse'],
      period: '2025.08 - 2026.03',
    },
    {
      name: '内容审核平台',
      desc: '对接多家机审服务商的内容安全中台，支持文本/图片审核流水线、人工复核工作台与申诉流程。',
      tech: ['Java', 'Spring Boot', 'MySQL', 'RabbitMQ'],
      period: '2024.11 - 2025.07',
    },
  ],
};
