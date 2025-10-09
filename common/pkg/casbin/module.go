package casbin

import (
	"github.com/casbin/casbin/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/databases/mysql"
)

// Module 提供了 Casbin Enforcer
var Module = fx.Module("casbin",
	// 提供一个函数来创建Casbin执行器，该函数接收MySQL管理器并返回执行器
	fx.Provide(func(manager mysql.ManagerInterface, logger *zap.Logger) (*casbin.SyncedCachedEnforcer, error) {
		// 获取主数据库客户端
		client, err := manager.Primary()
		if err != nil {
			return nil, err
		}

		// 使用主数据库客户端创建Casbin执行器
		return NewEnforcer(client, logger)
	}),
)
