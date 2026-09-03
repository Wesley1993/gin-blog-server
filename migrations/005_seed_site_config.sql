-- 站点配置种子数据
-- 幂等：仅当表为空时插入，可重复执行
-- 适用于 PostgreSQL 16

INSERT INTO blog_site_config (site_name, site_desc, copyright, founded_at, create_time, update_time)
SELECT 'WuzhiSpace', '记录技术与生活的个人博客', '© 2026 WuzhiSpace. All rights reserved.', '2026-08-29', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_site_config);
