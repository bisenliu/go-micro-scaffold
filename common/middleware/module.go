package middleware

import (
	"go.uber.org/fx"

	"common/interfaces"
)

// Module 中间件模块
var Module = fx.Module("middleware",
	fx.Provide(
		fx.Annotate(
			NewMiddlewareProvider,
			fx.As(new(interfaces.MiddlewareProvider)),
		),
	),
)