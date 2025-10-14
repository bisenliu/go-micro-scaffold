package config

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Config 应用配置结构
type Config struct {
	System         SystemConfig              `mapstructure:"system"`
	Token          TokenConfig               `mapstructure:"token"`
	SnowFlake      SnowFlakeConfig           `mapstructure:"snow_flake"`
	DatabaseCommon DatabaseConfig            `mapstructure:"database_common"`
	Databases      map[string]DatabaseConfig `mapstructure:"databases"`
	Redis          RedisConfig               `mapstructure:"redis"`
	Server         ServerConfig              `mapstructure:"server"`
	Zap            ZapConfig                 `mapstructure:"zap"`
	Validation     ValidationConfig          `mapstructure:"validation"`
	RateLimit      RateLimitConfig           `mapstructure:"rate_limit"`
}

type RateLimitConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	FillInterval    time.Duration `mapstructure:"fill_interval"`
	Capacity        int64         `mapstructure:"capacity"`
	Quantum         int64         `mapstructure:"quantum"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
	BucketExpiry    time.Duration `mapstructure:"bucket_expiry"`
}

type SystemConfig struct {
	Env        string `mapstructure:"env"`
	SecretKey  string `mapstructure:"secret_key"`
	ServerName string `mapstructure:"server_name"`
	Timezone   string `mapstructure:"timezone"` // 添加时区配置
}

type TokenConfig struct {
	ExpiredTime int `mapstructure:"expired_time"`
}

type SnowFlakeConfig struct {
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

type ServerConfig struct {
	Port           string        `mapstructure:"port"`
	Mode           string        `mapstructure:"mode"`
	EnableCORS     bool          `mapstructure:"enable_cors"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
}

type RedisConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	Database  int    `mapstructure:"database"`
	DefaultDB int    `mapstructure:"default_db"`
	PoolSize  int    `mapstructure:"pool_size"`
}

type ZapConfig struct {
	Level      string `mapstructure:"level"`
	Director   string `mapstructure:"director"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type ValidationConfig struct {
	Locale string `mapstructure:"locale"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // "mysql", "postgres", "sqlite" 等
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Charset  string `mapstructure:"charset"`

	// 连接池配置
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

// NewConfig 创建配置实例
func NewConfig() (*Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Module FX模块
var Module = fx.Provide(NewConfig)
