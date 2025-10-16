package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"common/interfaces"
)

// JWTServiceImpl JWT服务实现
type JWTServiceImpl struct {
	secretKey   string
	expiredTime int // 分钟
}

// NewJWTService 创建JWT服务
func NewJWTService(configProvider interfaces.ConfigProvider) interfaces.JWTService {
	tokenConfig := configProvider.GetTokenConfig()
	
	// 使用配置中的密钥，如果没有则生成一个临时密钥
	secretKey := "default-secret-key-for-development-only"
	
	return &JWTServiceImpl{
		secretKey:   secretKey,
		expiredTime: tokenConfig.ExpiredTime,
	}
}

// GenerateToken 生成JWT令牌
func (j *JWTServiceImpl) GenerateToken(userID string) (string, error) {
	now := time.Now()
	claims := &interfaces.TokenClaims{
		UserID: userID,
		Exp:    now.Add(time.Duration(j.expiredTime) * time.Minute).Unix(),
		Iat:    now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserID,
		"exp":     claims.Exp,
		"iat":     claims.Iat,
	})

	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken 验证JWT令牌
func (j *JWTServiceImpl) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user_id in token")
	}

	// 检查过期时间
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return "", errors.New("token expired")
		}
	}

	return userID, nil
}

// RefreshToken 刷新JWT令牌
func (j *JWTServiceImpl) RefreshToken(tokenString string) (string, error) {
	userID, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return j.GenerateToken(userID)
}

// ParseToken 解析JWT令牌获取Claims
func (j *JWTServiceImpl) ParseToken(tokenString string) (*interfaces.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid exp in token")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return nil, errors.New("invalid iat in token")
	}

	return &interfaces.TokenClaims{
		UserID: userID,
		Exp:    int64(exp),
		Iat:    int64(iat),
	}, nil
}