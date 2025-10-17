package response

import (
	"net/http"
	"sync"
)

// LazyErrorMapper 延迟初始化的错误映射器
type LazyErrorMapper struct {
	mappings map[ErrorType]*ErrorMapping
	once     sync.Once
	mu       sync.RWMutex
}

// NewLazyErrorMapper 创建新的延迟错误映射器
func NewLazyErrorMapper() *LazyErrorMapper {
	return &LazyErrorMapper{
		mappings: make(map[ErrorType]*ErrorMapping),
	}
}

// GetMapping 获取错误类型对应的映射信息
// 使用延迟初始化，只在首次访问时初始化映射
func (m *LazyErrorMapper) GetMapping(errorType ErrorType) (*ErrorMapping, bool) {
	// 使用 sync.Once 确保映射只初始化一次
	m.once.Do(m.initMappings)

	m.mu.RLock()
	defer m.mu.RUnlock()

	mapping, exists := m.mappings[errorType]
	return mapping, exists
}

// RegisterMapping 动态注册错误映射
func (m *LazyErrorMapper) RegisterMapping(errorType ErrorType, mapping *ErrorMapping) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mappings[errorType] = mapping
}

// initMappings 初始化默认错误映射
// 这个方法只会被 sync.Once 调用一次
func (m *LazyErrorMapper) initMappings() {
	// 使用写锁保护初始化过程
	m.mu.Lock()
	defer m.mu.Unlock()

	// 初始化所有默认映射
	m.mappings[ErrorTypeNotFound] = &ErrorMapping{
		BusinessCode:   CodeNotFound,
		HTTPStatus:     http.StatusNotFound,
		DefaultMessage: "资源不存在",
	}
	m.mappings[ErrorTypeValidationFailed] = &ErrorMapping{
		BusinessCode:   CodeValidation,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "验证失败",
	}
	m.mappings[ErrorTypeAlreadyExists] = &ErrorMapping{
		BusinessCode:   CodeAlreadyExists,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "资源已存在",
	}
	m.mappings[ErrorTypeUnauthorized] = &ErrorMapping{
		BusinessCode:   CodeUnauthorized,
		HTTPStatus:     http.StatusUnauthorized,
		DefaultMessage: "未授权访问",
	}
	m.mappings[ErrorTypeForbidden] = &ErrorMapping{
		BusinessCode:   CodeForbidden,
		HTTPStatus:     http.StatusForbidden,
		DefaultMessage: "禁止访问",
	}
	m.mappings[ErrorTypeBusinessRuleViolation] = &ErrorMapping{
		BusinessCode:   CodeBusinessError,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "业务规则违反",
	}
	m.mappings[ErrorTypeConcurrencyConflict] = &ErrorMapping{
		BusinessCode:   CodeConflict,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "并发冲突",
	}
	m.mappings[ErrorTypeResourceLocked] = &ErrorMapping{
		BusinessCode:   CodeConflict,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "资源已锁定",
	}
	m.mappings[ErrorTypeInvalidData] = &ErrorMapping{
		BusinessCode:   CodeValidation,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "无效的数据",
	}
	m.mappings[ErrorTypeCommandValidation] = &ErrorMapping{
		BusinessCode:   CodeValidation,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "命令验证失败",
	}
	m.mappings[ErrorTypeCommandExecution] = &ErrorMapping{
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "命令执行失败",
	}
	m.mappings[ErrorTypeQueryExecution] = &ErrorMapping{
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "查询执行失败",
	}
	m.mappings[ErrorTypeInternalServer] = &ErrorMapping{
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "内部服务器错误",
	}
	m.mappings[ErrorTypeInvalidRequest] = &ErrorMapping{
		BusinessCode:   CodeBadRequest,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "无效的请求",
	}
	m.mappings[ErrorTypeDatabaseConnection] = &ErrorMapping{
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "数据库连接失败",
	}
	m.mappings[ErrorTypeRecordNotFound] = &ErrorMapping{
		BusinessCode:   CodeNotFound,
		HTTPStatus:     http.StatusNotFound,
		DefaultMessage: "记录不存在",
	}
	m.mappings[ErrorTypeDuplicateKey] = &ErrorMapping{
		BusinessCode:   CodeAlreadyExists,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "重复键值",
	}
	m.mappings[ErrorTypeExternalServiceUnavailable] = &ErrorMapping{
		BusinessCode:   CodeThirdParty,
		HTTPStatus:     http.StatusBadGateway,
		DefaultMessage: "外部服务不可用",
	}
	m.mappings[ErrorTypeTimeout] = &ErrorMapping{
		BusinessCode:   CodeTimeout,
		HTTPStatus:     http.StatusRequestTimeout,
		DefaultMessage: "请求超时",
	}
	m.mappings[ErrorTypeNetworkError] = &ErrorMapping{
		BusinessCode:   CodeThirdParty,
		HTTPStatus:     http.StatusBadGateway,
		DefaultMessage: "网络错误",
	}
}

// GetAllMappings 获取所有映射（用于调试和测试）
func (m *LazyErrorMapper) GetAllMappings() map[ErrorType]*ErrorMapping {
	// 确保映射已初始化
	m.once.Do(m.initMappings)

	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回映射的副本以防止外部修改
	result := make(map[ErrorType]*ErrorMapping)
	for k, v := range m.mappings {
		result[k] = v
	}
	return result
}

// IsInitialized 检查映射是否已初始化（用于测试）
func (m *LazyErrorMapper) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.mappings) > 0
}

// 全局延迟错误映射器实例
var defaultLazyErrorMapper = NewLazyErrorMapper()

// GetDefaultLazyErrorMapper 获取默认的延迟错误映射器
func GetDefaultLazyErrorMapper() *LazyErrorMapper {
	return defaultLazyErrorMapper
}
