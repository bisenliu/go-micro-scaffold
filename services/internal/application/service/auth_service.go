package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	domainerrors "services/internal/domain/shared/errors"
	"services/internal/domain/user/repository"
)

// AuthServiceInterface 认证服务接口
type AuthServiceInterface interface {
	LoginByPassword(ctx context.Context, phoneNumber, password string) (string, string, error)
	LoginByWeChat(ctx context.Context, code string) (string, string, error)
}

// AuthService 认证服务
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository) AuthServiceInterface {
	return &AuthService{
		userRepo: userRepo,
	}
}

// LoginByPassword 账号密码登录
func (s *AuthService) LoginByPassword(ctx context.Context, phoneNumber, password string) (string, string, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return "", "", err // 可能是用户不存在或数据库错误
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(password)); err != nil {
		return "", "", domainerrors.NewDomainError(domainerrors.ErrUnauthorized, "手机号或密码错误")
	}

	// 3. 登录成功，返回用户ID和用户名
	return user.ID(), user.Name(), nil
}

// LoginByWeChat 微信登录
func (s *AuthService) LoginByWeChat(ctx context.Context, code string) (string, string, error) {
	// 1. Use the code to get user info from WeChat's API (e.g., openID).
	// 2. Check if a user with this openID exists in the database.
	//    - If not, create a new user.
	//    - If yes, retrieve the user.
	// 3. Return user ID and name.
	return "wechat_user_id", "wechat_user_name", nil
}
