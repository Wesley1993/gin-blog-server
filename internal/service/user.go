package service

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/pkg/logger"
	pkgredis "wuzhispace.com/pkg/redis"
)

// UserService 用户管理服务
type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
	db       *gorm.DB
	redis    *pkgredis.RedisClient
}

// NewUserService 创建用户管理服务实例
func NewUserService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	db *gorm.DB,
	redis *pkgredis.RedisClient,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		db:       db,
		redis:    redis,
	}
}

// Page 分页查询用户
func (s *UserService) Page(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.Page(page, pageSize)
}

// Create 创建用户
func (s *UserService) Create(req dto.CreateUserReq) error {
	// 检查用户名唯一
	existing, _ := s.userRepo.FindByUsername(req.Username)
	if existing != nil {
		return fmt.Errorf("10001") // ErrUsernameExists
	}

	// bcrypt 加密密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("创建用户加密密码失败", "username", req.Username, "error", err)
		return fmt.Errorf("500")
	}

	status := req.Status
	if status == 0 {
		status = 1 // 默认启用
	}

	now := time.Now()
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPwd),
		Nickname: req.Nickname,
		RoleID:   req.RoleID,
		Status:   status,
	}
	user.CreateTime = now
	user.UpdateTime = now

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("创建用户失败", "username", req.Username, "error", err)
		return err
	}
	return nil
}

// Update 更新用户
func (s *UserService) Update(req dto.UpdateUserReq) error {
	user, err := s.userRepo.FindByID(req.ID)
	if err != nil {
		logger.Error("更新用户查询失败", "userID", req.ID, "error", err)
		return fmt.Errorf("500")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.RoleID != 0 {
		user.RoleID = req.RoleID
	}
	if req.Status != 0 {
		user.Status = req.Status
	}
	user.UpdateTime = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		logger.Error("更新用户失败", "userID", user.ID, "error", err)
		return fmt.Errorf("500")
	}

	// 清除该用户 Redis 缓存
	ClearUserCache(s.redis, user.ID)
	return nil
}

// UpdateStatus 更新用户启用/禁用状态（允许置 0，与 Update 的非零判断区分）
func (s *UserService) UpdateStatus(userID int64, status int16) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Error("启停用用户查询失败", "userID", userID, "error", err)
		return fmt.Errorf("500")
	}

	if status != 1 {
		status = 0
	}
	user.Status = status
	user.UpdateTime = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		logger.Error("启停用用户保存失败", "userID", userID, "status", status, "error", err)
		return fmt.Errorf("500")
	}

	// 清除该用户 Redis 缓存
	ClearUserCache(s.redis, user.ID)
	logger.Info("启停用用户成功", "userID", userID, "status", status)
	return nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(userID int64, password string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Error("重置密码查询用户失败", "userID", userID, "error", err)
		return fmt.Errorf("500")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("重置密码加密失败", "userID", userID, "error", err)
		return fmt.Errorf("500")
	}

	user.Password = string(hashedPwd)
	user.UpdateTime = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		logger.Error("重置密码保存失败", "userID", userID, "error", err)
		return fmt.Errorf("500")
	}

	// 清除该用户缓存
	ClearUserCache(s.redis, user.ID)
	logger.Info("重置密码成功", "userID", userID)
	return nil
}

// Delete 删除用户
func (s *UserService) Delete(userID int64) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Error("删除用户查询失败", "userID", userID, "error", err)
		return fmt.Errorf("500")
	}

	// 检查是否超管
	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		logger.Error("删除用户查询角色失败", "userID", userID, "roleID", user.RoleID, "error", err)
		return fmt.Errorf("500")
	}
	if role.IsSuper == 1 {
		logger.Warn("删除用户拒绝，超管不可删", "userID", userID)
		return fmt.Errorf("10003") // ErrSuperAdminDel
	}

	if err := s.userRepo.Delete(userID); err != nil {
		logger.Error("删除用户失败", "userID", userID, "error", err)
		return fmt.Errorf("500")
	}

	// 清除该用户缓存
	ClearUserCache(s.redis, userID)
	logger.Info("删除用户成功", "userID", userID, "username", user.Username)
	return nil
}
