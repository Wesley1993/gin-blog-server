-- 常用网站增加 icon 图标字段（存储图标图片 URL）
-- 幂等：ADD COLUMN IF NOT EXISTS，可重复执行
-- 适用于 PostgreSQL 16

ALTER TABLE blog_link ADD COLUMN IF NOT EXISTS icon varchar(500) NOT NULL DEFAULT '';
