package repository

import (
	"context"
	"time"

	"services/internal/domain/user/entity"
)

// UserListFilter 用户列表过滤条件
type UserListFilter struct {
	Name      *string    // 姓名模糊查询
	Gender    *int       // 性别过滤
	StartTime *time.Time // 创建时间开始
	EndTime   *time.Time // 创建时间结束
}

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error

	// List 分页查询用户列表
	List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)

	// ListWithFilter 带过滤条件的分页查询用户列表
	ListWithFilter(ctx context.Context, filter *UserListFilter, offset, limit int) ([]*entity.User, int64, error)

	// ExistsByPhoneNumber 根据手机号查询用户是否存在
	ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *entity.User) error

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id string) (*entity.User, error)
}
