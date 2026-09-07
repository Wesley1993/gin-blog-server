-- 文章发布时间字段：支持自定义发布时间，公开查询/排序以此为准
-- 幂等：ADD COLUMN IF NOT EXISTS / CREATE INDEX IF NOT EXISTS，可重复执行
-- 适用于 PostgreSQL 16

ALTER TABLE blog_article ADD COLUMN IF NOT EXISTS published_at TIMESTAMP WITHOUT TIME ZONE;

-- 历史数据回填：已有文章用创建时间作为发布时间
UPDATE blog_article SET published_at = create_time WHERE published_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_article_published_at ON blog_article(published_at);

COMMENT ON COLUMN blog_article.published_at IS '发布时间（草稿可为空，公开列表按此排序与筛选）';
