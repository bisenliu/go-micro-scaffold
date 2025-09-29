package user

import "time"

// ListUsersQuery 用户列表查询
type ListUsersQuery struct {
	Page      int        `json:"page" validate:"min=1"`
	PageSize  int        `json:"page_size" validate:"min=1,max=100"`
	Name      *string    `json:"name,omitempty"`       // 姓名模糊查询
	Gender    *int       `json:"gender,omitempty"`     // 性别过滤
	StartTime *time.Time `json:"start_time,omitempty"` // 创建时间开始
	EndTime   *time.Time `json:"end_time,omitempty"`   // 创建时间结束
}
