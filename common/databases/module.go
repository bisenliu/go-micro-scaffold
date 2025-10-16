package databases

import (
	"go.uber.org/fx"

	// 暂时移除 redis 导入，等后续任务重构
	// "common/databases/redis"
)

// RedisModule provides Redis related functionalities
// var RedisModule = fx.Module("redis",
//     redis.RedisModule,
// )

// Module provides all database related functionalities
var Module = fx.Module("database",
	DatabaseModule,
	// 暂时禁用 Redis 模块，等后续任务重构
	// RedisModule,
)
