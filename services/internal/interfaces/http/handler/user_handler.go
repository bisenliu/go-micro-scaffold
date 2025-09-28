package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/pkg/jwt"
	"common/pkg/validation"
	"common/response"
	command "services/internal/application/command/user"
	"services/internal/application/commandhandler"
	"services/internal/application/queryhandler"
	requestdto "services/internal/interfaces/http/dto/request"
	responsedto "services/internal/interfaces/http/dto/response"
)

// UserHandler 用户HTTP处理器
type UserHandler struct {
	commandHandler *commandhandler.UserCommandHandler
	queryHandler   *queryhandler.UserQueryHandler
	logger         *zap.Logger
	validator      *validation.Validator
	jwtService     *jwt.JWT
}

// Ensure UserHandler implements Handler interface
var _ Handler = (*UserHandler)(nil)

// NewUserHandler 创建用户HTTP处理器
func NewUserHandler(
	commandHandler *commandhandler.UserCommandHandler,
	queryHandler *queryhandler.UserQueryHandler,
	zapLogger *zap.Logger,
	validator *validation.Validator,
	jwtService *jwt.JWT,
) *UserHandler {
	return &UserHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		logger:         zapLogger,
		validator:      validator,
		jwtService:     jwtService,
	}
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Creating user", zap.String("request_id", "create_user"))

	var req requestdto.CreateUserRequest
	if !h.validator.Verify(c, &req, validation.JSONBindAdapter) {
		return
	}

	command := &command.CreateUserCommand{
		OpenID:      req.OpenID,
		Name:        req.Name,
		Gender:      req.Gender.Int(),
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}
	user, err := h.commandHandler.HandleCreateUser(ctx, command)
	if err != nil {
		logger.Error(ctx, "创建用户失败", zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	logger.Info(ctx, "用户创建成功", zap.String("open_id", req.OpenID))

	response.Success(c, responsedto.ToUserInfoResponse(user))
}

// Login 用户登录示例
func (h *UserHandler) Login(c *gin.Context) {
	// 注意：这只是一个示例，实际登录需要验证用户名和密码
	// 假设已经验证了用户身份，用户ID为123，用户名为"testuser"

	// 生成token
	token, err := h.jwtService.GenToken("123", "testuser")
	if err != nil {
		response.BusinessError(c, response.CodeBusinessError, "Failed to generate token")
		return
	}

	response.Success(c, gin.H{
		"token": token,
	})
}
