package idgen

import (
	"go.uber.org/fx"

	"common/interfaces"
)

// NewGenerator 创建ID生成器实例
func NewGenerator(configProvider interfaces.ConfigProvider) (Generator, error) {
	return NewSnowflakeGenerator(configProvider)
}

// Module ID生成器模块
var Module = fx.Module("idgen",
	fx.Provide(NewGenerator),
)
