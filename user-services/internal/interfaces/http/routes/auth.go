package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"user-services/internal/interfaces/http/handler"
)

// SetupAuthRoutes 设置认证API路由
func SetupAuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware AuthMiddleware, logger *zap.Logger) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login/password", authHandler.LoginByPassword)
		auth.POST("/login/wechat", authHandler.LoginByWeChat)

		// 登出需要认证
		auth.POST("/logout", gin.HandlerFunc(authMiddleware), authHandler.Logout)
	}

	logger.Info("Auth API routes registered")
}
