package http

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	commonMiddleware "common/middleware"
	service "services/internal/application/service"
	"services/internal/interfaces/http/handler"
)

var InterfaceModuleFinal = fx.Module("interface_final",
	fx.Provide(
		// HTTP Handlers
		handler.NewUserHandler,
		handler.NewHealthHandler,

		// HTTP Server
		NewServer,

		// 创建 Casbin 中间件的 Provider
		func(permissionService service.PermissionServiceInterface, logger *zap.Logger) CasbinMiddleware {
			return CasbinMiddleware(commonMiddleware.CasbinMiddleware(permissionService.Enforce, logger))
		},
	),

	// 在应用启动时调用，用于设置路由
	fx.Invoke(RegisterServerLifecycle),
	fx.Invoke(SetupRoutesFinal),
)
