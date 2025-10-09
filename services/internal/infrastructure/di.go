package infrastructure

import (
	"go.uber.org/fx"

	commonDI "common/di"
	"services/internal/infrastructure/messaging"
	"services/internal/infrastructure/persistence/ent"
)

// InfrastructureModule 基础设施模块
var InfrastructureModule = fx.Module("infrastructure",
	// 包含 Ent 模块
	ent.Module,

	// Casbin 模块
	commonDI.CasbinModule,

	// 基础设施服务
	fx.Provide(
		// 消息发布
		messaging.NewRedisEventPublisher,
	),
)
