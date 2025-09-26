package databases

import (
	"go.uber.org/fx"

	"common/databases/mysql"
	"common/databases/redis"
)

// MySQLModule provides MySQL related functionalities
var MySQLModule = fx.Module("mysql",
	mysql.EntModule,
	mysql.DatabaseManagerModule,
	mysql.EntServiceModule,
)

// RedisModule provides Redis related functionalities
var RedisModule = fx.Module("redis",
	redis.RedisModule,
)

// Module provides all database related functionalities
var Module = fx.Module("database",
	MySQLModule,
	RedisModule,
)
