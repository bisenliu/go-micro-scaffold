package command

// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
	OpenID  string `json:"open_id" validate:"required"`
	UnionID string `json:"union_id"`
}
