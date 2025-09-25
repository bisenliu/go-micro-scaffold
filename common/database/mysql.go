package database

import (
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"

	"common/config"
)

// EntClient Ent ORM 客户端封装
type EntClient struct {
	config *config.Config
	driver *entsql.Driver
	db     *sql.DB
}

// EntClientParams 客户端依赖参数
type EntClientParams struct {
	fx.In
	Config *config.Config
}

// NewEntClient 创建新的 Ent 客户端
func NewEntClient(params EntClientParams) (*EntClient, error) {
	cfg := params.Config

	// 构建数据库连接字符串
	var dsn string
	var dbType string

	switch cfg.Database.Type {
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Database,
		)
		dbType = dialect.Postgres
	case "sqlite":
		dsn = cfg.Database.Database
		dbType = dialect.SQLite
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Database,
		)
		dbType = dialect.MySQL
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	// 打开数据库连接
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	if cfg.Database.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 创建 Ent 驱动
	driver := entsql.OpenDB(dialect.MySQL, db)

	return &EntClient{
		config: cfg,
		driver: driver,
		db:     db,
	}, nil
}

// Driver 获取 Ent 数据库驱动
func (c *EntClient) Driver() *entsql.Driver {
	return c.driver
}

// DB 获取原始数据库连接
func (c *EntClient) DB() *sql.DB {
	return c.db
}

// Close 关闭连接
func (c *EntClient) Close() error {
	return c.driver.Close()
}

// WithTx 在事务中执行操作（使用原生 SQL）
func (c *EntClient) WithTx(fn func(*sql.Tx) error) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// EntModule Ent ORM 模块
var EntModule = fx.Provide(NewEntClient)
