package request

// CreateUserRequest 创建用户请求DTO
type CreateUserRequest struct {
	OpenID    string `json:"open_id" binding:"required" label:"开放ID"`
	UnionID   string `json:"union_id" label:"联合ID"`
	Nickname  string `json:"nickname" binding:"required,max=50" label:"昵称"`
	AvatarURL string `json:"avatar_url" label:"头像URL"`
	Gender    int    `json:"gender" binding:"required,oneof=100 200 300" label:"性别"`
	Phone     string `json:"phone" binding:"required,e164" label:"手机号"`
	Birthday  string `json:"birthday" binding:"omitempty" label:"生日"`
	Location  string `json:"location" label:"位置"`
}

// UpdateUserRequest 更新用户请求DTO
type UpdateUserRequest struct {
	Nickname  string `json:"nickname" binding:"max=50" label:"昵称"`
	AvatarURL string `json:"avatar_url" label:"头像URL"`
	Gender    int    `json:"gender" binding:"oneof=100 200 300" label:"性别"`
	Phone     string `json:"phone" binding:"omitempty,e164" label:"手机号"`
	Birthday  string `json:"birthday" binding:"omitempty" label:"生日"`
	Location  string `json:"location" label:"位置"`
}

// UserListRequest 用户列表请求DTO
type UserListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1" label:"页码"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100" label:"每页数量"`
	Keyword  string `form:"keyword" binding:"omitempty,max=100" label:"搜索关键词"`
}

// UserLoginRequest 用户登录请求DTO
type UserLoginRequest struct {
	OpenID  string `json:"open_id" binding:"required" label:"开放ID"`
	Phone   string `json:"phone" binding:"omitempty,e164" label:"手机号"`
	Code    string `json:"code" binding:"required" label:"验证码"`
	Channel string `json:"channel" binding:"required,oneof=wechat mini_program" label:"登录渠道"`
}
