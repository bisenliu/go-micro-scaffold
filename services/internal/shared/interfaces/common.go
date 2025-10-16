// Package interfaces 重新导出 common 接口，避免 services 模块直接依赖 common 包
// 这样实现了依赖倒置，services 只依赖接口，不依赖具体实现
package interfaces

import "common/interfaces"

// 重新导出核心接口，避免直接依赖 common 包
type (
	// CommonServices 聚合所有 common 服务的接口
	CommonServices = interfaces.CommonServices
	
	// ConfigProvider 配置提供者接口
	ConfigProvider = interfaces.ConfigProvider
	
	// Logger 日志接口
	Logger = interfaces.Logger
	
	// DatabaseManager 数据库管理器接口
	DatabaseManager = interfaces.DatabaseManager
	
	// Database 数据库接口
	Database = interfaces.Database
	
	// CacheManager 缓存管理器接口
	CacheManager = interfaces.CacheManager
	
	// Validator 验证器接口
	Validator = interfaces.Validator
	
	// MiddlewareProvider 中间件提供者接口
	MiddlewareProvider = interfaces.MiddlewareProvider
	
	// JWTService JWT服务接口
	JWTService = interfaces.JWTService
)

// 重新导出配置类型
type (
	ServerConfig     = interfaces.ServerConfig
	DatabaseConfig   = interfaces.DatabaseConfig
	AuthConfig       = interfaces.AuthConfig
	LoggerConfig     = interfaces.LoggerConfig
	RedisConfig      = interfaces.RedisConfig
	TokenConfig      = interfaces.TokenConfig
	ValidationConfig = interfaces.ValidationConfig
)

// 重新导出其他类型
type (
	DatabaseConnection = interfaces.DatabaseConnection
	DatabaseStats      = interfaces.DatabaseStats
	HealthStatus       = interfaces.HealthStatus
	ValidationError    = interfaces.ValidationError
	ValidationErrors   = interfaces.ValidationErrors
	TokenClaims        = interfaces.TokenClaims
)