package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/pkg/jwt"
	"common/pkg/validation"
	"common/response"
	"services/internal/application/service"
	requestdto "services/internal/interfaces/http/dto/request"
)

// AuthHandler 认证HTTP处理器
type AuthHandler struct {
	authService service.AuthServiceInterface
	validator   *validation.Validator
	jwtService  *jwt.JWT
}

// NewAuthHandler 创建认证HTTP处理器
func NewAuthHandler(
	authService service.AuthServiceInterface,
	validator *validation.Validator,
	jwtService *jwt.JWT,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
		jwtService:  jwtService,
	}
}

// LoginByPassword 用户登录
func (h *AuthHandler) LoginByPassword(c *gin.Context) {
	ctx := c.Request.Context()
	var req requestdto.LoginRequest
	if !h.validator.Verify(c, &req, validation.JSONBindAdapter) {
		return
	}

	userID, userName, err := h.authService.LoginByPassword(ctx, req.PhoneNumber, req.Password)
	if err != nil {
		logger.Error(ctx, "Login failed", zap.Error(err))
		HandleWithLogging(c, nil, err)
		return
	}

	// 生成token
	token, err := h.jwtService.Generate(userID, userName)
	if err != nil {
		HandleWithLogging(c, nil, response.NewInternalServerError("Failed to generate token", err))
		return
	}

	response.OK(c, gin.H{
		"token": token,
	})
}

// LoginByWeChat 微信登录
func (h *AuthHandler) LoginByWeChat(c *gin.Context) {
	// 1. Get code from request
	// 2. Call authService.LoginByWeChat
	// 3. Generate JWT and return
	response.OK(c, gin.H{"message": "WeChat login placeholder"})
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	claimsValue, exists := c.Get("claims")
	if !exists {
		HandleWithLogging(c, nil, response.NewUnauthorizedError("无法获取用户信息"))
		return
	}

	claims, ok := claimsValue.(*jwt.CustomClaims)
	if !ok {
		HandleWithLogging(c, nil, response.NewInternalServerError("用户信息类型断言失败"))
		return
	}

	err := h.authService.Logout(ctx, claims.ID, claims.ExpiresAt.Time)
	if err != nil {
		logger.Error(ctx, "Logout failed", zap.Error(err))
	}

	// 使用带日志功能的统一API，自动判断成功或错误并记录DomainError上下文
	HandleWithLogging(c, "登出成功", err)
}
