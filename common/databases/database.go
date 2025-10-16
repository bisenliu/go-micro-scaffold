package databases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/interfaces"
)

// DatabaseImpl 数据库实现
type DatabaseImpl struct {
	db *sql.DB
}

// NewDatabase 创建数据库连接
func NewDatabase(config interfaces.DatabaseConnection, logger interfaces.Logger) (interfaces.Database, error) {
	dsn, err := buildDSN(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build DSN: %w", err)
	}

	driverName := config.Type
	if config.Type == "sqlite" {
		driverName = "sqlite3"
	}
	
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info(context.Background(), "Database connected successfully",
		zap.String("type", config.Type),
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
		zap.String("database", config.Database))

	return &DatabaseImpl{db: db}, nil
}

// GetDB 获取原始数据库连接
func (d *DatabaseImpl) GetDB() *sql.DB {
	return d.db
}

// Ping 检查数据库连接
func (d *DatabaseImpl) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close 关闭数据库连接
func (d *DatabaseImpl) Close() error {
	return d.db.Close()
}

// Stats 获取数据库连接统计信息
func (d *DatabaseImpl) Stats() sql.DBStats {
	return d.db.Stats()
}

// DatabaseManagerImpl 数据库管理器实现
type DatabaseManagerImpl struct {
	primaryDB  interfaces.Database
	readOnlyDB interfaces.Database
	logger     interfaces.Logger
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager(configProvider interfaces.ConfigProvider, logger interfaces.Logger) (interfaces.DatabaseManager, error) {
	config := configProvider.GetDatabaseConfig()

	// 创建主数据库连接
	primaryDB, err := NewDatabase(config.Primary, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create primary database: %w", err)
	}

	manager := &DatabaseManagerImpl{
		primaryDB: primaryDB,
		logger:    logger,
	}

	// 如果配置了只读数据库，创建只读连接
	if config.ReadOnly.Host != "" {
		readOnlyDB, err := NewDatabase(config.ReadOnly, logger)
		if err != nil {
			logger.Warn(context.Background(), "Failed to create read-only database, using primary database",
				zap.Error(err))
			manager.readOnlyDB = primaryDB
		} else {
			manager.readOnlyDB = readOnlyDB
		}
	} else {
		manager.readOnlyDB = primaryDB
	}

	return manager, nil
}

// GetPrimaryDB 获取主数据库连接
func (dm *DatabaseManagerImpl) GetPrimaryDB() interfaces.Database {
	return dm.primaryDB
}

// GetReadOnlyDB 获取只读数据库连接
func (dm *DatabaseManagerImpl) GetReadOnlyDB() interfaces.Database {
	return dm.readOnlyDB
}

// Health 检查所有数据库连接的健康状态
func (dm *DatabaseManagerImpl) Health(ctx context.Context) error {
	if err := dm.primaryDB.Ping(ctx); err != nil {
		return fmt.Errorf("primary database health check failed: %w", err)
	}

	// 如果只读数据库不是主数据库，也检查它的健康状态
	if dm.readOnlyDB != dm.primaryDB {
		if err := dm.readOnlyDB.Ping(ctx); err != nil {
			return fmt.Errorf("read-only database health check failed: %w", err)
		}
	}

	return nil
}

// Close 关闭所有数据库连接
func (dm *DatabaseManagerImpl) Close() error {
	var lastErr error

	if err := dm.primaryDB.Close(); err != nil {
		dm.logger.Error(context.Background(), "Failed to close primary database", zap.Error(err))
		lastErr = err
	}

	// 如果只读数据库不是主数据库，也关闭它
	if dm.readOnlyDB != dm.primaryDB {
		if err := dm.readOnlyDB.Close(); err != nil {
			dm.logger.Error(context.Background(), "Failed to close read-only database", zap.Error(err))
			lastErr = err
		}
	}

	return lastErr
}

// GetStats 获取所有数据库的统计信息
func (dm *DatabaseManagerImpl) GetStats() interfaces.DatabaseStats {
	stats := interfaces.DatabaseStats{
		Primary: dm.primaryDB.Stats(),
	}

	// 如果只读数据库不是主数据库，获取它的统计信息
	if dm.readOnlyDB != dm.primaryDB {
		stats.ReadOnly = dm.readOnlyDB.Stats()
	}

	return stats
}

// DatabaseFactoryImpl 数据库工厂实现
type DatabaseFactoryImpl struct {
	logger interfaces.Logger
}

// NewDatabaseFactory 创建数据库工厂
func NewDatabaseFactory(logger interfaces.Logger) interfaces.DatabaseFactory {
	return &DatabaseFactoryImpl{
		logger: logger,
	}
}

// CreateDatabase 根据配置创建数据库连接
func (f *DatabaseFactoryImpl) CreateDatabase(config interfaces.DatabaseConnection) (interfaces.Database, error) {
	return NewDatabase(config, f.logger)
}

// CreateManager 创建数据库管理器
func (f *DatabaseFactoryImpl) CreateManager(config interfaces.DatabaseConfig) (interfaces.DatabaseManager, error) {
	// 创建临时配置提供者
	tempProvider := &tempDatabaseConfigProvider{databaseConfig: config}
	return NewDatabaseManager(tempProvider, f.logger)
}

// buildDSN 构建数据库连接字符串
func buildDSN(config interfaces.DatabaseConnection) (string, error) {
	switch config.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.Database, config.Charset), nil
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.Username, config.Password, config.Database), nil
	case "sqlite", "sqlite3":
		return config.Database, nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", config.Type)
	}
}

// tempDatabaseConfigProvider 临时配置提供者，用于工厂方法
type tempDatabaseConfigProvider struct {
	databaseConfig interfaces.DatabaseConfig
}

func (t *tempDatabaseConfigProvider) GetDatabaseConfig() interfaces.DatabaseConfig {
	return t.databaseConfig
}

func (t *tempDatabaseConfigProvider) GetServerConfig() interfaces.ServerConfig         { return interfaces.ServerConfig{} }
func (t *tempDatabaseConfigProvider) GetAuthConfig() interfaces.AuthConfig             { return interfaces.AuthConfig{} }
func (t *tempDatabaseConfigProvider) GetLoggerConfig() interfaces.LoggerConfig         { return interfaces.LoggerConfig{} }
func (t *tempDatabaseConfigProvider) GetRedisConfig() interfaces.RedisConfig           { return interfaces.RedisConfig{} }
func (t *tempDatabaseConfigProvider) GetTokenConfig() interfaces.TokenConfig           { return interfaces.TokenConfig{} }
func (t *tempDatabaseConfigProvider) GetValidationConfig() interfaces.ValidationConfig { return interfaces.ValidationConfig{} }
func (t *tempDatabaseConfigProvider) GetRateLimitConfig() interfaces.RateLimitConfig   { return interfaces.RateLimitConfig{} }
func (t *tempDatabaseConfigProvider) Reload() error                                    { return nil }
func (t *tempDatabaseConfigProvider) GetEnv() string                                   { return "development" }

// DatabaseModule FX模块
var DatabaseModule = fx.Provide(
	fx.Annotate(
		NewDatabaseManager,
		fx.As(new(interfaces.DatabaseManager)),
	),
	fx.Annotate(
		NewDatabaseFactory,
		fx.As(new(interfaces.DatabaseFactory)),
	),
)