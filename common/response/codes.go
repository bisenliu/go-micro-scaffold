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
	CodeMissingParams   = 1003 // 缺少必要参数
	CodeInvalidFormat   = 1004 // 格式错误
	CodeInvalidToken    = 1005 // 无效的令牌
	CodeTokenExpired    = 1006 // 令牌已过期
)

// 认证授权错误状态码 (2000-2999)
const (
	CodeUnauthorized     = 2001 // 未授权
	CodeForbidden        = 2002 // 禁止访问
	CodeLoginRequired    = 2003 // 需要登录
	CodePermissionDenied = 2004 // 权限不足
	CodeAccountDisabled  = 2005 // 账户已禁用
	CodeAccountLocked    = 2006 // 账户已锁定
)

// 资源错误状态码 (3000-3999)
const (
	CodeNotFound        = 3001 // 资源不存在
	CodeResourceExists  = 3002 // 资源已存在
	CodeResourceDeleted = 3003 // 资源已删除
	CodeResourceLocked  = 3004 // 资源已锁定
)

// 业务逻辑错误状态码 (4000-4999)
const (
	CodeBusinessError       = 4001 // 通用业务错误
	CodeOperationFailed     = 4002 // 操作失败
	CodeDataInconsistent    = 4003 // 数据不一致
	CodeQuotaExceeded       = 4004 // 配额超限
	CodeFrequencyLimit      = 4005 // 频率限制
	CodeOperationNotAllowed = 4006 // 操作不被允许
)

// 系统错误状态码 (5000-5999)
const (
	CodeInternalError      = 5001 // 内部服务器错误
	CodeServiceUnavailable = 5002 // 服务不可用
	CodeDatabaseError      = 5003 // 数据库错误
	CodeNetworkError       = 5004 // 网络错误
	CodeTimeoutError       = 5005 // 超时错误
	CodeConfigError        = 5006 // 配置错误
)

// 第三方服务错误状态码 (6000-6999)
const (
	CodeThirdPartyError = 6001 // 第三方服务错误
	CodePaymentError    = 6002 // 支付错误
	CodeSMSError        = 6003 // 短信服务错误
	CodeEmailError      = 6004 // 邮件服务错误
	CodeStorageError    = 6005 // 存储服务错误
)

// 状态码消息映射
var CodeMessages = map[int]string{
	// 成功
	CodeSuccess: "操作成功",

	// 客户端错误
	CodeInvalidParams:   "参数错误",
	CodeValidationError: "验证错误",
	CodeMissingParams:   "缺少必要参数",
	CodeInvalidFormat:   "格式错误",
	CodeInvalidToken:    "无效的令牌",
	CodeTokenExpired:    "令牌已过期",

	// 认证授权错误
	CodeUnauthorized:     "未授权访问",
	CodeForbidden:        "禁止访问",
	CodeLoginRequired:    "需要登录",
	CodePermissionDenied: "权限不足",
	CodeAccountDisabled:  "账户已禁用",
	CodeAccountLocked:    "账户已锁定",

	// 资源错误
	CodeNotFound:        "资源不存在",
	CodeResourceExists:  "资源已存在",
	CodeResourceDeleted: "资源已删除",
	CodeResourceLocked:  "资源已锁定",

	// 业务逻辑错误
	CodeBusinessError:       "业务处理失败",
	CodeOperationFailed:     "操作失败",
	CodeDataInconsistent:    "数据不一致",
	CodeQuotaExceeded:       "配额超限",
	CodeFrequencyLimit:      "请求过于频繁",
	CodeOperationNotAllowed: "操作不被允许",

	// 系统错误
	CodeInternalError:      "内部服务器错误",
	CodeServiceUnavailable: "服务暂时不可用",
	CodeDatabaseError:      "数据库错误",
	CodeNetworkError:       "网络错误",
	CodeTimeoutError:       "请求超时",
	CodeConfigError:        "配置错误",

	// 第三方服务错误
	CodeThirdPartyError: "第三方服务错误",
	CodePaymentError:    "支付处理失败",
	CodeSMSError:        "短信发送失败",
	CodeEmailError:      "邮件发送失败",
	CodeStorageError:    "存储服务错误",
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

// IsClientError 判断是否为客户端错误
func IsClientError(code int) bool {
	return code >= 1000 && code < 2000
}

// IsAuthError 判断是否为认证授权错误
func IsAuthError(code int) bool {
	return code >= 2000 && code < 3000
}

// IsResourceError 判断是否为资源错误
func IsResourceError(code int) bool {
	return code >= 3000 && code < 4000
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
