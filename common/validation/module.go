package validation

import (
	"common/config"

	"go.uber.org/fx"
)

// NewValidatorForDI 为依赖注入创建验证器（从配置读取语言设置）
func NewValidatorForDI(cfg *config.Config) *Validator {
	validator, _ := NewValidator(cfg.Validation.Locale)
	return validator
}

// Module 验证模块
var Module = fx.Module("validation",
	fx.Provide(NewValidatorForDI), // 直接提供，不使用 fx.Annotated
)
