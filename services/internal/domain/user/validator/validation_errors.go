package validator

import "errors"

// 验证相关错误定义
var (
	// 手机号验证错误
	ErrInvalidPhoneNumber  = errors.New("手机号格式不正确")
	ErrPhoneNumberRequired = errors.New("手机号不能为空")
	ErrPhoneNumberTooLong  = errors.New("手机号长度不能超过11位")

	// 密码验证错误
	ErrPasswordRequired = errors.New("密码不能为空")
	ErrPasswordTooShort = errors.New("密码长度不能少于6位")
	ErrPasswordTooLong  = errors.New("密码长度不能超过20位")
	ErrPasswordTooWeak  = errors.New("密码强度不够，需要包含字母和数字")

	// 姓名验证错误
	ErrNameRequired      = errors.New("姓名不能为空")
	ErrNameTooLong       = errors.New("姓名长度不能超过50个字符")
	ErrInvalidNameFormat = errors.New("姓名格式不正确")

	// 性别验证错误
	ErrInvalidGender = errors.New("性别值不正确")
)
