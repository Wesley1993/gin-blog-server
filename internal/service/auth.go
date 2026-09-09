package service

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"wuzhispace.com/internal/dto"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/internal/repository"
	appjwt "wuzhispace.com/pkg/jwt"
	"wuzhispace.com/pkg/logger"
	pkgredis "wuzhispace.com/pkg/redis"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
	menuRepo *repository.MenuRepository
	db       *gorm.DB
	redis    *pkgredis.RedisClient
}

// NewAuthService 创建认证服务实例
func NewAuthService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	menuRepo *repository.MenuRepository,
	db *gorm.DB,
	redis *pkgredis.RedisClient,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		menuRepo: menuRepo,
		db:       db,
		redis:    redis,
	}
}

// Login 用户登录
func (s *AuthService) Login(req dto.LoginReq) (*dto.LoginResp, error) {
	// 1. 查用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		logger.Warn("用户登录失败", "username", req.Username, "error", err)
		return nil, fmt.Errorf("10005") // ErrLoginFailed
	}

	// 2. bcrypt 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn("用户登录失败，密码错误", "username", req.Username)
		return nil, fmt.Errorf("10005") // ErrLoginFailed
	}

	// 3. 检查状态
	if user.Status != 1 {
		logger.Warn("用户登录被拒绝，账号已禁用", "username", req.Username)
		return nil, fmt.Errorf("10004") // ErrAccountDisabled
	}

	// 4. 查角色信息
	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		logger.Error("登录查询角色失败", "userID", user.ID, "roleID", user.RoleID, "error", err)
		return nil, fmt.Errorf("10002") // ErrRoleNotFound
	}

	// 5. 生成 JWT token
	token, err := appjwt.GenerateToken(user.ID, user.Username, user.RoleID, role.IsSuper)
	if err != nil {
		logger.Error("登录生成JWT失败", "userID", user.ID, "error", err)
		return nil, fmt.Errorf("500")
	}

	// 6. 写入 Redis
	ttl := appjwt.GetExpiry()
	if err := s.redis.SetWithTTL(TokenKey(user.ID), token, ttl); err != nil {
		logger.Error("登录写入Redis token失败", "userID", user.ID, "error", err)
		return nil, fmt.Errorf("500")
	}

	// 7. 返回 token
	logger.Info("用户登录成功", "userID", user.ID, "username", user.Username)
	return &dto.LoginResp{Token: token}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(userID int64) error {
	s.redis.Del(TokenKey(userID))
	logger.Info("用户登出", "userID", userID)
	return nil
}

// ChangePassword 当前登录用户自助修改密码
func (s *AuthService) ChangePassword(userID int64, req dto.ChangePwdReq) error {
	// 1. 查用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Error("修改密码查询用户失败", "userID", userID, "error", err)
		return fmt.Errorf("500")
	}

	// 2. bcrypt 校验原密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		logger.Warn("修改密码失败，原密码错误", "userID", userID)
		return fmt.Errorf("10006") // ErrOldPwdIncorrect
	}

	// 3. 新密码不能与原密码相同
	if req.OldPassword == req.NewPassword {
		return fmt.Errorf("10007") // ErrPwdSame
	}

	// 4. 加密新密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("修改密码加密失败", "userID", userID, "error", err)
		return fmt.Errorf("500")
	}

	// 5. 落库：仅更新 password 与 update_time 两列，避免整行 Save 覆盖慢窗口内的并发修改
	if err := s.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password":    string(hashed),
		"update_time": time.Now(),
	}).Error; err != nil {
		logger.Error("修改密码保存失败", "userID", userID, "error", err)
		return fmt.Errorf("500")
	}

	// 6. 清除缓存，强制下线
	ClearUserCache(s.redis, userID)
	logger.Info("修改密码成功", "userID", userID)
	return nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID int64) (*dto.UserInfoResp, error) {
	// 1. 查用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Error("获取用户信息查询用户失败", "userID", userID, "error", err)
		return nil, fmt.Errorf("500")
	}

	// 2. 查角色信息
	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		logger.Error("获取用户信息查询角色失败", "userID", userID, "roleID", user.RoleID, "error", err)
		return nil, fmt.Errorf("10002")
	}

	// 3. 获取菜单（先查缓存）
	var menus []model.Menu
	menuCacheKey := RBACMenuKey(userID)
	cachedMenus, err := s.redis.Get(menuCacheKey)
	if err == nil && cachedMenus != "" {
		_ = json.Unmarshal([]byte(cachedMenus), &menus)
	} else {
		// 缓存未命中，查 DB
		if role.IsSuper == 1 {
			// 超管查全部菜单
			menus, err = s.menuRepo.FindAll()
			if err != nil {
				logger.Error("获取用户信息查询全部菜单失败", "userID", userID, "error", err)
				return nil, fmt.Errorf("500")
			}
		} else {
			// 按 role.menu_ids 查菜单
			var menuIDs []int64
			if len(role.MenuIDs) > 0 {
				if err := json.Unmarshal(role.MenuIDs, &menuIDs); err != nil {
					logger.Error("获取用户信息解析角色菜单ID失败", "userID", userID, "roleID", role.ID, "error", err)
					return nil, fmt.Errorf("500")
				}
			}
			if len(menuIDs) > 0 {
				menus, err = s.menuRepo.FindByIDs(menuIDs)
				if err != nil {
					logger.Error("获取用户信息按ID查询菜单失败", "userID", userID, "roleID", role.ID, "error", err)
					return nil, fmt.Errorf("500")
				}
			}
		}
		// 写入缓存
		if menuJSON, err := json.Marshal(menus); err == nil {
			_ = s.redis.SetWithTTL(menuCacheKey, string(menuJSON), CacheTTL)
		}
	}

	// 构建菜单树
	menuTree := s.menuRepo.BuildTree(menus)

	// 4. 获取权限标识集合（先查缓存）
	var permissions []string
	permCacheKey := RBACPermKey(userID)
	cachedPerms, err := s.redis.Get(permCacheKey)
	if err == nil && cachedPerms != "" {
		_ = json.Unmarshal([]byte(cachedPerms), &permissions)
	} else {
		// 缓存未命中，从 role.button_perms 解析
		if role.IsSuper == 1 {
			// 超管拥有所有权限，下发 '*' 通配标识（前端据此判定超管）
			permissions = []string{"*"}
		} else if len(role.ButtonPerms) > 0 {
			_ = json.Unmarshal(role.ButtonPerms, &permissions)
		}
		// 写入缓存
		if permJSON, err := json.Marshal(permissions); err == nil {
			_ = s.redis.SetWithTTL(permCacheKey, string(permJSON), CacheTTL)
		}
	}

	// 超管始终保证 '*' 通配标识存在（兼容修复前写入的旧缓存）
	if role.IsSuper == 1 && (len(permissions) == 0 || permissions[0] != "*") {
		permissions = []string{"*"}
	}

	// 5. 组装响应
	resp := &dto.UserInfoResp{
		User: dto.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			RoleID:   user.RoleID,
		},
		Menus:       menuTree,
		Permissions: permissions,
	}

	// 确保空切片返回 [] 而非 null
	if resp.Menus == nil {
		resp.Menus = []model.Menu{}
	}
	if resp.Permissions == nil {
		resp.Permissions = []string{}
	}

	return resp, nil
}
