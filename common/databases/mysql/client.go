package mysql

import (
	"context"
	"database/sql"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	"go.uber.org/zap"
)

// Client 数据库客户端接口
type ClientInterface interface {
	Name() string
	Driver() *entsql.Driver
	DB() *sql.DB
	Config() *DatabaseConfigWrapper
	Ping(ctx context.Context) error
	Close() error
	WithTx(ctx context.Context, fn func(*sql.Tx) error) error
}

// Client 数据库客户端实现
type Client struct {
	name   string
	driver *entsql.Driver
	db     *sql.DB
	config *DatabaseConfigWrapper
	logger *zap.Logger
}

// Name 获取客户端名称
func (c *Client) Name() string {
	return c.name
}

// Driver 获取 Ent 数据库驱动
func (c *Client) Driver() *entsql.Driver {
	return c.driver
}

// DB 获取原始数据库连接
func (c *Client) DB() *sql.DB {
	return c.db
}

// Config 获取配置信息
func (c *Client) Config() *DatabaseConfigWrapper {
	return c.config
}

// Ping 测试数据库连接
func (c *Client) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// Close 关闭连接
func (c *Client) Close() error {
	if err := c.driver.Close(); err != nil {
		c.logger.Error("Failed to close database driver",
			zap.String("client", c.name),
			zap.Error(err),
		)
		return err
	}

	c.logger.Info("Database client closed successfully",
		zap.String("client", c.name),
	)
	return nil
}

// WithTx 在事务中执行操作
func (c *Client) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				c.logger.Error("Failed to rollback transaction after panic",
					zap.String("client", c.name),
					zap.Error(rollbackErr),
				)
			}
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				c.logger.Error("Failed to rollback transaction",
					zap.String("client", c.name),
					zap.Error(rollbackErr),
				)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", commitErr)
			}
		}
	}()

	err = fn(tx)
	return err
}
