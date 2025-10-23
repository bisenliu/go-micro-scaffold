package swagger

// SwaggerError 统一的Swagger错误响应结构
// 用于Swagger文档中定义标准化的错误响应格式
type SwaggerError struct {
	Error   string      `json:"error" example:"Bad Request" description:"错误类型描述"`
	Message string      `json:"message" example:"Invalid input parameters" description:"详细错误信息"`
	Code    int         `json:"code" example:"400" description:"业务错误代码"`
	Details interface{} `json:"details,omitempty" description:"额外的错误详情信息"`
}

// SwaggerErrorTypes 错误类型常量
var SwaggerErrorTypes = struct {
	BadRequest          string
	Unauthorized        string
	Forbidden           string
	NotFound            string
	Conflict            string
	ValidationFailed    string
	InternalServerError string
	ServiceUnavailable  string
}{
	BadRequest:          "Bad Request",
	Unauthorized:        "Unauthorized",
	Forbidden:           "Forbidden",
	NotFound:            "Not Found",
	Conflict:            "Conflict",
	ValidationFailed:    "Validation Failed",
	InternalServerError: "Internal Server Error",
	ServiceUnavailable:  "Service Unavailable",
}
