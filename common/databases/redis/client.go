package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"

	"common/config"
)

// RedisClient Redis客户端
type RedisClient struct {
	*redis.Client
}

// RedisClientParams 客户端依赖参数
type RedisClientParams struct {
	fx.In
	Config *config.Config
}

// NewRedisClient 创建新的Redis客户端
func NewRedisClient(params RedisClientParams) (*RedisClient, error) {
	cfg := params.Config

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
		PoolSize: cfg.Redis.PoolSize,
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
