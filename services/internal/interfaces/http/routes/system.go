package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/response"
	"services/internal/interfaces/http/handler"
)

// SetupSystemRoutes 设置系统路由
func SetupSystemRoutes(engine *gin.Engine, healthHandler *handler.HealthHandler, logger *zap.Logger) {
	engine.GET("/health", healthHandler.Health)
	engine.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"message": "pong"})
	})

	logger.Info("System routes registered")
}
