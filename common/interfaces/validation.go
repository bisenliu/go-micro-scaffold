package interfaces

import (
	"context"
)

// Validator 验证器接口
// 定义了数据验证的统一接口
type Validator interface {
	// Validate 验证结构体
	Validate(ctx context.Context, data interface{}) error
	
	// ValidateStruct 验证结构体并返回详细错误信息
	ValidateStruct(ctx context.Context, data interface{}) ValidationErrors
	
	// ValidateField 验证单个字段
	ValidateField(ctx context.Context, field interface{}, tag string) error
	
	// RegisterValidation 注册自定义验证规则
	RegisterValidation(tag string, fn ValidationFunc) error
	
	// RegisterTranslation 注册错误消息翻译
	RegisterTranslation(tag string, message string) error
}

// ValidationFunc 自定义验证函数类型
type ValidationFunc func(value interface{}) bool

// ValidationError 单个验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors 验证错误集合
type ValidationErrors []ValidationError

// Error 实现error接口
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	return ve[0].Message
}

// HasErrors 检查是否有错误
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// GetFieldError 获取指定字段的错误
func (ve ValidationErrors) GetFieldError(field string) *ValidationError {
	for _, err := range ve {
		if err.Field == field {
			return &err
		}
	}
	return nil
}

// ValidatorFactory 验证器工厂接口
type ValidatorFactory interface {
	// CreateValidator 创建验证器
	CreateValidator(config ValidationConfig) (Validator, error)
}