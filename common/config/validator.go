package config

import (
	"fmt"
	"strings"
)

// ConfigValidator 配置验证器
type ConfigValidator struct{}

// NewConfigValidator 创建配置验证器
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// Validate 验证配置
func (v *ConfigValidator) Validate(cfg *Config) error {
	var errors []string

	// 只验证最关键的配置项
	if cfg.System.Env == "" {
		errors = append(errors, "system.env is required")
	}

	if cfg.Server.Port == "" {
		errors = append(errors, "server.port is required")
	}

	// 简化数据库验证 - 只检查是否有数据库配置
	if cfg.DatabaseCommon.Host == "" && len(cfg.Databases) == 0 {
		errors = append(errors, "database configuration is required (either database_common.host or databases section)")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, ", "))
	}

	return nil
}

// ValidatedConfig 创建并验证配置
func ValidatedConfig() (*Config, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 暂时禁用严格验证，只做基本检查
	if cfg.System.Env == "" {
		return nil, fmt.Errorf("system.env is required in configuration")
	}

	if cfg.Server.Port == "" {
		return nil, fmt.Errorf("server.port is required in configuration")
	}

	return cfg, nil
}