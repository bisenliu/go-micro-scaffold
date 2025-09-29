package mysql

import (
	"context"
	"database/sql"

	entsql "entgo.io/ent/dialect/sql"
	"go.uber.org/fx"
)

// Module 数据库模块，提供所有数据库相关的依赖注入
var Module = fx.Module("mysql_database",
	// 提供数据库管理器
	fx.Provide(NewManager),

	// 提供管理器接口的实现
	fx.Provide(func(manager *Manager) ManagerInterface {
		return manager
	}),

	// 生命周期管理
	fx.Invoke(func(lc fx.Lifecycle, manager *Manager) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return manager.Close()
			},
		})
	}),
)

// LegacyModule 向后兼容的模块，保持原有的 EntClient 接口
var LegacyModule = fx.Module("mysql_legacy",
	fx.Provide(NewLegacyEntClient),
)

// LegacyEntClient 向后兼容的 EntClient 结构
type LegacyEntClient struct {
	client ClientInterface
}

// NewLegacyEntClient 创建向后兼容的 EntClient
func NewLegacyEntClient(manager ManagerInterface) (*LegacyEntClient, error) {
	client, err := manager.Primary()
	if err != nil {
		return nil, err
	}

	return &LegacyEntClient{
		client: client,
	}, nil
}

// Driver 获取 Ent 数据库驱动（向后兼容）
func (c *LegacyEntClient) Driver() *entsql.Driver {
	return c.client.Driver()
}

// DB 获取原始数据库连接（向后兼容）
func (c *LegacyEntClient) DB() *sql.DB {
	return c.client.DB()
}

// Close 关闭连接（向后兼容）
func (c *LegacyEntClient) Close() error {
	return c.client.Close()
}
