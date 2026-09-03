-- 多云 OSS 支持：为站点配置表添加跳过 HTTPS 证书校验开关（自签名证书/私有端点场景）
ALTER TABLE blog_site_config ADD COLUMN IF NOT EXISTS oss_insecure SMALLINT NOT NULL DEFAULT 0;
