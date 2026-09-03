package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/pkg/logger"
	pkgredis "wuzhispace.com/pkg/redis"
)

// MenuService 菜单管理服务
type MenuService struct {
	menuRepo *repository.MenuRepository
	roleRepo *repository.RoleRepository
	userRepo *repository.UserRepository
	redis    *pkgredis.RedisClient
}

// NewMenuService 创建菜单管理服务实例
func NewMenuService(
	menuRepo *repository.MenuRepository,
	roleRepo *repository.RoleRepository,
	userRepo *repository.UserRepository,
	redis *pkgredis.RedisClient,
) *MenuService {
	return &MenuService{
		menuRepo: menuRepo,
		roleRepo: roleRepo,
		userRepo: userRepo,
		redis:    redis,
	}
}

// List 获取菜单树
func (s *MenuService) List() ([]model.Menu, error) {
	menus, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, err
	}

	tree := s.menuRepo.BuildTree(menus)
	if tree == nil {
		tree = []model.Menu{}
	}
	return tree, nil
}

// Create 创建菜单（有 path 时校验唯一性）
func (s *MenuService) Create(req *dto.CreateMenuReq) error {
	if req.Path != "" {
		if _, err := s.menuRepo.FindByPath(req.Path); err == nil {
			logger.Warn("创建菜单拒绝，path已存在", "path", req.Path)
			return fmt.Errorf("70003") // ErrMenuPathExists
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("创建菜单查重失败", "path", req.Path, "error", err)
			return fmt.Errorf("500")
		}
	}

	now := time.Now()
	status := req.Status
	if status == 0 {
		status = 1 // 默认启用
	}
	menu := &model.Menu{
		ParentID: req.ParentID,
		MenuName: req.MenuName,
		MenuType: int16(req.MenuType),
		Path:     req.Path,
		Perms:    req.Perms,
		Sort:     req.Sort,
		Status:   int16(status),
	}
	menu.CreateTime = now
	menu.UpdateTime = now

	if err := s.menuRepo.Create(menu); err != nil {
		logger.Error("创建菜单失败", "menuName", req.MenuName, "path", req.Path, "error", err)
		return err
	}
	// 菜单变更后清全体用户 RBAC 缓存，无需等 24h TTL 或重登即可看到新菜单/按钮
	ClearAllRBACCache(s.redis)
	logger.Info("创建菜单成功", "menuID", menu.ID, "menuName", menu.MenuName, "path", menu.Path)
	return nil
}

// Update 更新菜单（校验存在性，有 path 时校验唯一性）
func (s *MenuService) Update(req *dto.UpdateMenuReq) error {
	menu, err := s.menuRepo.FindByID(req.ID)
	if err != nil {
		logger.Error("更新菜单查询失败", "menuID", req.ID, "error", err)
		return fmt.Errorf("70001") // ErrMenuNotFound
	}

	if req.Path != "" && req.Path != menu.Path {
		existing, err := s.menuRepo.FindByPath(req.Path)
		if err == nil && existing.ID != menu.ID {
			logger.Warn("更新菜单拒绝，path已存在", "menuID", req.ID, "path", req.Path)
			return fmt.Errorf("70003") // ErrMenuPathExists
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("更新菜单查重失败", "menuID", req.ID, "path", req.Path, "error", err)
			return fmt.Errorf("500")
		}
	}

	if req.MenuName != "" {
		menu.MenuName = req.MenuName
	}
	menu.ParentID = req.ParentID
	if req.MenuType != 0 {
		menu.MenuType = int16(req.MenuType)
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Perms != "" {
		menu.Perms = req.Perms
	}
	menu.Sort = req.Sort
	if req.Status != 0 {
		menu.Status = int16(req.Status)
	}
	menu.UpdateTime = time.Now()

	if err := s.menuRepo.Update(menu); err != nil {
		logger.Error("更新菜单失败", "menuID", menu.ID, "error", err)
		return err
	}
	// 菜单变更后清全体用户 RBAC 缓存，无需等 24h TTL 或重登即可看到新菜单/按钮
	ClearAllRBACCache(s.redis)
	logger.Info("更新菜单成功", "menuID", menu.ID, "menuName", menu.MenuName, "path", menu.Path)
	return nil
}

// Delete 删除菜单：校验无子菜单 → 从引用角色中移除 → 清缓存 → 删除
func (s *MenuService) Delete(id int64) error {
	if _, err := s.menuRepo.FindByID(id); err != nil {
		logger.Error("删除菜单查询失败", "menuID", id, "error", err)
		return fmt.Errorf("70001") // ErrMenuNotFound
	}

	// 1. 存在子菜单时禁止删除
	children, err := s.menuRepo.FindByParentID(id)
	if err != nil {
		logger.Error("删除菜单查询子菜单失败", "menuID", id, "error", err)
		return fmt.Errorf("500")
	}
	if len(children) > 0 {
		logger.Warn("删除菜单拒绝，存在子菜单", "menuID", id, "childrenCount", len(children))
		return fmt.Errorf("70002") // ErrMenuHasChildren
	}

	// 2. 查询所有引用该菜单的角色，从 menu_ids 中移除该ID
	roles, err := s.roleRepo.FindByMenuID(id)
	if err != nil {
		logger.Error("删除菜单查询引用角色失败", "menuID", id, "error", err)
		return fmt.Errorf("500")
	}
	for i := range roles {
		var menuIDs []int64
		if len(roles[i].MenuIDs) > 0 {
			if err := json.Unmarshal(roles[i].MenuIDs, &menuIDs); err != nil {
				logger.Error("删除菜单解析角色菜单ID失败", "menuID", id, "roleID", roles[i].ID, "error", err)
				return fmt.Errorf("500")
			}
		}
		filtered := make([]int64, 0, len(menuIDs))
		for _, mid := range menuIDs {
			if mid != id {
				filtered = append(filtered, mid)
			}
		}
		raw, err := json.Marshal(filtered)
		if err != nil {
			logger.Error("删除菜单序列化角色菜单ID失败", "menuID", id, "roleID", roles[i].ID, "error", err)
			return fmt.Errorf("500")
		}
		roles[i].MenuIDs = raw
		roles[i].UpdateTime = time.Now()
		if err := s.roleRepo.Update(&roles[i]); err != nil {
			logger.Error("删除菜单清理角色引用失败", "menuID", id, "roleID", roles[i].ID, "error", err)
			return fmt.Errorf("500")
		}
		// 3. 清除受影响角色下所有用户的缓存
		ClearRoleUsersCache(s.redis, s.userRepo, roles[i].ID)
	}

	// 4. 执行删除
	if err := s.menuRepo.Delete(id); err != nil {
		logger.Error("删除菜单失败", "menuID", id, "error", err)
		return fmt.Errorf("500")
	}
	// 5. 删除后同样清全体用户 RBAC 缓存，保证菜单树立即可见（不受影响的用户缓存也不残留旧节点）
	ClearAllRBACCache(s.redis)
	logger.Info("删除菜单成功", "menuID", id)
	return nil
}
