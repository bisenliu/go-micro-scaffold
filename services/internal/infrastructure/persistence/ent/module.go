package ent

import (
	"go.uber.org/fx"

	entsql "entgo.io/ent/dialect/sql"
	"services/internal/infrastructure/persistence/ent/gen"
	"services/internal/shared/interfaces"
)

// Module Ent 模块
var Module = fx.Module("ent",
	// 提供 gen.Client，基于新的数据库接口
	fx.Provide(func(databaseManager interfaces.DatabaseManager, logger interfaces.Logger) (*gen.Client, error) {
		// 获取主数据库连接
		primaryDB := databaseManager.GetPrimaryDB()
		
		// 创建 Ent SQL 驱动
		driver := entsql.OpenDB("mysql", primaryDB.GetDB())
		
		// 构建 Ent 选项
		entOptions := []gen.Option{
			gen.Driver(driver),
		}

		// 在开发环境启用调试模式
		// TODO: 可以从配置中读取是否启用调试模式
		// if configProvider.GetEnv() == "development" {
		//     logger.Info(nil, "Enabling Ent debug mode")
		//     entOptions = append(entOptions, gen.Debug())
		// }

		// 创建 Ent 客户端
		entClient := gen.NewClient(entOptions...)

		return entClient, nil
	}),
)
