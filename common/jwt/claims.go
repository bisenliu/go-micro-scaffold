package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims 自定义声明类型并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 可根据需要自行添加字段
type CustomClaims struct {
	// 可根据需要自行添加字段
	UserID               string `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}
