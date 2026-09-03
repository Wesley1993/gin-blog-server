-- 站点配置表新增公安备案号字段（幂等）
ALTER TABLE blog_site_config ADD COLUMN IF NOT EXISTS police_icp VARCHAR(100) DEFAULT '';
