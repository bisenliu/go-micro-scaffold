package database

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/config"
	"common/logger"
)

// EntClientFactory Ent 客户端工厂接口
// 这个接口需要在具体的业务项目中实现，因为每个项目的 Ent 生成代码不同
type EntClientFactory interface {
	// NewClient 创建新的 Ent 客户端
	// driver: 数据库驱动
	// options: Ent 客户端选项
	NewClient(driver *EntClient, options ...interface{}) (interface{}, error)

	// Migrate 执行数据库迁移
	Migrate(ctx context.Context, client interface{}) error
}

// EntService Ent 服务封装
type EntService struct {
	client  *EntClient
	factory EntClientFactory
	logger  *zap.Logger
	config  *config.Config
}

// EntServiceParams Ent 服务依赖参数
type EntServiceParams struct {
	fx.In
	Client  *EntClient
	Factory EntClientFactory `optional:"true"`
	Logger  *zap.Logger
	Config  *config.Config
}

// NewEntService 创建 Ent 服务
func NewEntService(params EntServiceParams) *EntService {
	return &EntService{
		client:  params.Client,
		factory: params.Factory,
		logger:  params.Logger,
		config:  params.Config,
	}
}

// GetClient 获取 Ent 客户端
func (s *EntService) GetClient() *EntClient {
	return s.client
}

// CreateBusinessClient 创建业务 Ent 客户端
// 需要传入具体的工厂实现
func (s *EntService) CreateBusinessClient(factory EntClientFactory, options ...interface{}) (interface{}, error) {
	if factory == nil {
		return nil, fmt.Errorf("ent client factory is required")
	}

	return factory.NewClient(s.client, options...)
}

// Migrate 执行数据库迁移
func (s *EntService) Migrate(ctx context.Context, businessClient interface{}) error {
	if s.factory == nil {
		return fmt.Errorf("ent client factory is not configured")
	}

	logger.Info(s.logger, ctx, "Starting database migration")

	if err := s.factory.Migrate(ctx, businessClient); err != nil {
		logger.Error(s.logger, ctx, "Database migration failed", zap.Error(err))
		return fmt.Errorf("migration failed: %w", err)
	}

	logger.Info(s.logger, ctx, "Database migration completed successfully")
	return nil
}

// Close 关闭服务
func (s *EntService) Close() error {
	return s.client.Close()
}

// EntServiceModule Ent 服务模块
var EntServiceModule = fx.Module("ent_service",
	fx.Provide(NewEntService),
	fx.Invoke(func(lc fx.Lifecycle, service *EntService) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return service.Close()
			},
		})
	}),
)

// SimpleEntClientOptions Ent 客户端选项构建器
type SimpleEntClientOptions struct {
	Debug bool
}

// BuildEntOptions 构建 Ent 选项
func (o *SimpleEntClientOptions) BuildEntOptions() []interface{} {
	var options []interface{}

	if o.Debug {
		// 这里需要根据具体的 Ent 生成代码来调整
		// 通常是类似 ent.Debug() 的选项
		options = append(options, "debug")
	}

	return options
}

// NewSimpleEntClientOptions 创建简单的 Ent 客户端选项
func NewSimpleEntClientOptions(cfg *config.Config) *SimpleEntClientOptions {
	return &SimpleEntClientOptions{
		Debug: cfg.System.Env == "development",
	}
}
