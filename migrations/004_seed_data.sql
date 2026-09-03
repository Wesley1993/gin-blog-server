-- Blog 后台管理系统 范例数据种子迁移
-- 幂等：所有插入均带存在性检查，可重复执行
-- 适用于 PostgreSQL 16

-- 范例分类
INSERT INTO blog_category (name, parent_id, sort, status, create_time, update_time)
SELECT '技术', 0, 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_category WHERE name = '技术');

INSERT INTO blog_category (name, parent_id, sort, status, create_time, update_time)
SELECT '生活', 0, 2, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_category WHERE name = '生活');

-- 范例文章 1：欢迎文章（Markdown 正文，状态 1 已发布）
INSERT INTO blog_article (title, category_id, cover, content, tags, status, is_deleted, create_time, update_time)
SELECT '欢迎使用本博客系统',
       (SELECT id FROM blog_category WHERE name = '技术' LIMIT 1),
       'https://picsum.photos/seed/blog1/800/450',
       E'# 欢迎使用本博客系统\n\n这是一个基于 **Gin + PostgreSQL + Elasticsearch** 构建的博客系统。\n\n## 功能特性\n\n- 文章管理：支持 Markdown 编写\n- 全文检索：基于 ES 的中文分词搜索（IK 分词器）\n- RBAC 权限：用户、角色、菜单权限管理\n- 文件存储：阿里云 OSS 对象存储\n\n## 代码示例\n\n```go\nfunc main() {\n    fmt.Println("Hello, Blog!")\n}\n```\n\n> 提示：所有文章均以 Markdown 格式存储和渲染。\n\n享受写作吧！',
       '入门,教程',
       1, 0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_article WHERE title = '欢迎使用本博客系统');

-- 范例文章 2：Go 技术文章
INSERT INTO blog_article (title, category_id, cover, content, tags, status, is_deleted, create_time, update_time)
SELECT 'Go 优雅关闭服务的正确姿势',
       (SELECT id FROM blog_category WHERE name = '技术' LIMIT 1),
       'https://picsum.photos/seed/blog2/800/450',
       E'# Go 优雅关闭服务的正确姿势\n\n在生产环境中，直接 `kill -9` 进程会中断正在处理的请求。Go 标准库提供了完善的优雅关闭机制。\n\n## 核心流程\n\n1. 监听 `SIGINT` / `SIGTERM` 信号\n2. 收到信号后调用 `http.Server.Shutdown(ctx)` 停止接收新连接\n3. 等待存量请求处理完成（或超时强制退出）\n4. 关闭数据库、缓存等下游连接\n\n## 代码示例\n\n```go\nquit := make(chan os.Signal, 1)\nsignal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)\n<-quit\n\nctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\ndefer cancel()\nif err := srv.Shutdown(ctx); err != nil {\n    log.Printf("服务器强制关闭: %v", err)\n}\n```\n\n## 注意事项\n\n- **超时时间** 应根据最长请求耗时合理设置，避免连接长时间挂起。\n- `ListenAndServe` 返回 `http.ErrServerClosed` 属于正常关闭，不应视为错误。\n- 别忘了 `sql.DB.Close()` 与 Redis 客户端的 `Close()`。\n\n优雅关闭是服务稳定性的最后一道防线。',
       'Go,后端,运维',
       1, 0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_article WHERE title = 'Go 优雅关闭服务的正确姿势');

-- 范例文章 3：生活随笔
INSERT INTO blog_article (title, category_id, cover, content, tags, status, is_deleted, create_time, update_time)
SELECT '写代码之外的生活',
       (SELECT id FROM blog_category WHERE name = '生活' LIMIT 1),
       'https://picsum.photos/seed/blog3/800/450',
       E'# 写代码之外的生活\n\n程序员的生活不该只有屏幕和键盘。记录一些让状态变好的小事。\n\n## 我的充电方式\n\n- **散步**：不带目的地走上半小时，很多 bug 的思路就在路上冒出来。\n- **阅读**：每月至少读一本与技术无关的书。\n- **运动**：久坐是隐形杀手，每周三次运动雷打不动。\n\n## 一点感悟\n\n> 代码写得快不快，取决于思路清不清；思路清不清，往往是因为休息得不够。\n\n保持节奏，比偶尔冲刺更重要。愿大家都能写出好代码，也过得好生活。',
       '随笔,生活',
       1, 0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM blog_article WHERE title = '写代码之外的生活');
