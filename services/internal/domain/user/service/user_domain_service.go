package service

import (
	"context"

	"services/internal/domain/user/entity"
	userErrors "services/internal/domain/user/errors"
	"services/internal/domain/user/repository"
)

var (
	ErrPhoneAlreadyExists = userErrors.ErrPhoneAlreadyExists
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
func (s *UserDomainService) CreateUser(ctx context.Context, openID, name, phoneNumber, password string, gender int) (*entity.User, error) {

	// 检查手机号是否已绑定
	exists, err := s.userRepo.ExistsByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrPhoneAlreadyExists
	}

	user := entity.NewUser(openID, name, phoneNumber, password, gender)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
