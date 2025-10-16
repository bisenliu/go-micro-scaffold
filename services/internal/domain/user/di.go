package user

import (
	"go.uber.org/fx"

	"services/internal/domain/user/service"
	"services/internal/domain/user/validator"
)

// DomainModule 领域模块
// 只提供领域服务和验证器，不依赖任何外层实现
var DomainModule = fx.Module("user_domain",
	fx.Provide(
		// 验证器
		validator.NewUserValidator,

		// 领域服务
		service.NewUserDomainService,
	),
)
