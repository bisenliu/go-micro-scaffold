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
		if !sm.shouldAllowAccess() {
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

		// 在生产环境中需要额外的认证
		if sm.config.System.Env == "production" {
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
func (sm *SwaggerMiddleware) shouldAllowAccess() bool {
	// 检查Swagger是否启用
	if !sm.config.Swagger.Enabled {
		return false
	}

	// 检查环境配置
	env := strings.ToLower(sm.config.System.Env)

	switch env {
	case "production":
		// 生产环境需要明确启用
		return sm.config.Swagger.Enabled
	case "development", "dev", "test", "testing":
		// 开发和测试环境默认允许
		return true
	default:
		// 未知环境，谨慎处理
		return false
	}
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
