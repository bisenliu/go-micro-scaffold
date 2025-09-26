package errors

import "errors"

// 用户聚合根相关错误
var (
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserAlreadyExists  = errors.New("用户已存在")
	ErrUserInactive       = errors.New("用户已停用")
	ErrPhoneAlreadyExists = errors.New("手机号已存在")
	ErrInvalidUserData    = errors.New("无效的用户数据")
)

// 用户值对象相关错误
var (
	ErrInvalidPhone    = errors.New("无效的手机号")
	ErrInvalidGender   = errors.New("无效的性别")
	ErrInvalidNickname = errors.New("无效的昵称")
)

// 用户业务规则相关错误
var (
	ErrPhoneRequired         = errors.New("手机号必填")
	ErrNicknameRequired      = errors.New("昵称必填")
	ErrUserCannotJoinTeam    = errors.New("用户当前状态无法加入团队")
	ErrUserProfileIncomplete = errors.New("用户资料不完整")
)
