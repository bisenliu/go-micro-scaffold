package errors

import (
	domainerrors "services/internal/domain/shared/errors"
)

// 用户相关错误消息常量
const (
	MsgUserNotFound           = "用户不存在"
	MsgUserAlreadyExists      = "用户已存在"
	MsgPhoneAlreadyExists     = "手机号已存在"
	MsgInvalidUserID          = "无效的用户ID"
	MsgCreateUserFailed       = "创建用户失败"
	MsgUpdateUserFailed       = "更新用户失败"
	MsgQueryUserFailed        = "查询用户失败"
	MsgQueryUserListFailed    = "查询用户列表失败"
	MsgQueryUserCountFailed   = "查询用户总数失败"
	MsgCheckPhoneExistsFailed = "查询用户手机号是否存在失败"
	MsgFindUserByPhoneFailed  = "通过手机号查询用户失败"
)

// 用户相关错误
var (
	ErrUserInactive = domainerrors.NewInvalidDataError("用户已停用")
)

// 用户验证错误
var (
	// 手机号
	ErrInvalidPhone   = domainerrors.NewValidationError("无效的手机号")
	ErrPhoneRequired  = domainerrors.NewValidationError("手机号必填")
	ErrPhoneNotUnique = domainerrors.NewValidationError("手机号已被使用")
	// 姓名
	ErrInvalidNickname   = domainerrors.NewValidationError("无效的昵称")
	ErrNicknameRequired  = domainerrors.NewValidationError("昵称必填")
	ErrNameTooLong       = domainerrors.NewValidationError("姓名长度不能超过50个字符")
	ErrInvalidNameFormat = domainerrors.NewValidationError("姓名格式不正确")
	// 密码
	ErrPasswordRequired      = domainerrors.NewValidationError("密码不能为空")
	ErrPasswordTooShort      = domainerrors.NewValidationError("密码长度不能少于6位")
	ErrPasswordTooLong       = domainerrors.NewValidationError("密码长度不能超过20位")
	ErrPasswordTooWeak       = domainerrors.NewValidationError("密码强度不够，需要包含字母和数字")
	ErrPasswordHashingFailed = domainerrors.NewBusinessRuleViolationError("密码处理失败")
	// 性别
	ErrInvalidGender = domainerrors.NewValidationError("无效的性别")
)

// 用户业务规则错误
var (
	ErrUserCannotJoinTeam    = domainerrors.NewBusinessRuleViolationError("用户当前状态无法加入团队")
	ErrUserProfileIncomplete = domainerrors.NewBusinessRuleViolationError("用户资料不完整")
)
