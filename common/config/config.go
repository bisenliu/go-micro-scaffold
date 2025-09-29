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
}

type SystemConfig struct {
	Env        string `mapstructure:"env"`
	Port       string `mapstructure:"port"`
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
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
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
	Type     string `mapstructure:"type" json:"type" yaml:"type"` // "mysql", "postgres", "sqlite" 等
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Database string `mapstructure:"database" json:"database" yaml:"database"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Charset  string `mapstructure:"charset" json:"charset" yaml:"charset"`

	// 连接池配置
	MaxOpenConns    int           `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
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

	// 将数据库公共配置合并到各个数据库实例配置中
	for name, db := range config.Databases {
		// 如果某个字段没有在具体数据库中设置，则使用公共配置的值
		if db.Type == "" {
			db.Type = config.DatabaseCommon.Type
		}
		if db.Host == "" {
			db.Host = config.DatabaseCommon.Host
		}
		if db.Port == 0 {
			db.Port = config.DatabaseCommon.Port
		}
		if db.Username == "" {
			db.Username = config.DatabaseCommon.Username
		}
		if db.Password == "" {
			db.Password = config.DatabaseCommon.Password
		}
		if db.Charset == "" {
			db.Charset = config.DatabaseCommon.Charset
		}
		if db.MaxOpenConns == 0 {
			db.MaxOpenConns = config.DatabaseCommon.MaxOpenConns
		}
		if db.MaxIdleConns == 0 {
			db.MaxIdleConns = config.DatabaseCommon.MaxIdleConns
		}
		if db.ConnMaxLifetime == 0 {
			db.ConnMaxLifetime = config.DatabaseCommon.ConnMaxLifetime
		}
		if db.ConnMaxIdleTime == 0 {
			db.ConnMaxIdleTime = config.DatabaseCommon.ConnMaxIdleTime
		}

		// 更新配置
		config.Databases[name] = db
	}

	return &config, nil
}

// Module FX模块
var Module = fx.Provide(NewConfig)
