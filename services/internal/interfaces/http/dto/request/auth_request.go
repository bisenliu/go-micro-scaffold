package request

// LoginRequest 登录请求DTO
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" label:"手机号"`
	Password    string `json:"password" binding:"required" label:"密码"`
}

// WeChatLoginRequest 微信登录请求DTO
type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required" label:"微信授权码"`
}
