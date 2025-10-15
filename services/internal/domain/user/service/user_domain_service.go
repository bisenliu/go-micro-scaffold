package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"services/internal/domain/user/entity"
	userErrors "services/internal/domain/user/errors"
	"services/internal/domain/user/repository"
	"services/internal/domain/user/validator"
)

var (
	ErrPhoneAlreadyExists = userErrors.ErrPhoneAlreadyExists
)

// UserDomainService 用户领域服务
type UserDomainService struct {
	userRepo      repository.UserRepository
	userValidator validator.UserValidator
}

// NewUserDomainService 创建用户领域服务
func NewUserDomainService(userRepo repository.UserRepository, userValidator validator.UserValidator) *UserDomainService {
	return &UserDomainService{
		userRepo:      userRepo,
		userValidator: userValidator,
	}
}

// CreateUser 创建用户
func (s *UserDomainService) CreateUser(ctx context.Context, openID, name, phoneNumber, password string, gender int) (*entity.User, error) {
	if err := s.userValidator.ValidateForCreation(ctx, phoneNumber, password, name, gender); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, userErrors.ErrPasswordHashingFailed
	}

	user := entity.NewUser(openID, name, phoneNumber, string(hashedPassword), gender)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
