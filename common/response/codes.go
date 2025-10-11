package response

// BusinessCode 定义业务码和消息的结构体
type BusinessCode struct {
	Code    int
	Message string
	Type    ErrorType // 成功状态码不使用此字段
}

// IsSuccess 判断是否为成功状态码
func (bc BusinessCode) IsSuccess() bool {
	return bc.Code == 0
}

// 通用状态码 (0-999)
var (
	// 成功状态码
	CodeSuccess = BusinessCode{Code: 0, Message: "操作成功"}
)

// 客户端错误状态码 (1000-1999)
var (
	CodeInvalidParams   = BusinessCode{Code: 1001, Message: "参数错误", Type: ErrorTypeBusiness}
	CodeValidationError = BusinessCode{Code: 1002, Message: "验证错误", Type: ErrorTypeBusiness}
)

// 认证授权错误状态码 (2000-2999)
var (
	CodeUnauthorized = BusinessCode{Code: 2001, Message: "未授权访问", Type: ErrorTypeAuth}
	CodeForbidden    = BusinessCode{Code: 2002, Message: "禁止访问", Type: ErrorTypePermission}
)

// 业务逻辑错误状态码 (4000-4999)
var (
	CodeBusinessError = BusinessCode{Code: 4001, Message: "业务处理失败", Type: ErrorTypeBusiness}
	CodeNotFound      = BusinessCode{Code: 4004, Message: "资源未找到", Type: ErrorTypeBusiness}
)

// 系统错误状态码 (5000-5999)
var (
	CodeInternalError = BusinessCode{Code: 5001, Message: "内部服务器错误", Type: ErrorTypeSystem}
)

// 第三方服务错误状态码 (6000-6999)
var (
	CodeThirdPartyError = BusinessCode{Code: 6001, Message: "第三方服务错误", Type: ErrorTypeThirdParty}
)

// AllBusinessCodes 包含所有业务码的映射
var AllBusinessCodes = map[int]BusinessCode{
	CodeSuccess.Code:         CodeSuccess,
	CodeInvalidParams.Code:   CodeInvalidParams,
	CodeValidationError.Code: CodeValidationError,
	CodeUnauthorized.Code:    CodeUnauthorized,
	CodeForbidden.Code:       CodeForbidden,
	CodeBusinessError.Code:   CodeBusinessError,
	CodeNotFound.Code:        CodeNotFound,
	CodeInternalError.Code:   CodeInternalError,
	CodeThirdPartyError.Code: CodeThirdPartyError,
}

// GetCodeMessage 获取状态码对应的消息
func GetCodeMessage(code int) string {
	if status, exists := AllBusinessCodes[code]; exists {
		return status.Message
	}
	return "未知错误"
}

// GetBusinessCode 获取业务码对象
func GetBusinessCode(code int) BusinessCode {
	if status, exists := AllBusinessCodes[code]; exists {
		return status
	}
	return BusinessCode{Code: code, Message: "未知错误", Type: ErrorTypeSystem}
}
