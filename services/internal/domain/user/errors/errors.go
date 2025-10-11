package errors

import (
	"fmt"

	domainerrors "services/internal/domain/shared/errors"
)

// 用户相关错误
var (
	ErrUserNotFound       = fmt.Errorf("%w: %s", domainerrors.ErrNotFound, "用户不存在")
	ErrUserAlreadyExists  = fmt.Errorf("%w: %s", domainerrors.ErrAlreadyExists, "用户已存在")
	ErrUserInactive       = fmt.Errorf("%w: %s", domainerrors.ErrInvalidData, "用户已停用")
	ErrPhoneAlreadyExists = fmt.Errorf("%w: %s", domainerrors.ErrAlreadyExists, "手机号已存在")
)

// 用户验证错误
var (
	ErrInvalidPhone     = fmt.Errorf("%w: %s", domainerrors.ErrValidationFailed, "无效的手机号")
	ErrInvalidGender    = fmt.Errorf("%w: %s", domainerrors.ErrValidationFailed, "无效的性别")
	ErrInvalidNickname  = fmt.Errorf("%w: %s", domainerrors.ErrValidationFailed, "无效的昵称")
	ErrPhoneRequired    = fmt.Errorf("%w: %s", domainerrors.ErrValidationFailed, "手机号必填")
	ErrNicknameRequired = fmt.Errorf("%w: %s", domainerrors.ErrValidationFailed, "昵称必填")
)

// 用户业务规则错误
var (
	ErrUserCannotJoinTeam    = fmt.Errorf("%w: %s", domainerrors.ErrBusinessRuleViolation, "用户当前状态无法加入团队")
	ErrUserProfileIncomplete = fmt.Errorf("%w: %s", domainerrors.ErrBusinessRuleViolation, "用户资料不完整")
)
