-- 为已存在的范例文章补充占位封面图
-- 幂等：仅更新 cover 为空字符串的指定文章，可重复执行
-- 适用于 PostgreSQL 16

UPDATE blog_article
SET cover = 'https://picsum.photos/seed/blog1/800/450',
    update_time = NOW()
WHERE title = '欢迎使用本博客系统'
  AND (cover IS NULL OR cover = '');

UPDATE blog_article
SET cover = 'https://picsum.photos/seed/blog2/800/450',
    update_time = NOW()
WHERE title = 'Go 优雅关闭服务的正确姿势'
  AND (cover IS NULL OR cover = '');

UPDATE blog_article
SET cover = 'https://picsum.photos/seed/blog3/800/450',
    update_time = NOW()
WHERE title = '写代码之外的生活'
  AND (cover IS NULL OR cover = '');
