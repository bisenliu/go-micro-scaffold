package mapper

import (
	"github.com/google/uuid"

	"services/internal/domain/user/entity"
	"services/internal/infrastructure/persistence/ent/gen"
)

// UserMapperImpl 用户映射器实现
type UserMapperImpl struct{}

// NewUserMapper 创建用户映射器
func NewUserMapper() UserMapper {
	return &UserMapperImpl{}
}

// ToEntity 将Ent用户模型转换为领域实体（不设置ID等字段）
func (m *UserMapperImpl) ToEntity(entUser *gen.User) *entity.User {
	if entUser == nil {
		return nil
	}

	return entity.NewUser(
		entUser.OpenID,
		entUser.Name,
		entUser.PhoneNumber,
		entUser.Password,
		entUser.Gender,
	)
}

// ToEntityWithID 将Ent用户模型转换为领域实体并设置ID等字段
func (m *UserMapperImpl) ToEntityWithID(entUser *gen.User) *entity.User {
	if entUser == nil {
		return nil
	}

	user := m.ToEntity(entUser)
	
	// 设置ID和时间戳字段
	user.SetID(entUser.ID.String())
	user.SetCreatedAt(entUser.CreatedAt)
	user.SetUpdatedAt(entUser.UpdatedAt)

	return user
}

// ToData 将领域实体转换为Ent用户模型（通常不直接使用）
func (m *UserMapperImpl) ToData(userEntity *entity.User) *gen.User {
	if userEntity == nil {
		return nil
	}

	entUser := &gen.User{
		OpenID:      userEntity.OpenID(),
		Name:        userEntity.Name(),
		PhoneNumber: userEntity.PhoneNumber(),
		Password:    userEntity.Password(),
		Gender:      userEntity.Gender(),
	}

	// 如果有ID，尝试解析并设置
	if userEntity.ID() != "" {
		if id, err := uuid.Parse(userEntity.ID()); err == nil {
			entUser.ID = id
		}
	}

	return entUser
}

// ToCreateData 将领域实体转换为创建数据
func (m *UserMapperImpl) ToCreateData(userEntity *entity.User) *gen.UserCreate {
	if userEntity == nil {
		return nil
	}

	// 这里返回的是一个函数，用于配置UserCreate
	// 实际使用时需要在仓储中调用
	return nil // 这个方法需要在仓储中特殊处理
}

// ToUpdateData 将领域实体转换为更新数据
func (m *UserMapperImpl) ToUpdateData(userEntity *entity.User) map[string]interface{} {
	if userEntity == nil {
		return nil
	}

	return map[string]interface{}{
		"name":         userEntity.Name(),
		"phone_number": userEntity.PhoneNumber(),
		"gender":       userEntity.Gender(),
		// 注意：不包含password，因为密码更新通常需要特殊处理
		// 注意：不包含open_id，因为这通常是不可变的
	}
}

// ToEntities 批量转换为领域实体
func (m *UserMapperImpl) ToEntities(entUsers []*gen.User) []*entity.User {
	if entUsers == nil {
		return nil
	}

	users := make([]*entity.User, 0, len(entUsers))
	for _, entUser := range entUsers {
		if user := m.ToEntityWithID(entUser); user != nil {
			users = append(users, user)
		}
	}

	return users
}

// ToDataList 批量转换为数据模型
func (m *UserMapperImpl) ToDataList(entities []*entity.User) []*gen.User {
	if entities == nil {
		return nil
	}

	entUsers := make([]*gen.User, 0, len(entities))
	for _, userEntity := range entities {
		if entUser := m.ToData(userEntity); entUser != nil {
			entUsers = append(entUsers, entUser)
		}
	}

	return entUsers
}