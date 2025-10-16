package main

import (
	"os"
	"strings"
)

// AppConfig 应用配置
type AppConfig struct {
	Environment string
	Debug       bool
	EnableHTTP  bool
	EnableGraph bool
}

// NewAppConfig 创建应用配置
func NewAppConfig() *AppConfig {
	return &AppConfig{
		Environment: getEnv("APP_ENV", "development"),
		Debug:       getEnvBool("DEBUG", false),
		EnableHTTP:  getEnvBool("ENABLE_HTTP", true),
		EnableGraph: getEnvBool("ENABLE_GRAPH", false),
	}
}

// IsProduction 是否为生产环境
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment 是否为开发环境
func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsTest 是否为测试环境
func (c *AppConfig) IsTest() bool {
	return c.Environment == "test"
}

// ShouldLoadHTTPModule 是否应该加载HTTP模块
func (c *AppConfig) ShouldLoadHTTPModule() bool {
	return c.EnableHTTP && !c.IsTest()
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool 获取布尔类型环境变量
func getEnvBool(key string, defaultValue bool) bool {
	value := strings.ToLower(getEnv(key, ""))
	switch value {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}