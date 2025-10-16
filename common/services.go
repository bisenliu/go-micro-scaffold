package common

import (
	"go.uber.org/fx"
	
	"common/interfaces"
)

// CommonServicesImpl Common服务实现
// 实现 interfaces.CommonServices 接口，聚合所有 common 服务
type CommonServicesImpl struct {
	configProvider     interfaces.ConfigProvider
	logger            interfaces.Logger
	databaseManager   interfaces.DatabaseManager
	middlewareProvider interfaces.MiddlewareProvider
	jwtService        interfaces.JWTService
}

// CommonServicesParams 创建CommonServices所需的依赖
type CommonServicesParams struct {
	fx.In
	
	ConfigProvider     interfaces.ConfigProvider
	Logger            interfaces.Logger
	DatabaseManager   interfaces.DatabaseManager
	MiddlewareProvider interfaces.MiddlewareProvider `optional:"true"`
	JWTService        interfaces.JWTService         `optional:"true"`
}

// NewCommonServices 创建Common服务实例
func NewCommonServices(params CommonServicesParams) interfaces.CommonServices {
	return &CommonServicesImpl{
		configProvider:     params.ConfigProvider,
		logger:            params.Logger,
		databaseManager:   params.DatabaseManager,
		middlewareProvider: params.MiddlewareProvider,
		jwtService:        params.JWTService,
	}
}

// Config 获取配置提供者
func (cs *CommonServicesImpl) Config() interfaces.ConfigProvider {
	return cs.configProvider
}

// Logger 获取日志器
func (cs *CommonServicesImpl) Logger() interfaces.Logger {
	return cs.logger
}

// Database 获取数据库管理器
func (cs *CommonServicesImpl) Database() interfaces.DatabaseManager {
	return cs.databaseManager
}

// Cache 获取缓存管理器
func (cs *CommonServicesImpl) Cache() interfaces.CacheManager {
	// TODO: 在后续任务中实现
	return nil
}

// Validator 获取验证器
func (cs *CommonServicesImpl) Validator() interfaces.Validator {
	// TODO: 在后续任务中实现
	return nil
}

// Middleware 获取中间件提供者
func (cs *CommonServicesImpl) Middleware() interfaces.MiddlewareProvider {
	return cs.middlewareProvider
}

// JWT 获取JWT服务
func (cs *CommonServicesImpl) JWT() interfaces.JWTService {
	return cs.jwtService
}