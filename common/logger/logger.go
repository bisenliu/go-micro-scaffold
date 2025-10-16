package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"common/interfaces"
)

const (
	// TraceIDKey 是 context 中 traceID 的键名
	TraceIDKey = "traceID"
)

// LoggerImpl 日志实现
type LoggerImpl struct {
	logger *zap.Logger
}

// NewLogger 创建日志器
func NewLogger(configProvider interfaces.ConfigProvider) (interfaces.Logger, error) {
	config := configProvider.GetLoggerConfig()
	
	// 设置日志级别
	var level zapcore.Level
	switch config.Level {
	case "DEBUG", "debug":
		level = zapcore.DebugLevel
	case "INFO", "info":
		level = zapcore.InfoLevel
	case "WARN", "warn":
		level = zapcore.WarnLevel
	case "ERROR", "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 创建日志目录
	if err := os.MkdirAll(config.Director, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 通用文件输出 - 所有日志（按天分割）
	allWriter, err := getLogWriter(
		config.Director+"/app",
		time.Duration(config.MaxAge),
		int64(config.MaxSize),
		uint(config.MaxBackups),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create all log writer: %w", err)
	}

	// Info及以上等级文件输出（按天分割）
	infoWriter, err := getLogWriter(
		config.Director+"/info",
		time.Duration(config.MaxAge),
		int64(config.MaxSize),
		uint(config.MaxBackups),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create info log writer: %w", err)
	}

	// Error等级单独文件输出（按天分割）
	errorWriter, err := getLogWriter(
		config.Director+"/error",
		time.Duration(config.MaxAge),
		int64(config.MaxSize),
		uint(config.MaxBackups),
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

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &LoggerImpl{
		logger: zapLogger,
	}, nil
}

// Debug 记录调试级别日志
func (l *LoggerImpl) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Debug(msg, fields...)
}

// Info 记录信息级别日志
func (l *LoggerImpl) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Info(msg, fields...)
}

// Warn 记录警告级别日志
func (l *LoggerImpl) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Warn(msg, fields...)
}

// Error 记录错误级别日志
func (l *LoggerImpl) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Error(msg, fields...)
}

// Fatal 记录致命错误日志并退出程序
func (l *LoggerImpl) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Fatal(msg, fields...)
}

// With 创建带有预设字段的子日志器
func (l *LoggerImpl) With(fields ...zap.Field) interfaces.Logger {
	return &LoggerImpl{
		logger: l.logger.With(fields...),
	}
}

// WithContext 从上下文中提取字段创建子日志器
func (l *LoggerImpl) WithContext(ctx context.Context) interfaces.Logger {
	return &LoggerImpl{
		logger: l.withContext(ctx),
	}
}

// Sync 同步日志缓冲区
func (l *LoggerImpl) Sync() error {
	return l.logger.Sync()
}

// GetZapLogger 获取底层的zap.Logger实例
func (l *LoggerImpl) GetZapLogger() *zap.Logger {
	return l.logger
}

// withContext 从上下文中提取 traceID 并创建带有 traceID 的 logger
func (l *LoggerImpl) withContext(ctx context.Context) *zap.Logger {
	traceID := GetTraceID(ctx)
	if traceID != "" {
		return l.logger.With(zap.String("traceID", traceID))
	}
	return l.logger
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
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// LoggerFactoryImpl 日志工厂实现
type LoggerFactoryImpl struct{}

// NewLoggerFactory 创建日志工厂
func NewLoggerFactory() interfaces.LoggerFactory {
	return &LoggerFactoryImpl{}
}

// CreateLogger 创建日志器
func (f *LoggerFactoryImpl) CreateLogger(config interfaces.LoggerConfig) (interfaces.Logger, error) {
	// 创建临时配置提供者
	tempProvider := &tempConfigProvider{loggerConfig: config}
	return NewLogger(tempProvider)
}

// CreateDevelopmentLogger 创建开发环境日志器
func (f *LoggerFactoryImpl) CreateDevelopmentLogger() (interfaces.Logger, error) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	
	return &LoggerImpl{logger: zapLogger}, nil
}

// CreateProductionLogger 创建生产环境日志器
func (f *LoggerFactoryImpl) CreateProductionLogger(config interfaces.LoggerConfig) (interfaces.Logger, error) {
	return f.CreateLogger(config)
}

// tempConfigProvider 临时配置提供者，用于工厂方法
type tempConfigProvider struct {
	loggerConfig interfaces.LoggerConfig
}

func (t *tempConfigProvider) GetLoggerConfig() interfaces.LoggerConfig {
	return t.loggerConfig
}

func (t *tempConfigProvider) GetDatabaseConfig() interfaces.DatabaseConfig { return interfaces.DatabaseConfig{} }
func (t *tempConfigProvider) GetServerConfig() interfaces.ServerConfig     { return interfaces.ServerConfig{} }
func (t *tempConfigProvider) GetAuthConfig() interfaces.AuthConfig         { return interfaces.AuthConfig{} }
func (t *tempConfigProvider) GetRedisConfig() interfaces.RedisConfig       { return interfaces.RedisConfig{} }
func (t *tempConfigProvider) GetTokenConfig() interfaces.TokenConfig       { return interfaces.TokenConfig{} }
func (t *tempConfigProvider) GetValidationConfig() interfaces.ValidationConfig { return interfaces.ValidationConfig{} }
func (t *tempConfigProvider) GetRateLimitConfig() interfaces.RateLimitConfig { return interfaces.RateLimitConfig{} }
func (t *tempConfigProvider) Reload() error                                { return nil }
func (t *tempConfigProvider) GetEnv() string                               { return "development" }

// Module FX模块
var Module = fx.Module("logger",
	fx.Provide(NewLogger, NewLoggerFactory),
	fx.Invoke(func(logger interfaces.Logger) {
		// 设置全局logger实例
		SetGlobalLogger(logger)
	}),
)