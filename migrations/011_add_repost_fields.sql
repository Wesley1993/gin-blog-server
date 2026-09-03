-- 文章转载字段：转载标识、原文链接、原作者
-- 幂等：ADD COLUMN IF NOT EXISTS，可重复执行
-- 适用于 PostgreSQL 16

ALTER TABLE blog_article ADD COLUMN IF NOT EXISTS is_repost SMALLINT NOT NULL DEFAULT 0;
ALTER TABLE blog_article ADD COLUMN IF NOT EXISTS repost_url VARCHAR(500) DEFAULT '';
ALTER TABLE blog_article ADD COLUMN IF NOT EXISTS repost_author VARCHAR(100) DEFAULT '';

COMMENT ON COLUMN blog_article.is_repost IS '转载标识：0原创 1转载';
COMMENT ON COLUMN blog_article.repost_url IS '原文链接（转载时必填）';
COMMENT ON COLUMN blog_article.repost_author IS '原作者';
