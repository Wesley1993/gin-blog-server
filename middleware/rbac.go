package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"wuzhispace.com/internal/repository"
	"wuzhispace.com/internal/service"
	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/logger"
	pkgredis "wuzhispace.com/pkg/redis"
	"wuzhispace.com/pkg/response"
)

// RBACAuth RBAC 权限校验中间件
func RBACAuth(redisClient *pkgredis.RedisClient, roleRepo *repository.RoleRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Context 取 isSuper，超管直接放行
		isSuperVal, exists := c.Get("isSuper")
		if exists {
			if isSuper, ok := isSuperVal.(int16); ok && isSuper == 1 {
				c.Next()
				return
			}
		}

		// 获取当前路由所需权限标识
		requiredPerm := c.GetString("perm")
		if requiredPerm == "" {
			// 路由没有设置权限标识，默认放行（仅需登录即可）
			c.Next()
			return
		}

		// 获取 userID
		userIDVal, exists := c.Get("userID")
		if !exists {
			response.Fail(c, http.StatusUnauthorized, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
			c.Abort()
			return
		}
		userID := userIDVal.(int64)

		// 从 Redis 读权限集合
		permKey := service.RBACPermKey(userID)
		cachedPerms, err := redisClient.Get(permKey)
		var permissions []string

		if err == nil && cachedPerms != "" {
			// 缓存命中
			_ = json.Unmarshal([]byte(cachedPerms), &permissions)
		} else {
			// 缓存未命中，查 DB
			roleIDVal, _ := c.Get("roleID")
			roleID := roleIDVal.(int64)

			role, err := roleRepo.FindByID(roleID)
			if err != nil {
				logger.Error("权限校验查询角色失败", "userID", userID, "roleID", roleID, "error", err)
				response.Fail(c, http.StatusForbidden, errors.CodeForbidden, errors.GetMsg(errors.CodeForbidden))
				c.Abort()
				return
			}

			if len(role.ButtonPerms) > 0 {
				_ = json.Unmarshal(role.ButtonPerms, &permissions)
			}

			// 写入缓存
			if permJSON, err := json.Marshal(permissions); err == nil {
				_ = redisClient.SetWithTTL(permKey, string(permJSON), service.CacheTTL)
			}
		}

		// 比对权限
		hasPerm := false
		for _, p := range permissions {
			if p == requiredPerm {
				hasPerm = true
				break
			}
		}

		if !hasPerm {
			logger.Warn("权限校验失败", "userID", userID, "perm", requiredPerm)
			// 显式下发中文提示，供前端拦截器统一展示「账号权限不足」
			response.Fail(c, http.StatusForbidden, errors.CodeForbidden, "账号权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}
