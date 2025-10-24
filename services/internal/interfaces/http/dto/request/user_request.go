package request

import (
	"time"

	"common/pkg/pagination"
	uservo "services/internal/domain/user/valueobject"
)

// CreateUserRequest 创建用户请求DTO
type CreateUserRequest struct {
	OpenID      string        `json:"open_id" binding:"required" label:"开放ID" example:"wx_123456789"`    // 微信OpenID或其他第三方平台的唯一标识
	Name        string        `json:"name" binding:"required,max=50" label:"昵称" example:"张三"`            // 用户姓名，长度不超过50个字符
	Gender      uservo.Gender `json:"gender" binding:"required,enum" label:"性别" example:"100"`           // 性别：100-男性，200-女性，300-其他
	PhoneNumber string        `json:"phone_number" binding:"required" label:"手机号" example:"13800138000"` // 手机号码，需要符合中国大陆手机号格式
	Password    string        `json:"password" binding:"required" label:"密码" example:"password123"`      // 用户密码，长度至少6位
}

// ListUsersRequest 用户列表请求DTO
type ListUsersRequest struct {
	pagination.PageParams
	Name      string         `form:"name" binding:"omitempty,max=50" label:"姓名" example:"张三"`                                   // 用户姓名，支持模糊搜索
	Gender    *uservo.Gender `form:"gender" binding:"omitempty,enum" label:"性别" example:"100"`                                  // 性别过滤：100-男性，200-女性，300-其他
	StartTime *time.Time     `form:"start_time" binding:"omitempty" time_format:"2006-01-02" label:"开始时间" example:"2023-01-01"` // 创建时间范围的开始时间，格式：YYYY-MM-DD
	EndTime   *time.Time     `form:"end_time" binding:"omitempty" time_format:"2006-01-02" label:"结束时间" example:"2023-12-31"`   // 创建时间范围的结束时间，格式：YYYY-MM-DD
}
