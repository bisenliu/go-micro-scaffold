package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"common/interfaces"
)

// MiddlewareBuilder 中间件构建器
type MiddlewareBuilder struct {
	jwtService interfaces.JWTService
	logger     interfaces.Logger
	config     interfaces.ConfigProvider
}

// NewMiddlewareBuilder 创建中间件构建器
func NewMiddlewareBuilder(
	jwtService interfaces.JWTService,
	logger interfaces.Logger,
	config interfaces.ConfigProvider,
) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		jwtService: jwtService,
		logger:     logger,
		config:     config,
	}
}

// BuildAuthMiddleware 构建认证中间件
func (b *MiddlewareBuilder) BuildAuthMiddleware() gin.HandlerFunc {
	authConfig := b.config.GetAuthConfig()
	return AuthMiddleware(b.jwtService, authConfig, b.logger)
}

// BuildCORSMiddleware 构建CORS中间件
func (b *MiddlewareBuilder) BuildCORSMiddleware() gin.HandlerFunc {
	serverConfig := b.config.GetServerConfig()
	
	if !serverConfig.EnableCORS {
		// 如果禁用CORS，返回空中间件
		return func(c *gin.Context) {
			c.Next()
		}
	}

	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(config)
}

// BuildRateLimitMiddleware 构建限流中间件
func (b *MiddlewareBuilder) BuildRateLimitMiddleware() gin.HandlerFunc {
	rateLimitConfig := b.config.GetRateLimitConfig()
	return RateLimitMiddleware(rateLimitConfig, b.logger.GetZapLogger())
}

// BuildRecoveryMiddleware 构建恢复中间件
func (b *MiddlewareBuilder) BuildRecoveryMiddleware() gin.HandlerFunc {
	return RecoveryMiddleware()
}

// BuildRequestLogMiddleware 构建请求日志中间件
func (b *MiddlewareBuilder) BuildRequestLogMiddleware() gin.HandlerFunc {
	return RequestLogMiddleware()
}

// BuildClientIPMiddleware 构建客户端IP提取中间件
func (b *MiddlewareBuilder) BuildClientIPMiddleware() gin.HandlerFunc {
	return ExtractClientIPMiddleware()
}

// BuildIPWhitelistMiddleware 构建IP白名单中间件
func (b *MiddlewareBuilder) BuildIPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return IPWhitelistMiddleware(allowedIPs)
}

// BuildInternalIPMiddleware 构建内网IP限制中间件
func (b *MiddlewareBuilder) BuildInternalIPMiddleware() gin.HandlerFunc {
	return InternalIPMiddleware()
}

// BuildStandardChain 构建标准中间件链
func (b *MiddlewareBuilder) BuildStandardChain() *MiddlewareChain {
	chain := NewMiddlewareChain()
	
	// 1. 恢复中间件 - 最先执行，捕获panic
	chain.Add("recovery", b.BuildRecoveryMiddleware())
	
	// 2. 客户端IP提取 - 为后续中间件提供IP信息
	chain.Add("client_ip", b.BuildClientIPMiddleware())
	
	// 3. CORS中间件 - 处理跨域请求
	chain.Add("cors", b.BuildCORSMiddleware())
	
	// 4. 请求日志中间件 - 记录请求信息
	chain.Add("request_log", b.BuildRequestLogMiddleware())
	
	// 5. 限流中间件 - 控制请求频率
	chain.Add("rate_limit", b.BuildRateLimitMiddleware())
	
	return chain
}

// BuildSecureChain 构建安全中间件链（包含认证）
func (b *MiddlewareBuilder) BuildSecureChain() *MiddlewareChain {
	chain := b.BuildStandardChain()
	
	// 添加认证中间件
	chain.Add("auth", b.BuildAuthMiddleware())
	
	return chain
}

// BuildInternalChain 构建内网访问中间件链
func (b *MiddlewareBuilder) BuildInternalChain() *MiddlewareChain {
	chain := b.BuildStandardChain()
	
	// 添加内网IP限制
	chain.Add("internal_ip", b.BuildInternalIPMiddleware())
	
	return chain
}

// BuildCustomChain 根据配置构建自定义中间件链
func (b *MiddlewareBuilder) BuildCustomChain(config *MiddlewareConfig) *MiddlewareChain {
	chain := NewMiddlewareChain()
	
	// 恢复中间件总是第一个
	if config.Recovery == nil || config.Recovery.Enabled {
		chain.Add("recovery", b.BuildRecoveryMiddleware())
	}
	
	// 客户端IP提取
	chain.Add("client_ip", b.BuildClientIPMiddleware())
	
	// CORS中间件
	if config.CORS != nil && config.CORS.Enabled {
		chain.Add("cors", b.BuildCORSMiddleware())
	}
	
	// 请求日志中间件
	if config.RequestLog == nil || config.RequestLog.Enabled {
		chain.Add("request_log", b.BuildRequestLogMiddleware())
	}
	
	// IP访问控制
	if config.IPAccess != nil && config.IPAccess.Enabled {
		if config.IPAccess.InternalOnly {
			chain.Add("internal_ip", b.BuildInternalIPMiddleware())
		}
		if len(config.IPAccess.WhitelistIPs) > 0 {
			chain.Add("ip_whitelist", b.BuildIPWhitelistMiddleware(config.IPAccess.WhitelistIPs))
		}
	}
	
	// 限流中间件
	if config.RateLimit != nil && config.RateLimit.Enabled {
		chain.Add("rate_limit", b.BuildRateLimitMiddleware())
	}
	
	// 认证中间件
	if config.Auth != nil && config.Auth.Enabled {
		chain.Add("auth", b.BuildAuthMiddleware())
	}
	
	return chain
}