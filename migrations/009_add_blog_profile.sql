-- 站长个人资料表（单行表，展示端「关于页/联系站长」数据源）
-- 幂等：CREATE TABLE IF NOT EXISTS + WHERE NOT EXISTS 种子，可重复执行
-- 适用于 PostgreSQL 16

CREATE TABLE IF NOT EXISTS blog_profile (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) DEFAULT '',
    avatar VARCHAR(255) DEFAULT '',
    title VARCHAR(100) DEFAULT '',
    bio TEXT DEFAULT '',
    github VARCHAR(100) DEFAULT '',
    email VARCHAR(100) DEFAULT '',
    wechat VARCHAR(100) DEFAULT '',
    qq VARCHAR(100) DEFAULT '',
    skills JSONB DEFAULT '[]',
    projects JSONB DEFAULT '[]',
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- 种子数据：仅当表为空时插入
INSERT INTO blog_profile (name, avatar, title, bio, github, email, wechat, qq, skills, projects, create_time, update_time)
SELECT '吴之',
       '',
       '后端开发工程师 / 文字爱好者',
       '相信代码与文字是同一件事的两面：都是把混沌的想法整理成清晰的秩序。日常写 Go 与 TypeScript，偶尔写随笔。这个站点既是我沉淀技术笔记的地方，也是记录生活片段的角落。愿读到这里的你，也有所收获。',
       'https://github.com/wuzhi',
       'hello@wuzhispace.com',
       'wuzhispace',
       '85263741',
       '[{"name":"Go","level":90,"group":"后端"},{"name":"Gin / GORM","level":85,"group":"后端"},{"name":"PostgreSQL","level":80,"group":"后端"},{"name":"Redis","level":78,"group":"后端"},{"name":"Elasticsearch","level":70,"group":"后端"},{"name":"TypeScript / React","level":75,"group":"前端"},{"name":"Docker / CI","level":65,"group":"工程化"},{"name":"Linux","level":72,"group":"工程化"}]'::jsonb,
       '[{"name":"gin-blog-server","desc":"基于 Gin + GORM + PostgreSQL 的博客系统，包含管理后台与展示前台，集成 JWT 鉴权、RBAC 权限、Redis 缓存、Elasticsearch 全文搜索与阿里云 OSS 上传。","tech":["Go","Gin","GORM","PostgreSQL","Redis","Elasticsearch"],"start":"2026-05","end":"","link":""},{"name":"分布式短链服务","desc":"高并发短链接生成与跳转服务，基于发号器 + 布隆过滤器去重，支持统计埋点与灰度发布。","tech":["Go","Redis","Kafka","ClickHouse"],"start":"2025-08","end":"2026-03","link":""},{"name":"内容审核平台","desc":"对接多家机审服务商的内容安全中台，支持文本/图片审核流水线、人工复核工作台与申诉流程。","tech":["Java","Spring Boot","MySQL","RabbitMQ"],"start":"2024-11","end":"2025-07","link":""}]'::jsonb,
       NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_profile);
