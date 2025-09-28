package validation

import (
	"common/config"

	"go.uber.org/fx"
)

// NewValidator 为依赖注入创建验证器
func NewValidator(cfg *config.Config) (*Validator, error) {
	validator, err := NewLocalizedValidator(cfg.Validation.Locale)
	if err != nil {
		return nil, err
	}
	return validator, nil
}

// Module 验证模块
var Module = fx.Module("validation",
	fx.Provide(NewValidator),
)
