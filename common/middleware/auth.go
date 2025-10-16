package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/interfaces"
)

const (
	// AuthHeaderKey 认证头键名
	AuthHeaderKey = "Authorization"
	// TokenPrefix Token前缀
	TokenPrefix = "Bearer "
	// UserIDKey 用户ID在context中的键名
	UserIDKey = "user_id"
)

// AuthMiddleware 认证中间件
// 职责：只负责身份验证，不处理响应格式
func AuthMiddleware(jwtService interfaces.JWTService, authConfig interfaces.AuthConfig, logger interfaces.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 检查是否为白名单路径
		if isWhitelistedPath(c.Request.URL.Path, authConfig.Whitelist) {
			logger.Debug(ctx, "Whitelisted path, skipping authentication")
			c.Next()
			return
		}

		// 获取Authorization头
		authHeader := c.GetHeader(AuthHeaderKey)
		if authHeader == "" {
			logger.Warn(ctx, "Missing authorization header")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, TokenPrefix) {
			logger.Warn(ctx, "Invalid authorization header format")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 提取token
		token := strings.TrimPrefix(authHeader, TokenPrefix)
		if token == "" {
			logger.Warn(ctx, "Empty token")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 验证token并获取用户信息
		userID, err := jwtService.ValidateToken(token)
		if err != nil {
			logger.Error(ctx, "Token validation failed", 
				zap.Error(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 将用户ID存储到context中
		c.Set(UserIDKey, userID)
		ctx = context.WithValue(ctx, UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)

		logger.Debug(ctx, "Authentication successful")

		c.Next()
	}
}

// isWhitelistedPath 检查是否为白名单路径
func isWhitelistedPath(path string, whiteList []string) bool {
	for _, p := range whiteList {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// GetCurrentUserID 从context中获取当前用户ID
func GetCurrentUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	
	if uid, ok := userID.(string); ok {
		return uid, true
	}
	
	return "", false
}

// RequireAuth 要求认证的中间件工厂
func RequireAuth(jwtService interfaces.JWTService, authConfig interfaces.AuthConfig, logger interfaces.Logger) gin.HandlerFunc {
	return AuthMiddleware(jwtService, authConfig, logger)
}