package rdbms

import (
	"context"

	"go.uber.org/fx"
)

// ManagerInterface 数据库管理器接口
type ManagerInterface interface {
	GetClient(name string) (*Client, error)     // 通过数据库名获取客户端
	GetByAlias(alias string) (*Client, error)   // 通过别名获取客户端
	Default() (*Client, error)                  // 获取默认数据库
	ListClients() []string                      // 列出所有数据库名
	ListAliases() map[string]string             // 列出所有别名映射
	HasClient(name string) bool                 // 检查客户端是否存在
	Close() error                               // 关闭所有连接
}

// Module 数据库模块
var Module = fx.Module("mysql",

	fx.Provide(NewManager),
	fx.Provide(func(manager *Manager) ManagerInterface {
		return manager
	}),
	fx.Invoke(func(lc fx.Lifecycle, manager *Manager) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return manager.Close()
			},
		})
	}),
)
