package mysql

import (
	"fmt"

	"entgo.io/ent/dialect"

	"common/config"
)

// DatabaseConfigWrapper 数据库配置包装器，扩展原有配置
type DatabaseConfigWrapper struct {
	config.DatabaseConfig
	Name        string
	EnableDebug bool
}

// DSN 生成数据库连接字符串
func (cfg *DatabaseConfigWrapper) DSN() (string, string, error) {
	switch cfg.Type {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
		return dsn, dialect.Postgres, nil
	case "sqlite":
		return cfg.Database, dialect.SQLite, nil
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		return dsn, dialect.MySQL, nil
	default:
		return "", "", fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}

// ConfigManager 配置管理器
type ConfigManager struct {
	configs map[string]*DatabaseConfigWrapper
}

// NewConfigManager 创建配置管理器
func NewConfigManager(cfg *config.Config) *ConfigManager {
	manager := &ConfigManager{
		configs: make(map[string]*DatabaseConfigWrapper),
	}

	// 获取所有数据库配置
	dbConfigs := cfg.GetAllDatabaseConfigs()

	for name, dbConfig := range dbConfigs {
		manager.configs[name] = &DatabaseConfigWrapper{
			DatabaseConfig: dbConfig,
			Name:           name,
			EnableDebug:    cfg.System.Env == "development",
		}
	}

	return manager
}

// GetConfig 获取指定名称的数据库配置
func (cm *ConfigManager) GetConfig(name string) (*DatabaseConfigWrapper, error) {
	config, exists := cm.configs[name]
	if !exists {
		return nil, fmt.Errorf("database config '%s' not found", name)
	}
	return config, nil
}

// ListConfigs 列出所有配置名称
func (cm *ConfigManager) ListConfigs() []string {
	names := make([]string, 0, len(cm.configs))
	for name := range cm.configs {
		names = append(names, name)
	}
	return names
}

// HasConfig 检查是否存在指定配置
func (cm *ConfigManager) HasConfig(name string) bool {
	_, exists := cm.configs[name]
	return exists
}
