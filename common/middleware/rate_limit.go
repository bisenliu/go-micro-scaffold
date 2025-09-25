package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"go.uber.org/zap"

	"common/logger"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(fillInterval time.Duration, capacity, quantum int64, zapLogger *zap.Logger) gin.HandlerFunc {
	// 创建令牌桶
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, capacity, quantum)

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 尝试获取令牌
		if bucket.TakeAvailable(1) == 0 {
			logger.Warn(zapLogger, ctx, "Rate limit exceeded",
				zap.String("client_ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
