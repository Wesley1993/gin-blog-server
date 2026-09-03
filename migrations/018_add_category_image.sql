-- 分类支持设置图片：作为文章未上传封面时的默认封面来源
ALTER TABLE blog_category ADD COLUMN IF NOT EXISTS image VARCHAR(500) NOT NULL DEFAULT '';
