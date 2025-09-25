package ent

import (
	"context"
	"fmt"

	"common/database"
	"services/internal/infrastructure/persistence/ent/gen"
)

// EntClientFactoryImpl Ent 客户端工厂实现
type EntClientFactoryImpl struct{}

// NewEntClientFactory 创建 Ent 客户端工厂
func NewEntClientFactory() database.EntClientFactory {
	return &EntClientFactoryImpl{}
}

// NewClient 创建新的 Ent 客户端
func (f *EntClientFactoryImpl) NewClient(driver *database.EntClient, options ...interface{}) (interface{}, error) {
	// 将 common/database.EntClient 转换为 Ent 驱动
	// 这里需要根据你的具体实现来调整
	entOptions := []gen.Option{
		gen.Driver(driver.Driver()),
	}

	// 处理传入的选项
	for _, opt := range options {
		if entOpt, ok := opt.(gen.Option); ok {
			entOptions = append(entOptions, entOpt)
		}
	}

	// 创建 Ent 客户端
	client := gen.NewClient(entOptions...)
	return client, nil
}

// Migrate 执行数据库迁移
func (f *EntClientFactoryImpl) Migrate(ctx context.Context, client interface{}) error {
	// 将传入的客户端转换为 Ent 客户端
	entClient, ok := client.(*gen.Client)
	if !ok {
		return fmt.Errorf("invalid client type, expected *gen.Client")
	}

	// 执行迁移
	return entClient.Schema.Create(ctx)
}
