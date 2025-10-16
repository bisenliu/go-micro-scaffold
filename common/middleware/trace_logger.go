package middleware

import (
	"github.com/gin-gonic/gin"

	"common/interfaces"
	"common/logger"
	"common/pkg/contextutil"
)

// TraceLoggerMiddleware 创建 Logger 中间件，自动为每个请求注入带 traceID 的 logger
func TraceLoggerMiddleware(loggerInstance interfaces.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 尝试从请求头获取 traceID，如果没有则生成新的
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = logger.GenerateTraceID()
		}

		// 将 traceID 添加到 Gin context 和 Go context 中
		c.Set(contextutil.TraceIDKey, traceID)
		ctx = logger.WithTraceID(ctx, traceID)

		// 创建带有traceID的logger并存储到context中
		ctxLogger := loggerInstance.WithContext(ctx)
		ctx = logger.ToContext(ctx, ctxLogger)

		// 更新 request context
		c.Request = c.Request.WithContext(ctx)

		// 将 traceID 添加到响应头
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}
