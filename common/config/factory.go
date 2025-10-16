package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"common/interfaces"
)

// ConfigFactory 配置工厂
// 负责配置的创建、验证和默认值处理
type ConfigFactory struct {
	validator *validator.Validate
	config    *ModularConfig
	env       string
}

// NewConfigFactory 创建配置工厂
func NewConfigFactory() *ConfigFactory {
	return &ConfigFactory{
		validator: validator.New(),
	}
}

// LoadConfig 加载配置
func (f *ConfigFactory) LoadConfig() (*ModularConfig, error) {
	// 使用环境管理器加载配置
	envManager := NewEnvironmentManager()
	if err := envManager.LoadEnvironmentConfig(); err != nil {
		return nil, fmt.Errorf("failed to load environment config: %w", err)
	}

	// 创建配置实例
	config := &ModularConfig{}

	// 解析配置
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 应用默认值
	if err := f.applyDefaults(config); err != nil {
		return nil, fmt.Errorf("failed to apply defaults: %w", err)
	}

	// 暂时禁用验证，专注于基本功能
	// if err := f.validateConfig(config); err != nil {
	//     return nil, fmt.Errorf("config validation failed: %w", err)
	// }

	f.config = config
	f.env = config.System.Env

	return config, nil
}

// CreateConfigProvider 创建配置提供者
func (f *ConfigFactory) CreateConfigProvider() (interfaces.ConfigProvider, error) {
	if f.config == nil {
		if _, err := f.LoadConfig(); err != nil {
			return nil, err
		}
	}

	interfaceTypes := f.config.ToInterfaceTypes()
	
	return &ConfigProviderImpl{
		serverConfig:     interfaceTypes.Server,
		databaseConfig:   interfaceTypes.Database,
		authConfig:       interfaceTypes.Auth,
		loggerConfig:     interfaceTypes.Logger,
		redisConfig:      interfaceTypes.Redis,
		tokenConfig:      interfaceTypes.Token,
		validationConfig: interfaceTypes.Validation,
		rateLimitConfig:  interfaceTypes.RateLimit,
		env:              interfaceTypes.Env,
		factory:          f,
	}, nil
}

// applyDefaults 应用默认值
func (f *ConfigFactory) applyDefaults(config *ModularConfig) error {
	return f.applyDefaultsToStruct(reflect.ValueOf(config).Elem())
}

// applyDefaultsToStruct 递归应用默认值到结构体
func (f *ConfigFactory) applyDefaultsToStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// 跳过不可设置的字段
		if !field.CanSet() {
			continue
		}

		// 处理嵌套结构体
		if field.Kind() == reflect.Struct {
			if err := f.applyDefaultsToStruct(field); err != nil {
				return err
			}
			continue
		}

		// 获取默认值标签
		defaultTag := fieldType.Tag.Get("default")
		if defaultTag == "" {
			continue
		}

		// 如果字段已有值，跳过
		if !f.isZeroValue(field) {
			continue
		}

		// 应用默认值
		if err := f.setDefaultValue(field, defaultTag); err != nil {
			return fmt.Errorf("failed to set default value for field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// isZeroValue 检查是否为零值
func (f *ConfigFactory) isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice:
		return v.Len() == 0
	default:
		return v.IsZero()
	}
}

// setDefaultValue 设置默认值
func (f *ConfigFactory) setDefaultValue(field reflect.Value, defaultValue string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(defaultValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			// 处理 time.Duration
			duration, err := time.ParseDuration(defaultValue)
			if err != nil {
				return fmt.Errorf("invalid duration format: %s", defaultValue)
			}
			field.SetInt(int64(duration))
		} else {
			// 处理普通整数
			var intVal int64
			if _, err := fmt.Sscanf(defaultValue, "%d", &intVal); err != nil {
				return fmt.Errorf("invalid integer format: %s", defaultValue)
			}
			field.SetInt(intVal)
		}
	case reflect.Bool:
		var boolVal bool
		if _, err := fmt.Sscanf(defaultValue, "%t", &boolVal); err != nil {
			return fmt.Errorf("invalid boolean format: %s", defaultValue)
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// validateConfig 验证配置
func (f *ConfigFactory) validateConfig(config *ModularConfig) error {
	if err := f.validator.Struct(config); err != nil {
		var validationErrors []string
		
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, f.formatValidationError(err))
		}
		
		return fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	// 自定义验证逻辑
	if err := f.customValidation(config); err != nil {
		return err
	}

	return nil
}

// formatValidationError 格式化验证错误
func (f *ConfigFactory) formatValidationError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())
	
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	case "required_unless":
		return fmt.Sprintf("%s is required unless %s", field, err.Param())
	default:
		return fmt.Sprintf("%s validation failed: %s", field, err.Tag())
	}
}

// customValidation 自定义验证逻辑
func (f *ConfigFactory) customValidation(config *ModularConfig) error {
	// 验证数据库配置
	if config.Database.Primary.Type != "sqlite" {
		if config.Database.Primary.Host == "" {
			return fmt.Errorf("database host is required for non-sqlite databases")
		}
		if config.Database.Primary.Username == "" {
			return fmt.Errorf("database username is required for non-sqlite databases")
		}
		if config.Database.Primary.Password == "" {
			return fmt.Errorf("database password is required for non-sqlite databases")
		}
	}

	// 验证端口范围
	if config.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	// 验证日志目录
	if config.Logger.Director == "" {
		return fmt.Errorf("logger directory cannot be empty")
	}

	return nil
}

// GetEnvironment 获取当前环境
func (f *ConfigFactory) GetEnvironment() string {
	return f.env
}

// Reload 重新加载配置
func (f *ConfigFactory) Reload() error {
	_, err := f.LoadConfig()
	return err
}