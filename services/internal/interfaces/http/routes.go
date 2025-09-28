package http

import (
	"common/middleware"
	"common/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"services/internal/interfaces/http/handler"
)

// SetupRoutesFinal
func SetupRoutesFinal(
	engine *gin.Engine,
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	zapLogger *zap.Logger,
) {
	engine.Use(middleware.LoggerMiddleware(zapLogger))

	// 1. 系统路由（无需认证）
	setupSystemRoutesFinal(engine, healthHandler, zapLogger)

	// 2. API v1 路由组（需要认证）
	v1 := engine.Group("/api/v1")
	{
		// 这里可以添加认证中间件
		// v1.Use(middlewareManager.AuthMiddleware())

		setupUserAPIRoutes(v1, userHandler, zapLogger)
		// 后续添加其他模块
	}
	
	zapLogger.Info("All routes setup completed successfully")
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
	}

	logger.Info("User API routes registered", zap.Int("count", 5))
}

