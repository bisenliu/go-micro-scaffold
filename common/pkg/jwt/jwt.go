package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"common/interfaces"
)

// JWT JWT服务结构体
type JWT struct {
	configProvider interfaces.ConfigProvider
}

// NewJWT 创建JWT实例
func NewJWT(configProvider interfaces.ConfigProvider) interfaces.JWTService {
	return &JWT{
		configProvider: configProvider,
	}
}

// GenerateToken 生成JWT令牌
func (j *JWT) GenerateToken(userID string) (string, error) {
	tokenConfig := j.configProvider.GetTokenConfig()
	
	// 创建一个我们自己的声明
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenConfig.ExpiredTime) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-micro-scaffold", // 硬编码服务名
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用指定的secret签名并获得完整的编码后的字符串token
	secretKey := "default-secret-key" // 临时硬编码，后续可以从配置中获取
	return token.SignedString([]byte(secretKey))
}

// ValidateToken 验证JWT令牌
func (j *JWT) ValidateToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// RefreshToken 刷新JWT令牌
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	
	// 生成新的令牌
	return j.GenerateToken(claims.UserID)
}

// ParseToken 解析JWT令牌获取Claims
func (j *JWT) ParseToken(tokenString string) (*interfaces.TokenClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		secretKey := "default-secret-key" // 临时硬编码，后续可以从配置中获取
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return &interfaces.TokenClaims{
			UserID: claims.UserID,
			Exp:    claims.ExpiresAt.Unix(),
			Iat:    claims.IssuedAt.Unix(),
		}, nil
	}

	return nil, errors.New("invalid token")
}
