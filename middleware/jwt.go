package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/jwt"
	"wuzhispace.com/pkg/logger"
	pkgredis "wuzhispace.com/pkg/redis"
	"wuzhispace.com/pkg/response"
)

// JWTAuth JWT 鉴权中间件
func JWTAuth(redisClient *pkgredis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 提取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("鉴权失败，Authorization头缺失", "ip", c.ClientIP(), "path", c.Request.URL.Path)
			response.Fail(c, http.StatusUnauthorized, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
			c.Abort()
			return
		}

		// 解析 Bearer {token}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("鉴权失败，Authorization格式非法", "ip", c.ClientIP(), "path", c.Request.URL.Path)
			response.Fail(c, http.StatusUnauthorized, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 JWT
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			logger.Warn("JWT校验失败", "error", err)
			response.Fail(c, http.StatusUnauthorized, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
			c.Abort()
			return
		}

		// 检查 Redis 中 token 是否存在
		tokenKey := fmt.Sprintf("token:%d", claims.UserID)
		cachedToken, err := redisClient.Get(tokenKey)
		if err != nil {
			logger.Warn("鉴权失败，Redis查询token出错", "userID", claims.UserID, "ip", c.ClientIP(), "error", err)
			response.Fail(c, http.StatusUnauthorized, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
			c.Abort()
			return
		}
		if cachedToken == "" {
			logger.Warn("鉴权失败，Redis中token不存在", "userID", claims.UserID, "ip", c.ClientIP())
			response.Fail(c, http.StatusUnauthorized, errors.CodeUnauthorized, errors.GetMsg(errors.CodeUnauthorized))
			c.Abort()
			return
		}

		// 将用户信息写入 Context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)
		c.Set("isSuper", claims.IsSuper)

		c.Next()
	}
}
