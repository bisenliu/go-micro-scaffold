package jwt

import (
	"go.uber.org/fx"

	"common/config"
)

// NewJWTService 创建JWT实例
func NewJWTService(cfg *config.Config) *JWT {
	return NewJWT(cfg)
}

// Module JWT模块
var Module = fx.Module("jwt",
	fx.Provide(NewJWTService),
)
