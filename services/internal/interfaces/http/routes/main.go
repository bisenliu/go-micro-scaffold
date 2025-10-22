package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/config"
	commonMiddleware "common/middleware"
	"services/internal/interfaces/http/handler"
	"services/internal/interfaces/http/swagger"
)

// CasbinMiddleware 是一个具名类型，用于DI容器的类型安全注入
type (
	CasbinMiddleware gin.HandlerFunc
	AuthMiddleware   gin.HandlerFunc
)

// RoutesParams 定义了 SetupRoutesFinal 函数的依赖项
type RoutesParams struct {
	fx.In

	Engine           *gin.Engine
	UserHandler      *handler.UserHandler
	HealthHandler    *handler.HealthHandler
	AuthHandler      *handler.AuthHandler
	CasbinMiddleware CasbinMiddleware
	AuthMiddleware   AuthMiddleware
	Config           *config.Config
	ZapLogger        *zap.Logger
}

// SetupRoutesFinal
func SetupRoutesFinal(p RoutesParams) {

	// 1. 系统路由（无需认证）
	SetupSystemRoutes(p.Engine, p.HealthHandler, p.ZapLogger)

	// 2. Swagger API 文档路由（条件性启用）
	SetupSwaggerRoutesConditional(p.Engine, p.Config, p.ZapLogger)

	// 3. API v1 路由组
	v1 := p.Engine.Group("/api/v1")

	// 3.1 认证相关路由（部分需要Token）
	SetupAuthRoutes(v1, p.AuthHandler, p.AuthMiddleware, p.ZapLogger)

	// 3.2 业务路由（需要认证和授权）
	v1.Use(commonMiddleware.RequestLogMiddleware())
	v1.Use(gin.HandlerFunc(p.CasbinMiddleware))
	{
		SetupUserRoutes(v1, p.UserHandler, p.ZapLogger)
		// 后续添加其他模块
	}

	p.ZapLogger.Info("All routes setup completed successfully")
}

// SetupSwaggerRoutesConditional 条件性设置Swagger路由
func SetupSwaggerRoutesConditional(engine *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// 检查Swagger是否应该启用
	swaggerManager := swagger.NewSwaggerManager(cfg)

	if !swaggerManager.ShouldEnableInEnvironment(cfg.System.Env) {
		logger.Info("Swagger is disabled for current environment",
			zap.String("environment", cfg.System.Env),
			zap.Bool("config_enabled", cfg.Swagger.Enabled))
		return
	}

	// 启用Swagger路由
	swagger.SetupSwaggerRoutes(engine, cfg, logger)

	logger.Info("Swagger routes enabled and configured",
		zap.String("environment", cfg.System.Env),
		zap.String("swagger_url", "/swagger/index.html"),
		zap.String("api_docs_url", "/swagger/doc.json"))
}
