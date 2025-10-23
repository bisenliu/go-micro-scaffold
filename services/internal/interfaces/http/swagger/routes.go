package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"common/config"
	_ "services/docs" // 导入生成的Swagger文档
)

// SwaggerRoutes Swagger路由管理器
type SwaggerRoutes struct {
	manager *SwaggerManager
	logger  *zap.Logger
}

// NewSwaggerRoutes 创建Swagger路由管理器
func NewSwaggerRoutes(cfg *config.Config, logger *zap.Logger) *SwaggerRoutes {
	return &SwaggerRoutes{
		manager: NewSwaggerManager(cfg),
		logger:  logger,
	}
}

// SetupSwaggerRoutes 设置Swagger路由
func SetupSwaggerRoutes(engine *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	routes := NewSwaggerRoutes(cfg, logger)

	// 检查是否应该启用Swagger
	if !routes.shouldEnableSwagger(cfg.System.Env) {
		logger.Info("Swagger is disabled for current environment",
			zap.String("environment", cfg.System.Env))
		return
	}

	routes.setupRoutes(engine, cfg, logger)
	logger.Info("Swagger routes setup completed",
		zap.String("url", "/swagger/index.html"))
}

// shouldEnableSwagger 判断是否应该启用Swagger
func (sr *SwaggerRoutes) shouldEnableSwagger(env string) bool {
	return sr.manager.ShouldEnableInEnvironment(env)
}

// setupRoutes 设置具体的路由
func (sr *SwaggerRoutes) setupRoutes(engine *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	swaggerConfig := sr.manager.GetConfig()

	// 创建Swagger中间件（需要JWT服务，但这里我们简化处理）
	// 在实际使用中，如果需要JWT验证，应该通过DI注入
	middleware := NewSwaggerMiddleware(cfg, nil, logger)

	// 配置Swagger中间件，使用生成的文档
	url := ginSwagger.URL("/swagger/doc.json") // 指向生成的API定义

	// 创建Swagger路由组，应用访问控制中间件
	swaggerGroup := engine.Group("/swagger")
	swaggerGroup.Use(middleware.AccessControlMiddleware())

	// 添加重定向路由，方便访问（也需要中间件保护）
	engine.GET("/docs", middleware.AccessControlMiddleware(), sr.redirectToSwagger)
	engine.GET("/api-docs", middleware.AccessControlMiddleware(), sr.redirectToSwagger)

	// 设置Swagger路由 - 使用生成的文档
	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 添加Swagger健康检查路由（独立路由，避免与wildcard冲突）
	engine.GET("/swagger-health", middleware.AccessControlMiddleware(), sr.swaggerHealthCheck)

	sr.logger.Info("Swagger routes configured with access control",
		zap.String("title", swaggerConfig.Title),
		zap.String("version", swaggerConfig.Version),
		zap.String("base_path", swaggerConfig.BasePath))
}

// swaggerHealthCheck Swagger健康检查
func (sr *SwaggerRoutes) swaggerHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "swagger",
		"message": "Swagger documentation service is running",
	})
}

// redirectToSwagger 重定向到Swagger UI
func (sr *SwaggerRoutes) redirectToSwagger(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}

// GetSwaggerManager 获取Swagger管理器（用于其他模块）
func (sr *SwaggerRoutes) GetSwaggerManager() *SwaggerManager {
	return sr.manager
}
