package service

import (
	"fmt"
	"time"

	pkgredis "wuzhispace.com/pkg/redis"
	"wuzhispace.com/internal/repository"
	"wuzhispace.com/pkg/logger"
)

const (
	// TokenKeyPrefix token缓存前缀
	TokenKeyPrefix = "token:"
	// RBACPermKeyPrefix 权限缓存前缀
	RBACPermKeyPrefix = "rbac:perm:"
	// RBACMenuKeyPrefix 菜单缓存前缀
	RBACMenuKeyPrefix = "rbac:menu:"
	// CacheTTL 缓存过期时间
	CacheTTL = 24 * time.Hour
)

// TokenKey 生成token缓存key
func TokenKey(userID int64) string {
	return fmt.Sprintf("%s%d", TokenKeyPrefix, userID)
}

// RBACPermKey 生成权限缓存key
func RBACPermKey(userID int64) string {
	return fmt.Sprintf("%s%d", RBACPermKeyPrefix, userID)
}

// RBACMenuKey 生成菜单缓存key
func RBACMenuKey(userID int64) string {
	return fmt.Sprintf("%s%d", RBACMenuKeyPrefix, userID)
}

// ClearUserCache 清除用户所有缓存（token/perm/menu）
func ClearUserCache(redisClient *pkgredis.RedisClient, userID int64) {
	redisClient.Del(TokenKey(userID), RBACPermKey(userID), RBACMenuKey(userID))
	logger.Info("清除用户缓存", "userID", userID)
}

// ClearRoleUsersCache 批量清除某角色下所有用户的 perm/menu 缓存
func ClearRoleUsersCache(redisClient *pkgredis.RedisClient, userRepo *repository.UserRepository, roleID int64) {
	users, err := userRepo.FindByRoleID(roleID)
	if err != nil {
		return
	}
	for _, user := range users {
		redisClient.Del(RBACPermKey(user.ID), RBACMenuKey(user.ID))
	}
}

// ClearAllRBACCache 清除全体用户的 rbac:menu:* 与 rbac:perm:* 缓存。
// 菜单/按钮变更后调用，避免用户等待 24h TTL 或重新登录才能看到新菜单/权限。
// 采用 SCAN 前缀扫描而非遍历用户表，可覆盖缓存中存在但用户表已变更的孤儿键。
func ClearAllRBACCache(redisClient *pkgredis.RedisClient) {
	menuN, err := redisClient.ScanDel(RBACMenuKeyPrefix + "*")
	if err != nil {
		logger.Warn("清除全体菜单缓存失败", "error", err)
	}
	permN, err := redisClient.ScanDel(RBACPermKeyPrefix + "*")
	if err != nil {
		logger.Warn("清除全体权限缓存失败", "error", err)
	}
	logger.Info("清除全体用户RBAC缓存", "menuKeys", menuN, "permKeys", permN)
}
