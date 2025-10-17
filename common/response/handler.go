package response

import (
	"net/http"
)

// UnifiedErrorHandler 统一的错误处理器接口
type UnifiedErrorHandler interface {
	// HandleError 处理错误并返回错误结果
	HandleError(err error, options ...ErrorHandleOption) *ErrorResult
}

// ErrorHandleOption 错误处理选项
type ErrorHandleOption func(*ErrorHandleConfig)

// ErrorHandleConfig 错误处理配置
type ErrorHandleConfig struct {
	Code    *int   // 指定的业务码
	Message string // 自定义消息
	Data    any    // 额外数据
}

// WithCode 指定业务码选项
func WithCode(code int) ErrorHandleOption {
	return func(config *ErrorHandleConfig) {
		config.Code = &code
	}
}

// WithMessage 指定自定义消息选项
func WithMessage(message string) ErrorHandleOption {
	return func(config *ErrorHandleConfig) {
		config.Message = message
	}
}

// WithData 指定额外数据选项
func WithData(data any) ErrorHandleOption {
	return func(config *ErrorHandleConfig) {
		config.Data = data
	}
}

// unifiedErrorHandler 统一错误处理器实现
type unifiedErrorHandler struct {
	mapper  ErrorMapper
	factory ErrorFactory
}

// NewUnifiedErrorHandler 创建新的统一错误处理器
func NewUnifiedErrorHandler() UnifiedErrorHandler {
	return &unifiedErrorHandler{
		mapper:  NewLazyErrorMapper(),
		factory: GetDefaultErrorFactory(),
	}
}

// NewUnifiedErrorHandlerWithDependencies 使用指定依赖创建统一错误处理器
func NewUnifiedErrorHandlerWithDependencies(mapper ErrorMapper, factory ErrorFactory) UnifiedErrorHandler {
	return &unifiedErrorHandler{
		mapper:  mapper,
		factory: factory,
	}
}

// HandleError 统一的错误处理方法
func (h *unifiedErrorHandler) HandleError(err error, options ...ErrorHandleOption) *ErrorResult {
	// 解析选项
	config := &ErrorHandleConfig{}
	for _, option := range options {
		option(config)
	}

	// 如果没有错误且没有指定业务码，返回成功
	if err == nil && config.Code == nil {
		return &ErrorResult{
			Code:       CodeSuccess,
			Message:    "操作成功",
			HTTPStatus: http.StatusOK,
			Data:       config.Data,
		}
	}

	// 如果指定了业务码，使用业务码处理
	if config.Code != nil {
		return h.handleWithCode(*config.Code, config.Message, config.Data)
	}

	// 处理错误对象
	return h.handleErrorObject(err, config)
}

// handleWithCode 使用指定业务码处理
func (h *unifiedErrorHandler) handleWithCode(code int, message string, data any) *ErrorResult {
	// 查找代码信息
	if codeInfo, exists := GetCodeInfo(code); exists {
		finalMessage := message
		if finalMessage == "" {
			finalMessage = codeInfo.Message
		}
		return &ErrorResult{
			Code:       code,
			Message:    finalMessage,
			HTTPStatus: codeInfo.HTTPStatus,
			Data:       data,
		}
	}

	// 如果代码信息不存在，使用默认处理
	finalMessage := message
	if finalMessage == "" {
		finalMessage = "未知错误"
	}
	return &ErrorResult{
		Code:       code,
		Message:    finalMessage,
		HTTPStatus: http.StatusInternalServerError,
		Data:       data,
	}
}

// handleErrorObject 处理错误对象
func (h *unifiedErrorHandler) handleErrorObject(err error, config *ErrorHandleConfig) *ErrorResult {
	if err == nil {
		return &ErrorResult{
			Code:       CodeInternalError,
			Message:    "内部错误：空错误对象",
			HTTPStatus: http.StatusInternalServerError,
			Data:       config.Data,
		}
	}

	// 检查是否为DomainError
	if domainErr, ok := err.(*DomainError); ok {
		return h.handleDomainError(domainErr, config)
	}

	// 默认处理为内部服务器错误
	message := config.Message
	if message == "" {
		message = err.Error()
	}

	return &ErrorResult{
		Code:       CodeInternalError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Data:       config.Data,
	}
}

// handleDomainError 处理领域错误
func (h *unifiedErrorHandler) handleDomainError(err *DomainError, config *ErrorHandleConfig) *ErrorResult {
	// 查找错误映射
	if mapping, exists := h.mapper.GetMapping(err.Type); exists {
		message := config.Message
		if message == "" {
			message = err.Message
		}
		if message == "" {
			message = mapping.DefaultMessage
		}

		return &ErrorResult{
			Code:       mapping.BusinessCode,
			Message:    message,
			HTTPStatus: mapping.HTTPStatus,
			Data:       config.Data,
		}
	}

	// 如果没有找到映射，使用默认处理
	message := config.Message
	if message == "" {
		message = err.Message
	}
	if message == "" {
		message = "业务处理失败"
	}

	return &ErrorResult{
		Code:       CodeBusinessError,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Data:       config.Data,
	}
}

// 全局统一错误处理器实例
var defaultUnifiedErrorHandler = NewUnifiedErrorHandler()

// GetDefaultUnifiedErrorHandler 获取默认的统一错误处理器
func GetDefaultUnifiedErrorHandler() UnifiedErrorHandler {
	return defaultUnifiedErrorHandler
}

// 便捷函数：使用默认统一处理器

// HandleError 使用默认统一处理器处理错误
func HandleError(err error, options ...ErrorHandleOption) *ErrorResult {
	return defaultUnifiedErrorHandler.HandleError(err, options...)
}

// HandleErrorWithCode 使用指定业务码处理错误
func HandleErrorWithCode(code int, message string) *ErrorResult {
	return defaultUnifiedErrorHandler.HandleError(nil, WithCode(code), WithMessage(message))
}

// HandleErrorWithData 处理带有额外数据的错误
func HandleErrorWithData(err error, data any) *ErrorResult {
	return defaultUnifiedErrorHandler.HandleError(err, WithData(data))
}
