package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"services/internal/interfaces/http/handler"
)

	// SetupAuthRoutes 设置认证API路由
func SetupAuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, logger *zap.Logger) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login/password", authHandler.LoginByPassword)
		auth.POST("/login/wechat", authHandler.LoginByWeChat)
	}

	logger.Info("Auth API routes registered", zap.Int("count", 2))
}
