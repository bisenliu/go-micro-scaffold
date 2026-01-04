package validator

import (
	"regexp"
	"strings"

	userErrors "user-services/internal/domain/user/errors"
)

// PasswordValidator 密码验证器接口
type PasswordValidator interface {
	Validate(password string) error
}

// passwordValidator 密码验证器实现
type passwordValidator struct{}

// NewPasswordValidator 创建密码验证器
func NewPasswordValidator() PasswordValidator {
	return &passwordValidator{}
}

// Validate 验证密码
func (v *passwordValidator) Validate(password string) error {
	// 去除空格
	password = strings.TrimSpace(password)

	// 检查是否为空
	if password == "" {
		return userErrors.ErrPasswordRequired
	}

	// 检查长度
	if len(password) < 6 {
		return userErrors.ErrPasswordTooShort
	}
	if len(password) > 20 {
		return userErrors.ErrPasswordTooLong
	}

	// 检查强度（示例：必须包含字母和数字）
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasNumber {
		return userErrors.ErrPasswordTooWeak
	}

	return nil
}
