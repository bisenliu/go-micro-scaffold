package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"services/internal/domain/errors"
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

// Create 保存用户
func (r *UserRepositoryImpl) Create(ctx context.Context, userEntity *entity.User) error {
	// 检查手机号是否已存在
	exists, err := r.ExistsByPhoneNumber(ctx, userEntity.PhoneNumber())
	if err != nil {
		return fmt.Errorf("failed to check phone number: %w", err)
	}
	if exists {
		return errors.ErrPhoneAlreadyExists
	}

	user, err := r.client.User.Create().
		SetOpenID(userEntity.OpenID()).
		SetName(userEntity.Name()).
		SetPhoneNumber(userEntity.PhoneNumber()).
		SetPassword(userEntity.Password()).
		SetGender(userEntity.Gender()).
		Save(ctx)

	if err != nil {
		if gen.IsConstraintError(err) {
			return errors.ErrUserAlreadyExists
		}
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

func (r *UserRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	// 查询用户列表
	entUsers, err := r.client.User.Query().
		Offset(offset).
		Limit(limit).
		Order(gen.Desc(entuser.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}

	// 查询总数
	total, err := r.client.User.Query().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// 转换为领域实体
	users := make([]*entity.User, 0, len(entUsers))
	for _, entUser := range entUsers {
		user := r.entUserToEntity(entUser)
		users = append(users, user)
	}

	return users, int64(total), nil
}

// List 获取用户列表
func (r *UserRepositoryImpl) ListWithFilter(ctx context.Context, filter *repository.UserListFilter, offset, limit int) ([]*entity.User, int64, error) {
	// 构建基础查询（只构建一次）
	baseQuery := r.buildUserQuery(filter)

	// 先查询总数（使用 Clone 避免修改原始查询）
	total, err := baseQuery.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users with filter: %w", err)
	}

	// 再查询分页数据（复用相同的查询条件）
	entUsers, err := baseQuery.
		Offset(offset).
		Limit(limit).
		Order(gen.Desc(entuser.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users with filter: %w", err)
	}

	// 转换为领域实体
	users := make([]*entity.User, 0, len(entUsers))
	for _, entUser := range entUsers {
		user := r.entUserToEntity(entUser)
		users = append(users, user)
	}

	return users, int64(total), nil
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

// 提取一个私有方法来构建查询条件
func (r *UserRepositoryImpl) buildUserQuery(filter *repository.UserListFilter) *gen.UserQuery {
	query := r.client.User.Query()

	if filter.Name != nil && *filter.Name != "" {
		query = query.Where(entuser.NameContains(*filter.Name))
	}
	if filter.Gender != nil {
		query = query.Where(entuser.Gender(*filter.Gender))
	}
	if filter.StartTime != nil {
		query = query.Where(entuser.CreatedAtGTE(*filter.StartTime))
	}
	if filter.EndTime != nil {
		query = query.Where(entuser.CreatedAtLT(*filter.EndTime))
	}

	return query
}

// entUserToEntity 将Ent用户实体转换为领域用户实体
func (r *UserRepositoryImpl) entUserToEntity(entUser *gen.User) *entity.User {
	if entUser == nil {
		return nil
	}

	// 创建领域用户实体
	user := entity.NewUser(
		entUser.OpenID,
		entUser.Name,
		entUser.PhoneNumber,
		entUser.Password,
		entUser.Gender,
	)

	// 设置ID和其他字段
	user.SetID(entUser.ID.String())
	user.SetCreatedAt(entUser.CreatedAt)
	user.SetUpdatedAt(entUser.UpdatedAt)

	return user
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	entUser, err := r.client.User.Get(ctx, userID)
	if err != nil {
		if gen.IsNotFound(err) {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	user := r.entUserToEntity(entUser)
	return user, nil
}
