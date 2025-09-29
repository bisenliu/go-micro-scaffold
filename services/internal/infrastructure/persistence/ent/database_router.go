package ent

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"common/databases/mysql"
)

// RouterInterface 数据库路由器接口
type RouterInterface interface {
	GetClient(name string) (*BusinessClient, error)
	Primary() (*BusinessClient, error)
	Close() error
}

// Router 数据库路由器
type Router struct {
	clients sync.Map
	builder *ClientBuilder
	logger  *zap.Logger
}

// NewRouter 创建数据库路由器
func NewRouter(manager mysql.ManagerInterface, logger *zap.Logger) (*Router, error) {
	builder := NewClientBuilder(manager, logger)

	router := &Router{
		builder: builder,
		logger:  logger,
	}

	// 为每个数据库连接创建业务Ent客户端
	clientNames := manager.ListClients()
	for _, name := range clientNames {
		businessClient, err := builder.BuildClient(name)
		if err != nil {
			router.Close()
			return nil, fmt.Errorf("failed to create business client for '%s': %w", name, err)
		}

		router.clients.Store(name, businessClient)
	}

	logger.Info("Database router initialized successfully",
		zap.Int("client_count", len(clientNames)),
		zap.Strings("clients", clientNames))

	return router, nil
}

// GetClient 获取指定名称的业务数据库客户端
func (r *Router) GetClient(name string) (*BusinessClient, error) {
	if client, ok := r.clients.Load(name); ok {
		return client.(*BusinessClient), nil
	}
	return nil, fmt.Errorf("business database client '%s' not found", name)
}

// Primary 获取主数据库客户端
func (r *Router) Primary() (*BusinessClient, error) {
	return r.GetClient(mysql.DB1)
}

// Close 关闭所有客户端连接
func (r *Router) Close() error {
	var lastErr error
	r.clients.Range(func(key, value interface{}) bool {
		client := value.(*BusinessClient)
		if err := client.Close(); err != nil {
			r.logger.Error("Failed to close business client",
				zap.String("name", key.(string)),
				zap.Error(err))
			lastErr = err
		}
		return true
	})
	return lastErr
}

// MigrateAll 执行所有客户端的数据库迁移
func (r *Router) MigrateAll(ctx context.Context) error {
	var clients []*BusinessClient
	r.clients.Range(func(key, value interface{}) bool {
		clients = append(clients, value.(*BusinessClient))
		return true
	})

	for _, client := range clients {
		if err := client.Migrate(ctx); err != nil {
			return fmt.Errorf("migration failed for client '%s': %w", client.Name(), err)
		}
	}

	r.logger.Info("All database migrations completed successfully")
	return nil
}
