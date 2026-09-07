-- 只读「测试」角色 + 测试账号 种子数据迁移
-- 适用于 PostgreSQL 16
-- 所有插入均使用 WHERE NOT EXISTS / ON CONFLICT 保证幂等，可重复执行
--
-- 说明：本项目 sys_role 实际字段为 (id, role_name, menu_ids, button_perms,
--       is_super, create_time, update_time)，并无 code/description/status 列，
--       故角色以 role_name='测试' 唯一标识，只读语义通过 button_perms='[]'
--       （无任何按钮/写操作权限）+ is_super=0 实现。

-- ============================================================
-- 1. 只读「测试」角色
--    menu_ids: 聚合所有「目录(menu_type=1)」与「页面(menu_type=2)」菜单 ID，
--              即拥有全部菜单的可见/可浏览权限；
--              刻意排除「按钮(menu_type=3)」，配合空 button_perms 实现只读。
--    button_perms: '[]' 空数组，无任何写操作权限标识。
--    is_super: 0 非超级管理员，走常规权限校验分支。
-- ============================================================
INSERT INTO sys_role (role_name, menu_ids, button_perms, is_super, create_time, update_time)
SELECT '测试',
       COALESCE((SELECT jsonb_agg(id ORDER BY id) FROM sys_menu
                 WHERE menu_type IN (1, 2) AND status = 1), '[]'::jsonb),
       '[]'::jsonb,
       0, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_role WHERE role_name = '测试'
);

-- ============================================================
-- 2. 测试账号（关联只读「测试」角色）
--    username: test
--    password: test123 （bcrypt 加密，与项目 internal/service/auth.go 一致）
--    role_id : 子查询回填「测试」角色实际 ID
--    status  : 1 启用
-- ============================================================
INSERT INTO sys_user (username, password, nickname, role_id, status, create_time, update_time)
SELECT 'test',
       '$2a$10$3uxXJpBLW6W7gsAWE/Z00uJfuEHTKpxzcJeEHNW6hFkr5q2Q.XUQ2',
       '测试账号',
       (SELECT id FROM sys_role WHERE role_name = '测试'),
       1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sys_user WHERE username = 'test'
);
