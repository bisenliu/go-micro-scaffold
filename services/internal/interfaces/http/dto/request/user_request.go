package request

import uservo "services/internal/domain/user/valueobject"

// CreateUserRequest 创建用户请求DTO
type CreateUserRequest struct {
	OpenID      string        `json:"open_id" binding:"required" label:"开放ID"`
	Name        string        `json:"name" binding:"required,max=50" label:"昵称"`
	Gender      uservo.Gender `json:"gender" binding:"required,enum" label:"性别"`
	PhoneNumber string        `json:"phone_number" binding:"required" label:"手机号"`
	Password    string        `json:"password" binding:"required" label:"密码"`
}
