package repository

import (
	"context"

	"services/internal/domain/user/entity"
	"services/internal/domain/user/repository"
	"services/internal/infrastructure/persistence/ent/gen"
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
func (r *UserRepositoryImpl) Save(ctx context.Context, userEntity *entity.User) error {
	// 使用 Ent ORM 提供的创建方法

	return nil
}

// List 获取用户列表
func (r *UserRepositoryImpl) List(ctx context.Context, page, size int) ([]*entity.User, int64, error) {
	// 使用 Ent ORM 提供的查询方法

	return nil, 0, nil
}
