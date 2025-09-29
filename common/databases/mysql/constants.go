package mysql

// 数据库实例名称常量定义
const (
	// 主数据库实例名称
	DB1 = "db1"

	// 分析数据库实例名称
	DB2 = "db2"

	// 读数据库实例名称（用于读写分离）
	ReadDB = "read"

	// 写数据库实例名称（用于读写分离）
	WriteDB = "write"
)
