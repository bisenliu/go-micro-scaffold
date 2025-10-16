package mapper

import (
	"services/internal/domain/user/entity"
	"services/internal/infrastructure/persistence/ent/gen"
)

// Mapper 通用映射器接口
type Mapper[TDomain any, TData any] interface {
	// ToEntity 将数据模型转换为领域实体
	ToEntity(data *TData) *TDomain
	
	// ToData 将领域实体转换为数据模型
	ToData(entity *TDomain) *TData
	
	// ToEntities 批量转换为领域实体
	ToEntities(dataList []*TData) []*TDomain
	
	// ToDataList 批量转换为数据模型
	ToDataList(entities []*TDomain) []*TData
}

// UserMapper 用户映射器接口
type UserMapper interface {
	Mapper[entity.User, gen.User]
	
	// ToEntityWithID 转换为领域实体并设置ID等字段
	ToEntityWithID(entUser *gen.User) *entity.User
	
	// ToCreateData 转换为创建数据（不包含ID等自动生成字段）
	ToCreateData(userEntity *entity.User) *gen.UserCreate
	
	// ToUpdateData 转换为更新数据
	ToUpdateData(userEntity *entity.User) map[string]interface{}
}