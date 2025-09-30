package validator

import (
	"context"

	"services/internal/domain/user/repository"
)

// UserValidator 用户验证器接口
type UserValidator interface {
	ValidateForCreation(ctx context.Context, phoneNumber, password, name string, gender int) error
	ValidateForUpdate(ctx context.Context, userID string, updates map[string]interface{}) error
	ValidatePhoneNumber(phoneNumber string) error
	ValidatePassword(password string) error
	ValidateName(name string) error
}

// userValidator 用户验证器实现
type userValidator struct {
	phoneValidator    PhoneValidator
	passwordValidator PasswordValidator
	nameValidator     NameValidator
}

// NewUserValidator 创建用户验证器
func NewUserValidator(userRepo repository.UserRepository) UserValidator {
	return &userValidator{
		phoneValidator:    NewPhoneValidator(userRepo),
		passwordValidator: NewPasswordValidator(),
		nameValidator:     NewNameValidator(),
	}
}

// ValidateForCreation 验证用户创建
func (v *userValidator) ValidateForCreation(ctx context.Context, phoneNumber, password, name string, gender int) error {
	// 验证手机号格式
	if err := v.phoneValidator.Validate(phoneNumber); err != nil {
		return err
	}

	// 验证手机号唯一性
	if err := v.phoneValidator.CheckUniqueness(ctx, phoneNumber); err != nil {
		return err
	}

	// 验证密码
	// if err := v.passwordValidator.Validate(password); err != nil {
	// 	return err
	// }

	// 验证姓名
	if err := v.nameValidator.Validate(name); err != nil {
		return err
	}

	return nil
}

// ValidateForUpdate 验证用户更新
func (v *userValidator) ValidateForUpdate(ctx context.Context, userID string, updates map[string]interface{}) error {
	// 根据更新字段进行相应验证
	if phoneNumber, ok := updates["phone_number"].(string); ok {
		if err := v.phoneValidator.Validate(phoneNumber); err != nil {
			return err
		}
		// 注意：更新时需要排除当前用户的手机号
		if err := v.phoneValidator.CheckUniqueness(ctx, phoneNumber); err != nil {
			return err
		}
	}

	if password, ok := updates["password"].(string); ok {
		if err := v.passwordValidator.Validate(password); err != nil {
			return err
		}
	}

	if name, ok := updates["name"].(string); ok {
		if err := v.nameValidator.Validate(name); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePhoneNumber 验证手机号
func (v *userValidator) ValidatePhoneNumber(phoneNumber string) error {
	return v.phoneValidator.Validate(phoneNumber)
}

// ValidatePassword 验证密码
func (v *userValidator) ValidatePassword(password string) error {
	return v.passwordValidator.Validate(password)
}

// ValidateName 验证姓名
func (v *userValidator) ValidateName(name string) error {
	return v.nameValidator.Validate(name)
}
