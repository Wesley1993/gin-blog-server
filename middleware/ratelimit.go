package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wuzhispace.com/pkg/errors"
	"wuzhispace.com/pkg/logger"
	pkgredis "wuzhispace.com/pkg/redis"
	"wuzhispace.com/pkg/response"
)

// RateLimit 登录接口限流中间件（5次/分钟/IP）
func RateLimit(redisClient *pkgredis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate:login:%s", clientIP)

		count, err := redisClient.Incr(key)
		if err != nil {
			// Redis 故障不阻断请求，静默放行但记录告警
			logger.Warn("限流器Redis故障，静默放行", "ip", clientIP, "error", err)
			c.Next()
			return
		}

		// 首次访问设置过期时间
		if count == 1 {
			_, _ = redisClient.Expire(key, 60*time.Second)
		}

		// 超过 5 次返回 429
		if count > 5 {
			logger.Warn("登录限流触发", "ip", clientIP, "count", count)
			response.Fail(c, http.StatusTooManyRequests, errors.CodeTooManyReq, errors.GetMsg(errors.CodeTooManyReq))
			c.Abort()
			return
		}

		c.Next()
	}
}
