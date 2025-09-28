package di

import (
	"go.uber.org/fx"

	"common/config"
	"common/databases"
	"common/http"
	"common/logger"
	"common/pkg/idgen"
	"common/pkg/jwt"
	"common/pkg/timezone"
	"common/pkg/validation"
)

// ConfigModule 配置模块
var ConfigModule = fx.Module("config",
	config.Module,
)

// LoggerModule 日志模块
var LoggerModule = fx.Module("logger",
	logger.Module,
)

// DatabasesModule 数据库模块
var DatabasesModule = fx.Module("databases",
	databases.Module,
)

// ValidationModule 验证模块
var ValidationModule = fx.Module("validation",
	validation.Module,
)

// IDGenModule ID生成器模块
var IDGenModule = fx.Module("idgen",
	idgen.Module,
)

// JWTModule JWT模块
var JWTModule = fx.Module("jwt",
	jwt.Module,
)

// HTTPModule HTTP模块
var HTTPModule = fx.Module("http",
	http.Module,
)

// TimezoneModule 时区模块
var TimezoneModule = fx.Module("timezone",
	timezone.Module,
)

// GetCoreModules 获取核心模块，用于CLI和其他应用
func GetCoreModules() fx.Option {
	return fx.Options(
		// 基础设施模块
		ConfigModule,
		LoggerModule,
		DatabasesModule,
		ValidationModule,
		IDGenModule,
		JWTModule,
		TimezoneModule,
	)
}

// GetWebModules 获取Web应用模块
func GetWebModules() fx.Option {
	return fx.Options(
		GetCoreModules(),
		HTTPModule,
	)
}
