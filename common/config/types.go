package config

import (
	"time"

	"common/interfaces"
)

// ModularConfig 模块化配置结构
// 将配置按功能模块分离，便于管理和扩展
type ModularConfig struct {
	System     SystemConfig     `mapstructure:"system"`
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Token      TokenConfig      `mapstructure:"token"`
	Validation ValidationConfig `mapstructure:"validation"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	Env        string `mapstructure:"env" validate:"required,oneof=development staging production" default:"development"`
	SecretKey  string `mapstructure:"secret_key"`
	ServerName string `mapstructure:"server_name" default:"go-micro-scaffold"`
	Timezone   string `mapstructure:"timezone" default:"Asia/Shanghai"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port           string        `mapstructure:"port" validate:"required" default:"8080"`
	Host           string        `mapstructure:"host" default:"localhost"`
	Mode           string        `mapstructure:"mode" validate:"oneof=debug release test" default:"debug"`
	EnableCORS     bool          `mapstructure:"enable_cors" default:"true"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout" default:"5s"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout" default:"10s"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout" default:"120s"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes" default:"1048576"` // 1MB
	EnableMetrics  bool          `mapstructure:"enable_metrics" default:"false"`
	EnableTracing  bool          `mapstructure:"enable_tracing" default:"false"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Primary  DatabaseConnection `mapstructure:"primary"`
	ReadOnly DatabaseConnection `mapstructure:"read_only"`
}

// DatabaseConnection 数据库连接配置
type DatabaseConnection struct {
	Type            string        `mapstructure:"type" validate:"required,oneof=mysql postgres sqlite" default:"mysql"`
	Host            string        `mapstructure:"host" default:"localhost"`
	Port            int           `mapstructure:"port" validate:"min=1,max=65535" default:"3306"`
	Database        string        `mapstructure:"database" validate:"required"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Charset         string        `mapstructure:"charset" default:"utf8mb4"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" validate:"min=1" default:"100"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" validate:"min=1" default:"10"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" default:"60m"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" default:"10m"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Whitelist []string `mapstructure:"whitelist"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level" validate:"oneof=debug info warn error" default:"info"`
	Director   string `mapstructure:"director" validate:"required" default:"./logs"`
	MaxAge     int    `mapstructure:"max_age" validate:"min=1" default:"7"`
	MaxSize    int    `mapstructure:"max_size" validate:"min=1" default:"100"`
	MaxBackups int    `mapstructure:"max_backups" validate:"min=1" default:"10"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host      string `mapstructure:"host" default:"localhost"`
	Port      int    `mapstructure:"port" validate:"min=1,max=65535" default:"6379"`
	Password  string `mapstructure:"password"`
	Database  int    `mapstructure:"database" validate:"min=0,max=15" default:"0"`
	DefaultDB int    `mapstructure:"default_db" validate:"min=0,max=15" default:"0"`
	PoolSize  int    `mapstructure:"pool_size" validate:"min=1" default:"10"`
}

// TokenConfig Token配置
type TokenConfig struct {
	ExpiredTime int `mapstructure:"expired_time" validate:"min=1" default:"30"` // 分钟
}

// ValidationConfig 验证配置
type ValidationConfig struct {
	Locale string `mapstructure:"locale" validate:"oneof=zh en" default:"zh"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled         bool          `mapstructure:"enabled" default:"false"`
	FillInterval    time.Duration `mapstructure:"fill_interval" default:"1s"`
	Capacity        int64         `mapstructure:"capacity" default:"100"`
	Quantum         int64         `mapstructure:"quantum" default:"1"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval" default:"5m"`
	BucketExpiry    time.Duration `mapstructure:"bucket_expiry" default:"10m"`
}

// ToInterfaceTypes 转换为接口类型
func (c *ModularConfig) ToInterfaceTypes() ConfigInterfaces {
	return ConfigInterfaces{
		Server: interfaces.ServerConfig{
			Port:           c.Server.Port,
			Host:           c.Server.Host,
			Mode:           c.Server.Mode,
			EnableCORS:     c.Server.EnableCORS,
			ReadTimeout:    c.Server.ReadTimeout,
			WriteTimeout:   c.Server.WriteTimeout,
			IdleTimeout:    c.Server.IdleTimeout,
			MaxHeaderBytes: c.Server.MaxHeaderBytes,
			EnableMetrics:  c.Server.EnableMetrics,
			EnableTracing:  c.Server.EnableTracing,
		},
		Database: interfaces.DatabaseConfig{
			Primary: interfaces.DatabaseConnection{
				Type:            c.Database.Primary.Type,
				Host:            c.Database.Primary.Host,
				Port:            c.Database.Primary.Port,
				Database:        c.Database.Primary.Database,
				Username:        c.Database.Primary.Username,
				Password:        c.Database.Primary.Password,
				Charset:         c.Database.Primary.Charset,
				MaxOpenConns:    c.Database.Primary.MaxOpenConns,
				MaxIdleConns:    c.Database.Primary.MaxIdleConns,
				ConnMaxLifetime: c.Database.Primary.ConnMaxLifetime,
				ConnMaxIdleTime: c.Database.Primary.ConnMaxIdleTime,
			},
			ReadOnly: interfaces.DatabaseConnection{
				Type:            c.Database.ReadOnly.Type,
				Host:            c.Database.ReadOnly.Host,
				Port:            c.Database.ReadOnly.Port,
				Database:        c.Database.ReadOnly.Database,
				Username:        c.Database.ReadOnly.Username,
				Password:        c.Database.ReadOnly.Password,
				Charset:         c.Database.ReadOnly.Charset,
				MaxOpenConns:    c.Database.ReadOnly.MaxOpenConns,
				MaxIdleConns:    c.Database.ReadOnly.MaxIdleConns,
				ConnMaxLifetime: c.Database.ReadOnly.ConnMaxLifetime,
				ConnMaxIdleTime: c.Database.ReadOnly.ConnMaxIdleTime,
			},
		},
		Auth: interfaces.AuthConfig{
			Whitelist: c.Auth.Whitelist,
		},
		Logger: interfaces.LoggerConfig{
			Level:      c.Logger.Level,
			Director:   c.Logger.Director,
			MaxAge:     c.Logger.MaxAge,
			MaxSize:    c.Logger.MaxSize,
			MaxBackups: c.Logger.MaxBackups,
		},
		Redis: interfaces.RedisConfig{
			Host:      c.Redis.Host,
			Port:      c.Redis.Port,
			Password:  c.Redis.Password,
			Database:  c.Redis.Database,
			DefaultDB: c.Redis.DefaultDB,
			PoolSize:  c.Redis.PoolSize,
		},
		Token: interfaces.TokenConfig{
			ExpiredTime: c.Token.ExpiredTime,
		},
		Validation: interfaces.ValidationConfig{
			Locale: c.Validation.Locale,
		},
		RateLimit: interfaces.RateLimitConfig{
			Enabled:         c.RateLimit.Enabled,
			FillInterval:    c.RateLimit.FillInterval,
			Capacity:        c.RateLimit.Capacity,
			Quantum:         c.RateLimit.Quantum,
			CleanupInterval: c.RateLimit.CleanupInterval,
			BucketExpiry:    c.RateLimit.BucketExpiry,
		},
		Env: c.System.Env,
	}
}

// ConfigInterfaces 接口类型配置集合
type ConfigInterfaces struct {
	Server     interfaces.ServerConfig
	Database   interfaces.DatabaseConfig
	Auth       interfaces.AuthConfig
	Logger     interfaces.LoggerConfig
	Redis      interfaces.RedisConfig
	Token      interfaces.TokenConfig
	Validation interfaces.ValidationConfig
	RateLimit  interfaces.RateLimitConfig
	Env        string
}