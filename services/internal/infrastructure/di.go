package infrastructure

import (
	"go.uber.org/fx"

	"services/internal/infrastructure/messaging"
	"services/internal/infrastructure/persistence/ent"
	"services/internal/infrastructure/persistence/ent/gen"
)

// InfrastructureModule 基础设施模块
var InfrastructureModule = fx.Module("infrastructure",
	// 包含 Ent 模块
	ent.Module,

	// 提供向后兼容的主数据库客户端
	fx.Provide(
		func(router ent.RouterInterface) (*gen.Client, error) {
			client, err := router.Primary()
			if err != nil {
				return nil, err
			}
			return client.Query(), nil
		},
	),

	// 基础设施服务
	fx.Provide(
		// 消息发布
		messaging.NewRedisEventPublisher,
	),
)
