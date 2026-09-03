package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/pkg/logger"
	pkgredis "wuzhispace.com/pkg/redis"
)

// RoleService 角色管理服务
type RoleService struct {
	roleRepo *repository.RoleRepository
	userRepo *repository.UserRepository
	db       *gorm.DB
	redis    *pkgredis.RedisClient
}

// NewRoleService 创建角色管理服务实例
func NewRoleService(
	roleRepo *repository.RoleRepository,
	userRepo *repository.UserRepository,
	db *gorm.DB,
	redis *pkgredis.RedisClient,
) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
		db:       db,
		redis:    redis,
	}
}

// List 角色列表
func (s *RoleService) List() ([]model.Role, error) {
	return s.roleRepo.FindAll()
}

// Create 创建角色
func (s *RoleService) Create(req dto.CreateRoleReq) error {
	now := time.Now()
	role := &model.Role{
		RoleName:    req.RoleName,
		MenuIDs:     req.MenuIDs,
		ButtonPerms: req.ButtonPerms,
		IsSuper:     0,
	}
	role.CreateTime = now
	role.UpdateTime = now

	if err := s.roleRepo.Create(role); err != nil {
		logger.Error("创建角色失败", "roleName", req.RoleName, "error", err)
		return err
	}
	logger.Info("创建角色成功", "roleID", role.ID, "roleName", role.RoleName)
	return nil
}

// Update 更新角色
func (s *RoleService) Update(req dto.UpdateRoleReq) error {
	role, err := s.roleRepo.FindByID(req.ID)
	if err != nil {
		logger.Error("更新角色查询失败", "roleID", req.ID, "error", err)
		return fmt.Errorf("10002") // ErrRoleNotFound
	}

	role.RoleName = req.RoleName
	if req.MenuIDs != nil {
		role.MenuIDs = req.MenuIDs
	}
	if req.ButtonPerms != nil {
		role.ButtonPerms = req.ButtonPerms
	}
	role.UpdateTime = time.Now()

	if err := s.roleRepo.Update(role); err != nil {
		logger.Error("更新角色失败", "roleID", role.ID, "error", err)
		return fmt.Errorf("500")
	}

	// 批量清除该角色下所有用户的 perm/menu 缓存
	ClearRoleUsersCache(s.redis, s.userRepo, role.ID)
	logger.Info("更新角色成功", "roleID", role.ID, "roleName", role.RoleName)
	return nil
}

// Delete 删除角色
func (s *RoleService) Delete(id int64) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		logger.Error("删除角色查询失败", "roleID", id, "error", err)
		return fmt.Errorf("10002") // ErrRoleNotFound
	}

	// 超管角色不可删
	if role.IsSuper == 1 {
		logger.Warn("删除角色拒绝，超管角色不可删", "roleID", id)
		return fmt.Errorf("10003") // ErrSuperAdminDel
	}

	if err := s.roleRepo.Delete(id); err != nil {
		logger.Error("删除角色失败", "roleID", id, "error", err)
		return fmt.Errorf("500")
	}

	// 清除该角色下所有用户的缓存
	ClearRoleUsersCache(s.redis, s.userRepo, id)
	logger.Info("删除角色成功", "roleID", id, "roleName", role.RoleName)
	return nil
}
