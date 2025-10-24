package request

// LoginRequest 登录请求DTO
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" label:"手机号" example:"13800138000"` // 用户手机号码
	Password    string `json:"password" binding:"required" label:"密码" example:"password123"`      // 用户密码
}

// WeChatLoginRequest 微信登录请求DTO
type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required" label:"微信授权码" example:"wx_auth_code_123456"` // 微信授权后获得的临时授权码
}
