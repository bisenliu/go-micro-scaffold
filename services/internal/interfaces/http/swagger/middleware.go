package swagger

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/config"
	"common/pkg/jwt"
	"common/response"
)

// SwaggerMiddleware Swagger访问控制中间件
type SwaggerMiddleware struct {
	config     *config.Config
	jwtService *jwt.JWT
	logger     *zap.Logger
}

// NewSwaggerMiddleware 创建Swagger中间件
func NewSwaggerMiddleware(cfg *config.Config, jwtService *jwt.JWT, logger *zap.Logger) *SwaggerMiddleware {
	return &SwaggerMiddleware{
		config:     cfg,
		jwtService: jwtService,
		logger:     logger,
	}
}

// AccessControlMiddleware 访问控制中间件
func (sm *SwaggerMiddleware) AccessControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查环境和配置
		if !sm.shouldAllowAccess(c) {
			sm.logger.Warn("Swagger access denied",
				zap.String("path", c.Request.URL.Path),
				zap.String("environment", sm.config.System.Env),
				zap.String("client_ip", c.ClientIP()))

			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Swagger documentation is not available in this environment",
			})
			c.Abort()
			return
		}

		// 在生产环境中可能需要额外的认证
		if sm.config.System.Env == "production" && sm.requiresAuthentication() {
			if !sm.validateAccess(c) {
				sm.logger.Warn("Swagger authentication failed",
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()))

				err := response.NewUnauthorizedError("Authentication required for Swagger access")
				response.Handle(c, nil, err)
				c.Abort()
				return
			}
		}

		// 添加安全头
		sm.addSecurityHeaders(c)

		c.Next()
	}
}

// shouldAllowAccess 检查是否应该允许访问
func (sm *SwaggerMiddleware) shouldAllowAccess(c *gin.Context) bool {
	// 检查Swagger是否启用
	if !sm.config.Swagger.Enabled {
		return false
	}

	// 检查环境配置
	env := strings.ToLower(sm.config.System.Env)

	switch env {
	case "production":
		// 生产环境需要明确启用
		return sm.isExplicitlyEnabled()
	case "development", "dev":
		// 开发环境默认允许
		return true
	case "test", "testing":
		// 测试环境默认允许
		return true
	default:
		// 未知环境，谨慎处理
		return false
	}
}

// isExplicitlyEnabled 检查是否明确启用
func (sm *SwaggerMiddleware) isExplicitlyEnabled() bool {
	// 可以通过环境变量或配置文件明确启用
	return sm.config.Swagger.Enabled
}

// requiresAuthentication 检查是否需要认证
func (sm *SwaggerMiddleware) requiresAuthentication() bool {
	// 生产环境默认需要认证
	// 可以通过配置或环境变量控制
	return sm.config.System.Env == "production"
}

// validateAccess 验证访问权限
func (sm *SwaggerMiddleware) validateAccess(c *gin.Context) bool {
	// 获取Authorization头
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false
	}

	// 检查Bearer前缀
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	// 提取token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return false
	}

	// 验证JWT token
	_, err := sm.jwtService.ParseToken(token)
	if err != nil {
		sm.logger.Debug("Swagger JWT validation failed", zap.Error(err))
		return false
	}

	return true
}

// addSecurityHeaders 添加安全头
func (sm *SwaggerMiddleware) addSecurityHeaders(c *gin.Context) {
	// 防止在iframe中加载
	c.Header("X-Frame-Options", "DENY")

	// 防止MIME类型嗅探
	c.Header("X-Content-Type-Options", "nosniff")

	// XSS保护
	c.Header("X-XSS-Protection", "1; mode=block")

	// 内容安全策略
	c.Header("Content-Security-Policy", "default-src 'self' 'unsafe-inline' 'unsafe-eval'; img-src 'self' data:; font-src 'self' data:")

	// 引用策略
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
}

// IPWhitelistMiddleware IP白名单中间件（可选功能）
func (sm *SwaggerMiddleware) IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(allowedIPs) == 0 {
			// 如果没有配置白名单，则跳过检查
			c.Next()
			return
		}

		clientIP := c.ClientIP()

		// 检查IP是否在白名单中
		for _, allowedIP := range allowedIPs {
			if clientIP == allowedIP || allowedIP == "*" {
				c.Next()
				return
			}
		}

		sm.logger.Warn("Swagger access denied - IP not in whitelist",
			zap.String("client_ip", clientIP),
			zap.Strings("allowed_ips", allowedIPs))

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "Access denied - IP not in whitelist",
		})
		c.Abort()
	}
}

// RateLimitMiddleware 简单的速率限制中间件（可选功能）
func (sm *SwaggerMiddleware) RateLimitMiddleware() gin.HandlerFunc {
	// 这里可以实现基于IP的简单速率限制
	// 或者集成现有的rate limit中间件
	return func(c *gin.Context) {
		// 目前直接通过，后续可以根据需要实现具体的限制逻辑
		c.Next()
	}
}

// SecurityConfig Swagger安全配置
type SecurityConfig struct {
	RequireAuth     bool     `json:"require_auth"`      // 是否需要认证
	AllowedIPs      []string `json:"allowed_ips"`       // IP白名单
	EnableRateLimit bool     `json:"enable_rate_limit"` // 是否启用速率限制
}

// GetSecurityConfig 获取安全配置
func (sm *SwaggerMiddleware) GetSecurityConfig() SecurityConfig {
	return SecurityConfig{
		RequireAuth:     sm.requiresAuthentication(),
		AllowedIPs:      []string{}, // 可以从配置文件读取
		EnableRateLimit: false,      // 可以从配置文件读取
	}
}
