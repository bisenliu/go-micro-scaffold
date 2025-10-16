package response

import (
	"net/http"
)

// ErrorMapping 错误映射结构
type ErrorMapping struct {
	BusinessCode   int
	HTTPStatus     int
	DefaultMessage string
}

// 预定义的通用错误映射
// 各个服务可以根据需要扩展这些映射
var CommonErrorMappings = map[string]*ErrorMapping{
	// 通用错误类型
	"not_found": {
		BusinessCode:   CodeNotFound,
		HTTPStatus:     http.StatusNotFound,
		DefaultMessage: "资源不存在",
	},
	"validation_failed": {
		BusinessCode:   CodeValidation,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "验证失败",
	},
	"already_exists": {
		BusinessCode:   CodeAlreadyExists,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "资源已存在",
	},
	"unauthorized": {
		BusinessCode:   CodeUnauthorized,
		HTTPStatus:     http.StatusUnauthorized,
		DefaultMessage: "未授权访问",
	},
	"forbidden": {
		BusinessCode:   CodeForbidden,
		HTTPStatus:     http.StatusForbidden,
		DefaultMessage: "禁止访问",
	},
	"business_rule_violation": {
		BusinessCode:   CodeBusinessError,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "业务规则违反",
	},
	"concurrency_conflict": {
		BusinessCode:   CodeConflict,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "并发冲突",
	},
	"resource_locked": {
		BusinessCode:   CodeConflict,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "资源已锁定",
	},
	"invalid_data": {
		BusinessCode:   CodeValidation,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "无效的数据",
	},
	
	// 应用层错误
	"command_validation": {
		BusinessCode:   CodeValidation,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "命令验证失败",
	},
	"command_execution": {
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "命令执行失败",
	},
	"query_execution": {
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "查询执行失败",
	},
	
	// 基础设施层错误
	"internal_server": {
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "内部服务器错误",
	},
	"invalid_request": {
		BusinessCode:   CodeBadRequest,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "无效的请求",
	},
	"database_connection": {
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "数据库连接失败",
	},
	"record_not_found": {
		BusinessCode:   CodeNotFound,
		HTTPStatus:     http.StatusNotFound,
		DefaultMessage: "记录不存在",
	},
	"duplicate_key": {
		BusinessCode:   CodeAlreadyExists,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "重复键值",
	},
	"external_service_unavailable": {
		BusinessCode:   CodeThirdParty,
		HTTPStatus:     http.StatusBadGateway,
		DefaultMessage: "外部服务不可用",
	},
	"timeout": {
		BusinessCode:   CodeTimeout,
		HTTPStatus:     http.StatusRequestTimeout,
		DefaultMessage: "请求超时",
	},
	"network_error": {
		BusinessCode:   CodeThirdParty,
		HTTPStatus:     http.StatusBadGateway,
		DefaultMessage: "网络错误",
	},
}

// GetErrorMapping 根据错误类型字符串获取错误映射
func GetErrorMapping(errorType string) (*ErrorMapping, bool) {
	mapping, exists := CommonErrorMappings[errorType]
	return mapping, exists
}

// RegisterErrorMapping 注册自定义错误映射
func RegisterErrorMapping(errorType string, mapping *ErrorMapping) {
	CommonErrorMappings[errorType] = mapping
}