package user

import (
	"go.uber.org/fx"

	domainrepo "user-services/internal/domain/user/repository"
	"user-services/internal/domain/user/service"
	"user-services/internal/domain/user/validator"
	entrepo "user-services/internal/infrastructure/persistence/ent/repository"
)

// DomainModule 领域模块
var DomainModule = fx.Module("domain",
	fx.Provide(
		// 验证器
		validator.NewUserValidator,

		// 领域服务
		service.NewUserDomainService,

		// 仓储实现
		fx.Annotate(
			entrepo.NewUserRepository,
			fx.As(new(domainrepo.UserRepository)),
		),
	),
)
