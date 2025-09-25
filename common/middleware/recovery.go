package middleware

import (
	"common/logger"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GracefulRecoveryMiddleware 优雅恢复中间件
func GracefulRecoveryMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		ctx := c.Request.Context()

		// 记录panic信息
		logger.Error(zapLogger, ctx, "Panic recovered",
			zap.Any("error", recovered),
			zap.String("stack", string(debug.Stack())),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()))

		// 返回统一的错误响应
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal server error",
			"error":   fmt.Sprintf("%v", recovered),
		})
	})
}
