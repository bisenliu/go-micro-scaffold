package config

import (
	"go.uber.org/fx"

	"common/interfaces"
)

// ConfigProviderImpl 配置提供者实现
type ConfigProviderImpl struct {
	serverConfig     interfaces.ServerConfig
	databaseConfig   interfaces.DatabaseConfig
	authConfig       interfaces.AuthConfig
	loggerConfig     interfaces.LoggerConfig
	redisConfig      interfaces.RedisConfig
	tokenConfig      interfaces.TokenConfig
	validationConfig interfaces.ValidationConfig
	rateLimitConfig  interfaces.RateLimitConfig
	env              string
	factory          *ConfigFactory
}

// NewConfigProvider 创建配置提供者
func NewConfigProvider() (interfaces.ConfigProvider, error) {
	factory := NewConfigFactory()
	return factory.CreateConfigProvider()
}

// GetDatabaseConfig 获取数据库配置
func (p *ConfigProviderImpl) GetDatabaseConfig() interfaces.DatabaseConfig {
	return p.databaseConfig
}

// GetServerConfig 获取服务器配置
func (p *ConfigProviderImpl) GetServerConfig() interfaces.ServerConfig {
	return p.serverConfig
}

// GetAuthConfig 获取认证配置
func (p *ConfigProviderImpl) GetAuthConfig() interfaces.AuthConfig {
	return p.authConfig
}

// GetLoggerConfig 获取日志配置
func (p *ConfigProviderImpl) GetLoggerConfig() interfaces.LoggerConfig {
	return p.loggerConfig
}

// GetRedisConfig 获取Redis配置
func (p *ConfigProviderImpl) GetRedisConfig() interfaces.RedisConfig {
	return p.redisConfig
}

// GetTokenConfig 获取Token配置
func (p *ConfigProviderImpl) GetTokenConfig() interfaces.TokenConfig {
	return p.tokenConfig
}

// GetValidationConfig 获取验证配置
func (p *ConfigProviderImpl) GetValidationConfig() interfaces.ValidationConfig {
	return p.validationConfig
}

// GetRateLimitConfig 获取限流配置
func (p *ConfigProviderImpl) GetRateLimitConfig() interfaces.RateLimitConfig {
	return p.rateLimitConfig
}

// Reload 重新加载配置
func (p *ConfigProviderImpl) Reload() error {
	if p.factory != nil {
		return p.factory.Reload()
	}
	return nil
}

// GetEnv 获取当前环境
func (p *ConfigProviderImpl) GetEnv() string {
	return p.env
}

// Module FX模块
var Module = fx.Provide(NewConfigProvider)