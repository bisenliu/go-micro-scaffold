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
		HandleErrorResponse(c, err)
		return
	}

	// 生成token
	token, err := h.jwtService.Generate(userID, userName)
	if err != nil {
		response.FailWithCode(c, response.CodeInternalError, "Failed to generate token")
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
		response.FailWithCode(c, response.CodeUnauthorized, "无法获取用户信息")
		return
	}

	claims, ok := claimsValue.(*jwt.CustomClaims)
	if !ok {
		response.FailWithCode(c, response.CodeInternalError, "用户信息类型断言失败")
		return
	}

	if err := h.authService.Logout(ctx, claims.ID, claims.ExpiresAt.Time); err != nil {
		logger.Error(ctx, "Logout failed", zap.Error(err))
		HandleErrorResponse(c, err)
		return
	}

	response.OK(c, "登出成功")
}
