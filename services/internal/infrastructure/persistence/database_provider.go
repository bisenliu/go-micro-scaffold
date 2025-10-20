package persistence

import (
	"go.uber.org/zap"

	"common/databases/rdbms"
	"services/internal/infrastructure/persistence/ent/gen"
)

// DatabaseProvider 数据库提供者，统一管理数据库访问
type DatabaseProvider struct {
	manager rdbms.ManagerInterface
	logger  *zap.Logger
}

// NewDatabaseProvider 创建数据库提供者
func NewDatabaseProvider(manager rdbms.ManagerInterface, logger *zap.Logger) *DatabaseProvider {
	return &DatabaseProvider{
		manager: manager,
		logger:  logger,
	}
}

// GetEntClient 获取Ent客户端使用的数据库
func (p *DatabaseProvider) GetEntClient() (*rdbms.Client, error) {
	// 可以通过别名配置指定Ent使用的数据库
	if client, err := p.manager.GetByAlias("ent"); err == nil {
		return client, nil
	}
	// 回退到默认数据库
	return p.manager.Default()
}

// GetCasbinClient 获取Casbin使用的数据库
func (p *DatabaseProvider) GetCasbinClient() (*rdbms.Client, error) {
	// 可以通过别名配置指定Casbin使用的数据库
	if client, err := p.manager.GetByAlias("casbin"); err == nil {
		return client, nil
	}
	// 回退到默认数据库
	return p.manager.Default()
}

// GetHealthCheckClient 获取健康检查使用的数据库
func (p *DatabaseProvider) GetHealthCheckClient() (*rdbms.Client, error) {
	// 可以通过别名配置指定健康检查使用的数据库
	if client, err := p.manager.GetByAlias("health_check"); err == nil {
		return client, nil
	}
	// 回退到默认数据库
	return p.manager.Default()
}

// CreateEntClient 创建Ent客户端
func (p *DatabaseProvider) CreateEntClient() (*gen.Client, error) {
	dbClient, err := p.GetEntClient()
	if err != nil {
		return nil, err
	}

	// 构建 Ent 选项
	entOptions := []gen.Option{
		gen.Driver(dbClient.Driver()),
	}

	p.logger.Info("Creating Ent client",
		zap.String("database", dbClient.Name()))

	// 创建 Ent 客户端
	entClient := gen.NewClient(entOptions...)

	return entClient, nil
}
