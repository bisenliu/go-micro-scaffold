package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/response"
	"services/internal/application/service"
)

// CasbinAuthzMiddleware 创建基于Casbin的授权中间件
func CasbinAuthzMiddleware(permissionService *service.PermissionService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		allowed, err := permissionService.Enforce(c.Request.Context(), userID, resource, action)
		if err != nil {
			logger.Error("Failed to enforce policy",
				zap.String("user_id", userID),
				zap.String("resource", resource),
				zap.String("action", action),
				zap.Error(err))
			response.InternalServerError(c, "Authorization check failed")
			c.Abort()
			return
		}

		if !allowed {
			logger.Warn("Access denied",
				zap.String("user_id", userID),
				zap.String("resource", resource),
				zap.String("action", action))
			response.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		// 权限检查通过，继续处理请求
		logger.Debug("Access granted",
			zap.String("user_id", userID),
			zap.String("resource", resource),
			zap.String("action", action))
		c.Next()
	}
}
