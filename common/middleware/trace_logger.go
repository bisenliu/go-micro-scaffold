package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
)

// TraceLoggerMiddleware 创建 Logger 中间件，自动为每个请求注入带 traceID 的 logger
func TraceLoggerMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 尝试从请求头获取 traceID，如果没有则生成新的
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = logger.GenerateTraceID()
		}

		// 将 traceID 添加到 Gin context 和 Go context 中
		c.Set("traceID", traceID)
		ctx = logger.WithTraceID(ctx, traceID)

		// 创建带 traceID 的 logger 并存入 context
		ctxLogger := logger.WithContext(zapLogger, ctx)
		ctx = logger.ToContext(ctx, ctxLogger)

		// 更新 request context
		c.Request = c.Request.WithContext(ctx)

		// 将 traceID 添加到响应头
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}
