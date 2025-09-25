package service

import (
	"context"

	"services/internal/domain/user/entity"
	userErrors "services/internal/domain/user/errors"
	"services/internal/domain/user/repository"
)

var (
	ErrUserNotFound      = userErrors.ErrUserNotFound
	ErrUserAlreadyExists = userErrors.ErrUserAlreadyExists
	ErrInvalidUserData   = userErrors.ErrInvalidUserData
)

// UserDomainService 用户领域服务
type UserDomainService struct {
	userRepo repository.UserRepository
}

// NewUserDomainService 创建用户领域服务
func NewUserDomainService(userRepo repository.UserRepository) *UserDomainService {
	return &UserDomainService{
		userRepo: userRepo,
	}
}

// CreateUser 创建用户
func (s *UserDomainService) CreateUser(ctx context.Context, openID, unionID string) (*entity.User, error) {

	return nil, nil
}
