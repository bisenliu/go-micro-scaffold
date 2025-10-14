package middleware

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/response"
)

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.RecoveryWithWriter(nil, func(c *gin.Context, recovered interface{}) {
		ctx := c.Request.Context()

		// 记录错误日志
		logger.Error(ctx, "Panic recovered",
			zap.Any("recovered", recovered),
			zap.String("stack", string(debug.Stack())),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()))

		// 返回统一的错误响应
		response.InternalServerError(c, "Internal server error")
	})
}
