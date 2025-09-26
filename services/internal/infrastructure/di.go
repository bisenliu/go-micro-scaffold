package infrastructure

import (
	"go.uber.org/fx"

	"services/internal/infrastructure/messaging"
	"services/internal/infrastructure/persistence/ent"
	"services/internal/infrastructure/persistence/ent/gen"
	"services/internal/infrastructure/validation"
)

// InfrastructureModule 基础设施模块
var InfrastructureModule = fx.Module("infrastructure",
	// 提供Ent客户端工厂
	fx.Provide(
		// Ent客户端工厂
		ent.NewEntClientFactory,

		// 数据库路由器
		ent.NewDatabaseRouter,

		// 向后兼容：提供主数据库客户端
		func(router *ent.DatabaseRouter) (*gen.Client, error) {
			return router.Primary()
		},
	),

	// 基础设施服务
	fx.Provide(
		// 验证器
		validation.NewUserInfrastructureValidator,

		// 消息发布
		messaging.NewRedisEventPublisher,
	),
)
