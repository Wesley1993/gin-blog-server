-- 多云 OSS 支持：为站点配置表添加存储厂商和区域字段
ALTER TABLE blog_site_config ADD COLUMN IF NOT EXISTS oss_provider varchar(20) NOT NULL DEFAULT 'aliyun';
ALTER TABLE blog_site_config ADD COLUMN IF NOT EXISTS oss_region varchar(100) NOT NULL DEFAULT '';
