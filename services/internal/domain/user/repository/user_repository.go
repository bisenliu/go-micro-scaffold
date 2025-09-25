package repository

import (
	"context"

	"services/internal/domain/user/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Save 保存用户
	Save(ctx context.Context, user *entity.User) error

	// List 分页查询用户列表
	List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)
}
