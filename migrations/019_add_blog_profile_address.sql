-- 站点个人信息新增地址字段
ALTER TABLE blog_profile ADD COLUMN IF NOT EXISTS address VARCHAR(200) DEFAULT '';
