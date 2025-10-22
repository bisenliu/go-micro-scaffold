package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Config 应用配置结构
type Config struct {
	// 1. 系统与服务器配置
	System SystemConfig `mapstructure:"system"`
	Server ServerConfig `mapstructure:"server"`

	// 2. 中间件配置
	Auth      AuthConfig      `mapstructure:"auth"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`

	// 3. 业务逻辑相关配置
	Token      TokenConfig      `mapstructure:"token"`
	SnowFlake  SnowFlakeConfig  `mapstructure:"snow_flake"`
	Validation ValidationConfig `mapstructure:"validation"`

	// 4. 外部服务依赖配置
	DatabaseCommon  DatabaseConfig            `mapstructure:"database_common"`
	Databases       map[string]DatabaseConfig `mapstructure:"databases"`
	DatabaseAliases map[string]string         `mapstructure:"database_aliases"`
	Redis           RedisConfig               `mapstructure:"redis"`

	// 5. Swagger API 文档配置
	Swagger SwaggerConfig `mapstructure:"swagger"`

	// 6. 日志配置
	Zap ZapConfig `mapstructure:"zap"`
}

// --- 1. 系统与服务器配置 ---

type SystemConfig struct {
	Env        string `mapstructure:"env"`
	SecretKey  string `mapstructure:"secret_key"`
	ServerName string `mapstructure:"server_name"`
	Timezone   string `mapstructure:"timezone"`
}

type ServerConfig struct {
	Port           string        `mapstructure:"port"`
	Mode           string        `mapstructure:"mode"`
	EnableCORS     bool          `mapstructure:"enable_cors"` // CORS开关也属于中间件配置
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
}

// --- 2. 中间件配置 ---

type AuthConfig struct {
	Whitelist []string `mapstructure:"whitelist"`
}

type RateLimitConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	FillInterval    time.Duration `mapstructure:"fill_interval"`
	Capacity        int64         `mapstructure:"capacity"`
	Quantum         int64         `mapstructure:"quantum"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
	BucketExpiry    time.Duration `mapstructure:"bucket_expiry"`
}

// --- 3. 业务逻辑相关配置 ---

type TokenConfig struct {
	ExpiredTime int `mapstructure:"expired_time"`
}

type SnowFlakeConfig struct {
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

type ValidationConfig struct {
	Locale string `mapstructure:"locale"`
}

// --- 4. 外部服务依赖配置 ---

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

type RedisConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	Database  int    `mapstructure:"database"`
	DefaultDB int    `mapstructure:"default_db"`
	PoolSize  int    `mapstructure:"pool_size"`
}

// --- 5. Swagger API 文档配置 ---

// SwaggerConfig Swagger配置结构
type SwaggerConfig struct {
	Enabled     bool          `mapstructure:"enabled"`     // 是否启用Swagger
	Title       string        `mapstructure:"title"`       // API文档标题
	Description string        `mapstructure:"description"` // API描述
	Version     string        `mapstructure:"version"`     // API版本
	Host        string        `mapstructure:"host"`        // API主机地址
	BasePath    string        `mapstructure:"base_path"`   // API基础路径
	Contact     ContactConfig `mapstructure:"contact"`     // 联系信息
	License     LicenseConfig `mapstructure:"license"`     // 许可证信息
}

// ContactConfig 联系信息配置
type ContactConfig struct {
	Name  string `mapstructure:"name"`  // 联系人姓名
	Email string `mapstructure:"email"` // 联系人邮箱
	URL   string `mapstructure:"url"`   // 联系人URL
}

// LicenseConfig 许可证配置
type LicenseConfig struct {
	Name string `mapstructure:"name"` // 许可证名称
	URL  string `mapstructure:"url"`  // 许可证URL
}

// --- 6. 日志配置 ---

type ZapConfig struct {
	Level      string `mapstructure:"level"`
	Director   string `mapstructure:"director"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// --- 配置加载 ---

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

	// 验证数据库别名配置
	if err := validateDatabaseAliases(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateDatabaseAliases 验证数据库别名配置的有效性
func validateDatabaseAliases(config *Config) error {
	if config.DatabaseAliases == nil {
		return nil // 别名配置是可选的
	}

	// 检查每个别名是否指向存在的数据库
	for alias, dbName := range config.DatabaseAliases {
		if _, exists := config.Databases[dbName]; !exists {
			return fmt.Errorf("database alias '%s' points to non-existent database '%s'", alias, dbName)
		}
	}

	return nil
}

// Module FX模块
var Module = fx.Provide(NewConfig)
