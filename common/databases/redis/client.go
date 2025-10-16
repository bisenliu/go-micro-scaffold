package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"

	"common/interfaces"
)

// RedisClient Redis客户端
type RedisClient struct {
	*redis.Client
}

// NewRedisClient 创建新的Redis客户端
func NewRedisClient(configProvider interfaces.ConfigProvider) (*RedisClient, error) {
	cfg := configProvider.GetRedisConfig()

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Database,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisClient{
		Client: rdb,
	}, nil
}

// Close 关闭客户端
func (c *RedisClient) Close() error {
	return c.Client.Close()
}

// Module FX模块
var RedisModule = fx.Module("redis",
	fx.Provide(NewRedisClient),
	fx.Invoke(func(lc fx.Lifecycle, client *RedisClient) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return client.Close()
			},
		})
	}),
)
