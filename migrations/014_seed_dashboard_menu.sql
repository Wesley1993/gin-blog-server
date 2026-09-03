-- 仪表盘菜单种子数据迁移
-- 适用于 PostgreSQL 16
-- 使用 WHERE NOT EXISTS 保证幂等，可重复执行

-- 仪表盘（顶层页面，sort = 0 使其排在菜单首位）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT 0, '仪表盘', 2, '/dashboard', '', 0, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu WHERE menu_name = '仪表盘' AND parent_id = 0
);
