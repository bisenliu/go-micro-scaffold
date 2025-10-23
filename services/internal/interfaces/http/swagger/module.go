package swagger

import (
	"common/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module Swagger模块的FX依赖注入配置
var Module = fx.Module("swagger",
	// 提供简化后的Swagger核心服务
	fx.Provide(
		NewSwaggerManager,    // 配置管理器 - 简化的配置管理
		NewSwaggerRoutes,     // 路由管理器 - 集成访问控制中间件
		NewSwaggerMiddleware, // 访问控制中间件 - 保留核心安全功能
	),
)

// SwaggerModule 完整的Swagger模块配置
// 提供完整的Swagger功能集成
var SwaggerModule = fx.Options(
	Module,
	// Swagger模块是自包含的，不需要额外的模块依赖
)

// SwaggerIntegrationParams 定义Swagger集成所需的参数
type SwaggerIntegrationParams struct {
	fx.In

	Engine *gin.Engine `optional:"true"`
	Config *config.Config
	Logger *zap.Logger
}

// EnableSwaggerIntegration 启用Swagger集成的便捷函数
// 可以在其他模块中调用此函数来集成Swagger功能
func EnableSwaggerIntegration(params SwaggerIntegrationParams) {
	if params.Engine != nil {
		SetupSwaggerRoutes(params.Engine, params.Config, params.Logger)
	}
}
