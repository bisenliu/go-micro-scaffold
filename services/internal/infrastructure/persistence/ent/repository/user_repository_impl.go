package repository

import (
	"context"

	"github.com/google/uuid"

	domainerrors "services/internal/domain/shared/errors"
	"services/internal/domain/user/entity"
	usererrors "services/internal/domain/user/errors"
	"services/internal/domain/user/repository"
	"services/internal/infrastructure/mapper"
	"services/internal/infrastructure/persistence/ent/gen"
	entuser "services/internal/infrastructure/persistence/ent/gen/user"
)

// UserRepositoryImpl Ent用户仓储实现
type UserRepositoryImpl struct {
	client       *gen.Client
	userMapper   mapper.UserMapper
	mapperHelper *mapper.UserMapperHelper
}

// NewUserRepository 创建用户仓储
func NewUserRepository(
	client *gen.Client,
	userMapper mapper.UserMapper,
	mapperHelper *mapper.UserMapperHelper,
) repository.UserRepository {
	return &UserRepositoryImpl{
		client:       client,
		userMapper:   userMapper,
		mapperHelper: mapperHelper,
	}
}

// Create 保存用户
func (r *UserRepositoryImpl) Create(ctx context.Context, userEntity *entity.User) error {
	// 检查手机号是否已存在
	exists, err := r.ExistsByPhoneNumber(ctx, userEntity.PhoneNumber())
	if err != nil {
		return err // 直接返回错误，因为ExistsByPhoneNumber已经包装过
	}
	if exists {
		return usererrors.ErrPhoneAlreadyExists
	}

	// 使用映射器辅助类配置创建操作
	createOp := r.client.User.Create()
	createOp = r.mapperHelper.ConfigureCreate(createOp, userEntity)
	
	user, err := createOp.Save(ctx)
	if err != nil {
		if gen.IsConstraintError(err) {
			return usererrors.ErrUserAlreadyExists
		}
		return domainerrors.NewInternalServerError("创建用户失败")
	}

	// 使用映射器辅助类应用数据库生成的字段到领域实体
	r.mapperHelper.ApplyEntityChanges(userEntity, user)

	return nil
}

// Update 更新用户信息
func (r *UserRepositoryImpl) Update(ctx context.Context, userEntity *entity.User) error {
	// 将字符串类型的ID转换为uuid.UUID类型
	userID, err := uuid.Parse(userEntity.ID())
	if err != nil {
		return domainerrors.NewInvalidDataError("无效的用户ID")
	}

	// 使用映射器辅助类配置更新操作
	updateOp := r.client.User.UpdateOneID(userID)
	updateOp = r.mapperHelper.ConfigureUpdate(updateOp, userEntity)
	
	// 更新用户时，updated_at 字段会自动更新为当前时间
	// 因为在数据库层面已经配置了 UpdateDefault(time.Now)
	_, err = updateOp.Save(ctx)

	if err != nil {
		if gen.IsNotFound(err) {
			return usererrors.ErrUserNotFound
		}
		if gen.IsConstraintError(err) {
			return usererrors.ErrPhoneAlreadyExists
		}
		return domainerrors.NewInternalServerError("更新用户失败")
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
		return nil, 0, domainerrors.NewInternalServerError("查询用户列表失败")
	}

	// 查询总数
	total, err := r.client.User.Query().Count(ctx)
	if err != nil {
		return nil, 0, domainerrors.NewInternalServerError("查询用户总数失败")
	}

	// 使用映射器转换为领域实体
	users := r.userMapper.ToEntities(entUsers)

	return users, int64(total), nil
}

// ListWithFilter 获取用户列表
func (r *UserRepositoryImpl) ListWithFilter(ctx context.Context, filter *repository.UserListFilter, offset, limit int) ([]*entity.User, int64, error) {
	// 构建基础查询（只构建一次）
	baseQuery := r.buildUserQuery(filter)

	// 先查询总数（使用 Clone 避免修改原始查询）
	total, err := baseQuery.Clone().Count(ctx)
	if err != nil {
		return nil, 0, domainerrors.NewInternalServerError("查询用户总数失败")
	}

	// 再查询分页数据（复用相同的查询条件）
	entUsers, err := baseQuery.
		Offset(offset).
		Limit(limit).
		Order(gen.Desc(entuser.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, domainerrors.NewInternalServerError("查询用户列表失败")
	}

	// 使用映射器转换为领域实体
	users := r.userMapper.ToEntities(entUsers)

	return users, int64(total), nil
}

func (r *UserRepositoryImpl) ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	exists, err := r.client.User.
		Query().
		Where(entuser.PhoneNumber(phoneNumber)).
		Exist(ctx)
	if err != nil {
		return false, domainerrors.NewInternalServerError("查询用户手机号是否存在失败")
	}

	return exists, nil
}

// 提取一个私有方法来构建查询条件
func (r *UserRepositoryImpl) buildUserQuery(filter *repository.UserListFilter) *gen.UserQuery {
	query := r.client.User.Query()

	if filter == nil {
		return query
	}

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



func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, domainerrors.NewInvalidDataError("无效的用户ID")
	}

	entUser, err := r.client.User.Get(ctx, userID)
	if err != nil {
		if gen.IsNotFound(err) {
			return nil, usererrors.ErrUserNotFound
		}
		return nil, domainerrors.NewInternalServerError("查询用户失败")
	}

	// 使用映射器转换为领域实体
	user := r.userMapper.ToEntityWithID(entUser)
	return user, nil
}

// FindByPhoneNumber 根据手机号获取用户
func (r *UserRepositoryImpl) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error) {
	entUser, err := r.client.User.
		Query().
		Where(entuser.PhoneNumberEQ(phoneNumber)).
		Only(ctx)
	if err != nil {
		if gen.IsNotFound(err) {
			return nil, usererrors.ErrUserNotFound
		}
		return nil, domainerrors.NewInternalServerError("通过手机号查询用户失败")
	}
	
	// 使用映射器转换为领域实体
	return r.userMapper.ToEntityWithID(entUser), nil
}
