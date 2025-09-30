package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"go.uber.org/zap"

	"common/logger"
	"common/response"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(fillInterval int64, capacity int64, zapLogger *zap.Logger) gin.HandlerFunc {
	// 创建令牌桶
	bucket := ratelimit.NewBucketWithRate(float64(fillInterval), capacity)

	zapLogger.Info("Rate limit middleware initialized",
		zap.Int64("fill_interval", fillInterval),
		zap.Int64("capacity", capacity))

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 尝试获取令牌
		if bucket.TakeAvailable(1) == 0 {
			logger.Warn(ctx, "Rate limit exceeded",
				zap.String("client_ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))

			response.BadRequest(c, "Rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
