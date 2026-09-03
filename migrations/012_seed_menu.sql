-- Blog 后台管理系统 菜单种子数据迁移
-- 适用于 PostgreSQL 16
-- 所有插入均使用 WHERE NOT EXISTS 保证幂等，可重复执行

-- ============================================================
-- 一级菜单（目录 / 顶层页面，parent_id = 0）
-- ============================================================

-- 内容管理（目录）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT 0, '内容管理', 1, '/content', '', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0
);

-- 权限管理（目录）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT 0, '权限管理', 1, '/system', '', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0
);

-- 站点设置（顶层页面）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT 0, '站点设置', 2, '/site', '', 3, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu WHERE menu_name = '站点设置' AND parent_id = 0
);

-- 个人资料（顶层页面）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT 0, '个人资料', 2, '/profile-config', '', 4, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu WHERE menu_name = '个人资料' AND parent_id = 0
);

-- 常用网站（顶层页面）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT 0, '常用网站', 2, '/links', '', 5, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu WHERE menu_name = '常用网站' AND parent_id = 0
);

-- ============================================================
-- 二级菜单（页面）：父级 ID 通过子查询获取
-- ============================================================

-- 文章管理（页面，归属内容管理）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0),
       '文章管理', 2, '/article', '', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '文章管理'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)
);

-- 分类管理（页面，归属内容管理）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0),
       '分类管理', 2, '/category', '', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '分类管理'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)
);

-- 用户管理（页面，归属权限管理）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0),
       '用户管理', 2, '/user', '', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '用户管理'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)
);

-- 角色管理（页面，归属权限管理）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0),
       '角色管理', 2, '/role', '', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '角色管理'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)
);

-- 菜单管理（页面，归属权限管理）
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0),
       '菜单管理', 2, '/menu', '', 3, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '菜单管理'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)
);

-- ============================================================
-- 三级菜单（按钮权限）：父级为对应页面，ID 通过子查询获取
-- ============================================================

-- 文章管理下的按钮：新增 / 编辑 / 删除
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)),
       '新增文章', 3, '', 'article:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '新增文章'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)),
       '编辑文章', 3, '', 'article:edit', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '编辑文章'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)),
       '删除文章', 3, '', 'article:delete', 3, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '删除文章'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '文章管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0))
);

-- 分类管理下的按钮：新增 / 编辑 / 删除
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '分类管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)),
       '新增分类', 3, '', 'category:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '新增分类'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '分类管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '分类管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)),
       '编辑分类', 3, '', 'category:edit', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '编辑分类'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '分类管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '分类管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0)),
       '删除分类', 3, '', 'category:delete', 3, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '删除分类'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '分类管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '内容管理' AND parent_id = 0))
);

-- 用户管理下的按钮：重置密码
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '重置密码', 3, '', 'user:resetPwd', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '重置密码'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

-- 角色管理下的按钮：新增 / 编辑 / 删除
INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '角色管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '新增角色', 3, '', 'role:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '新增角色'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '角色管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '角色管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '编辑角色', 3, '', 'role:edit', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '编辑角色'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '角色管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '角色管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '删除角色', 3, '', 'role:delete', 3, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '删除角色'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '角色管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

-- ============================================================
-- 演示角色：编辑（非超级管理员）
-- 仅授权"内容管理"分支（含目录本身，否则菜单树无法挂载子页面），
-- menu_ids 通过子查询聚合实际菜单 ID 为 JSONB 数组
-- ============================================================
INSERT INTO sys_role (role_name, menu_ids, button_perms, is_super, create_time, update_time)
SELECT '编辑',
       (SELECT jsonb_agg(id ORDER BY id) FROM sys_menu
        WHERE menu_name IN ('内容管理', '文章管理', '分类管理')),
       '["article:add", "article:edit", "article:delete", "category:add", "category:edit", "category:delete"]'::jsonb,
       0, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_role WHERE role_name = '编辑'
);
