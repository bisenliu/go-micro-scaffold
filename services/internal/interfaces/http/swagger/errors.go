package swagger

// ErrorResponse 统一的Swagger错误响应结构
// 用于Swagger文档中定义标准化的错误响应格式
type ErrorResponse struct {
	Error   string      `json:"error" example:"Bad Request" description:"错误类型描述"`
	Message string      `json:"message" example:"Invalid input parameters" description:"详细错误信息"`
	Code    int         `json:"code" example:"400" description:"业务错误代码"`
	Details interface{} `json:"details,omitempty" description:"额外的错误详情信息"`
}

// ValidationErrorResponse 参数验证错误响应
// 用于表示请求参数验证失败的错误响应
type ValidationErrorResponse struct {
	Error   string                 `json:"error" example:"Validation Failed" description:"错误类型"`
	Message string                 `json:"message" example:"Request validation failed" description:"错误信息"`
	Code    int                    `json:"code" example:"400" description:"业务错误代码"`
	Details ValidationErrorDetails `json:"details" description:"验证错误详情"`
}

// ValidationErrorDetails 验证错误详情
type ValidationErrorDetails struct {
	Fields []FieldError `json:"fields" description:"字段验证错误列表"`
}

// FieldError 字段验证错误
type FieldError struct {
	Field   string `json:"field" example:"name" description:"出错的字段名"`
	Message string `json:"message" example:"Name is required" description:"字段错误信息"`
	Value   string `json:"value,omitempty" example:"" description:"字段的值"`
}

// UnauthorizedErrorResponse 未授权错误响应
type UnauthorizedErrorResponse struct {
	Error   string `json:"error" example:"Unauthorized" description:"错误类型"`
	Message string `json:"message" example:"Authentication required" description:"错误信息"`
	Code    int    `json:"code" example:"401" description:"业务错误代码"`
}

// ForbiddenErrorResponse 禁止访问错误响应
type ForbiddenErrorResponse struct {
	Error   string `json:"error" example:"Forbidden" description:"错误类型"`
	Message string `json:"message" example:"Access denied" description:"错误信息"`
	Code    int    `json:"code" example:"403" description:"业务错误代码"`
}

// NotFoundErrorResponse 资源不存在错误响应
type NotFoundErrorResponse struct {
	Error   string `json:"error" example:"Not Found" description:"错误类型"`
	Message string `json:"message" example:"Resource not found" description:"错误信息"`
	Code    int    `json:"code" example:"404" description:"业务错误代码"`
}

// ConflictErrorResponse 资源冲突错误响应
type ConflictErrorResponse struct {
	Error   string `json:"error" example:"Conflict" description:"错误类型"`
	Message string `json:"message" example:"Resource already exists" description:"错误信息"`
	Code    int    `json:"code" example:"409" description:"业务错误代码"`
}

// InternalServerErrorResponse 服务器内部错误响应
type InternalServerErrorResponse struct {
	Error   string `json:"error" example:"Internal Server Error" description:"错误类型"`
	Message string `json:"message" example:"An internal server error occurred" description:"错误信息"`
	Code    int    `json:"code" example:"500" description:"业务错误代码"`
}

// ServiceUnavailableErrorResponse 服务不可用错误响应
type ServiceUnavailableErrorResponse struct {
	Error   string `json:"error" example:"Service Unavailable" description:"错误类型"`
	Message string `json:"message" example:"Service is temporarily unavailable" description:"错误信息"`
	Code    int    `json:"code" example:"503" description:"业务错误代码"`
}

// SwaggerErrorDefinitions Swagger错误定义常量
// 用于在Swagger注释中引用标准错误响应
var SwaggerErrorDefinitions = struct {
	BadRequest          string
	Unauthorized        string
	Forbidden           string
	NotFound            string
	Conflict            string
	ValidationFailed    string
	InternalServerError string
	ServiceUnavailable  string
}{
	BadRequest:          "ErrorResponse",
	Unauthorized:        "UnauthorizedErrorResponse",
	Forbidden:           "ForbiddenErrorResponse",
	NotFound:            "NotFoundErrorResponse",
	Conflict:            "ConflictErrorResponse",
	ValidationFailed:    "ValidationErrorResponse",
	InternalServerError: "InternalServerErrorResponse",
	ServiceUnavailable:  "ServiceUnavailableErrorResponse",
}

// CommonSwaggerErrors 常用的Swagger错误响应注释
// 可以在Handler方法的Swagger注释中直接使用
var CommonSwaggerErrors = struct {
	BadRequest          string
	Unauthorized        string
	Forbidden           string
	NotFound            string
	Conflict            string
	ValidationFailed    string
	InternalServerError string
	ServiceUnavailable  string
}{
	BadRequest:          "400 {object} swagger.ErrorResponse \"请求参数错误\"",
	Unauthorized:        "401 {object} swagger.UnauthorizedErrorResponse \"未授权访问\"",
	Forbidden:           "403 {object} swagger.ForbiddenErrorResponse \"禁止访问\"",
	NotFound:            "404 {object} swagger.NotFoundErrorResponse \"资源不存在\"",
	Conflict:            "409 {object} swagger.ConflictErrorResponse \"资源冲突\"",
	ValidationFailed:    "400 {object} swagger.ValidationErrorResponse \"参数验证失败\"",
	InternalServerError: "500 {object} swagger.InternalServerErrorResponse \"服务器内部错误\"",
	ServiceUnavailable:  "503 {object} swagger.ServiceUnavailableErrorResponse \"服务不可用\"",
}

// ErrorResponseConverter 错误响应转换器
// 用于将系统内部错误转换为Swagger标准错误响应格式
type ErrorResponseConverter struct{}

// NewErrorResponseConverter 创建错误响应转换器
func NewErrorResponseConverter() *ErrorResponseConverter {
	return &ErrorResponseConverter{}
}

// ConvertToSwaggerError 将通用错误转换为Swagger错误响应
func (c *ErrorResponseConverter) ConvertToSwaggerError(err error, httpStatus int) interface{} {
	if err == nil {
		return nil
	}

	// 根据HTTP状态码返回对应的错误响应结构
	switch httpStatus {
	case 400:
		return &ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
			Code:    400,
		}
	case 401:
		return &UnauthorizedErrorResponse{
			Error:   "Unauthorized",
			Message: err.Error(),
			Code:    401,
		}
	case 403:
		return &ForbiddenErrorResponse{
			Error:   "Forbidden",
			Message: err.Error(),
			Code:    403,
		}
	case 404:
		return &NotFoundErrorResponse{
			Error:   "Not Found",
			Message: err.Error(),
			Code:    404,
		}
	case 409:
		return &ConflictErrorResponse{
			Error:   "Conflict",
			Message: err.Error(),
			Code:    409,
		}
	case 500:
		return &InternalServerErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
			Code:    500,
		}
	case 503:
		return &ServiceUnavailableErrorResponse{
			Error:   "Service Unavailable",
			Message: err.Error(),
			Code:    503,
		}
	default:
		return &ErrorResponse{
			Error:   "Error",
			Message: err.Error(),
			Code:    httpStatus,
		}
	}
}

// ConvertValidationError 转换验证错误为Swagger格式
func (c *ErrorResponseConverter) ConvertValidationError(fieldErrors []FieldError) *ValidationErrorResponse {
	return &ValidationErrorResponse{
		Error:   "Validation Failed",
		Message: "Request validation failed",
		Code:    400,
		Details: ValidationErrorDetails{
			Fields: fieldErrors,
		},
	}
}

// CreateFieldError 创建字段错误
func CreateFieldError(field, message, value string) FieldError {
	return FieldError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// SwaggerErrorHelper Swagger错误辅助工具
type SwaggerErrorHelper struct {
	converter *ErrorResponseConverter
}

// NewSwaggerErrorHelper 创建Swagger错误辅助工具
func NewSwaggerErrorHelper() *SwaggerErrorHelper {
	return &SwaggerErrorHelper{
		converter: NewErrorResponseConverter(),
	}
}

// GetConverter 获取错误转换器
func (h *SwaggerErrorHelper) GetConverter() *ErrorResponseConverter {
	return h.converter
}

// GetCommonErrorAnnotations 获取常用错误注释
// 返回可以直接在Swagger注释中使用的错误响应定义
func (h *SwaggerErrorHelper) GetCommonErrorAnnotations() map[string]string {
	return map[string]string{
		"400": CommonSwaggerErrors.BadRequest,
		"401": CommonSwaggerErrors.Unauthorized,
		"403": CommonSwaggerErrors.Forbidden,
		"404": CommonSwaggerErrors.NotFound,
		"409": CommonSwaggerErrors.Conflict,
		"500": CommonSwaggerErrors.InternalServerError,
		"503": CommonSwaggerErrors.ServiceUnavailable,
	}
}

// GetValidationErrorAnnotation 获取验证错误注释
func (h *SwaggerErrorHelper) GetValidationErrorAnnotation() string {
	return CommonSwaggerErrors.ValidationFailed
}
