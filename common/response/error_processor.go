package response

import (
	"net/http"
)

// ErrorProcessor 错误处理器
type ErrorProcessor struct {
	errorMappings map[string]*ErrorMapping
}

// NewErrorProcessor 创建新的错误处理器
func NewErrorProcessor() *ErrorProcessor {
	return &ErrorProcessor{
		errorMappings: CommonErrorMappings,
	}
}

// Process 处理错误并返回错误结果
func (ep *ErrorProcessor) Process(err error) *ErrorResult {
	if err == nil {
		return &ErrorResult{
			Code:       CodeSuccess,
			Message:    "操作成功",
			HTTPStatus: http.StatusOK,
		}
	}

	// 检查是否为自定义错误接口
	if customErr, ok := err.(CustomError); ok {
		return ep.processCustomError(customErr)
	}

	// 默认处理为内部服务器错误
	return &ErrorResult{
		Code:       CodeInternalError,
		Message:    err.Error(),
		HTTPStatus: http.StatusInternalServerError,
	}
}

// ProcessWithCode 使用指定业务码处理错误
func (ep *ErrorProcessor) ProcessWithCode(code int, message string) *ErrorResult {
	// 获取业务码信息
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

	// 如果业务码不存在，使用默认处理
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

// ProcessWithData 处理带有额外数据的错误
func (ep *ErrorProcessor) ProcessWithData(err error, data interface{}) *ErrorResult {
	result := ep.Process(err)
	result.Data = data
	return result
}

// processCustomError 处理自定义错误
func (ep *ErrorProcessor) processCustomError(customErr CustomError) *ErrorResult {
	errorType := customErr.ErrorType()
	message := customErr.Error()

	// 查找错误映射
	if mapping, exists := ep.errorMappings[errorType]; exists {
		// 如果自定义错误没有消息，使用映射的默认消息
		if message == "" {
			message = mapping.DefaultMessage
		}

		return &ErrorResult{
			Code:       mapping.BusinessCode,
			Message:    message,
			HTTPStatus: mapping.HTTPStatus,
		}
	}

	// 如果没有找到映射，使用默认的业务错误处理
	if message == "" {
		message = "业务处理失败"
	}

	return &ErrorResult{
		Code:       CodeBusinessError,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// CustomError 自定义错误接口
// 各个服务可以实现这个接口来提供错误类型信息
type CustomError interface {
	error
	ErrorType() string // 返回错误类型字符串，用于映射查找
}

// SimpleCustomError 简单的自定义错误实现
type SimpleCustomError struct {
	Type    string
	Message string
}

func (e *SimpleCustomError) Error() string {
	return e.Message
}

func (e *SimpleCustomError) ErrorType() string {
	return e.Type
}

// NewCustomError 创建新的自定义错误
func NewCustomError(errorType, message string) *SimpleCustomError {
	return &SimpleCustomError{
		Type:    errorType,
		Message: message,
	}
}
