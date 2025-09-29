package ent

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module 业务层 Ent 模块
var Module = fx.Module("business_ent",
	// 提供客户端构建器
	fx.Provide(NewClientBuilder),

	// 提供数据库路由器
	fx.Provide(NewRouter),

	// 提供路由器接口的实现
	fx.Provide(func(router *Router) RouterInterface {
		return router
	}),

	// 提供主业务客户端（常用的快捷方式）
	fx.Provide(func(router RouterInterface) (*BusinessClient, error) {
		return router.Primary()
	}),

	// 生命周期管理
	fx.Invoke(func(lc fx.Lifecycle, router *Router, logger *zap.Logger) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				// 可选：在启动时执行数据库迁移
				// return router.MigrateAll(ctx)
				logger.Info("Business Ent module started successfully")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Shutting down business Ent module")
				return router.Close()
			},
		})
	}),
)

// MigrationModule 数据库迁移模块（可选）
var MigrationModule = fx.Module("ent_migration",
	fx.Invoke(func(lc fx.Lifecycle, router RouterInterface, logger *zap.Logger) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("Starting database migration")
				if err := router.(*Router).MigrateAll(ctx); err != nil {
					logger.Error("Database migration failed", zap.Error(err))
					return err
				}
				logger.Info("Database migration completed successfully")
				return nil
			},
		})
	}),
)
