package routes

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"common/config"
	_ "user-services/docs" // 导入生成的Swagger文档
)

// SetupSwaggerRoutes 设置Swagger相关路由
func SetupSwaggerRoutes(engine *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// 检查Swagger是否应该启用
	if !isSwaggerEnabled(cfg.System.Env) {
		logger.Info("Swagger is disabled for current environment",
			zap.String("environment", cfg.System.Env),
			zap.Bool("enabled", false))
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

// isSwaggerEnabled 检查Swagger是否应该启用
func isSwaggerEnabled(env string) bool {
	// 首先检查环境变量 SWAGGER_ENABLED
	if enabledEnv := os.Getenv("SWAGGER_ENABLED"); enabledEnv != "" {
		enabled, err := strconv.ParseBool(enabledEnv)
		if err == nil {
			return enabled
		}
	}

	// 根据环境自动判断
	envLower := strings.ToLower(env)
	switch envLower {
	case "production", "prod":
		// 生产环境默认禁用
		return false
	case "development", "dev", "local":
		// 开发环境默认启用
		return true
	case "testing", "test", "staging":
		// 测试环境默认启用
		return true
	default:
		// 未知环境默认禁用
		return false
	}
}
