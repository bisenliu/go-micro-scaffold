package mysql

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/config"
)

// Manager 数据库管理器接口
type ManagerInterface interface {
	GetClient(name string) (ClientInterface, error)
	Primary() (ClientInterface, error)
	Read() (ClientInterface, error)
	Write() (ClientInterface, error)
	ListClients() []string
	HasClient(name string) bool
	Close() error
}

// Manager 数据库管理器，管理多个数据库连接
type Manager struct {
	clients       map[string]*Client
	configManager *ConfigManager
	logger        *zap.Logger
	mu            sync.RWMutex
}

// ManagerParams 数据库管理器依赖参数
type ManagerParams struct {
	fx.In
	Config *config.Config
	Logger *zap.Logger
}

// NewManager 创建数据库管理器
func NewManager(params ManagerParams) (*Manager, error) {
	configManager := NewConfigManager(params.Config)

	manager := &Manager{
		clients:       make(map[string]*Client),
		configManager: configManager,
		logger:        params.Logger,
	}

	// 为每个数据库配置创建客户端
	configNames := configManager.ListConfigs()
	for _, name := range configNames {
		dbConfig, err := configManager.GetConfig(name)
		if err != nil {
			manager.closeAllClients()
			return nil, fmt.Errorf("failed to get config for '%s': %w", name, err)
		}

		builder := NewClientBuilder(dbConfig, params.Logger)
		client, err := builder.Build()
		if err != nil {
			manager.closeAllClients()
			return nil, fmt.Errorf("failed to create database client for '%s': %w", name, err)
		}

		manager.clients[name] = client
	}

	params.Logger.Info("Database manager initialized successfully",
		zap.Int("client_count", len(manager.clients)),
		zap.Strings("clients", configNames),
	)

	return manager, nil
}

// GetClient 获取指定名称的数据库客户端
func (m *Manager) GetClient(name string) (ClientInterface, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("database client '%s' not found", name)
	}
	return client, nil
}

// Primary 获取主数据库客户端
func (m *Manager) Primary() (ClientInterface, error) {
	return m.GetClient("primary")
}

// Read 获取读数据库客户端（读写分离场景）
func (m *Manager) Read() (ClientInterface, error) {
	// 优先使用read客户端，如果没有则使用primary
	if client, err := m.GetClient("read"); err == nil {
		return client, nil
	}
	return m.Primary()
}

// Write 获取写数据库客户端（读写分离场景）
func (m *Manager) Write() (ClientInterface, error) {
	// 优先使用write客户端，如果没有则使用primary
	if client, err := m.GetClient("write"); err == nil {
		return client, nil
	}
	return m.Primary()
}

// ListClients 列出所有数据库客户端名称
func (m *Manager) ListClients() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	return names
}

// HasClient 检查是否存在指定名称的客户端
func (m *Manager) HasClient(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.clients[name]
	return exists
}

// Close 关闭所有数据库连接
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.closeAllClients()
}

// closeAllClients 关闭所有客户端连接（内部方法，不加锁）
func (m *Manager) closeAllClients() error {
	var lastErr error
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			m.logger.Error("Failed to close database client",
				zap.String("name", name),
				zap.Error(err),
			)
			lastErr = err
		}
	}
	m.clients = make(map[string]*Client)
	return lastErr
}

// ManagerModule 数据库管理器模块
var ManagerModule = fx.Module("database_manager",
	fx.Provide(NewManager),
	fx.Invoke(func(lc fx.Lifecycle, manager *Manager) {
		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return manager.Close()
			},
		})
	}),
)
