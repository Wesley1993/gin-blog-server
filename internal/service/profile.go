package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/repository"
	pkgredis "wuzhispace.com/pkg/redis"
)

// ProfileService 个人信息服务
type ProfileService struct {
	userRepo *repository.UserRepository
	db       *gorm.DB
	redis    *pkgredis.RedisClient
}

// NewProfileService 创建个人信息服务实例
func NewProfileService(userRepo *repository.UserRepository, db *gorm.DB, redis *pkgredis.RedisClient) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
		db:       db,
		redis:    redis,
	}
}

// GetProfile 获取当前用户个人信息
func (s *ProfileService) GetProfile(userID int64) (*dto.ProfileDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("500")
	}

	return &dto.ProfileDTO{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		RoleID:   user.RoleID,
	}, nil
}

// UpdateProfile 更新个人信息（昵称、头像、简介），更新后清除用户缓存
func (s *ProfileService) UpdateProfile(userID int64, req dto.UpdateProfileReq) (*dto.ProfileDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("500")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	user.Avatar = req.Avatar
	user.Bio = req.Bio
	user.UpdateTime = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("500")
	}

	// 清除用户缓存（保留 token，避免更新个人信息后被踢出登录）
	if s.redis != nil {
		s.redis.Del(RBACPermKey(userID), RBACMenuKey(userID))
	}

	return &dto.ProfileDTO{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		RoleID:   user.RoleID,
	}, nil
}
