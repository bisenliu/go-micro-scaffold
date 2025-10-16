package http

import (
	"go.uber.org/fx"

	"services/internal/interfaces/http/handler"
	"services/internal/interfaces/http/routes"
)

var InterfaceModuleFinal = fx.Module("interface_final",
	fx.Provide(
		// HTTP Handlers
		handler.NewUserHandler,
		handler.NewHealthHandler,
		handler.NewAuthHandler,

		// HTTP Server
		NewServer,

		// Middleware
		NewCasbinMiddleware,
		NewAuthMiddleware,
	),

	// 在应用启动时调用，用于设置路由
	fx.Invoke(RegisterServerLifecycle),
	fx.Invoke(routes.SetupRoutesFinal),
)
