package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 错误类型枚举
type ErrorType int

const (
	ErrorTypeValidation ErrorType = iota // 验证错误
	ErrorTypeBusiness                    // 业务错误
	ErrorTypeSystem                      // 系统错误
	ErrorTypeAuth                        // 认证错误
)

// BaseResponse 基础响应结构 (保持与您的版本一致，已包含 Details)
type BaseResponse struct {
	Code    int         `json:"code"`             // 业务状态码
	Message string      `json:"message"`          // 响应消息
	Data    interface{} `json:"data,omitempty"`   // 响应数据
	Errors  interface{} `json:"errors,omitempty"` // 错误详情(用于验证错误等)
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`        // 当前页
	PageSize   int   `json:"page_size"`   // 每页大小
	Total      int64 `json:"total"`       // 总数量
	TotalPages int   `json:"total_pages"` // 总页数
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// AppError 统一错误结构
type AppError struct {
	Type    ErrorType   `json:"-"`                // 错误类型（不序列化）
	Code    int         `json:"code"`             // 业务错误码
	Message string      `json:"message"`          // 错误消息
	Errors  interface{} `json:"errors,omitempty"` // 错误详情
}

func (e *AppError) Error() string {
	return e.Message
}

// 预定义错误
var (
	ErrInvalidParams = &AppError{
		Type: ErrorTypeValidation,
		Code: CodeInvalidParams,
	}

	ErrUnauthorized = &AppError{
		Type: ErrorTypeAuth,
		Code: CodeUnauthorized,
	}

	ErrForbidden = &AppError{
		Type: ErrorTypeAuth,
		Code: CodeForbidden,
	}

	ErrNotFound = &AppError{
		Type: ErrorTypeBusiness,
		Code: CodeNotFound,
	}

	ErrInternalError = &AppError{
		Type: ErrorTypeSystem,
		Code: CodeInternalError,
	}

	ErrServiceUnavailable = &AppError{
		Type: ErrorTypeSystem,
		Code: CodeServiceUnavailable,
	}
)

// NewAppError 创建自定义错误
// 如果 message 为空，将自动从 codes.go 中获取对应的消息
func NewAppError(errorType ErrorType, code int, message string, errors interface{}) *AppError {
	if message == "" {
		message = GetCodeMessage(code)
	}

	return &AppError{
		Type:    errorType,
		Code:    code,
		Message: message,
		Errors:  errors,
	}
}

// getHTTPStatusByErrorType 根据错误类型获取HTTP状态码
func getHTTPStatusByErrorType(errorType ErrorType) int {
	switch errorType {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeAuth:
		return http.StatusUnauthorized
	case ErrorTypeBusiness:
		return http.StatusOK // 业务错误通常返回200，通过业务码区分
	case ErrorTypeSystem:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Success 统一的成功响应方法
func Success(c *gin.Context, data interface{}) {
	response := &BaseResponse{
		Code:    CodeSuccess,
		Message: GetCodeMessage(CodeSuccess),
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// SuccessWithPagination 分页成功响应
func SuccessWithPagination(c *gin.Context, data interface{}, pagination *Pagination) {
	response := &PageResponse{
		Code: CodeSuccess,
		// 优化：使用 GetCodeMessage 替代硬编码 "success"
		Message:    GetCodeMessage(CodeSuccess),
		Data:       data,
		Pagination: pagination,
	}
	c.JSON(http.StatusOK, response)
}

// Error 统一的错误响应方法（使用AppError）
func Error(c *gin.Context, err *AppError) {
	// 注意：NewAppError 已经处理了 Message 为空的情况，
	// 所以这里不再需要重复检查和获取 Message，只需确保 AppError 构造正确。

	httpStatus := getHTTPStatusByErrorType(err.Type)

	response := &BaseResponse{
		Code:    err.Code,
		Message: err.Message,
		Errors:  err.Errors,
	}

	c.JSON(httpStatus, response)
}

// --- 快捷错误响应函数

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeValidation, CodeInvalidParams, message, nil))
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeAuth, CodeUnauthorized, message, nil))
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeAuth, CodeForbidden, message, nil))
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeBusiness, CodeNotFound, message, nil))
}

// InternalServerError 500错误响应
func InternalServerError(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeSystem, CodeInternalError, message, nil))
}

// BusinessError 业务错误响应
func BusinessError(c *gin.Context, code int, message string) {
	Error(c, NewAppError(ErrorTypeBusiness, code, message, nil))
}

// ValidationError 验证错误响应
func ValidationError(c *gin.Context, message string, errors map[string]string) {
	Error(c, NewAppError(ErrorTypeValidation, CodeValidationError, message, errors))
}
