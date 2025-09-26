package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
)

const (
	// UserIDKey 用户ID在context中的键名
	UserIDKey = "userID"
	// AuthHeaderKey 认证头名称
	AuthHeaderKey = "Authorization"
	// TokenPrefix Token前缀
	TokenPrefix = "Bearer "
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 检查是否为白名单路径
		if isWhitelistedPath(c.Request.URL.Path) {
			logger.Debug(ctx, "Whitelisted path, skipping auth",
				zap.String("path", c.Request.URL.Path))
			c.Next()
			return
		}

		// 获取Authorization头
		authHeader := c.GetHeader(AuthHeaderKey)
		if authHeader == "" {
			logger.Warn(ctx, "Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Missing authorization header",
			})
			c.Abort()
			return
		}

		var token string
		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, TokenPrefix) {
			logger.Warn(ctx, "Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// 提取token
		token = strings.TrimPrefix(authHeader, TokenPrefix)
		if token == "" {
			logger.Warn(ctx, "Empty token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Empty token",
			})
			c.Abort()
			return
		}

		// 验证token并获取用户信息
		userID, err := validateToken(ctx, token, zapLogger)
		if err != nil {
			logger.Error(ctx, "Token validation failed", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// 将用户ID存储到context中
		c.Set(UserIDKey, userID)
		ctx = context.WithValue(ctx, UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)

		logger.Debug(ctx, "Authentication successful",
			zap.Int64("user_id", userID))

		c.Next()
	}
}

// isWhitelistedPath 检查是否为白名单路径
func isWhitelistedPath(path string) bool {
	// 白名单路径列表
	whiteList := []string{
		"/health",
		"/ping",
	}

	for _, p := range whiteList {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// validateToken 验证token并返回用户ID
// TODO: 这里需要根据实际的JWT实现来调整
func validateToken(ctx context.Context, token string, logger *zap.Logger) (int64, error) {
	// 这里应该实现实际的JWT验证逻辑
	// 包括：
	// 1. 解析JWT token
	// 2. 验证签名
	// 3. 检查过期时间
	// 4. 从Redis中验证token是否有效
	// 5. 返回用户ID

	// 临时实现，实际使用时需要替换
	logger.Debug("Token validation not implemented yet", zap.String("token", token[:10]+"..."))

	// 返回模拟的用户ID，实际使用时需要从JWT中解析
	return 123456, nil
}

// GetUserIDFromContext 从context中获取用户ID
func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	if id, ok := userID.(int64); ok {
		return id, true
	}

	return 0, false
}
