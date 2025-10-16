package jwt

import (
	"go.uber.org/fx"

	"common/interfaces"
)

// Module JWT模块
var Module = fx.Module("jwt",
	fx.Provide(
		fx.Annotate(
			NewJWTService,
			fx.As(new(interfaces.JWTService)),
		),
	),
)