-- Blog 后台管理系统 初始化建表 SQL
-- 适用于 PostgreSQL 16

-- 用户表
CREATE TABLE IF NOT EXISTS sys_user (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    nickname VARCHAR(50) DEFAULT '',
    avatar VARCHAR(255) DEFAULT '',
    bio VARCHAR(500) DEFAULT '',
    role_id BIGINT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- 角色表
CREATE TABLE IF NOT EXISTS sys_role (
    id BIGSERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL,
    menu_ids JSONB DEFAULT '[]',
    button_perms JSONB DEFAULT '[]',
    is_super SMALLINT NOT NULL DEFAULT 0,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- 菜单权限表
CREATE TABLE IF NOT EXISTS sys_menu (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NOT NULL DEFAULT 0,
    menu_name VARCHAR(50) NOT NULL,
    menu_type SMALLINT NOT NULL DEFAULT 1,
    path VARCHAR(100) DEFAULT '',
    perms VARCHAR(100) DEFAULT '',
    sort INT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- 分类表
CREATE TABLE IF NOT EXISTS blog_category (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(60) NOT NULL,
    parent_id BIGINT NOT NULL DEFAULT 0,
    sort INT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- 文章表
CREATE TABLE IF NOT EXISTS blog_article (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(120) NOT NULL,
    category_id BIGINT NOT NULL DEFAULT 0,
    cover VARCHAR(255) DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    tags VARCHAR(200) DEFAULT '',
    status SMALLINT NOT NULL DEFAULT 0,
    is_deleted SMALLINT NOT NULL DEFAULT 0,
    is_repost SMALLINT NOT NULL DEFAULT 0,
    repost_url VARCHAR(500) DEFAULT '',
    repost_author VARCHAR(100) DEFAULT '',
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_article_category ON blog_article(category_id);
CREATE INDEX IF NOT EXISTS idx_article_status ON blog_article(status);
CREATE INDEX IF NOT EXISTS idx_article_create_time ON blog_article(create_time);

-- 站点配置表
CREATE TABLE IF NOT EXISTS blog_site_config (
    id BIGSERIAL PRIMARY KEY,
    site_name VARCHAR(100) DEFAULT '',
    site_desc VARCHAR(255) DEFAULT '',
    copyright VARCHAR(255) DEFAULT '',
    founded_at VARCHAR(20) DEFAULT '',
    icp VARCHAR(100) DEFAULT '',
    police_icp VARCHAR(100) DEFAULT '',
    oss_access_key VARCHAR(200) DEFAULT '',
    oss_secret_key VARCHAR(200) DEFAULT '',
    oss_bucket VARCHAR(100) DEFAULT '',
    oss_endpoint VARCHAR(200) DEFAULT '',
    oss_domain VARCHAR(200) DEFAULT '',
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- 站长个人资料表（单行表）
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

-- 常用网站（友情链接）表
CREATE TABLE IF NOT EXISTS blog_link (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(255) NOT NULL,
    description VARCHAR(200) DEFAULT '',
    sort INT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_link_status ON blog_link(status);
CREATE INDEX IF NOT EXISTS idx_link_sort ON blog_link(sort);

-- 初始数据：超级管理员角色（sys_role 无唯一约束，用 WHERE NOT EXISTS 保证幂等）
INSERT INTO sys_role (role_name, menu_ids, button_perms, is_super, create_time, update_time)
SELECT '超级管理员', '[]', '[]', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM sys_role WHERE role_name = '超级管理员');

-- 初始数据：超级管理员用户（密码 admin123，bcrypt 加密）
INSERT INTO sys_user (username, password, nickname, role_id, status, create_time, update_time)
VALUES ('admin', '$2a$10$wkjv4BSUNfOivbVmMxfqce8AAzStKv3T.WP30dPAHDy5DZZ6iGysa', '超级管理员', 1, 1, NOW(), NOW())
ON CONFLICT (username) DO NOTHING;
