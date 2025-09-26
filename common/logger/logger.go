package logger

import (
	"common/config"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// TraceIDKey 是 context 中 traceID 的键名
	TraceIDKey = "traceID"
)

// NewLogger 创建logger实例
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	// 设置日志级别
	var level zapcore.Level
	switch cfg.Zap.Level {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "WARN":
		level = zapcore.WarnLevel
	case "ERROR":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 创建日志目录
	if err := os.MkdirAll(cfg.Zap.Director, 0755); err != nil {
		return nil, err
	}

	// 通用文件输出 - 所有日志（按天分割）
	allWriter, err := getLogWriter(
		cfg.Zap.Director+"/app",
		time.Duration(cfg.Zap.MaxAge),
		int64(cfg.Zap.MaxSize),
		uint(cfg.Zap.MaxBackups),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create all log writer: %w", err)
	}

	// Info及以上等级文件输出（按天分割）
	infoWriter, err := getLogWriter(
		cfg.Zap.Director+"/info",
		time.Duration(cfg.Zap.MaxAge),
		int64(cfg.Zap.MaxSize),
		uint(cfg.Zap.MaxBackups),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create info log writer: %w", err)
	}

	// Error等级单独文件输出（按天分割）
	errorWriter, err := getLogWriter(
		cfg.Zap.Director+"/error",
		time.Duration(cfg.Zap.MaxAge),
		int64(cfg.Zap.MaxSize),
		uint(cfg.Zap.MaxBackups),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create error log writer: %w", err)
	}

	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 创建不同等级的过滤器
	infoLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	errorLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// 创建核心 - 多个输出目标
	core := zapcore.NewTee(
		// 所有日志写入通用文件
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			allWriter,
			level,
		),
		// Info及以上等级写入info文件
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			infoWriter,
			infoLevelEnabler,
		),
		// Error等级写入error文件
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			errorWriter,
			errorLevelEnabler,
		),
		// 控制台输出
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleWriter,
			level,
		),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

// getLogWriter 创建按天分割的日志写入器
func getLogWriter(filename string, maxAge time.Duration, maxSize int64, maxBackup uint) (zapcore.WriteSyncer, error) {
	// 按天分割日志，每天创建一个新的日志文件
	logPath := fmt.Sprintf("%s.%%Y-%%m-%%d.log", filename)
	rotationLogs, err := rotatelogs.New(
		logPath, // 日志文件路径
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(maxAge*24*time.Hour),     // 日志文件保留时间
		rotatelogs.WithRotationTime(24*time.Hour),      // 每24小时分割一次日志
		rotatelogs.WithRotationSize(maxSize*1024*1024), // 每个文件保存的最大尺寸
	)
	if err != nil {
		return nil, fmt.Errorf("rotatelogs log failed: %w", err)
	}
	return zapcore.AddSync(rotationLogs), nil
}

// GenerateTraceID 生成新的 traceID
func GenerateTraceID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Failed to generate UUID: %v", err)
		return "unknown"
	}
	return id.String()
}

// WithTraceID 在 context 中设置 traceID
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceID 从 context 中获取 traceID
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// WithContext 创建带有 traceID 的 logger
func WithContext(logger *zap.Logger, ctx context.Context) *zap.Logger {
	traceID := GetTraceID(ctx)
	return logger.With(zap.String("traceID", traceID))
}

// ToContext 将 logger 存入 context
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, "contextLogger", logger)
}

// FromContext 从 context 中获取 logger
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value("contextLogger").(*zap.Logger); ok {
		return logger
	}
	// 返回 noop logger，避免 panic
	return zap.NewNop()
}

// Info 从 context 获取 logger 记录 Info 日志
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Info(msg, fields...)
}

// Error 从 context 获取 logger 记录 Error 日志
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Error(msg, fields...)
}

// Warn 从 context 获取 logger 记录 Warn 日志
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Warn(msg, fields...)
}

// Debug 从 context 获取 logger 记录 Debug 日志
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Debug(msg, fields...)
}

// Module FX模块
var Module = fx.Provide(NewLogger)
