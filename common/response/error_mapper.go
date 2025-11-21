package response

import (
	"net/http"
	"sync"
)

// errorMapper 错误映射实现
type errorMapper struct {
	mappings map[ErrorType]*ErrorMapping
	mu       sync.RWMutex
}

// newErrorMapper 创建新的错误映射器
func newErrorMapper() ErrorMapper {
	mapper := &errorMapper{
		mappings: make(map[ErrorType]*ErrorMapping),
	}
	
	// 初始化默认错误映射
	mapper.initDefaultMappings()
	
	return mapper
}

// initDefaultMappings 初始化默认错误映射
func (em *errorMapper) initDefaultMappings() {
	defaultMappings := map[ErrorType]*ErrorMapping{
		ErrorTypeNotFound: {
			BusinessCode:   CodeNotFound,
			HTTPStatus:     http.StatusNotFound,
			DefaultMessage: "资源不存在",
		},
		ErrorTypeValidationFailed: {
			BusinessCode:   CodeValidation,
			HTTPStatus:     http.StatusBadRequest,
			DefaultMessage: "验证失败",
		},
		ErrorTypeAlreadyExists: {
			BusinessCode:   CodeAlreadyExists,
			HTTPStatus:     http.StatusConflict,
			DefaultMessage: "资源已存在",
		},
		ErrorTypeUnauthorized: {
			BusinessCode:   CodeUnauthorized,
			HTTPStatus:     http.StatusUnauthorized,
			DefaultMessage: "未授权访问",
		},
		ErrorTypeForbidden: {
			BusinessCode:   CodeForbidden,
			HTTPStatus:     http.StatusForbidden,
			DefaultMessage: "禁止访问",
		},
		ErrorTypeBusinessRuleViolation: {
			BusinessCode:   CodeBusinessError,
			HTTPStatus:     http.StatusBadRequest,
			DefaultMessage: "业务规则违反",
		},
		ErrorTypeConcurrencyConflict: {
			BusinessCode:   CodeConflict,
			HTTPStatus:     http.StatusConflict,
			DefaultMessage: "并发冲突",
		},
		ErrorTypeResourceLocked: {
			BusinessCode:   CodeConflict,
			HTTPStatus:     http.StatusConflict,
			DefaultMessage: "资源已锁定",
		},
		ErrorTypeInvalidData: {
			BusinessCode:   CodeValidation,
			HTTPStatus:     http.StatusBadRequest,
			DefaultMessage: "无效的数据",
		},
		ErrorTypeCommandValidation: {
			BusinessCode:   CodeValidation,
			HTTPStatus:     http.StatusBadRequest,
			DefaultMessage: "命令验证失败",
		},
		ErrorTypeCommandExecution: {
			BusinessCode:   CodeInternalError,
			HTTPStatus:     http.StatusInternalServerError,
			DefaultMessage: "命令执行失败",
		},
		ErrorTypeQueryExecution: {
			BusinessCode:   CodeInternalError,
			HTTPStatus:     http.StatusInternalServerError,
			DefaultMessage: "查询执行失败",
		},
		ErrorTypeInternalServer: {
			BusinessCode:   CodeInternalError,
			HTTPStatus:     http.StatusInternalServerError,
			DefaultMessage: "内部服务器错误",
		},
		ErrorTypeInvalidRequest: {
			BusinessCode:   CodeBadRequest,
			HTTPStatus:     http.StatusBadRequest,
			DefaultMessage: "无效的请求",
		},
		ErrorTypeDatabaseConnection: {
			BusinessCode:   CodeInternalError,
			HTTPStatus:     http.StatusInternalServerError,
			DefaultMessage: "数据库连接失败",
		},
		ErrorTypeRecordNotFound: {
			BusinessCode:   CodeNotFound,
			HTTPStatus:     http.StatusNotFound,
			DefaultMessage: "记录不存在",
		},
		ErrorTypeDuplicateKey: {
			BusinessCode:   CodeAlreadyExists,
			HTTPStatus:     http.StatusConflict,
			DefaultMessage: "重复键值",
		},
		ErrorTypeExternalServiceUnavailable: {
			BusinessCode:   CodeThirdParty,
			HTTPStatus:     http.StatusBadGateway,
			DefaultMessage: "外部服务不可用",
		},
		ErrorTypeTimeout: {
			BusinessCode:   CodeTimeout,
			HTTPStatus:     http.StatusRequestTimeout,
			DefaultMessage: "请求超时",
		},
		ErrorTypeNetworkError: {
			BusinessCode:   CodeThirdParty,
			HTTPStatus:     http.StatusBadGateway,
			DefaultMessage: "网络错误",
		},
	}
	
	for errorType, mapping := range defaultMappings {
		em.mappings[errorType] = mapping
	}
}

// GetErrorMapping 获取错误类型对应的映射信息
func (em *errorMapper) GetErrorMapping(errorType ErrorType) (*ErrorMapping, bool) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	mapping, exists := em.mappings[errorType]
	return mapping, exists
}

// RegisterErrorMapping 注册错误映射
func (em *errorMapper) RegisterErrorMapping(errorType ErrorType, mapping *ErrorMapping) {
	if mapping == nil {
		return
	}

	em.mu.Lock()
	defer em.mu.Unlock()

	em.mappings[errorType] = mapping
}