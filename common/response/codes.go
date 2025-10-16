package response

import "net/http"

// 业务状态码定义
const (
	// 成功码 (0)
	CodeSuccess = 0

	// 客户端错误 (1000-1999)
	CodeBadRequest   = 1000
	CodeValidation   = 1001
	CodeUnauthorized = 1002
	CodeForbidden    = 1003
	CodeNotFound     = 1004

	// 业务错误 (2000-2999)
	CodeBusinessError = 2000
	CodeAlreadyExists = 2001
	CodeConflict      = 2002

	// 系统错误 (5000-5999)
	CodeInternalError = 5000
	CodeTimeout       = 5001
	CodeRateLimit     = 5002
	CodeThirdParty    = 5003
)

// CodeInfo 业务码信息
type CodeInfo struct {
	Code       int
	Message    string
	HTTPStatus int
	Category   string
}

// 业务码信息映射表
var CodeInfoMap = map[int]*CodeInfo{
	CodeSuccess: {
		Code:       CodeSuccess,
		Message:    "操作成功",
		HTTPStatus: http.StatusOK,
		Category:   "success",
	},

	// 客户端错误
	CodeBadRequest: {
		Code:       CodeBadRequest,
		Message:    "请求参数错误",
		HTTPStatus: http.StatusBadRequest,
		Category:   "client_error",
	},
	CodeValidation: {
		Code:       CodeValidation,
		Message:    "数据验证失败",
		HTTPStatus: http.StatusBadRequest,
		Category:   "client_error",
	},
	CodeUnauthorized: {
		Code:       CodeUnauthorized,
		Message:    "未授权访问",
		HTTPStatus: http.StatusUnauthorized,
		Category:   "client_error",
	},
	CodeForbidden: {
		Code:       CodeForbidden,
		Message:    "禁止访问",
		HTTPStatus: http.StatusForbidden,
		Category:   "client_error",
	},
	CodeNotFound: {
		Code:       CodeNotFound,
		Message:    "资源不存在",
		HTTPStatus: http.StatusNotFound,
		Category:   "client_error",
	},

	// 业务错误
	CodeBusinessError: {
		Code:       CodeBusinessError,
		Message:    "业务处理失败",
		HTTPStatus: http.StatusBadRequest,
		Category:   "business_error",
	},
	CodeAlreadyExists: {
		Code:       CodeAlreadyExists,
		Message:    "资源已存在",
		HTTPStatus: http.StatusConflict,
		Category:   "business_error",
	},
	CodeConflict: {
		Code:       CodeConflict,
		Message:    "资源冲突",
		HTTPStatus: http.StatusConflict,
		Category:   "business_error",
	},

	// 系统错误
	CodeInternalError: {
		Code:       CodeInternalError,
		Message:    "内部服务器错误",
		HTTPStatus: http.StatusInternalServerError,
		Category:   "system_error",
	},
	CodeTimeout: {
		Code:       CodeTimeout,
		Message:    "请求超时",
		HTTPStatus: http.StatusRequestTimeout,
		Category:   "system_error",
	},
	CodeRateLimit: {
		Code:       CodeRateLimit,
		Message:    "请求过于频繁",
		HTTPStatus: http.StatusTooManyRequests,
		Category:   "system_error",
	},
	CodeThirdParty: {
		Code:       CodeThirdParty,
		Message:    "第三方服务错误",
		HTTPStatus: http.StatusBadGateway,
		Category:   "system_error",
	},
}

// GetCodeInfo 获取业务码信息
func GetCodeInfo(code int) (*CodeInfo, bool) {
	info, exists := CodeInfoMap[code]
	return info, exists
}

// GetCodeMessage 获取业务码对应的消息
func GetCodeMessage(code int) string {
	if info, exists := CodeInfoMap[code]; exists {
		return info.Message
	}
	return "未知错误"
}

// GetHTTPStatus 获取业务码对应的HTTP状态码
func GetHTTPStatus(code int) int {
	if info, exists := CodeInfoMap[code]; exists {
		return info.HTTPStatus
	}
	return http.StatusInternalServerError
}