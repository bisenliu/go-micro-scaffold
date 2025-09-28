package jwt

import (
	"errors"
	"time"

	"common/config"

	"github.com/golang-jwt/jwt/v4"
)

// JWT JWT服务结构体
type JWT struct {
	config *config.Config
}

// NewJWT 创建JWT实例
func NewJWT(cfg *config.Config) *JWT {
	return &JWT{
		config: cfg,
	}
}

// GenToken 生成JWT
// @param userID 用户ID
// @param username 用户名
// @return string token
// @return error 生成失败异常
func (j *JWT) GenToken(userID uint64, username string) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.config.Token.ExpiredTime) * time.Minute)),
			Issuer:    j.config.System.ServerName,
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(j.config.System.SecretKey))
}

// ParseToken 解析JWT
// @param tokenString token
// @return *CustomClaims 自定义声明
// @return error 解析失败异常
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(j.config.System.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
