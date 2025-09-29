package ent

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"common/databases/mysql"
	"services/internal/infrastructure/persistence/ent/gen"
)

// BusinessClient 业务层 Ent 客户端
type BusinessClient struct {
	client *gen.Client
	name   string
	logger *zap.Logger
}

// ClientBuilder 业务客户端构建器
type ClientBuilder struct {
	manager mysql.ManagerInterface
	logger  *zap.Logger
}

// NewClientBuilder 创建业务客户端构建器
func NewClientBuilder(manager mysql.ManagerInterface, logger *zap.Logger) *ClientBuilder {
	return &ClientBuilder{
		manager: manager,
		logger:  logger,
	}
}

// BuildClient 构建指定名称的业务客户端
func (b *ClientBuilder) BuildClient(name string) (*BusinessClient, error) {
	// 获取底层数据库客户端
	dbClient, err := b.manager.GetClient(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get database client '%s': %w", name, err)
	}

	// 构建 Ent 选项
	entOptions := []gen.Option{
		gen.Driver(dbClient.Driver()),
	}

	// 根据配置启用调试模式
	if dbClient.Config().EnableDebug {
		b.logger.Info("Enabling Ent debug mode for client", zap.String("client", name))
		entOptions = append(entOptions, gen.Debug())
	}

	// 创建 Ent 客户端
	entClient := gen.NewClient(entOptions...)

	return &BusinessClient{
		client: entClient,
		name:   name,
		logger: b.logger,
	}, nil
}

// BuildPrimary 构建主数据库客户端
func (b *ClientBuilder) BuildPrimary() (*BusinessClient, error) {
	return b.BuildClient("primary")
}

// Query 获取查询客户端
func (bc *BusinessClient) Query() *gen.Client {
	return bc.client
}

// Tx 开始事务
func (bc *BusinessClient) Tx(ctx context.Context) (*gen.Tx, error) {
	return bc.client.Tx(ctx)
}

// Close 关闭客户端
func (bc *BusinessClient) Close() error {
	if err := bc.client.Close(); err != nil {
		bc.logger.Error("Failed to close business client",
			zap.String("client", bc.name),
			zap.Error(err),
		)
		return err
	}

	bc.logger.Info("Business client closed successfully",
		zap.String("client", bc.name),
	)
	return nil
}

// Migrate 执行数据库迁移
func (bc *BusinessClient) Migrate(ctx context.Context) error {
	bc.logger.Info("Starting database migration", zap.String("client", bc.name))

	if err := bc.client.Schema.Create(ctx); err != nil {
		bc.logger.Error("Database migration failed",
			zap.String("client", bc.name),
			zap.Error(err),
		)
		return fmt.Errorf("migration failed for client '%s': %w", bc.name, err)
	}

	bc.logger.Info("Database migration completed successfully", zap.String("client", bc.name))
	return nil
}

// Name 获取客户端名称
func (bc *BusinessClient) Name() string {
	return bc.name
}
