package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"common/interfaces"
)

var (
	globalLogger interfaces.Logger
	globalMutex  sync.RWMutex
)

// SetGlobalLogger 设置全局logger实例
func SetGlobalLogger(logger interfaces.Logger) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	globalLogger = logger
}

// GetGlobalLogger 获取全局logger实例
func GetGlobalLogger() interfaces.Logger {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return globalLogger
}

// Debug 全局调试日志
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Debug(ctx, msg, fields...)
	}
}

// Info 全局信息日志
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Info(ctx, msg, fields...)
	}
}

// Warn 全局警告日志
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Warn(ctx, msg, fields...)
	}
}

// Error 全局错误日志
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Error(ctx, msg, fields...)
	}
}

// Fatal 全局致命错误日志
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Fatal(ctx, msg, fields...)
	}
}

// WithContext 从上下文创建logger
func WithContext(ctx context.Context) interfaces.Logger {
	if logger := GetGlobalLogger(); logger != nil {
		return logger.WithContext(ctx)
	}
	return nil
}

// ToContext 将logger存储到context中
func ToContext(ctx context.Context, logger interfaces.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

// FromContext 从context中获取logger
func FromContext(ctx context.Context) interfaces.Logger {
	if logger, ok := ctx.Value("logger").(interfaces.Logger); ok {
		return logger
	}
	return GetGlobalLogger()
}