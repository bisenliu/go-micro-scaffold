package handler

import (
	"common/logger"
	"common/validation"
	"net/http"
	"services/internal/application/commandhandler"
	"services/internal/application/queryhandler"
	requestdto "services/internal/interfaces/http/dto/request"
	responsedto "services/internal/interfaces/http/dto/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	// 这里应该调用应用层的命令处理器来创建用户
	// 例如:
	// command := command.NewCreateUserCommand(req.OpenID, req.Nickname, ...)
	// if err := h.commandHandler.HandleCreateUser(ctx, command); err != nil {
	//     logger.Error(ctx, "创建用户失败", zap.Error(err))
	//     c.JSON(http.StatusInternalServerError, responsedto.ErrorResponse(500, "创建用户失败"))
	//     return
	// }

	logger.Info(ctx, "用户创建成功", zap.String("open_id", req.OpenID))
	c.JSON(http.StatusOK, responsedto.SuccessResponse(nil))
}
