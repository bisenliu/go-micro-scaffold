package interfaces

import (
	"context"
	"time"
)

// CacheManager 缓存管理器接口
// 定义了缓存操作的统一接口
type CacheManager interface {
	// Get 获取缓存值
	Get(ctx context.Context, key string) (string, error)
	
	// Set 设置缓存值
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	
	// Delete 删除缓存
	Delete(ctx context.Context, key string) error
	
	// Exists 检查键是否存在
	Exists(ctx context.Context, key string) (bool, error)
	
	// Expire 设置键的过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) error
	
	// TTL 获取键的剩余生存时间
	TTL(ctx context.Context, key string) (time.Duration, error)
	
	// Keys 获取匹配模式的所有键
	Keys(ctx context.Context, pattern string) ([]string, error)
	
	// FlushDB 清空当前数据库
	FlushDB(ctx context.Context) error
	
	// Ping 检查连接
	Ping(ctx context.Context) error
	
	// Close 关闭连接
	Close() error
}

// CacheFactory 缓存工厂接口
type CacheFactory interface {
	// CreateCacheManager 创建缓存管理器
	CreateCacheManager(config RedisConfig) (CacheManager, error)
}