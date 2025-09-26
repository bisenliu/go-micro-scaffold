package handler

import (
	"common/logger"
	"common/validation"
	command "services/internal/application/command/user"
	"services/internal/application/commandhandler"
	"services/internal/application/queryhandler"
	requestdto "services/internal/interfaces/http/dto/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/response"
)

// UserHandler 用户HTTP处理器
type UserHandler struct {
	commandHandler *commandhandler.UserCommandHandler
	queryHandler   *queryhandler.UserQueryHandler
	logger         *zap.Logger
	validator      *validation.Validator
}

// Ensure UserHandler implements Handler interface
var _ Handler = (*UserHandler)(nil)

// NewUserHandler 创建用户HTTP处理器
func NewUserHandler(
	commandHandler *commandhandler.UserCommandHandler,
	queryHandler *queryhandler.UserQueryHandler,
	zapLogger *zap.Logger,
	validator *validation.Validator,
) *UserHandler {
	return &UserHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		logger:         zapLogger,
		validator:      validator,
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
	response.Success(c, user)
}
