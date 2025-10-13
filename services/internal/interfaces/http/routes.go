package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	commonMiddleware "common/middleware"
	"common/response"
	"services/internal/interfaces/http/handler"
)

// CasbinMiddleware 是一个具名类型，用于DI容器的类型安全注入
type CasbinMiddleware gin.HandlerFunc

// RoutesParams 定义了 SetupRoutesFinal 函数的依赖项
type RoutesParams struct {
	fx.In

	Engine           *gin.Engine
	UserHandler      *handler.UserHandler
	HealthHandler    *handler.HealthHandler
	CasbinMiddleware CasbinMiddleware
	ZapLogger        *zap.Logger
}

// SetupRoutesFinal
func SetupRoutesFinal(p RoutesParams) {
	p.Engine.Use(commonMiddleware.LoggerMiddleware(p.ZapLogger))

	// 1. 系统路由（无需认证）
	setupSystemRoutesFinal(p.Engine, p.HealthHandler, p.ZapLogger)

	// 2. API v1 路由组（需要认证和授权）
	v1 := p.Engine.Group("/api/v1")
	// v1.Use(gin.HandlerFunc(p.CasbinMiddleware))
	{
		// 添加认证中间件
		setupUserAPIRoutes(v1, p.UserHandler, p.ZapLogger)
		// 后续添加其他模块
	}

	p.ZapLogger.Info("All routes setup completed successfully")
}

// setupSystemRoutesFinal 设置系统路由
func setupSystemRoutesFinal(engine *gin.Engine, healthHandler *handler.HealthHandler, logger *zap.Logger) {
	engine.GET("/health", healthHandler.Health)
	engine.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"message": "pong"})
	})

	logger.Info("System routes registered", zap.Int("count", 2))
}

// setupUserAPIRoutes 设置用户API路由
func setupUserAPIRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, logger *zap.Logger) {
	users := rg.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.POST("/login", userHandler.Login)
	}

	logger.Info("User API routes registered", zap.Int("count", 3))
}
