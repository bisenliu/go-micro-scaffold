package errors

import (
	domainerrors "services/internal/domain/shared/errors"
)

// 用户相关错误
var (
	ErrUserNotFound       = domainerrors.NewDomainError(domainerrors.ErrNotFound, "用户不存在")
	ErrUserAlreadyExists  = domainerrors.NewDomainError(domainerrors.ErrAlreadyExists, "用户已存在")
	ErrUserInactive       = domainerrors.NewDomainError(domainerrors.ErrInvalidData, "用户已停用")
	ErrPhoneAlreadyExists = domainerrors.NewDomainError(domainerrors.ErrAlreadyExists, "手机号已存在")
)

// 用户验证错误
var (
	ErrInvalidPhone     = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "无效的手机号")
	ErrInvalidGender    = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "无效的性别")
	ErrInvalidNickname  = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "无效的昵称")
	ErrPhoneRequired    = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "手机号必填")
	ErrNicknameRequired = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "昵称必填")
)

// 用户业务规则错误
var (
	ErrUserCannotJoinTeam    = domainerrors.NewDomainError(domainerrors.ErrBusinessRuleViolation, "用户当前状态无法加入团队")
	ErrUserProfileIncomplete = domainerrors.NewDomainError(domainerrors.ErrBusinessRuleViolation, "用户资料不完整")
)