package validation

import (
	"github.com/go-playground/validator/v10"
)

// Enum 可验证枚举接口
type Enum interface {
	IsValid() bool
}

// ValidateEnum 枚举验证函数
func ValidateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(Enum)
	return value.IsValid()
}
