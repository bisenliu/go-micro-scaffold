package validation

import (
	"go.uber.org/fx"

	"common/interfaces"
)

// NewValidator 为依赖注入创建验证器
func NewValidator(configProvider interfaces.ConfigProvider) (*Validator, error) {
	cfg := configProvider.GetValidationConfig()
	validator, err := NewLocalizedValidator(cfg.Locale)
	if err != nil {
		return nil, err
	}
	return validator, nil
}

// Module 验证模块
var Module = fx.Module("validation",
	fx.Provide(NewValidator),
)
