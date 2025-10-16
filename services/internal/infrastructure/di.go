package infrastructure

import (
	"go.uber.org/fx"

	domainrepo "services/internal/domain/user/repository"
	"services/internal/infrastructure/mapper"
	"services/internal/infrastructure/messaging"
	"services/internal/infrastructure/persistence/ent"
	entrepo "services/internal/infrastructure/persistence/ent/repository"
)

// InfrastructureModule 基础设施模块
// 在基础设施层注册所有仓储实现，实现领域层定义的接口
var InfrastructureModule = fx.Module("infrastructure",
	// 包含映射器模块
	mapper.Module,
	
	// 包含 Ent 模块
	ent.Module,

	// 暂时禁用 Casbin 模块，等后续任务重构
	// commonDI.CasbinModule,

	// 基础设施服务
	fx.Provide(
		// 消息发布器（简化版本）
		messaging.NewSimpleEventPublisher,

		// 仓储实现 - 实现领域层定义的接口
		fx.Annotate(
			entrepo.NewUserRepository,
			fx.As(new(domainrepo.UserRepository)),
		),
	),
)
