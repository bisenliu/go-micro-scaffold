package response

// 业务状态码定义文件
// 这里定义了所有可能的业务状态码，便于统一管理和维护

// 通用状态码 (0-999)
const (
	// 成功状态码
	CodeSuccess = 0 // 成功
)

// 客户端错误状态码 (1000-1999)
const (
	CodeInvalidParams   = 1001 // 参数错误
	CodeValidationError = 1002 // 验证错误
)

// 认证授权错误状态码 (2000-2999)
const (
	CodeUnauthorized = 2001 // 未授权
	CodeForbidden    = 2002 // 禁止访问
)

// 业务逻辑错误状态码 (4000-4999)
const (
	CodeBusinessError = 4001 // 通用业务错误
)

// 系统错误状态码 (5000-5999)
const (
	CodeInternalError = 5001 // 内部服务器错误
)

// 第三方服务错误状态码 (6000-6999)
const (
	CodeThirdPartyError = 6001 // 第三方服务错误
)

// 状态码消息映射
var CodeMessages = map[int]string{
	// 成功
	CodeSuccess: "操作成功",

	// 客户端错误
	CodeInvalidParams:   "参数错误",
	CodeValidationError: "验证错误",

	// 认证授权错误
	CodeUnauthorized: "未授权访问",

	// 业务逻辑错误
	CodeBusinessError: "业务处理失败",

	// 系统错误
	CodeInternalError: "内部服务器错误",

	// 第三方服务错误
	CodeThirdPartyError: "第三方服务错误",
}

// GetCodeMessage 获取状态码对应的消息
func GetCodeMessage(code int) string {
	if message, exists := CodeMessages[code]; exists {
		return message
	}
	return "未知错误"
}

// IsSuccessCode 判断是否为成功状态码
func IsSuccessCode(code int) bool {
	return code == CodeSuccess
}

// IsClientError 判断是否为客户端错误 (参数验证错误)
func IsClientError(code int) bool {
	return code >= 1000 && code < 2000
}

// IsAuthError 判断是否为认证授权错误
func IsAuthError(code int) bool {
	return code >= 2000 && code < 3000
}

// IsBusinessError 判断是否为业务错误
func IsBusinessError(code int) bool {
	return code >= 4000 && code < 5000
}

// IsSystemError 判断是否为系统错误
func IsSystemError(code int) bool {
	return code >= 5000 && code < 6000
}

// IsThirdPartyError 判断是否为第三方服务错误
func IsThirdPartyError(code int) bool {
	return code >= 6000 && code < 7000
}
