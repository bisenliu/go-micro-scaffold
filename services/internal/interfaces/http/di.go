package http

import (
	"go.uber.org/fx"

	"services/internal/interfaces/http/handler"
)

// InterfaceModuleFinal
var InterfaceModuleFinal = fx.Module("interface_final",
	// 处理器
	fx.Provide(
		handler.NewUserHandler,
		handler.NewHealthHandler,
		// 后续添加其他处理器
	),

	// HTTP服务器相关组件
	fx.Provide(NewServer),
	fx.Invoke(RegisterServerLifecycle),

	// 启动器
	fx.Invoke(SetupRoutesFinal),
)
