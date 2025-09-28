package jwt

import (
	"common/config"

	"go.uber.org/fx"
)

// NewJWTService 创建JWT实例
func NewJWTService(cfg *config.Config) *JWT {
	return NewJWT(cfg)
}

// Module JWT模块
var Module = fx.Module("jwt",
	fx.Provide(NewJWTService),
)
