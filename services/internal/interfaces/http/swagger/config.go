package swagger

import (
	"os"
	"strconv"
	"strings"

	"common/config"
)

// SwaggerManager Swagger配置管理器
type SwaggerManager struct {
	config *config.SwaggerConfig
}

// NewSwaggerManager 创建Swagger配置管理器
func NewSwaggerManager(cfg *config.Config) *SwaggerManager {
	swaggerConfig := cfg.Swagger

	// 应用环境变量覆盖
	applyEnvironmentOverrides(&swaggerConfig)

	// 设置默认值
	setDefaults(&swaggerConfig)

	return &SwaggerManager{
		config: &swaggerConfig,
	}
}

// GetConfig 获取Swagger配置
func (sm *SwaggerManager) GetConfig() *config.SwaggerConfig {
	return sm.config
}

// IsEnabled 检查Swagger是否启用
func (sm *SwaggerManager) IsEnabled() bool {
	return sm.config.Enabled
}

// ShouldEnableInEnvironment 根据环境判断是否应该启用Swagger
func (sm *SwaggerManager) ShouldEnableInEnvironment(env string) bool {
	// 生产环境默认禁用
	if strings.ToLower(env) == "production" {
		// 除非明确通过环境变量启用
		if enabledEnv := os.Getenv("SWAGGER_ENABLED"); enabledEnv != "" {
			enabled, _ := strconv.ParseBool(enabledEnv)
			return enabled
		}
		return false
	}

	// 其他环境（开发、测试）默认启用
	return sm.config.Enabled
}

// applyEnvironmentOverrides 应用环境变量覆盖配置
func applyEnvironmentOverrides(cfg *config.SwaggerConfig) {
	// SWAGGER_ENABLED 环境变量
	if enabledEnv := os.Getenv("SWAGGER_ENABLED"); enabledEnv != "" {
		if enabled, err := strconv.ParseBool(enabledEnv); err == nil {
			cfg.Enabled = enabled
		}
	}
}

// setDefaults 设置默认值
func setDefaults(cfg *config.SwaggerConfig) {
	// 设置默认标题
	if cfg.Title == "" {
		cfg.Title = "API Documentation"
	}

	// 设置默认版本
	if cfg.Version == "" {
		cfg.Version = "1.0.0"
	}
}
