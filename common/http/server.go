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

	// 核心中间件：顺序非常重要
	// 1. TraceLogger: 必须在最前面，为所有请求注入 traceID
	engine.Use(middleware.TraceLoggerMiddleware(params.Logger))
	// 2. Recovery: 必须在 Logger 之后，以便在 panic 时可以记录带有 traceID 的日志
	engine.Use(middleware.RecoveryMiddleware())

	// 3. ExtractClientIP: 尽早获取客户端IP，为后续的日志、限流、访问控制提供支持
	engine.Use(middleware.ExtractClientIPMiddleware())

	// 4. CORS: 处理跨域请求(通过配置启用/禁用)，对于 OPTIONS 预检请求会直接中断后续中间件
	engine.Use(middleware.CORSMiddleware(params.Config.Server))

	// 5. RateLimit: 基于IP进行限流(通过配置启用/禁用)，保护后端服务
	engine.Use(middleware.RateLimitMiddleware(params.Config.RateLimit, params.Logger))

	return engine
}

// Module 提供了一个配置好的 *gin.Engine
var Module = fx.Module("http-engine",
	fx.Provide(NewBaseEngine),
)
