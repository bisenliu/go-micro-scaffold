package errors

import (
	"errors"
	"fmt"
)

// DomainError 是我们自定义的领域错误结构体
// 它实现了 Go 的 error 接口
type DomainError struct {
	// BaseErr 是被包装的底层错误，用于 errors.Is 判断
	BaseErr error
	// Message 是纯净的、用于API响应的顶层消息
	Message string
}

// Error 返回完整的错误信息，通常用于日志
func (e *DomainError) Error() string {
	if e.BaseErr != nil {
		return fmt.Sprintf("%s: %s", e.BaseErr.Error(), e.Message)
	}
	return e.Message
}

// Unwrap 返回被包装的底层错误，以支持 errors.Is
func (e *DomainError) Unwrap() error {
	return e.BaseErr
}

// NewDomainError 创建一个新的 DomainError 实例
func NewDomainError(baseErr error, message string) *DomainError {
	return &DomainError{BaseErr: baseErr, Message: message}
}

// 定义基础错误类型，它们将作为 BaseErr 被包装
var (
	ErrNotFound            = errors.New("资源不存在")
	ErrAlreadyExists       = errors.New("资源已存在")
	ErrInvalidData         = errors.New("无效的数据")
	ErrValidationFailed    = errors.New("验证失败")
	ErrBusinessRuleViolation = errors.New("业务规则违反")
	ErrConcurrencyConflict   = errors.New("并发冲突")
	ErrResourceLocked        = errors.New("资源已锁定")
	ErrUnauthorized        = errors.New("未授权")
	ErrForbidden           = errors.New("禁止访问")
)
