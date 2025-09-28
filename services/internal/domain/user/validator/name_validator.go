package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// NameValidator 姓名验证器接口
type NameValidator interface {
	Validate(name string) error
}

// nameValidator 姓名验证器实现
type nameValidator struct{}

// NewNameValidator 创建姓名验证器
func NewNameValidator() NameValidator {
	return &nameValidator{}
}

// Validate 验证姓名
func (v *nameValidator) Validate(name string) error {
	// 去除首尾空格
	name = strings.TrimSpace(name)

	// 检查是否为空
	if name == "" {
		return ErrNameRequired
	}

	// 检查长度（按字符数计算，支持中文）
	length := utf8.RuneCountInString(name)
	if length > 50 {
		return ErrNameTooLong
	}

	// 检查格式：只允许中文、英文字母、数字和常见符号
	nameRegex := regexp.MustCompile(`^[\p{Han}a-zA-Z0-9\s\-_.·]+$`)
	if !nameRegex.MatchString(name) {
		return ErrInvalidNameFormat
	}

	return nil
}