package http

import (
	"go.uber.org/fx"

	"services/internal/interfaces/http/handler"
)

// InterfaceModuleFinal 最终推荐的接口模块
var InterfaceModuleFinal = fx.Module("interface_final",
	// 处理器
	fx.Provide(
		handler.NewUserHandler,
		handler.NewHealthHandler,
		// 后续添加其他处理器
		// handler.NewOrderHandler,
		// handler.NewProductHandler,
		// handler.NewPaymentHandler,
		// handler.NewTeamHandler,
	),

	// 启动器 - 使用最终推荐的路由设置方案
	fx.Invoke(SetupRoutesFinal),
)
