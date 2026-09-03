-- 站点配置扩展：建站日期
-- 适用于 PostgreSQL 16

ALTER TABLE blog_site_config ADD COLUMN IF NOT EXISTS founded_at VARCHAR(20) DEFAULT '';
