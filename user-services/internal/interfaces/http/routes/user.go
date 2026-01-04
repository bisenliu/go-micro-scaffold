package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"user-services/internal/interfaces/http/handler"
)

// SetupUserRoutes 设置用户API路由
func SetupUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, logger *zap.Logger) {
	users := rg.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
	}

	logger.Info("User API routes registered")
}
