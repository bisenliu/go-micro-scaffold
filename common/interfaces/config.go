package interfaces

import (
	"time"
)

// ConfigProvider 配置提供者接口
// 定义了获取各种配置的统一接口，支持配置的动态重载
type ConfigProvider interface {
	// GetDatabaseConfig 获取数据库配置
	GetDatabaseConfig() DatabaseConfig
	
	// GetServerConfig 获取服务器配置
	GetServerConfig() ServerConfig
	
	// GetAuthConfig 获取认证配置
	GetAuthConfig() AuthConfig
	
	// GetLoggerConfig 获取日志配置
	GetLoggerConfig() LoggerConfig
	
	// GetRedisConfig 获取Redis配置
	GetRedisConfig() RedisConfig
	
	// GetTokenConfig 获取Token配置
	GetTokenConfig() TokenConfig
	
	// GetValidationConfig 获取验证配置
	GetValidationConfig() ValidationConfig
	
	// GetRateLimitConfig 获取限流配置
	GetRateLimitConfig() RateLimitConfig
	
	// Reload 重新加载配置
	Reload() error
	
	// GetEnv 获取当前环境
	GetEnv() string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Primary  DatabaseConnection `json:"primary"`
	ReadOnly DatabaseConnection `json:"read_only,omitempty"`
}

// DatabaseConnection 数据库连接配置
type DatabaseConnection struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
	
	// 连接池配置
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port           string        `json:"port"`
	Host           string        `json:"host"`
	Mode           string        `json:"mode"`
	EnableCORS     bool          `json:"enable_cors"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
	IdleTimeout    time.Duration `json:"idle_timeout"`
	MaxHeaderBytes int           `json:"max_header_bytes"`
	EnableMetrics  bool          `json:"enable_metrics"`
	EnableTracing  bool          `json:"enable_tracing"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Whitelist []string `json:"whitelist"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `json:"level"`
	Director   string `json:"director"`
	MaxAge     int    `json:"max_age"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Password  string `json:"password"`
	Database  int    `json:"database"`
	DefaultDB int    `json:"default_db"`
	PoolSize  int    `json:"pool_size"`
}

// TokenConfig Token配置
type TokenConfig struct {
	ExpiredTime int `json:"expired_time"`
}

// ValidationConfig 验证配置
type ValidationConfig struct {
	Locale string `json:"locale"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled         bool          `json:"enabled"`
	FillInterval    time.Duration `json:"fill_interval"`
	Capacity        int64         `json:"capacity"`
	Quantum         int64         `json:"quantum"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
	BucketExpiry    time.Duration `json:"bucket_expiry"`
}