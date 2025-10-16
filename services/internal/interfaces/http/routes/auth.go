package routes

import (
	"context"

	"github.com/gin-gonic/gin"

	"common/interfaces"
	"services/internal/interfaces/http/handler"
)

// SetupAuthRoutes 设置认证API路由
func SetupAuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware AuthMiddleware, logger interfaces.Logger) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login/password", authHandler.LoginByPassword)
		auth.POST("/login/wechat", authHandler.LoginByWeChat)

		// 登出需要认证
		auth.POST("/logout", gin.HandlerFunc(authMiddleware), authHandler.Logout)
	}

	ctx := context.Background()
	logger.Info(ctx, "Auth API routes registered")
}
