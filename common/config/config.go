package config

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Config 应用配置结构
type Config struct {
	System     SystemConfig              `mapstructure:"system"`
	Token      TokenConfig               `mapstructure:"token"`
	SnowFlake  SnowFlakeConfig           `mapstructure:"snow_flake"`
	Database   DatabaseConfig            `mapstructure:"database"`
	Databases  map[string]DatabaseConfig `mapstructure:"databases"` // 多数据库配置
	Redis      RedisConfig               `mapstructure:"redis"`
	Server     ServerConfig              `mapstructure:"server"`
	Zap        ZapConfig                 `mapstructure:"zap"`
	WeChat     WeChatConfig              `mapstructure:"wechat"`
	AliyunSMS  AliyunSMSConfig           `mapstructure:"aliyun_sms"`
	Validation ValidationConfig          `mapstructure:"validation"`
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

type WeChatConfig struct {
	AppIDMini              string `mapstructure:"app_id_mini"`
	SecretMini             string `mapstructure:"secret_mini"`
	AppIDServe             string `mapstructure:"app_id_serve"`
	SecretServe            string `mapstructure:"secret_serve"`
	TeamStatusMessageID    string `mapstructure:"team_status_message_id"`
	PaymentNoticeMessageID string `mapstructure:"payment_notice_message_id"`
	MerchID                string `mapstructure:"merch_id"`
	PaymentNotifyURL       string `mapstructure:"payment_notify_url"`
	RefundNotifyURL        string `mapstructure:"refund_notify_url"`
	GetOpenIDURL           string `mapstructure:"get_openid_url"`
	GetAccessTokenURL      string `mapstructure:"get_access_token_url"`
	GetPhoneNumberURL      string `mapstructure:"get_phone_number_url"`
	SendTemplateMessageURL string `mapstructure:"send_template_message_url"`
}

type AliyunSMSConfig struct {
	AccessKeyID       string `mapstructure:"access_key_id"`
	AccessKeySecret   string `mapstructure:"access_key_secret"`
	RegionID          string `mapstructure:"region_id"`
	Endpoint          string `mapstructure:"endpoint"`
	SignName          string `mapstructure:"sign_name"`
	TeamNoticeCode    string `mapstructure:"team_notice_code"`
	PaymentNoticeCode string `mapstructure:"payment_notice_code"`
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

	return &config, nil
}

// GetDatabaseConfig 获取数据库配置，支持向后兼容
func (c *Config) GetDatabaseConfig(name string) (DatabaseConfig, bool) {
	// 如果指定了databases配置，优先使用
	if c.Databases != nil {
		if config, exists := c.Databases[name]; exists {
			return config, true
		}
	}

	// 向后兼容：如果请求primary且没有databases配置，使用database配置
	if name == "primary" && c.Databases == nil {
		return c.Database, true
	}

	return DatabaseConfig{}, false
}

// GetAllDatabaseConfigs 获取所有数据库配置
func (c *Config) GetAllDatabaseConfigs() map[string]DatabaseConfig {
	if c.Databases != nil {
		return c.Databases
	}

	// 向后兼容：如果没有databases配置，返回primary数据库
	return map[string]DatabaseConfig{
		"primary": c.Database,
	}
}

// Module FX模块
var Module = fx.Provide(NewConfig)
