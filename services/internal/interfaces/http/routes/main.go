package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	commonMiddleware "common/middleware"
	"services/internal/interfaces/http/handler"
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
	ZapLogger        *zap.Logger
}

// SetupRoutesFinal
func SetupRoutesFinal(p RoutesParams) {

	// 1. 系统路由（无需认证）
	SetupSystemRoutes(p.Engine, p.HealthHandler, p.ZapLogger)

	// 2. API v1 路由组
	v1 := p.Engine.Group("/api/v1")

	// 2.1 认证相关路由（部分需要Token）
	SetupAuthRoutes(v1, p.AuthHandler, p.AuthMiddleware, p.ZapLogger)

	// 2.2 业务路由（需要认证和授权）
	v1.Use(commonMiddleware.RequestLogMiddleware())
	// v1.Use(gin.HandlerFunc(p.CasbinMiddleware))
	{
		SetupUserRoutes(v1, p.UserHandler, p.ZapLogger)
		// 后续添加其他模块
	}

	p.ZapLogger.Info("All routes setup completed successfully")
}
