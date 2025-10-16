package routes

import (
	"context"

	"github.com/gin-gonic/gin"

	"common/interfaces"
	"services/internal/interfaces/http/handler"
)

// SetupUserRoutes 设置用户API路由
func SetupUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, logger interfaces.Logger) {
	users := rg.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
	}

	ctx := context.Background()
	logger.Info(ctx, "User API routes registered")
}
