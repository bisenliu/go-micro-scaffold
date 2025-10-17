package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/config"
	"common/logger"
	"common/pkg/contextutil"
	"common/pkg/jwt"
	"common/response"
)

const (
	// AuthHeaderKey 认证头键名
	AuthHeaderKey = "Authorization"
	// TokenPrefix Token前缀
	TokenPrefix = "Bearer "
)

// AuthMiddleware 认证中间件
func AuthMiddleware(jwtService *jwt.JWT, cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 检查是否为白名单路径
		if isWhitelistedPath(c.Request.URL.Path, cfg.Whitelist) {
			logger.Debug(ctx, "Whitelisted path, skipping authentication",
				zap.String("path", c.Request.URL.Path))
			c.Next()
			return
		}

		// 获取Authorization头
		authHeader := c.GetHeader(AuthHeaderKey)
		if authHeader == "" {
			logger.Warn(ctx, "Missing authorization header")
			err := response.NewUnauthorizedError("Missing authorization header")
			response.Handle(c, nil, err)
			c.Abort()
			return
		}

		var token string
		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, TokenPrefix) {
			logger.Warn(ctx, "Invalid authorization header format")
			err := response.NewUnauthorizedError("Invalid authorization header format")
			response.Handle(c, nil, err)
			c.Abort()
			return
		}

		// 提取token
		token = strings.TrimPrefix(authHeader, TokenPrefix)
		if token == "" {
			logger.Warn(ctx, "Empty token")
			err := response.NewUnauthorizedError("Empty token")
			response.Handle(c, nil, err)
			c.Abort()
			return
		}

		// 验证token并获取用户信息
		userID, err := validateToken(ctx, token, jwtService)
		if err != nil {
			logger.Error(ctx, "Token validation failed", zap.Error(err))
			authErr := response.NewUnauthorizedError("Invalid token")
			response.Handle(c, nil, authErr)
			c.Abort()
			return
		}

		// 将用户ID存储到context中
		c.Set(contextutil.UserIDKey, userID)
		ctx = context.WithValue(ctx, contextutil.UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)

		logger.Debug(ctx, "Authentication successful",
			zap.String("user_id", userID))

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

// validateToken 验证token并返回用户ID
func validateToken(ctx context.Context, tokenString string, jwtService *jwt.JWT) (string, error) {
	// 解析JWT token
	claims, err := jwtService.ParseToken(tokenString)
	if err != nil {
		logger.Debug(ctx, "Token parsing failed", zap.Error(err), zap.String("token", tokenString[:10]+"..."))
		return "", err
	}

	// 返回用户ID
	return claims.UserID, nil
}
