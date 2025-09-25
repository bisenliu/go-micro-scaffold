package user

import (
	"go.uber.org/fx"

	domainrepo "services/internal/domain/user/repository"
	"services/internal/domain/user/service"
	entrepo "services/internal/infrastructure/persistence/ent/repository"
)

// DomainModule 领域模块
var DomainModule = fx.Module("domain",
	fx.Provide(
		// 领域服务
		service.NewUserDomainService,

		// 仓储实现
		fx.Annotate(
			entrepo.NewUserRepository,
			fx.As(new(domainrepo.UserRepository)),
		),
	),
)
