package mysql

import (
	"database/sql"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// ClientBuilder 客户端构建器
type ClientBuilder struct {
	config *DatabaseConfigWrapper
	logger *zap.Logger
}

// NewClientBuilder 创建客户端构建器
func NewClientBuilder(config *DatabaseConfigWrapper, logger *zap.Logger) *ClientBuilder {
	return &ClientBuilder{
		config: config,
		logger: logger,
	}
}

// Build 构建数据库客户端
func (b *ClientBuilder) Build() (*Client, error) {
	// 获取DSN和数据库类型
	dsn, dbType, err := b.config.DSN()
	if err != nil {
		return nil, fmt.Errorf("failed to generate DSN: %w", err)
	}

	// 打开数据库连接
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 配置连接池
	b.configureConnectionPool(db)

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 创建 Ent 驱动
	driver := entsql.OpenDB(dbType, db)

	client := &Client{
		name:   b.config.Name,
		driver: driver,
		db:     db,
		config: b.config,
		logger: b.logger,
	}

	b.logger.Info("Database client created successfully",
		zap.String("name", b.config.Name),
		zap.String("type", b.config.Type),
		zap.String("host", b.config.Host),
		zap.Int("port", b.config.Port),
		zap.String("database", b.config.Database),
	)

	return client, nil
}

// configureConnectionPool 配置连接池
func (b *ClientBuilder) configureConnectionPool(db *sql.DB) {
	db.SetMaxOpenConns(b.config.MaxOpenConns)
	db.SetMaxIdleConns(b.config.MaxIdleConns)
	db.SetConnMaxLifetime(b.config.ConnMaxLifetime)
	if b.config.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(b.config.ConnMaxIdleTime)
	}
}
