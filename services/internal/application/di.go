package application

import (
	"go.uber.org/fx"

	"services/internal/application/commandhandler"
	"services/internal/application/queryhandler"
)

// ApplicationModule 应用模块
var ApplicationModule = fx.Module("application",
	fx.Provide(
		// 命令处理器
		commandhandler.NewUserCommandHandler,

		// 查询处理器
		queryhandler.NewUserQueryHandler,
	),
)
