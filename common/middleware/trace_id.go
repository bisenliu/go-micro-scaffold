package middleware

import (
	"common/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TraceIDMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取 traceID
		traceID := c.GetHeader("X-Trace-ID")

		// 如果请求头中没有 traceID，则生成一个新的
		if traceID == "" {
			traceID = logger.GenerateTraceID()
		}

		// 将 traceID 添加到 Gin context 中
		c.Set("traceID", traceID)

		// 将 traceID 添加到 Go context 中
		ctx := logger.WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// 将 traceID 添加到响应头
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}
