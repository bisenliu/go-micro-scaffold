package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/config"
	"common/middleware"
)

// EngineParams 定义了创建Gin引擎所需的依赖
type EngineParams struct {
	fx.In
	Config *config.Config
	Logger *zap.Logger
}

// NewBaseEngine 创建一个带有通用配置和中间件的Gin引擎
func NewBaseEngine(params EngineParams) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(params.Config.Server.Mode)

	// 创建Gin引擎
	engine := gin.New()

	// 添加CORS 中间件
	if params.Config.Server.EnableCORS {
		engine.Use(middleware.CORSMiddleware())
	}

	engine.Use(middleware.ExtractClientIPMiddleware(params.Logger))
	engine.Use(middleware.RecoveryMiddleware(params.Logger))
	engine.Use(middleware.TraceLoggerMiddleware(params.Logger))

	return engine
}

// Module 提供了一个配置好的 *gin.Engine
var Module = fx.Module("http-engine",
	fx.Provide(NewBaseEngine),
)
