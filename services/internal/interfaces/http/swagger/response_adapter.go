package swagger

import (
	"net/http"

	"common/response"

	"github.com/gin-gonic/gin"
)

// ResponseAdapter 响应适配器
// 将系统内部的响应格式转换为Swagger文档中定义的标准格式
type ResponseAdapter struct {
	converter *ErrorResponseConverter
}

// NewResponseAdapter 创建响应适配器
func NewResponseAdapter() *ResponseAdapter {
	return &ResponseAdapter{
		converter: NewErrorResponseConverter(),
	}
}

// AdaptErrorResponse 适配错误响应
// 将系统内部错误转换为Swagger标准错误响应格式
func (a *ResponseAdapter) AdaptErrorResponse(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 检查是否为DomainError
	if domainErr, ok := err.(*response.DomainError); ok {
		a.handleDomainError(c, domainErr)
		return
	}

	// 默认处理为内部服务器错误
	errorResp := &InternalServerErrorResponse{
		Error:   "Internal Server Error",
		Message: err.Error(),
		Code:    500,
	}
	c.JSON(http.StatusInternalServerError, errorResp)
}

// handleDomainError 处理领域错误
func (a *ResponseAdapter) handleDomainError(c *gin.Context, domainErr *response.DomainError) {
	switch domainErr.Type {
	case response.ErrorTypeValidationFailed:
		a.handleValidationError(c, domainErr)
	case response.ErrorTypeNotFound, response.ErrorTypeRecordNotFound:
		a.handleNotFoundError(c, domainErr)
	case response.ErrorTypeUnauthorized:
		a.handleUnauthorizedError(c, domainErr)
	case response.ErrorTypeForbidden:
		a.handleForbiddenError(c, domainErr)
	case response.ErrorTypeAlreadyExists, response.ErrorTypeDuplicateKey:
		a.handleConflictError(c, domainErr)
	case response.ErrorTypeExternalServiceUnavailable:
		a.handleServiceUnavailableError(c, domainErr)
	default:
		a.handleGenericError(c, domainErr)
	}
}

// handleValidationError 处理验证错误
func (a *ResponseAdapter) handleValidationError(c *gin.Context, domainErr *response.DomainError) {
	// 尝试从上下文中获取字段错误信息
	var fieldErrors []FieldError
	if domainErr.HasContext() {
		if fields, exists := domainErr.GetContextValue("fields"); exists {
			if fieldList, ok := fields.([]FieldError); ok {
				fieldErrors = fieldList
			}
		}
	}

	// 如果没有具体的字段错误，创建一个通用的验证错误
	if len(fieldErrors) == 0 {
		fieldErrors = []FieldError{
			{
				Field:   "request",
				Message: domainErr.Message,
				Value:   "",
			},
		}
	}

	errorResp := &ValidationErrorResponse{
		Error:   "Validation Failed",
		Message: domainErr.Message,
		Code:    400,
		Details: ValidationErrorDetails{
			Fields: fieldErrors,
		},
	}
	c.JSON(http.StatusBadRequest, errorResp)
}

// handleNotFoundError 处理资源不存在错误
func (a *ResponseAdapter) handleNotFoundError(c *gin.Context, domainErr *response.DomainError) {
	errorResp := &NotFoundErrorResponse{
		Error:   "Not Found",
		Message: domainErr.Message,
		Code:    404,
	}
	c.JSON(http.StatusNotFound, errorResp)
}

// handleUnauthorizedError 处理未授权错误
func (a *ResponseAdapter) handleUnauthorizedError(c *gin.Context, domainErr *response.DomainError) {
	errorResp := &UnauthorizedErrorResponse{
		Error:   "Unauthorized",
		Message: domainErr.Message,
		Code:    401,
	}
	c.JSON(http.StatusUnauthorized, errorResp)
}

// handleForbiddenError 处理禁止访问错误
func (a *ResponseAdapter) handleForbiddenError(c *gin.Context, domainErr *response.DomainError) {
	errorResp := &ForbiddenErrorResponse{
		Error:   "Forbidden",
		Message: domainErr.Message,
		Code:    403,
	}
	c.JSON(http.StatusForbidden, errorResp)
}

// handleConflictError 处理资源冲突错误
func (a *ResponseAdapter) handleConflictError(c *gin.Context, domainErr *response.DomainError) {
	errorResp := &ConflictErrorResponse{
		Error:   "Conflict",
		Message: domainErr.Message,
		Code:    409,
	}
	c.JSON(http.StatusConflict, errorResp)
}

// handleServiceUnavailableError 处理服务不可用错误
func (a *ResponseAdapter) handleServiceUnavailableError(c *gin.Context, domainErr *response.DomainError) {
	errorResp := &ServiceUnavailableErrorResponse{
		Error:   "Service Unavailable",
		Message: domainErr.Message,
		Code:    503,
	}
	c.JSON(http.StatusServiceUnavailable, errorResp)
}

// handleGenericError 处理通用错误
func (a *ResponseAdapter) handleGenericError(c *gin.Context, domainErr *response.DomainError) {
	// 根据错误类型确定HTTP状态码
	var httpStatus int
	var errorType string
	var code int

	switch domainErr.Type {
	case response.ErrorTypeBusinessRuleViolation, response.ErrorTypeInvalidData, response.ErrorTypeInvalidRequest:
		httpStatus = http.StatusBadRequest
		errorType = "Bad Request"
		code = 400
	case response.ErrorTypeTimeout:
		httpStatus = http.StatusRequestTimeout
		errorType = "Request Timeout"
		code = 408
	case response.ErrorTypeConcurrencyConflict, response.ErrorTypeResourceLocked:
		httpStatus = http.StatusConflict
		errorType = "Conflict"
		code = 409
	default:
		httpStatus = http.StatusInternalServerError
		errorType = "Internal Server Error"
		code = 500
	}

	errorResp := &ErrorResponse{
		Error:   errorType,
		Message: domainErr.Message,
		Code:    code,
	}
	c.JSON(httpStatus, errorResp)
}

// AdaptSuccessResponse 适配成功响应
// 确保成功响应格式与Swagger文档一致
func (a *ResponseAdapter) AdaptSuccessResponse(c *gin.Context, data any) {
	// 成功响应保持原有格式，因为它们已经符合文档要求
	c.JSON(http.StatusOK, data)
}

// AdaptPagingResponse 适配分页响应
// 确保分页响应格式与Swagger文档一致
func (a *ResponseAdapter) AdaptPagingResponse(c *gin.Context, data any, page, pageSize int, total int64) {
	// 分页响应保持原有格式，因为它们已经符合文档要求
	response.HandlePaging(c, data, page, pageSize, total, nil)
}

// CreateValidationFieldError 创建验证字段错误
// 辅助函数，用于创建符合Swagger格式的字段验证错误
func CreateValidationFieldError(field, message, value string) FieldError {
	return FieldError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// CreateValidationErrorWithFields 创建带字段错误的验证错误
// 辅助函数，用于创建包含具体字段错误信息的DomainError
func CreateValidationErrorWithFields(message string, fieldErrors []FieldError) *response.DomainError {
	domainErr := response.CreateError(response.ErrorTypeValidationFailed, message)
	if len(fieldErrors) > 0 {
		domainErr = domainErr.WithContext("fields", fieldErrors)
	}
	return domainErr
}

// SwaggerResponseHelper Swagger响应辅助工具
type SwaggerResponseHelper struct {
	adapter *ResponseAdapter
}

// NewSwaggerResponseHelper 创建Swagger响应辅助工具
func NewSwaggerResponseHelper() *SwaggerResponseHelper {
	return &SwaggerResponseHelper{
		adapter: NewResponseAdapter(),
	}
}

// HandleWithSwaggerFormat 使用Swagger格式处理响应
// 这是一个便捷方法，确保响应格式符合Swagger文档
func (h *SwaggerResponseHelper) HandleWithSwaggerFormat(c *gin.Context, data any, err error) {
	if err != nil {
		h.adapter.AdaptErrorResponse(c, err)
		return
	}
	h.adapter.AdaptSuccessResponse(c, data)
}

// HandlePagingWithSwaggerFormat 使用Swagger格式处理分页响应
func (h *SwaggerResponseHelper) HandlePagingWithSwaggerFormat(c *gin.Context, data any, page, pageSize int, total int64, err error) {
	if err != nil {
		h.adapter.AdaptErrorResponse(c, err)
		return
	}
	h.adapter.AdaptPagingResponse(c, data, page, pageSize, total)
}

// GetAdapter 获取响应适配器
func (h *SwaggerResponseHelper) GetAdapter() *ResponseAdapter {
	return h.adapter
}

// 全局Swagger响应辅助工具实例
var defaultSwaggerHelper = NewSwaggerResponseHelper()

// HandleWithSwaggerFormat 全局便捷函数
func HandleWithSwaggerFormat(c *gin.Context, data any, err error) {
	defaultSwaggerHelper.HandleWithSwaggerFormat(c, data, err)
}

// HandlePagingWithSwaggerFormat 全局分页便捷函数
func HandlePagingWithSwaggerFormat(c *gin.Context, data any, page, pageSize int, total int64, err error) {
	defaultSwaggerHelper.HandlePagingWithSwaggerFormat(c, data, page, pageSize, total, err)
}

// GetDefaultSwaggerHelper 获取默认Swagger响应辅助工具
func GetDefaultSwaggerHelper() *SwaggerResponseHelper {
	return defaultSwaggerHelper
}

// SetDefaultSwaggerHelper 设置默认Swagger响应辅助工具
func SetDefaultSwaggerHelper(helper *SwaggerResponseHelper) {
	if helper != nil {
		defaultSwaggerHelper = helper
	}
}
