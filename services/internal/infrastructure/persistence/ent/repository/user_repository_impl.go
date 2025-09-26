package repository

import (
	"context"
	"fmt"

	"services/internal/domain/user/entity"
	"services/internal/domain/user/repository"
	"services/internal/infrastructure/persistence/ent/gen"
	entuser "services/internal/infrastructure/persistence/ent/gen/user"
)

// UserRepositoryImpl Ent用户仓储实现
type UserRepositoryImpl struct {
	client *gen.Client
}

// NewUserRepository 创建用户仓储
func NewUserRepository(client *gen.Client) repository.UserRepository {
	return &UserRepositoryImpl{
		client: client,
	}
}

// Save 保存用户
func (r *UserRepositoryImpl) Create(ctx context.Context, userEntity *entity.User) error {
	_, err := r.client.User.Create().
		SetOpenID(userEntity.OpenID()).
		SetName(userEntity.Name()).
		SetPhoneNumber(userEntity.PhoneNumber()).
		SetPassword(userEntity.Password()).
		SetGender(userEntity.Gender()).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// List 获取用户列表
func (r *UserRepositoryImpl) List(ctx context.Context, page, size int) ([]*entity.User, int64, error) {
	// 使用 Ent ORM 提供的查询方法

	return nil, 0, nil
}

func (r *UserRepositoryImpl) ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {

	exists, err := r.client.User.
		Query().
		Where(entuser.PhoneNumber(phoneNumber)).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check phone existence: %w", err)
	}

	return exists, nil
}
