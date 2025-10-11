package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	commonMiddleware "common/middleware"
	"common/response"
	serviceInterface "services/internal/application/service"
	"services/internal/interfaces/http/handler"
)

// SetupRoutesFinal
func SetupRoutesFinal(
	engine *gin.Engine,
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	permissionService *serviceInterface.PermissionService,
	zapLogger *zap.Logger,
) {
	engine.Use(commonMiddleware.LoggerMiddleware(zapLogger))

	// 1. 系统路由（无需认证）
	setupSystemRoutesFinal(engine, healthHandler, zapLogger)

	// 2. API v1 路由组（需要认证和授权）
	v1 := engine.Group("/api/v1")
	// v1.Use(commonMiddleware.CasbinMiddleware(permissionService.Enforce, zapLogger))
	{
		// 添加认证中间件
		setupUserAPIRoutes(v1, userHandler, zapLogger)
		// 后续添加其他模块
	}

	// 3. 初始化一些示例权限策略
	initExamplePolicies(permissionService, zapLogger)

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
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.POST("/login", userHandler.Login)
	}

	logger.Info("User API routes registered", zap.Int("count", 3))
}

// initExamplePolicies 初始化示例权限策略
func initExamplePolicies(permissionService *serviceInterface.PermissionService, logger *zap.Logger) {
	ctx := context.Background()

	// 添加一些示例策略
	// 1. 管理员可以访问所有用户接口
	if err := permissionService.AddPolicy(ctx, "admin", "/api/v1/users", "*"); err != nil {
		logger.Error("Failed to add admin policy", zap.Error(err))
	}

	// 2. 普通用户只能查看用户列表
	if err := permissionService.AddPolicy(ctx, "user", "/api/v1/users", "GET"); err != nil {
		logger.Error("Failed to add user policy", zap.Error(err))
	}

	// 3. 添加角色关系：alice 是管理员
	if err := permissionService.AddRoleForUser(ctx, "alice", "admin"); err != nil {
		logger.Error("Failed to add role for alice", zap.Error(err))
	}

	// 4. 添加角色关系：bob 是普通用户
	if err := permissionService.AddRoleForUser(ctx, "bob", "user"); err != nil {
		logger.Error("Failed to add role for bob", zap.Error(err))
	}

	logger.Info("Example policies initialized")
}
