-- 补齐菜单管理与用户管理的按钮权限种子数据
-- 适用于 PostgreSQL 16
-- 所有插入均使用 WHERE NOT EXISTS 保证幂等，可重复执行
-- 注意：「编辑」演示角色授权保持不动（不授予权限管理模块按钮，符合最小权限演示意图）

-- ============================================================
-- 菜单管理下的按钮：新增 / 编辑 / 删除
-- ============================================================

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '菜单管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '新增菜单', 3, '', 'menu:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '新增菜单'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '菜单管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '菜单管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '编辑菜单', 3, '', 'menu:edit', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '编辑菜单'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '菜单管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '菜单管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '删除菜单', 3, '', 'menu:delete', 3, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '删除菜单'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '菜单管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

-- ============================================================
-- 用户管理下的按钮：新增 / 编辑 / 删除 / 启停用
-- （重置密码 user:resetPwd 已在 012 迁移中种入，此处从 sort=2 起排）
-- ============================================================

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '新增用户', 3, '', 'user:add', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '新增用户'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '编辑用户', 3, '', 'user:edit', 3, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '编辑用户'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '删除用户', 3, '', 'user:delete', 4, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '删除用户'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);

INSERT INTO sys_menu (parent_id, menu_name, menu_type, path, perms, sort, status, create_time, update_time)
SELECT (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
          AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0)),
       '启停用用户', 3, '', 'user:status', 5, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_menu
    WHERE menu_name = '启停用用户'
      AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '用户管理'
                         AND parent_id = (SELECT id FROM sys_menu WHERE menu_name = '权限管理' AND parent_id = 0))
);
