package di

import (
	"go.uber.org/fx"

	"common/cache"
	"common/config"
	"common/database"
	"common/http"
	"common/logger"
	"common/validation"
)

// ConfigModule 配置模块
var ConfigModule = fx.Module("config",
	config.Module,
)

// LoggerModule 日志模块
var LoggerModule = fx.Module("logger",
	logger.Module,
)

// DatabaseModule 数据库模块 (Ent ORM)
var DatabaseModule = fx.Module("database",
	// 保持向后兼容的单数据库支持
	database.EntModule,
	database.EntServiceModule,
	// 新的多数据库支持
	database.DatabaseManagerModule,
)

// CacheModule 缓存模块
var CacheModule = fx.Module("cache",
	cache.RedisModule,
)

// ValidationModule 验证模块
var ValidationModule = fx.Module("validation",
	validation.Module,
)

// HTTPModule HTTP模块
var HTTPModule = fx.Module("http",
	http.Module,
)

// GetCoreModules 获取核心模块，用于CLI和其他应用
func GetCoreModules() fx.Option {
	return fx.Options(
		// 基础设施模块
		ConfigModule,
		LoggerModule,
		DatabaseModule,
		CacheModule,
		ValidationModule,
	)
}

// GetWebModules 获取Web应用模块
func GetWebModules() fx.Option {
	return fx.Options(
		GetCoreModules(),
		HTTPModule,
	)
}
