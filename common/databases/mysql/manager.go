package mysql

import (
	"fmt"
	"sync"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/config"
)

// Manager 数据库管理器
type Manager struct {
	clients       sync.Map
	configManager *ConfigManager
	logger        *zap.Logger
}

// ManagerParams 管理器依赖参数
type ManagerParams struct {
	fx.In
	Config *config.Config
	Logger *zap.Logger
}

// NewManager 创建数据库管理器
func NewManager(params ManagerParams) (*Manager, error) {
	configManager := NewConfigManager(params.Config)

	manager := &Manager{
		configManager: configManager,
		logger:        params.Logger,
	}

	// 初始化所有客户端
	if err := manager.initializeClients(params); err != nil {
		return nil, err
	}

	return manager, nil
}

// initializeClients 初始化所有客户端
func (m *Manager) initializeClients(params ManagerParams) error {
	configNames := m.configManager.ListConfigs()

	for _, name := range configNames {
		dbConfig, err := m.configManager.GetConfig(name)
		if err != nil {
			return fmt.Errorf("failed to get config for '%s': %w", name, err)
		}

		builder := NewClientBuilder(dbConfig, params.Logger)
		client, err := builder.Build()
		if err != nil {
			return fmt.Errorf("failed to create database client for '%s': %w", name, err)
		}

		m.clients.Store(name, client)
	}

	m.logger.Info("Database manager initialized successfully",
		zap.Int("client_count", len(configNames)),
		zap.Strings("clients", configNames))

	return nil
}

// GetClient 获取指定名称的数据库客户端
func (m *Manager) GetClient(name string) (ClientInterface, error) {
	if client, ok := m.clients.Load(name); ok {
		return client.(*Client), nil
	}
	return nil, fmt.Errorf("database client '%s' not found", name)
}

// Primary 获取主数据库客户端
func (m *Manager) Primary() (ClientInterface, error) {
	return m.GetClient(DB1)
}

// ListClients 列出所有数据库客户端名称
func (m *Manager) ListClients() []string {
	var names []string
	m.clients.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

// HasClient 检查是否存在指定名称的客户端
func (m *Manager) HasClient(name string) bool {
	_, ok := m.clients.Load(name)
	return ok
}

// Close 关闭所有数据库连接
func (m *Manager) Close() error {
	var lastErr error
	m.clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		if err := client.Close(); err != nil {
			m.logger.Error("Failed to close database client",
				zap.String("name", key.(string)),
				zap.Error(err))
			lastErr = err
		}
		return true
	})
	return lastErr
}
