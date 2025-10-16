package interfaces

import (
	"context"
	
	"go.uber.org/zap"
)

// Logger 日志接口
// 定义了统一的日志记录接口，支持结构化日志和上下文传递
type Logger interface {
	// Debug 记录调试级别日志
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	
	// Info 记录信息级别日志
	Info(ctx context.Context, msg string, fields ...zap.Field)
	
	// Warn 记录警告级别日志
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	
	// Error 记录错误级别日志
	Error(ctx context.Context, msg string, fields ...zap.Field)
	
	// Fatal 记录致命错误日志并退出程序
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	
	// With 创建带有预设字段的子日志器
	With(fields ...zap.Field) Logger
	
	// WithContext 从上下文中提取字段创建子日志器
	WithContext(ctx context.Context) Logger
	
	// Sync 同步日志缓冲区
	Sync() error
	
	// GetZapLogger 获取底层的zap.Logger实例
	GetZapLogger() *zap.Logger
}

// LoggerFactory 日志工厂接口
// 用于创建不同配置的日志器实例
type LoggerFactory interface {
	// CreateLogger 创建日志器
	CreateLogger(config LoggerConfig) (Logger, error)
	
	// CreateDevelopmentLogger 创建开发环境日志器
	CreateDevelopmentLogger() (Logger, error)
	
	// CreateProductionLogger 创建生产环境日志器
	CreateProductionLogger(config LoggerConfig) (Logger, error)
}