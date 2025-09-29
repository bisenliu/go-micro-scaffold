package mysql

import (
	"context"

	"go.uber.org/fx"
)

// ManagerInterface 数据库管理器接口
type ManagerInterface interface {
	GetClient(name string) (ClientInterface, error)
	Primary() (ClientInterface, error)
	Read() (ClientInterface, error)
	Write() (ClientInterface, error)
	ListClients() []string
	HasClient(name string) bool
	Close() error
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
