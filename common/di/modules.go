package di

import (
	"go.uber.org/fx"

	"common"
	"common/config"
	"common/databases"
	"common/databases/redis"
	"common/interfaces"
	"common/jwt"
	"common/logger"
	"common/middleware"
	"common/pkg/casbin"
	"common/pkg/validation"
)

// ConfigModule 配置模块
var ConfigModule = fx.Module("config",
	fx.Provide(
		fx.Annotate(
			config.NewConfigProvider,
			fx.As(new(interfaces.ConfigProvider)),
		),
	),
)

// LoggerModule 日志模块
var LoggerModule = fx.Module("logger",
	fx.Provide(
		fx.Annotate(
			logger.NewLogger,
			fx.As(new(interfaces.Logger)),
		),
		fx.Annotate(
			logger.NewLoggerFactory,
			fx.As(new(interfaces.LoggerFactory)),
		),
	),
)

// DatabasesModule 数据库模块
var DatabasesModule = fx.Module("databases",
	databases.Module,
	redis.RedisModule,
)

// JWTModule JWT模块
var JWTModule = fx.Module("jwt",
	jwt.Module,
)

// MiddlewareModule 中间件模块
var MiddlewareModule = fx.Module("middleware",
	middleware.Module,
)

// ValidationModule 验证模块
var ValidationModule = fx.Module("validation",
	validation.Module,
)

// IDGenModule ID生成器模块
// var IDGenModule = fx.Module("idgen",
//     idgen.Module,
// )

// HTTPModule 基础gin.Engine
// var HTTPModule = fx.Module("http",
//     http.Module,
// )

// TimezoneModule 时区模块
// var TimezoneModule = fx.Module("timezone",
//     timezone.Module,
// )

// CasbinModule Casbin权限模块
var CasbinModule = fx.Module("casbin",
	casbin.Module,
)

// CommonServicesModule Common服务聚合模块
var CommonServicesModule = fx.Module("common_services",
	fx.Provide(
		fx.Annotate(
			common.NewCommonServices,
			fx.As(new(interfaces.CommonServices)),
		),
	),
)

// GetCoreModules 获取核心模块，用于CLI和其他应用
func GetCoreModules() fx.Option {
	return fx.Options(
		// 基础设施模块 - 按依赖顺序
		ConfigModule,
		LoggerModule,
		DatabasesModule,
		JWTModule,
		MiddlewareModule,
		ValidationModule,
		CasbinModule,
		// Common服务聚合模块
		CommonServicesModule,
	)
}

// GetWebModules 获取Web应用模块
func GetWebModules() fx.Option {
	return fx.Options(
		GetCoreModules(),
		// 暂时禁用 HTTP 模块，等后续任务重构
		// HTTPModule,
	)
}
