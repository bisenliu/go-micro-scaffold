package application

import (
	"go.uber.org/fx"

	"user-services/internal/application/commandhandler"
	"user-services/internal/application/queryhandler"
	"user-services/internal/application/service"
)

// ApplicationModule 应用模块
var ApplicationModule = fx.Module("application",
	fx.Provide(
		// 命令处理器
		commandhandler.NewUserCommandHandler,

		// 查询处理器
		queryhandler.NewUserQueryHandler,

		// 应用服务
		service.NewPermissionService,
		service.NewAuthService,
	),
)
