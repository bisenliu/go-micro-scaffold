package http

import (
	"go.uber.org/fx"

	commonMiddleware "common/middleware"
	service "services/internal/application/service"
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

		// 创建 Casbin 中间件的 Provider
		func(permissionService service.PermissionServiceInterface) routes.CasbinMiddleware {
			return routes.CasbinMiddleware(commonMiddleware.CasbinMiddleware(permissionService.Enforce))
		},
	),

	// 在应用启动时调用，用于设置路由
	fx.Invoke(RegisterServerLifecycle),
	fx.Invoke(routes.SetupRoutesFinal),
)
