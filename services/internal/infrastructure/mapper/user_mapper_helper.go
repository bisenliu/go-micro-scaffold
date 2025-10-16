package mapper

import (
	"services/internal/domain/user/entity"
	"services/internal/infrastructure/persistence/ent/gen"
)

// UserMapperHelper 用户映射器辅助类
// 提供更实用的映射方法，特别是针对Ent的创建和更新操作
type UserMapperHelper struct {
	mapper UserMapper
}

// NewUserMapperHelper 创建用户映射器辅助类
func NewUserMapperHelper(mapper UserMapper) *UserMapperHelper {
	return &UserMapperHelper{
		mapper: mapper,
	}
}

// ConfigureCreate 配置用户创建操作
func (h *UserMapperHelper) ConfigureCreate(create *gen.UserCreate, userEntity *entity.User) *gen.UserCreate {
	if create == nil || userEntity == nil {
		return create
	}

	return create.
		SetOpenID(userEntity.OpenID()).
		SetName(userEntity.Name()).
		SetPhoneNumber(userEntity.PhoneNumber()).
		SetPassword(userEntity.Password()).
		SetGender(userEntity.Gender())
}

// ConfigureUpdate 配置用户更新操作
func (h *UserMapperHelper) ConfigureUpdate(update *gen.UserUpdateOne, userEntity *entity.User) *gen.UserUpdateOne {
	if update == nil || userEntity == nil {
		return update
	}

	return update.
		SetName(userEntity.Name()).
		SetPhoneNumber(userEntity.PhoneNumber()).
		SetGender(userEntity.Gender())
		// 注意：不更新OpenID和Password，这些字段通常需要特殊处理
}

// ConfigureBulkUpdate 配置批量更新操作
func (h *UserMapperHelper) ConfigureBulkUpdate(update *gen.UserUpdate, userEntity *entity.User) *gen.UserUpdate {
	if update == nil || userEntity == nil {
		return update
	}

	return update.
		SetName(userEntity.Name()).
		SetPhoneNumber(userEntity.PhoneNumber()).
		SetGender(userEntity.Gender())
}

// ApplyEntityChanges 将数据库结果应用到领域实体
func (h *UserMapperHelper) ApplyEntityChanges(userEntity *entity.User, entUser *gen.User) {
	if userEntity == nil || entUser == nil {
		return
	}

	// 设置数据库生成的字段
	userEntity.SetID(entUser.ID.String())
	userEntity.SetCreatedAt(entUser.CreatedAt)
	userEntity.SetUpdatedAt(entUser.UpdatedAt)
}