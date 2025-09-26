package mysql

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/config"
)

// DatabaseManager 数据库管理器，管理多个数据库连接
type DatabaseManager struct {
	clients map[string]*EntClient
	config  *config.Config
	logger  *zap.Logger
	mu      sync.RWMutex
}

// DatabaseManagerParams 数据库管理器依赖参数
type DatabaseManagerParams struct {
	fx.In
	Config *config.Config
	Logger *zap.Logger
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager(params DatabaseManagerParams) (*DatabaseManager, error) {
	manager := &DatabaseManager{
		clients: make(map[string]*EntClient),
		config:  params.Config,
		logger:  params.Logger,
	}

	// 获取所有数据库配置
	dbConfigs := params.Config.GetAllDatabaseConfigs()

	// 为每个数据库创建连接
	for name, dbConfig := range dbConfigs {
		client, err := createEntClientWithConfig(dbConfig)
		if err != nil {
			// 清理已创建的连接
			manager.closeAllClients()
			return nil, fmt.Errorf("failed to create database client for '%s': %w", name, err)
		}
		manager.clients[name] = client

		params.Logger.Info("Database client created successfully",
			zap.String("name", name),
			zap.String("host", dbConfig.Host),
			zap.Int("port", dbConfig.Port),
			zap.String("database", dbConfig.Database),
		)
	}

	return manager, nil
}

// GetClient 获取指定名称的数据库客户端
func (dm *DatabaseManager) GetClient(name string) (*EntClient, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	client, exists := dm.clients[name]
	if !exists {
		return nil, fmt.Errorf("database client '%s' not found", name)
	}
	return client, nil
}

// GetPrimaryClient 获取主数据库客户端（向后兼容）
func (dm *DatabaseManager) GetPrimaryClient() (*EntClient, error) {
	return dm.GetClient("primary")
}

// ListClients 列出所有数据库客户端名称
func (dm *DatabaseManager) ListClients() []string {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	names := make([]string, 0, len(dm.clients))
	for name := range dm.clients {
		names = append(names, name)
	}
	return names
}

// Close 关闭所有数据库连接
func (dm *DatabaseManager) Close() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	return dm.closeAllClients()
}

// closeAllClients 关闭所有客户端连接（内部方法，不加锁）
func (dm *DatabaseManager) closeAllClients() error {
	var lastErr error
	for name, client := range dm.clients {
		if err := client.Close(); err != nil {
			dm.logger.Error("Failed to close database client",
				zap.String("name", name),
				zap.Error(err),
			)
			lastErr = err
		}
	}
	dm.clients = make(map[string]*EntClient)
	return lastErr
}

// createEntClientWithConfig 使用指定配置创建EntClient
func createEntClientWithConfig(dbConfig config.DatabaseConfig) (*EntClient, error) {
	// 创建临时的参数结构
	tempConfig := &config.Config{
		Database: dbConfig,
	}

	params := EntClientParams{
		Config: tempConfig,
	}

	return NewEntClient(params)
}

// DatabaseManagerModule 数据库管理器模块
var DatabaseManagerModule = fx.Module("database_manager",
	fx.Provide(NewDatabaseManager),
	fx.Invoke(func(lc fx.Lifecycle, manager *DatabaseManager) {
		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return manager.Close()
			},
		})
	}),
)
