package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 错误类型枚举
type ErrorType int

const (
	ErrorTypeValidation ErrorType = iota // 验证错误
	ErrorTypeSystem                      // 系统错误
	ErrorTypeAuth                        // 认证错误
	ErrorTypeThirdParty                  // 第三方服务错误
)

// BaseResponse 基础响应结构
type BaseResponse struct {
	Code    int         `json:"code"`             // 业务状态码
	Message string      `json:"message"`          // 响应消息
	Data    interface{} `json:"data,omitempty"`   // 响应数据
	Errors  interface{} `json:"errors,omitempty"` // 错误详情
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`        // 当前页
	PageSize   int   `json:"page_size"`   // 每页大小
	Total      int64 `json:"total"`       // 总数量
	TotalPages int   `json:"total_pages"` // 总页数
}

// PageData 分页数据结构
type PageData struct {
	Items interface{} `json:"items"`
	*Pagination
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
		return http.StatusBadRequest // 400
	case ErrorTypeAuth:
		return http.StatusUnauthorized // 401
	case ErrorTypeSystem:
		return http.StatusInternalServerError // 500
	case ErrorTypeThirdParty:
		return http.StatusBadGateway // 502
	default:
		return http.StatusInternalServerError // 500
	}
}

// buildSuccessResponse 构建成功响应的通用方法
func buildSuccessResponse(data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:    CodeSuccess,
		Message: GetCodeMessage(CodeSuccess),
		Data:    data,
	}
}

// Success 统一的成功响应方法
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, buildSuccessResponse(data))
}

// SuccessWithPagination 分页成功响应
func SuccessWithPagination(c *gin.Context, data interface{}, pagination *Pagination) {
	pageData := PageData{
		Items:      data,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, buildSuccessResponse(pageData))
}

// Error 统一的错误响应方法（使用AppError）
func Error(c *gin.Context, err *AppError) {
	httpStatus := getHTTPStatusByErrorType(err.Type)
	c.JSON(httpStatus, &BaseResponse{
		Code:    err.Code,
		Message: err.Message,
		Errors:  err.Errors,
	})
}

// --- 快捷错误响应函数

// BadRequest 400错误响应（参数格式错误、缺失必填字段等）
func BadRequest(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeValidation, CodeInvalidParams, message, nil))
}

// ValidationError 验证错误响应
func ValidationError(c *gin.Context, message string, errors interface{}) {
	Error(c, NewAppError(ErrorTypeValidation, CodeValidationError, message, errors))
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeAuth, CodeUnauthorized, message, nil))
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeAuth, CodeForbidden, message, nil))
}

// InternalServerError 500错误响应
func InternalServerError(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeSystem, CodeInternalError, message, nil))
}

// ThirdPartyError 502错误响应
func ThirdPartyError(c *gin.Context, message string) {
	Error(c, NewAppError(ErrorTypeThirdParty, CodeThirdPartyError, message, nil))
}
