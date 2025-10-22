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
	setDefaults(&swaggerConfig, &cfg.System)

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

// GetSwaggerURL 获取Swagger UI访问URL
func (sm *SwaggerManager) GetSwaggerURL() string {
	return "/swagger/*any"
}

// GetAPIDocsURL 获取API文档JSON URL
func (sm *SwaggerManager) GetAPIDocsURL() string {
	return "/swagger/doc.json"
}

// applyEnvironmentOverrides 应用环境变量覆盖配置
func applyEnvironmentOverrides(cfg *config.SwaggerConfig) {
	// SWAGGER_ENABLED 环境变量
	if enabledEnv := os.Getenv("SWAGGER_ENABLED"); enabledEnv != "" {
		if enabled, err := strconv.ParseBool(enabledEnv); err == nil {
			cfg.Enabled = enabled
		}
	}

	// SWAGGER_TITLE 环境变量
	if title := os.Getenv("SWAGGER_TITLE"); title != "" {
		cfg.Title = title
	}

	// SWAGGER_DESCRIPTION 环境变量
	if description := os.Getenv("SWAGGER_DESCRIPTION"); description != "" {
		cfg.Description = description
	}

	// SWAGGER_VERSION 环境变量
	if version := os.Getenv("SWAGGER_VERSION"); version != "" {
		cfg.Version = version
	}

	// SWAGGER_HOST 环境变量
	if host := os.Getenv("SWAGGER_HOST"); host != "" {
		cfg.Host = host
	}

	// SWAGGER_BASE_PATH 环境变量
	if basePath := os.Getenv("SWAGGER_BASE_PATH"); basePath != "" {
		cfg.BasePath = basePath
	}
}

// setDefaults 设置默认值
func setDefaults(cfg *config.SwaggerConfig, systemCfg *config.SystemConfig) {
	// 根据环境设置默认启用状态
	if systemCfg.Env == "production" {
		// 生产环境默认禁用，除非明确配置启用
		if cfg.Title == "" { // 如果没有配置过，则设置默认禁用
			cfg.Enabled = false
		}
	}

	// 设置默认标题
	if cfg.Title == "" {
		cfg.Title = "API Documentation"
	}

	// 设置默认描述
	if cfg.Description == "" {
		cfg.Description = "API documentation for " + systemCfg.ServerName
	}

	// 设置默认版本
	if cfg.Version == "" {
		cfg.Version = "1.0.0"
	}

	// 设置默认基础路径
	if cfg.BasePath == "" {
		cfg.BasePath = "/api/v1"
	}

	// 设置默认联系信息
	if cfg.Contact.Name == "" {
		cfg.Contact.Name = "API Support"
	}

	// 设置默认许可证
	if cfg.License.Name == "" {
		cfg.License.Name = "MIT"
	}
}
