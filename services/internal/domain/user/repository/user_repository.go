package repository

import (
	"context"

	"services/internal/domain/user/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户 	保存用户
	Create(ctx context.Context, user *entity.User) error

	// List 分页查询用户列表
	List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)

	// 根据手机号查询用户是否存在
	ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
}
