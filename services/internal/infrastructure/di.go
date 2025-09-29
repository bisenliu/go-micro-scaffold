package infrastructure

import (
	"go.uber.org/fx"

	"services/internal/infrastructure/messaging"
	"services/internal/infrastructure/persistence/ent"
)

// InfrastructureModule 基础设施模块
var InfrastructureModule = fx.Module("infrastructure",
	// 包含 Ent 模块
	ent.Module,

	// 基础设施服务
	fx.Provide(
		// 消息发布
		messaging.NewRedisEventPublisher,
	),
)
