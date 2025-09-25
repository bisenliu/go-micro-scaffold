package entity

// User 用户聚合根
type User struct {
	openID string
}

// NewUser 创建新用户
func NewUser(openID, unionID string) *User {
	return &User{
		openID: openID,
	}
}
