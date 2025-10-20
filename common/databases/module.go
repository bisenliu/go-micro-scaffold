package databases

import (
	"go.uber.org/fx"

	"common/databases/rdbms"
	"common/databases/redis"
)

// RDBMSModule provides relational database management system functionalities
var RDBMSModule = fx.Module("rdbms",
	rdbms.Module,
)

// RedisModule provides Redis related functionalities
var RedisModule = fx.Module("redis",
	redis.RedisModule,
)

// Module provides all database related functionalities
var Module = fx.Module("database",
	RDBMSModule,
	RedisModule,
)
