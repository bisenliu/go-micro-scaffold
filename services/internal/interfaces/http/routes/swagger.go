package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"common/config"
	_ "services/docs" // 导入生成的Swagger文档
	"services/internal/interfaces/http/swagger"
)

// SetupSwaggerRoutes 设置Swagger相关路由
func SetupSwaggerRoutes(engine *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// 检查Swagger是否应该启用
	swaggerManager := swagger.NewSwaggerManager(cfg)

	if !swaggerManager.ShouldEnableInEnvironment(cfg.System.Env) {
		logger.Info("Swagger is disabled for current environment",
			zap.String("environment", cfg.System.Env),
			zap.Bool("config_enabled", cfg.Swagger.Enabled))
		return
	}

	// 设置Swagger路由 - 简单直接的配置
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 添加便捷的重定向路由
	engine.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	engine.GET("/api-docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// 添加Swagger健康检查
	engine.GET("/swagger-health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "swagger",
			"message": "Swagger documentation service is running",
		})
	})

	logger.Info("Swagger routes enabled and configured",
		zap.String("environment", cfg.System.Env),
		zap.String("swagger_url", "/swagger/index.html"),
		zap.String("api_docs_url", "/swagger/doc.json"))
}
