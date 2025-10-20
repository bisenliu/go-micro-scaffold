package casbin

import (
	"github.com/casbin/casbin/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/databases/rdbms"
)

// Module 提供了 Casbin Enforcer
var Module = fx.Module("casbin",
	// 提供一个函数来创建Casbin执行器，该函数接收数据库管理器并返回执行器
	fx.Provide(func(manager rdbms.ManagerInterface, logger *zap.Logger) (*casbin.SyncedCachedEnforcer, error) {
		// 尝试通过别名获取Casbin专用数据库，如果没有则使用默认数据库
		var client *rdbms.Client
		var err error

		if casbinClient, aliasErr := manager.GetByAlias("casbin"); aliasErr == nil {
			client = casbinClient
		} else {
			client, err = manager.Default()
			if err != nil {
				return nil, err
			}
		}

		// 使用数据库客户端创建Casbin执行器
		return NewEnforcer(client, logger)
	}),
)
