package interfaces

import (
	"github.com/gin-gonic/gin"
)

// MiddlewareProvider 中间件提供者接口
// 定义了创建各种中间件的统一接口
type MiddlewareProvider interface {
	// CreateAuthMiddleware 创建认证中间件
	CreateAuthMiddleware() gin.HandlerFunc
	
	// CreateCORSMiddleware 创建CORS中间件
	CreateCORSMiddleware() gin.HandlerFunc
	
	// CreateRateLimitMiddleware 创建限流中间件
	CreateRateLimitMiddleware() gin.HandlerFunc
	
	// CreateRequestLogMiddleware 创建请求日志中间件
	CreateRequestLogMiddleware() gin.HandlerFunc
	
	// CreateRecoveryMiddleware 创建恢复中间件
	CreateRecoveryMiddleware() gin.HandlerFunc
	
	// GetBuilder 获取中间件构建器
	GetBuilder() MiddlewareBuilder
}

// MiddlewareBuilder 中间件构建器接口
// 提供更灵活的中间件构建能力
type MiddlewareBuilder interface {
	// BuildStandardChain 构建标准中间件链
	BuildStandardChain() []gin.HandlerFunc
	
	// BuildSecureChain 构建安全中间件链（包含认证）
	BuildSecureChain() []gin.HandlerFunc
	
	// BuildInternalChain 构建内网访问中间件链
	BuildInternalChain() []gin.HandlerFunc
	
	// BuildCustomChain 根据配置构建自定义中间件链
	BuildCustomChain(config interface{}) []gin.HandlerFunc
}

// JWTService JWT服务接口
// 定义了JWT令牌的创建和验证接口
type JWTService interface {
	// GenerateToken 生成JWT令牌
	GenerateToken(userID string) (string, error)
	
	// ValidateToken 验证JWT令牌
	ValidateToken(token string) (string, error)
	
	// RefreshToken 刷新JWT令牌
	RefreshToken(token string) (string, error)
	
	// ParseToken 解析JWT令牌获取Claims
	ParseToken(token string) (*TokenClaims, error)
}

// TokenClaims JWT令牌声明
type TokenClaims struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

// RateLimiter 限流器接口
type RateLimiter interface {
	// Allow 检查是否允许请求
	Allow(key string) bool
	
	// Take 获取令牌
	Take(key string, count int64) bool
	
	// Reset 重置限流器
	Reset(key string)
}