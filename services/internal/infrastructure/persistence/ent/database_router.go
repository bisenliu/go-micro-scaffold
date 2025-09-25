package ent

import (
	"fmt"

	"common/database"
	"services/internal/infrastructure/persistence/ent/gen"
)

// DatabaseRouter 数据库路由器，管理多个业务Ent客户端
type DatabaseRouter struct {
	clients map[string]*gen.Client
	factory database.EntClientFactory
}

// NewDatabaseRouter 创建数据库路由器
func NewDatabaseRouter(manager *database.DatabaseManager, factory database.EntClientFactory) (*DatabaseRouter, error) {
	router := &DatabaseRouter{
		clients: make(map[string]*gen.Client),
		factory: factory,
	}

	// 为每个数据库连接创建业务Ent客户端
	clientNames := manager.ListClients()
	for _, name := range clientNames {
		commonClient, err := manager.GetClient(name)
		if err != nil {
			return nil, fmt.Errorf("failed to get database client '%s': %w", name, err)
		}

		businessClient, err := factory.NewClient(commonClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create business client for '%s': %w", name, err)
		}

		genClient, ok := businessClient.(*gen.Client)
		if !ok {
			return nil, fmt.Errorf("invalid client type for '%s', expected *gen.Client", name)
		}

		router.clients[name] = genClient
	}

	return router, nil
}

// GetClient 获取指定名称的业务数据库客户端
func (dr *DatabaseRouter) GetClient(name string) (*gen.Client, error) {
	client, exists := dr.clients[name]
	if !exists {
		return nil, fmt.Errorf("business database client '%s' not found", name)
	}
	return client, nil
}

// Primary 获取主数据库客户端
func (dr *DatabaseRouter) Primary() (*gen.Client, error) {
	return dr.GetClient("primary")
}

// Analytics 获取分析数据库客户端（如果配置了的话）
func (dr *DatabaseRouter) Analytics() (*gen.Client, error) {
	return dr.GetClient("analytics")
}


// Read 获取读数据库客户端（读写分离场景）
func (dr *DatabaseRouter) Read() (*gen.Client, error) {
	// 优先使用read客户端，如果没有则使用primary
	if client, err := dr.GetClient("read"); err == nil {
		return client, nil
	}
	return dr.Primary()
}

// Write 获取写数据库客户端（读写分离场景）
func (dr *DatabaseRouter) Write() (*gen.Client, error) {
	// 优先使用write客户端，如果没有则使用primary
	if client, err := dr.GetClient("write"); err == nil {
		return client, nil
	}
	return dr.Primary()
}

// ListClients 列出所有可用的数据库客户端名称
func (dr *DatabaseRouter) ListClients() []string {
	names := make([]string, 0, len(dr.clients))
	for name := range dr.clients {
		names = append(names, name)
	}
	return names
}

// HasClient 检查是否存在指定名称的客户端
func (dr *DatabaseRouter) HasClient(name string) bool {
	_, exists := dr.clients[name]
	return exists
}

// Close 关闭所有客户端连接
func (dr *DatabaseRouter) Close() error {
	var lastErr error
	for name, client := range dr.clients {
		if err := client.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close client '%s': %w", name, err)
		}
	}
	dr.clients = make(map[string]*gen.Client)
	return lastErr
}
