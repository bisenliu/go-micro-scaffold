package validator

import (
	"context"
	"regexp"
	"strings"

	userErrors "user-services/internal/domain/user/errors"
	"user-services/internal/domain/user/repository"
)

// PhoneValidator 手机号验证器接口
type PhoneValidator interface {
	Validate(phoneNumber string) error
	CheckUniqueness(ctx context.Context, phoneNumber string) error
}

// phoneValidator 手机号验证器实现
type phoneValidator struct {
	userRepo repository.UserRepository
}

// NewPhoneValidator 创建手机号验证器
func NewPhoneValidator(userRepo repository.UserRepository) PhoneValidator {
	return &phoneValidator{
		userRepo: userRepo,
	}
}

// Validate 验证手机号格式
func (v *phoneValidator) Validate(phoneNumber string) error {
	originalPhone := phoneNumber // 保留原始输入用于上下文

	// 去除空格
	phoneNumber = strings.TrimSpace(phoneNumber)

	// 检查是否为空
	if phoneNumber == "" {
		return userErrors.ErrPhoneRequired.WithContext("input", originalPhone)
	}

	// 检查长度
	if len(phoneNumber) > 11 {
		return userErrors.ErrInvalidPhone.
			WithContext("input", originalPhone).
			WithContext("length", len(originalPhone)).
			WithContext("max_length", 11).
			WithContext("rule", "length")
	}

	// 检查格式（中国大陆手机号）
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(phoneNumber) {
		return userErrors.ErrInvalidPhone.
			WithContext("input", originalPhone).
			WithContext("rule", "format").
			WithContext("expected", "1[3-9]xxxxxxxxx")
	}

	return nil
}

// CheckUniqueness 检查手机号唯一性
func (v *phoneValidator) CheckUniqueness(ctx context.Context, phoneNumber string) error {
	exists, err := v.userRepo.ExistsByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return err
	}

	if exists {
		return userErrors.ErrPhoneNotUnique.
			WithContext("input", phoneNumber)
	}

	return nil
}
