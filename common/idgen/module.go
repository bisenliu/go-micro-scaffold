package idgen

import (
	"common/config"

	"go.uber.org/fx"
)

// NewGenerator 创建ID生成器实例
func NewGenerator(cfg *config.Config) (Generator, error) {
	return NewSnowflakeGenerator(cfg)
}

// Module ID生成器模块
var Module = fx.Module("idgen",
	fx.Provide(NewGenerator),
)
