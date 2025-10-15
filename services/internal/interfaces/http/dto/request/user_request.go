package request

import (
	"time"

	"common/pkg/pagination"
	uservo "services/internal/domain/user/valueobject"
)

// CreateUserRequest 创建用户请求DTO
type CreateUserRequest struct {
	OpenID      string        `json:"open_id" binding:"required" label:"开放ID"`
	Name        string        `json:"name" binding:"required,max=50" label:"昵称"`
	Gender      uservo.Gender `json:"gender" binding:"required,enum" label:"性别"`
	PhoneNumber string        `json:"phone_number" binding:"required" label:"手机号"`
	Password    string        `json:"password" binding:"required" label:"密码"`
}

// ListUsersRequest 用户列表请求DTO
type ListUsersRequest struct {
	pagination.PageParams
	Name      string     `form:"name" binding:"omitempty,max=50" label:"姓名"`
	Gender    *int       `form:"gender" binding:"omitempty,oneof=100 200 300" label:"性别"`
	StartTime *time.Time `form:"start_time" binding:"omitempty" time_format:"2006-01-02" label:"开始时间"`
	EndTime   *time.Time `form:"end_time" binding:"omitempty" time_format:"2006-01-02" label:"结束时间"`
}
