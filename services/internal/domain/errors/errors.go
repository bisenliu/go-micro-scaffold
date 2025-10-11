package errors

import (
	"errors"

	infraerrors "services/internal/infrastructure/errors"
)

// 核心业务错误
var (
	// 用户相关错误
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserAlreadyExists  = errors.New("用户已存在")
	ErrUserInactive       = errors.New("用户已停用")
	ErrPhoneAlreadyExists = errors.New("手机号已存在")
	ErrInvalidUserData    = errors.New("无效的用户数据")

	// 验证相关错误
	ErrValidationFailed = errors.New("验证失败")
	ErrInvalidPhone     = errors.New("无效的手机号")
	ErrInvalidGender    = errors.New("无效的性别")
	ErrInvalidNickname  = errors.New("无效的昵称")
	ErrPhoneRequired    = errors.New("手机号必填")
	ErrNicknameRequired = errors.New("昵称必填")

	// 业务规则错误
	ErrBusinessRuleViolation = errors.New("业务规则违反")
	ErrConcurrencyConflict   = errors.New("并发冲突")
	ErrResourceLocked        = errors.New("资源已锁定")

	// 认证授权错误
	ErrUnauthorized = errors.New("未授权")
	ErrForbidden    = errors.New("禁止访问")
)

// InfrastructureError 包装基础设施错误，转换为业务错误
func InfrastructureError(err error) error {
	if err == nil {
		return nil
	}

	// 将基础设施错误映射到业务错误
	switch {
	case errors.Is(err, infraerrors.ErrRecordNotFound):
		return ErrUserNotFound
	case errors.Is(err, infraerrors.ErrDuplicateKey):
		return ErrUserAlreadyExists
	default:
		return err
	}
}

// ValidationError 包装验证错误
func ValidationError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Join(ErrValidationFailed, err)
}
