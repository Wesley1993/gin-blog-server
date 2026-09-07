-- 站点信息 / 个人信息 / 常用网站 / 上传 / 索引重建 写操作按钮权限种子数据
-- 适用于 PostgreSQL 16
-- 所有插入均使用 WHERE NOT EXISTS 保证幂等，可重复执行
--
-- 背景：router.go 的 writePermTable 新登记了以下写路由权限标识，
--       非超管用户须由角色 button_perms 显式授予才能操作，否则一律拒绝
--       （只读「测试」角色 button_perms 为空，因此无法修改站点信息/个人信息）。
--       本迁移将这些权限标识作为「按钮(menu_type=3)」节点种入 sys_menu，
--       使其在角色编辑页的按钮权限选择器中可见、可勾选授予。
--
-- 权限标识与父级菜单归属：
--   site:edit         -> 站点设置（保存基础配置 / 保存 OSS / 连通测试同属站点配置写操作）
--   siteProfile:edit  -> 个人资料（保存站长资料）
--   profile:edit      -> 个人资料（编辑当前登录用户个人信息）
--   link:add/edit/del -> 常用网站
--   article:rebuild   -> 文章管理（重建 ES 索引）
--   upload:create     -> 文章管理（文件上传统一入口）

-- ============================================================
-- 站点设置 下的按钮：保存站点配置
-- ============================================================
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '站点设置' AND parent_id = 0),
       '保存站点配置', 3, '', 'site:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '保存站点配置'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '站点设置' AND parent_id = 0)
);

-- ============================================================
-- 个人资料 下的按钮：保存站长资料 / 编辑个人信息
-- ============================================================
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '个人资料' AND parent_id = 0),
       '保存站长资料', 3, '', 'siteProfile:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '保存站长资料'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '个人资料' AND parent_id = 0)
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '个人资料' AND parent_id = 0),
       '编辑个人信息', 3, '', 'profile:edit', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '编辑个人信息'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '个人资料' AND parent_id = 0)
);

-- ============================================================
-- 常用网站 下的按钮：新增 / 编辑 / 删除
-- ============================================================
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '常用网站' AND parent_id = 0),
       '新增网站', 3, '', 'link:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '新增网站'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '常用网站' AND parent_id = 0)
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '常用网站' AND parent_id = 0),
       '编辑网站', 3, '', 'link:edit', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '编辑网站'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '常用网站' AND parent_id = 0)
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '常用网站' AND parent_id = 0),
       '删除网站', 3, '', 'link:delete', 3, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '删除网站'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '常用网站' AND parent_id = 0)
);

-- ============================================================
-- 文章管理 下的按钮：重建索引 / 文件上传
-- （article:add/edit/delete 已在 012 迁移中种入，此处从 sort=4 起排）
-- ============================================================
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)),
       '重建索引', 3, '', 'article:rebuild', 4, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '重建索引'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)),
       '文件上传', 3, '', 'upload:create', 5, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '文件上传'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0))
);

-- ============================================================
-- 「编辑」演示角色补授权：文件上传 + 重建索引 + 编辑个人信息
-- 编辑角色可创建/编辑文章与分类，需保留封面/图片上传能力（upload:create）、
-- 文章页「重建索引」按钮能力（article:rebuild），以及在个人中心修改自己
-- 昵称/头像/简介的能力（profile:edit），避免本次收紧写权限后既有功能回退。
-- 通过 JSONB 数组拼接 + DISTINCT 去重实现幂等，重复执行不会产生重复标识。
-- 站点信息/站长资料/常用网站相关权限刻意不授予（编辑角色无权管理站点）。
-- ============================================================
UPDATE sys_role
SET button_perms = (
        SELECT COALESCE(jsonb_agg(DISTINCT elem), '[]'::jsonb)
        FROM jsonb_array_elements_text(
            COALESCE(button_perms, '[]'::jsonb) || '["upload:create", "article:rebuild", "profile:edit"]'::jsonb
        ) AS elem
    ),
    update_time = NOW()
WHERE role_name = '编辑';
