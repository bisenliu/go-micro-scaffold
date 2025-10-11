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
	// 手机号
	ErrInvalidPhone     = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "无效的手机号")
	ErrPhoneRequired    = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "手机号必填")
	// 姓名
	ErrInvalidNickname  = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "无效的昵称")
	ErrNicknameRequired = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "昵称必填")
	ErrNameTooLong      = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "姓名长度不能超过50个字符")
	ErrInvalidNameFormat= domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "姓名格式不正确")
	// 密码
	ErrPasswordRequired = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "密码不能为空")
	ErrPasswordTooShort = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "密码长度不能少于6位")
	ErrPasswordTooLong  = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "密码长度不能超过20位")
	ErrPasswordTooWeak  = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "密码强度不够，需要包含字母和数字")
	// 性别
	ErrInvalidGender    = domainerrors.NewDomainError(domainerrors.ErrValidationFailed, "无效的性别")
)

// 用户业务规则错误
var (
	ErrUserCannotJoinTeam    = domainerrors.NewDomainError(domainerrors.ErrBusinessRuleViolation, "用户当前状态无法加入团队")
	ErrUserProfileIncomplete = domainerrors.NewDomainError(domainerrors.ErrBusinessRuleViolation, "用户资料不完整")
)