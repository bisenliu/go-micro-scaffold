package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// EnvironmentManager 环境管理器
// 支持环境特定的配置覆盖
type EnvironmentManager struct {
	env        string
	configPath string
}

// NewEnvironmentManager 创建环境管理器
func NewEnvironmentManager() *EnvironmentManager {
	return &EnvironmentManager{}
}

// LoadEnvironmentConfig 加载环境特定配置
func (em *EnvironmentManager) LoadEnvironmentConfig() error {
	// 获取环境变量
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}
	em.env = env

	// 设置配置文件搜索路径
	configPaths := []string{
		"./configs",
		".",
		"../configs",
	}

	// 尝试加载基础配置文件
	baseConfigFound := false
	for _, path := range configPaths {
		viper.AddConfigPath(path)
		em.configPath = path
		
		// 尝试加载 app.yaml
		viper.SetConfigName("app")
		viper.SetConfigType("yaml")
		
		if err := viper.ReadInConfig(); err == nil {
			baseConfigFound = true
			break
		}
	}

	if !baseConfigFound {
		return fmt.Errorf("base config file 'app.yaml' not found in paths: %v", configPaths)
	}

	// 尝试加载环境特定配置文件
	if err := em.loadEnvironmentSpecificConfig(); err != nil {
		// 环境特定配置文件不存在是可以接受的
		fmt.Printf("Environment specific config not found (this is optional): %v\n", err)
	}

	// 加载环境变量覆盖
	em.loadEnvironmentVariables()

	return nil
}

// loadEnvironmentSpecificConfig 加载环境特定配置
func (em *EnvironmentManager) loadEnvironmentSpecificConfig() error {
	envConfigFile := fmt.Sprintf("app.%s.yaml", em.env)
	envConfigPath := filepath.Join(em.configPath, envConfigFile)

	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("environment config file %s not found", envConfigFile)
	}

	// 创建新的 viper 实例来读取环境配置
	envViper := viper.New()
	envViper.SetConfigFile(envConfigPath)
	envViper.SetConfigType("yaml")

	if err := envViper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read environment config: %w", err)
	}

	// 合并环境配置到主配置
	if err := viper.MergeConfigMap(envViper.AllSettings()); err != nil {
		return fmt.Errorf("failed to merge environment config: %w", err)
	}

	return nil
}

// loadEnvironmentVariables 加载环境变量覆盖
func (em *EnvironmentManager) loadEnvironmentVariables() {
	// 设置环境变量前缀
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	
	// 设置环境变量键名替换规则
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 绑定特定的环境变量
	envMappings := map[string]string{
		"APP_ENV":              "system.env",
		"APP_SECRET_KEY":       "system.secret_key",
		"APP_SERVER_PORT":      "server.port",
		"APP_SERVER_HOST":      "server.host",
		"APP_DB_HOST":          "database.primary.host",
		"APP_DB_PORT":          "database.primary.port",
		"APP_DB_NAME":          "database.primary.database",
		"APP_DB_USER":          "database.primary.username",
		"APP_DB_PASSWORD":      "database.primary.password",
		"APP_REDIS_HOST":       "redis.host",
		"APP_REDIS_PORT":       "redis.port",
		"APP_REDIS_PASSWORD":   "redis.password",
		"APP_LOG_LEVEL":        "logger.level",
		"APP_LOG_DIR":          "logger.director",
	}

	for envKey, configKey := range envMappings {
		if envValue := os.Getenv(envKey); envValue != "" {
			viper.Set(configKey, envValue)
		}
	}
}

// GetEnvironment 获取当前环境
func (em *EnvironmentManager) GetEnvironment() string {
	return em.env
}

// IsProduction 是否为生产环境
func (em *EnvironmentManager) IsProduction() bool {
	return em.env == "production"
}

// IsDevelopment 是否为开发环境
func (em *EnvironmentManager) IsDevelopment() bool {
	return em.env == "development"
}

// IsStaging 是否为预发布环境
func (em *EnvironmentManager) IsStaging() bool {
	return em.env == "staging"
}

// GetConfigPath 获取配置文件路径
func (em *EnvironmentManager) GetConfigPath() string {
	return em.configPath
}