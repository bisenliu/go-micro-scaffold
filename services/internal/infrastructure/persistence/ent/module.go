package ent

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/databases/mysql"
	"services/internal/infrastructure/persistence/ent/gen"
)

// Module Ent 模块
var Module = fx.Module("ent",
	// 提供 gen.Client，基于 mysql.ManagerInterface
	fx.Provide(func(manager mysql.ManagerInterface, logger *zap.Logger) (*gen.Client, error) {
		// 获取主数据库客户端
		dbClient, err := manager.Primary()
		if err != nil {
			return nil, err
		}

		// 构建 Ent 选项
		entOptions := []gen.Option{
			gen.Driver(dbClient.Driver()),
		}

		// 根据配置启用调试模式
		if dbClient.Config().EnableDebug {
			logger.Info("Enabling Ent debug mode for primary client")
			entOptions = append(entOptions, gen.Debug())
		}

		// 创建 Ent 客户端
		entClient := gen.NewClient(entOptions...)

		return entClient, nil
	}),
)
