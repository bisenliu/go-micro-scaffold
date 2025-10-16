package interfaces

import (
	"context"
	"database/sql"
)

// Database 数据库接口
// 定义了数据库操作的基本接口
type Database interface {
	// GetDB 获取原始数据库连接
	GetDB() *sql.DB
	
	// Ping 检查数据库连接
	Ping(ctx context.Context) error
	
	// Close 关闭数据库连接
	Close() error
	
	// Stats 获取数据库连接统计信息
	Stats() sql.DBStats
}

// DatabaseManager 数据库管理器接口
// 管理多个数据库连接，支持主从分离
type DatabaseManager interface {
	// GetPrimaryDB 获取主数据库连接
	GetPrimaryDB() Database
	
	// GetReadOnlyDB 获取只读数据库连接
	// 如果没有配置只读数据库，返回主数据库
	GetReadOnlyDB() Database
	
	// Health 检查所有数据库连接的健康状态
	Health(ctx context.Context) error
	
	// Close 关闭所有数据库连接
	Close() error
	
	// GetStats 获取所有数据库的统计信息
	GetStats() DatabaseStats
}

// DatabaseStats 数据库统计信息
type DatabaseStats struct {
	Primary  sql.DBStats `json:"primary"`
	ReadOnly sql.DBStats `json:"read_only,omitempty"`
}

// DatabaseFactory 数据库工厂接口
// 用于创建不同类型的数据库连接
type DatabaseFactory interface {
	// CreateDatabase 根据配置创建数据库连接
	CreateDatabase(config DatabaseConnection) (Database, error)
	
	// CreateManager 创建数据库管理器
	CreateManager(config DatabaseConfig) (DatabaseManager, error)
}