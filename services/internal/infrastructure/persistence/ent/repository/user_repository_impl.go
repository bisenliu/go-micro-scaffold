package repository

import (
	"context"
	"fmt"

	"services/internal/domain/user/entity"
	"services/internal/domain/user/repository"
	"services/internal/infrastructure/persistence/ent/gen"
	entuser "services/internal/infrastructure/persistence/ent/gen/user"

	"github.com/google/uuid"
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

// Create 保存用户
func (r *UserRepositoryImpl) Create(ctx context.Context, userEntity *entity.User) error {

	user, err := r.client.User.Create().
		SetOpenID(userEntity.OpenID()).
		SetName(userEntity.Name()).
		SetPhoneNumber(userEntity.PhoneNumber()).
		SetPassword(userEntity.Password()).
		SetGender(userEntity.Gender()).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// 将数据库生成的ID/时间戳设置给领域实体，并返回给调用方回领域实体
	userEntity.SetID(user.ID.String())
	userEntity.SetUpdatedAt(user.UpdatedAt)
	userEntity.SetCreatedAt(user.CreatedAt)

	return nil
}

// Update 更新用户信息
func (r *UserRepositoryImpl) Update(ctx context.Context, userEntity *entity.User) error {
	// 将字符串类型的ID转换为uuid.UUID类型
	userID, err := uuid.Parse(userEntity.ID())
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// 更新用户时，updated_at 字段会自动更新为当前时间
	// 因为在数据库层面已经配置了 UpdateDefault(time.Now)
	_, err = r.client.User.UpdateOneID(userID).
		SetName(userEntity.Name()).
		SetPhoneNumber(userEntity.PhoneNumber()).
		SetGender(userEntity.Gender()).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
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
