package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/response"
)

// PermissionEnforceFunc 权限检查函数类型
type PermissionEnforceFunc func(ctx context.Context, sub, obj, act string) (bool, error)

// CasbinMiddleware 创建基于Casbin的授权中间件
func CasbinMiddleware(enforceFunc PermissionEnforceFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// 这里应该从 JWT token 或 session 中获取用户信息
		// 为了演示，我们使用一个示例用户 ID
		userID := c.GetString("user_id")
		if userID == "" {
			// 如果没有用户 ID，可以尝试从其他地方获取，比如 JWT claims
			userID = "anonymous"
		}

		// 获取请求的资源和操作
		resource := c.Request.URL.Path
		action := c.Request.Method

		// 检查权限
		allowed, err := enforceFunc(c.Request.Context(), userID, resource, action)
		if err != nil {
			logger.Error(ctx, "Failed to enforce policy",
				zap.String("user_id", userID),
				zap.String("resource", resource),
				zap.String("action", action),
				zap.Error(err))
			response.InternalServerError(c, "Authorization check failed")
			c.Abort()
			return
		}

		if !allowed {
			logger.Warn(ctx, "Access denied",
				zap.String("user_id", userID),
				zap.String("resource", resource),
				zap.String("action", action))
			response.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		// 权限检查通过，继续处理请求
		logger.Debug(ctx, "Access granted",
			zap.String("user_id", userID),
			zap.String("resource", resource),
			zap.String("action", action))
		c.Next()
	}
}
