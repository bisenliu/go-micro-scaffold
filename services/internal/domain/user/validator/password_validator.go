package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// PasswordValidator 密码验证器接口
type PasswordValidator interface {
	Validate(password string) error
	ValidateStrength(password string) error
}

// passwordValidator 密码验证器实现
type passwordValidator struct{}

// NewPasswordValidator 创建密码验证器
func NewPasswordValidator() PasswordValidator {
	return &passwordValidator{}
}

// Validate 验证密码基本规则
func (v *passwordValidator) Validate(password string) error {
	// 去除空格
	password = strings.TrimSpace(password)

	// 检查是否为空
	if password == "" {
		return ErrPasswordRequired
	}

	// 检查长度
	length := utf8.RuneCountInString(password)
	if length < 6 {
		return ErrPasswordTooShort
	}
	if length > 20 {
		return ErrPasswordTooLong
	}

	// 验证密码强度
	if err := v.ValidateStrength(password); err != nil {
		return err
	}

	return nil
}

// ValidateStrength 验证密码强度
func (v *passwordValidator) ValidateStrength(password string) error {
	// 至少包含一个字母
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	// 至少包含一个数字
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)

	if !hasLetter || !hasNumber {
		return ErrPasswordTooWeak
	}

	return nil
}
