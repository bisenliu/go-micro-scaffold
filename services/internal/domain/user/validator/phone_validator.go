package validator

import (
	"context"
	"regexp"
	"strings"

	userErrors "services/internal/domain/user/errors"
	"services/internal/domain/user/repository"
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
	// 去除空格
	phoneNumber = strings.TrimSpace(phoneNumber)

	// 检查是否为空
	if phoneNumber == "" {
		return ErrPhoneNumberRequired
	}

	// 检查长度
	if len(phoneNumber) > 11 {
		return ErrPhoneNumberTooLong
	}

	// 检查格式（中国大陆手机号）
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(phoneNumber) {
		return ErrInvalidPhoneNumber
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
		return userErrors.ErrPhoneAlreadyExists
	}

	return nil
}
