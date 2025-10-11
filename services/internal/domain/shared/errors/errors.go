package errors

import (
	"errors"
	"fmt"
)

// 领域错误定义
var (
	// 通用错误
	ErrNotFound            = errors.New("资源不存在")
	ErrAlreadyExists       = errors.New("资源已存在")
	ErrInvalidData         = errors.New("无效的数据")
	ErrValidationFailed    = errors.New("验证失败")
	ErrBusinessRuleViolation = errors.New("业务规则违反")
	ErrConcurrencyConflict   = errors.New("并发冲突")
	ErrResourceLocked        = errors.New("资源已锁定")

	// 认证授权错误
	ErrUnauthorized        = errors.New("未授权")
	ErrForbidden          = errors.New("禁止访问")
)

// WrapErr 包装错误信息
func WrapErr(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// IsNotFound 检查错误是否为"不存在"错误
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists 检查错误是否为"已存在"错误
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsValidationError 检查错误是否为验证错误
func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidationFailed)
}