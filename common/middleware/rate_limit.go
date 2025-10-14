package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"

	"common/logger"
	"common/response"
)

// 使用 sync.Map 来存储每个IP对应的令牌桶
var ipBuckets sync.Map

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var ip string
		// 优先从 context 获取真实 IP
		if clientIP, exists := c.Get(ClientIPContextKey); exists {
			ip, _ = clientIP.(string)
		}

		// 如果 context 中没有，则回退到 gin 的默认方法
		if ip == "" {
			logger.Warn(ctx, "Client IP not found in context, falling back to c.ClientIP() for rate limiting")
			ip = c.ClientIP()
		}

		if ip == "" {
			logger.Error(ctx, "Could not determine client IP for rate limiting")
			response.InternalServerError(c, "Could not determine client IP")
			c.Abort()
			return
		}

		// 获取或创建一个令牌桶
		bucket, _ := ipBuckets.LoadOrStore(ip, ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum))

		// 强制类型转换
		ipBucket := bucket.(*ratelimit.Bucket)

		// 检查是否有可用令牌
		if ipBucket.TakeAvailable(1) < 1 {
			response.RateLimit(c, "Rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}

}
