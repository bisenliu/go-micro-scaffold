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
// @Summary 用户密码登录
// @Description 使用手机号和密码进行用户登录，成功后返回JWT Token
// @Tags 认证授权
// @Accept json
// @Produce json
// @Param request body requestdto.LoginRequest true "登录请求"
// @Success 200 {object} map[string]string "登录成功，返回token"
// @Failure 400 {object} services_internal_interfaces_http_swagger.ValidationErrorResponse "请求参数验证失败"
// @Failure 401 {object} services_internal_interfaces_http_swagger.UnauthorizedErrorResponse "用户名或密码错误"
// @Failure 500 {object} services_internal_interfaces_http_swagger.InternalServerErrorResponse "服务器内部错误"
// @Router /auth/login [post]
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

	response.Handle(c, gin.H{
		"token": token,
	}, nil)
}

// LoginByWeChat 微信登录
// @Summary 微信登录
// @Description 使用微信授权码进行用户登录，成功后返回JWT Token
// @Tags 认证授权
// @Accept json
// @Produce json
// @Param request body requestdto.WeChatLoginRequest true "微信登录请求"
// @Success 200 {object} map[string]string "登录成功，返回token"
// @Failure 400 {object} services_internal_interfaces_http_swagger.ValidationErrorResponse "请求参数验证失败"
// @Failure 401 {object} services_internal_interfaces_http_swagger.UnauthorizedErrorResponse "微信授权失败"
// @Failure 500 {object} services_internal_interfaces_http_swagger.InternalServerErrorResponse "服务器内部错误"
// @Router /auth/wechat [post]
func (h *AuthHandler) LoginByWeChat(c *gin.Context) {
	// 1. Get code from request
	// 2. Call authService.LoginByWeChat
	// 3. Generate JWT and return
	response.Handle(c, gin.H{"message": "WeChat login placeholder"}, nil)
}

// Logout 登出
// @Summary 用户登出
// @Description 用户登出，使当前Token失效
// @Tags 认证授权
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "登出成功"
// @Failure 401 {object} services_internal_interfaces_http_swagger.UnauthorizedErrorResponse "未授权或Token无效"
// @Failure 500 {object} services_internal_interfaces_http_swagger.InternalServerErrorResponse "服务器内部错误"
// @Security BearerAuth
// @Router /auth/logout [post]
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
