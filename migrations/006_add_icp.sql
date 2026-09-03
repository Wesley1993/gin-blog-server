-- 站点配置表新增 ICP 备案号字段（幂等）
ALTER TABLE blog_site_config ADD COLUMN IF NOT EXISTS icp VARCHAR(100) DEFAULT '';
