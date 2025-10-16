package mapper

// MapperFactory 映射器工厂
// 提供统一的映射器访问接口
type MapperFactory struct {
	userMapper       UserMapper
	userMapperHelper *UserMapperHelper
	baseMapper       *BaseMapper
}

// NewMapperFactory 创建映射器工厂
func NewMapperFactory(
	userMapper UserMapper,
	userMapperHelper *UserMapperHelper,
) *MapperFactory {
	return &MapperFactory{
		userMapper:       userMapper,
		userMapperHelper: userMapperHelper,
		baseMapper:       NewBaseMapper(),
	}
}

// User 获取用户映射器
func (f *MapperFactory) User() UserMapper {
	return f.userMapper
}

// UserHelper 获取用户映射器辅助类
func (f *MapperFactory) UserHelper() *UserMapperHelper {
	return f.userMapperHelper
}

// Base 获取基础映射器
func (f *MapperFactory) Base() *BaseMapper {
	return f.baseMapper
}

// MapperFactoryProvider 映射器工厂提供者接口
type MapperFactoryProvider interface {
	GetMapperFactory() *MapperFactory
}

// MapperFactoryProviderImpl 映射器工厂提供者实现
type MapperFactoryProviderImpl struct {
	factory *MapperFactory
}

// NewMapperFactoryProvider 创建映射器工厂提供者
func NewMapperFactoryProvider(factory *MapperFactory) MapperFactoryProvider {
	return &MapperFactoryProviderImpl{
		factory: factory,
	}
}

// GetMapperFactory 获取映射器工厂
func (p *MapperFactoryProviderImpl) GetMapperFactory() *MapperFactory {
	return p.factory
}