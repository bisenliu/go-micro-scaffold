package swagger

import (
	"go.uber.org/fx"
)

// Module Swagger模块的FX依赖注入配置
var Module = fx.Module("swagger",
	// 提供Swagger相关的服务
	fx.Provide(
		NewSwaggerManager,
		NewSwaggerRoutes,
		NewSwaggerMiddleware,
		NewErrorResponseConverter,
		NewSwaggerErrorHelper,
		NewSwaggerResponseHelper,
		NewResponseAdapter,
	),
)

// SwaggerModule 完整的Swagger模块配置
var SwaggerModule = fx.Options(
	Module,
	// 可以在这里添加其他相关的模块依赖
)
