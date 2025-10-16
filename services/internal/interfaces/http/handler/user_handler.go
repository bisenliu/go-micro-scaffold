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
	"services/internal/application/query/user"
	"services/internal/application/queryhandler"
	requestdto "services/internal/interfaces/http/dto/request"
	responsedto "services/internal/interfaces/http/dto/response"
)

// UserHandler 用户HTTP处理器
type UserHandler struct {
	commandHandler *commandhandler.UserCommandHandler
	queryHandler   *queryhandler.UserQueryHandler
	validator      *validation.Validator
	jwtService     *jwt.JWT
}

// Ensure UserHandler implements Handler interface
var _ Handler = (*UserHandler)(nil)

// NewUserHandler 创建用户HTTP处理器
func NewUserHandler(
	commandHandler *commandhandler.UserCommandHandler,
	queryHandler *queryhandler.UserQueryHandler,
	validator *validation.Validator,
) *UserHandler {
	return &UserHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
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
		logger.Error(ctx, "Failed to create user", zap.Error(err))
		HandleErrorResponse(c, err)
		return
	}

	logger.Info(ctx, "User created successfully", zap.String("open_id", req.OpenID))

	response.OK(c, responsedto.ToUserInfoResponse(user))
}

// ListUsers 获取用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Listing users", zap.String("request_id", "list_users"))

	var req requestdto.ListUsersRequest
	if !h.validator.Verify(c, &req, validation.QueryBindAdapter) {
		return
	}

	// 构建查询对象
	query := &user.ListUsersQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 设置过滤条件
	if req.Name != "" {
		query.Name = &req.Name
	}
	if req.Gender != nil {
		query.Gender = req.Gender
	}
	if req.StartTime != nil {
		query.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		query.EndTime = req.EndTime
	}

	// 调用查询处理器
	users, total, err := h.queryHandler.HandleListUsers(ctx, query)
	if err != nil {
		logger.Error(ctx, "Failed to list users", zap.Error(err))
		HandleErrorResponse(c, err)
		return
	}

	logger.Info(ctx, "User list retrieved successfully",
		zap.Int64("total", total),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize))

	// 将领域实体转换为DTO
	userResponses := responsedto.ToUserListResponse(users)

	// 返回分页响应
	response.OKWithPaging(c, userResponses, req.Page, req.PageSize, total)
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	// 获取用户ID
	userID := c.Param("id")

	// 构建查询对象
	query := &user.GetUserQuery{ID: userID}

	// 获取用户信息
	userInfo, err := h.queryHandler.HandleGetUser(ctx, query)
	if err != nil {
		logger.Error(ctx, "Failed to get user info", zap.Error(err), zap.String("user_id", userID))
		HandleErrorResponse(c, err)
		return
	}

	response.OK(c, responsedto.ToUserInfoResponse(userInfo))
}
