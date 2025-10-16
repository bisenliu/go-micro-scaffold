package middleware

import (
	"github.com/gin-gonic/gin"

	"common/interfaces"
)

// MiddlewareProviderImpl 中间件提供者实现
type MiddlewareProviderImpl struct {
	builder *MiddlewareBuilder
}

// NewMiddlewareProvider 创建中间件提供者
func NewMiddlewareProvider(
	jwtService interfaces.JWTService,
	logger interfaces.Logger,
	config interfaces.ConfigProvider,
) interfaces.MiddlewareProvider {
	builder := NewMiddlewareBuilder(jwtService, logger, config)
	return &MiddlewareProviderImpl{
		builder: builder,
	}
}

// CreateAuthMiddleware 创建认证中间件
func (p *MiddlewareProviderImpl) CreateAuthMiddleware() gin.HandlerFunc {
	return p.builder.BuildAuthMiddleware()
}

// CreateCORSMiddleware 创建CORS中间件
func (p *MiddlewareProviderImpl) CreateCORSMiddleware() gin.HandlerFunc {
	return p.builder.BuildCORSMiddleware()
}

// CreateRateLimitMiddleware 创建限流中间件
func (p *MiddlewareProviderImpl) CreateRateLimitMiddleware() gin.HandlerFunc {
	return p.builder.BuildRateLimitMiddleware()
}

// CreateRequestLogMiddleware 创建请求日志中间件
func (p *MiddlewareProviderImpl) CreateRequestLogMiddleware() gin.HandlerFunc {
	return p.builder.BuildRequestLogMiddleware()
}

// CreateRecoveryMiddleware 创建恢复中间件
func (p *MiddlewareProviderImpl) CreateRecoveryMiddleware() gin.HandlerFunc {
	return p.builder.BuildRecoveryMiddleware()
}

// GetBuilder 获取中间件构建器
func (p *MiddlewareProviderImpl) GetBuilder() interfaces.MiddlewareBuilder {
	return &MiddlewareBuilderAdapter{builder: p.builder}
}