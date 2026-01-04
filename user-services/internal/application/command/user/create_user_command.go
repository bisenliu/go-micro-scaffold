package command

// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
	OpenID      string
	Name        string
	PhoneNumber string
	Password    string
	Gender      int
}
