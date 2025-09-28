package entity

import (
	"time"
)

// User 用户聚合根
type User struct {
	id          string
	openID      string
	name        string
	gender      int
	phoneNumber string
	password    string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewUser 创建新用户
func NewUser(openID, name, phoneNumber, password string, gender int) *User {
	return &User{
		openID:      openID,
		name:        name,
		gender:      gender,
		phoneNumber: phoneNumber,
		password:    password,
	}
}

func (u *User) ID() string {
	return u.id
}

func (u *User) OpenID() string {
	return u.openID
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Gender() int {
	return u.gender
}

func (u *User) PhoneNumber() string {
	return u.phoneNumber
}

func (u *User) Password() string {
	return u.password
}

func (u *User) GetCreatedAt() int64 {
	return u.createdAt.UnixMilli()
}

func (u *User) GetUpdatedAt() int64 {
	return u.updatedAt.UnixMilli()
}

func (u *User) SetID(id string) {
	u.id = id
}

func (u *User) SetCreatedAt(createdAt time.Time) {
	u.createdAt = createdAt
}

func (u *User) SetUpdatedAt(updatedAt time.Time) {
	u.updatedAt = updatedAt
}
