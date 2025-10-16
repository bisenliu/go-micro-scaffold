package routes

import (
	"context"

	"github.com/gin-gonic/gin"

	"common/interfaces"
	"common/response"
	"services/internal/interfaces/http/handler"
)

// SetupSystemRoutes 设置系统路由
func SetupSystemRoutes(engine *gin.Engine, healthHandler *handler.HealthHandler, logger interfaces.Logger) {
	engine.GET("/health", healthHandler.Health)
	engine.GET("/ping", func(c *gin.Context) {
		response.OK(c, gin.H{"message": "pong"})
	})

	ctx := context.Background()
	logger.Info(ctx, "System routes registered")
}
