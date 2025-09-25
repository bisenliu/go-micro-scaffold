package user

// ListUsersQuery 用户列表查询
type ListUsersQuery struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}
