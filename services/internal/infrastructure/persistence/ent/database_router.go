package ent

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"common/databases/mysql"
	"services/internal/infrastructure/persistence/ent/gen"
)

// RouterInterface 数据库路由器接口
type RouterInterface interface {
	GetClient(name string) (*BusinessClient, error)
	Primary() (*BusinessClient, error)
	Read() (*BusinessClient, error)
	Write() (*BusinessClient, error)
	Analytics() (*BusinessClient, error)
	ListClients() []string
	HasClient(name string) bool
	Close() error
}

// Router 数据库路由器，管理多个业务Ent客户端
type Router struct {
	clients map[string]*BusinessClient
	builder *ClientBuilder
	logger  *zap.Logger
	mu      sync.RWMutex
}

// NewRouter 创建数据库路由器
func NewRouter(manager mysql.ManagerInterface, logger *zap.Logger) (*Router, error) {
	builder := NewClientBuilder(manager, logger)

	router := &Router{
		clients: make(map[string]*BusinessClient),
		builder: builder,
		logger:  logger,
	}

	// 为每个数据库连接创建业务Ent客户端
	clientNames := manager.ListClients()
	for _, name := range clientNames {
		businessClient, err := builder.BuildClient(name)
		if err != nil {
			// 清理已创建的客户端
			router.closeAllClients()
			return nil, fmt.Errorf("failed to create business client for '%s': %w", name, err)
		}

		router.clients[name] = businessClient
	}

	logger.Info("Database router initialized successfully",
		zap.Int("client_count", len(router.clients)),
		zap.Strings("clients", clientNames))

	return router, nil
}

// GetClient 获取指定名称的业务数据库客户端
func (r *Router) GetClient(name string) (*BusinessClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, exists := r.clients[name]
	if !exists {
		return nil, fmt.Errorf("business database client '%s' not found", name)
	}
	return client, nil
}

// Primary 获取主数据库客户端
func (r *Router) Primary() (*BusinessClient, error) {
	return r.GetClient(mysql.DB1)
}

// Read 获取读数据库客户端（读写分离场景）
func (r *Router) Read() (*BusinessClient, error) {
	// 优先使用read客户端，如果没有则使用db1
	if client, err := r.GetClient(mysql.ReadDB); err == nil {
		return client, nil
	}
	return r.Primary()
}

// Write 获取写数据库客户端（读写分离场景）
func (r *Router) Write() (*BusinessClient, error) {
	// 优先使用write客户端，如果没有则使用db1
	if client, err := r.GetClient(mysql.WriteDB); err == nil {
		return client, nil
	}
	return r.Primary()
}

// Analytics 获取分析数据库客户端（如果配置了的话）
func (r *Router) Analytics() (*BusinessClient, error) {
	return r.GetClient(mysql.DB2)
}

// ListClients 列出所有可用的数据库客户端名称
func (r *Router) ListClients() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.clients))
	for name := range r.clients {
		names = append(names, name)
	}
	return names
}

// HasClient 检查是否存在指定名称的客户端
func (r *Router) HasClient(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.clients[name]
	return exists
}

// Close 关闭所有客户端连接
func (r *Router) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.closeAllClients()
}

// closeAllClients 关闭所有客户端连接（内部方法，不加锁）
func (r *Router) closeAllClients() error {
	var lastErr error
	for name, client := range r.clients {
		if err := client.Close(); err != nil {
			r.logger.Error("Failed to close business client",
				zap.String("name", name),
				zap.Error(err))
			lastErr = err
		}
	}
	r.clients = make(map[string]*BusinessClient)
	return lastErr
}

// MigrateAll 执行所有客户端的数据库迁移
func (r *Router) MigrateAll(ctx context.Context) error {
	r.mu.RLock()
	clients := make([]*BusinessClient, 0, len(r.clients))
	for _, client := range r.clients {
		clients = append(clients, client)
	}
	r.mu.RUnlock()

	for _, client := range clients {
		if err := client.Migrate(ctx); err != nil {
			return fmt.Errorf("migration failed for client '%s': %w", client.Name(), err)
		}
	}

	r.logger.Info("All database migrations completed successfully")
	return nil
}

// GetGenClient 获取原生 gen.Client（用于向后兼容）
func (r *Router) GetGenClient(name string) (*gen.Client, error) {
	businessClient, err := r.GetClient(name)
	if err != nil {
		return nil, err
	}
	return businessClient.Query(), nil
}

// PrimaryGenClient 获取主数据库的原生 gen.Client（用于向后兼容）
func (r *Router) PrimaryGenClient() (*gen.Client, error) {
	return r.GetGenClient(mysql.DB1)
}
