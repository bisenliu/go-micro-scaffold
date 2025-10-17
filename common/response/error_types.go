package response

import (
	"fmt"
	"net/http"
)

// ErrorType 错误类型枚举
type ErrorType int

const (
	ErrorTypeNotFound ErrorType = iota
	ErrorTypeValidationFailed
	ErrorTypeAlreadyExists
	ErrorTypeUnauthorized
	ErrorTypeForbidden
	ErrorTypeBusinessRuleViolation
	ErrorTypeConcurrencyConflict
	ErrorTypeResourceLocked
	ErrorTypeInvalidData
	ErrorTypeCommandValidation
	ErrorTypeCommandExecution
	ErrorTypeQueryExecution
	ErrorTypeInternalServer
	ErrorTypeInvalidRequest
	ErrorTypeDatabaseConnection
	ErrorTypeRecordNotFound
	ErrorTypeDuplicateKey
	ErrorTypeExternalServiceUnavailable
	ErrorTypeTimeout
	ErrorTypeNetworkError
)

// ErrorMapping 错误映射结构
type ErrorMapping struct {
	BusinessCode   int
	HTTPStatus     int
	DefaultMessage string
}

// DomainError 领域错误结构
type DomainError struct {
	Type    ErrorType
	Message string
	BaseErr error
	Context map[string]interface{}
}

// Error 实现error接口
func (e *DomainError) Error() string {
	if e.BaseErr != nil {
		return fmt.Sprintf("%s: %s", e.BaseErr.Error(), e.Message)
	}
	return e.Message
}

// Unwrap 返回被包装的底层错误
func (e *DomainError) Unwrap() error {
	return e.BaseErr
}

func (e *DomainError) WithContext(key string, value interface{}) *DomainError {
	// 创建新实例，避免修改原始错误
	newErr := &DomainError{
		Type:    e.Type,
		Message: e.Message,
		BaseErr: e.BaseErr,
		Context: make(map[string]interface{}),
	}

	// 复制现有上下文
	if e.Context != nil {
		for k, v := range e.Context {
			newErr.Context[k] = v
		}
	}

	// 添加新的上下文
	newErr.Context[key] = value

	return newErr
}

// WithContextMap 批量添加上下文信息
func (e *DomainError) WithContextMap(contextMap map[string]interface{}) *DomainError {
	newErr := &DomainError{
		Type:    e.Type,
		Message: e.Message,
		BaseErr: e.BaseErr,
		Context: make(map[string]interface{}),
	}

	// 复制现有上下文
	if e.Context != nil {
		for k, v := range e.Context {
			newErr.Context[k] = v
		}
	}

	// 添加新的上下文
	for k, v := range contextMap {
		newErr.Context[k] = v
	}

	return newErr
}

// NewDomainError 创建新的领域错误
func NewDomainError(errorType ErrorType, message string) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
		BaseErr: nil,
		Context: make(map[string]interface{}),
	}
}

// NewDomainErrorWithCause 创建带有原因的领域错误
func NewDomainErrorWithCause(errorType ErrorType, message string, cause error) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
		BaseErr: cause,
		Context: make(map[string]interface{}),
	}
}

// ErrorMapper 错误映射器
type ErrorMapper struct {
	mappings map[ErrorType]*ErrorMapping
}

// NewErrorMapper 创建新的错误映射器
func NewErrorMapper() *ErrorMapper {
	mapper := &ErrorMapper{
		mappings: make(map[ErrorType]*ErrorMapping),
	}
	mapper.initDefaultMappings()
	return mapper
}

// initDefaultMappings 初始化默认错误映射
func (m *ErrorMapper) initDefaultMappings() {
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

// GetMapping 根据错误类型获取错误映射
func (m *ErrorMapper) GetMapping(errorType ErrorType) (*ErrorMapping, bool) {
	mapping, exists := m.mappings[errorType]
	return mapping, exists
}

// ErrorHandler 统一错误处理器
type ErrorHandler struct {
	mapper *ErrorMapper
}

// NewErrorHandler 创建新的错误处理器
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		mapper: NewErrorMapper(),
	}
}

// Handle 处理错误并返回错误结果
func (h *ErrorHandler) Handle(err error) *ErrorResult {
	if err == nil {
		return &ErrorResult{
			Code:       CodeSuccess,
			Message:    "操作成功",
			HTTPStatus: http.StatusOK,
		}
	}

	// 检查是否为DomainError
	if domainErr, ok := err.(*DomainError); ok {
		return h.handleDomainError(domainErr)
	}

	// 默认处理为内部服务器错误
	return &ErrorResult{
		Code:       CodeInternalError,
		Message:    err.Error(),
		HTTPStatus: http.StatusInternalServerError,
	}
}

// HandleWithCode 使用指定业务码处理错误
func (h *ErrorHandler) HandleWithCode(code int, message string) *ErrorResult {
	if codeInfo, exists := GetCodeInfo(code); exists {
		finalMessage := message
		if finalMessage == "" {
			finalMessage = codeInfo.Message
		}
		return &ErrorResult{
			Code:       code,
			Message:    finalMessage,
			HTTPStatus: codeInfo.HTTPStatus,
		}
	}

	finalMessage := message
	if finalMessage == "" {
		finalMessage = "未知错误"
	}
	return &ErrorResult{
		Code:       code,
		Message:    finalMessage,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// HandleWithData 处理带有额外数据的错误
func (h *ErrorHandler) HandleWithData(err error, data interface{}) *ErrorResult {
	result := h.Handle(err)
	result.Data = data
	return result
}

// handleDomainError 处理领域错误
func (h *ErrorHandler) handleDomainError(err *DomainError) *ErrorResult {
	if mapping, exists := h.mapper.GetMapping(err.Type); exists {
		message := err.Message
		if message == "" {
			message = mapping.DefaultMessage
		}
		return &ErrorResult{
			Code:       mapping.BusinessCode,
			Message:    message,
			HTTPStatus: mapping.HTTPStatus,
		}
	}

	message := err.Message
	if message == "" {
		message = "业务处理失败"
	}
	return &ErrorResult{
		Code:       CodeBusinessError,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}
