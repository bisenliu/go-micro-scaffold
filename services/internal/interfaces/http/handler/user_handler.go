package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/pkg/validation"
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
// @Summary 创建新用户
// @Description 创建一个新的用户账户，需要提供用户的基本信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body requestdto.CreateUserRequest true "创建用户请求"
// @Success 200 {object} responsedto.UserInfoResponse "创建成功"
// @Failure 400 {object} services_internal_interfaces_http_swagger.SwaggerError "请求参数验证失败"
// @Failure 409 {object} services_internal_interfaces_http_swagger.SwaggerError "用户已存在"
// @Failure 500 {object} services_internal_interfaces_http_swagger.SwaggerError "服务器内部错误"
// @Security BearerAuth
// @Router /users [post]
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
	} else {
		logger.Info(ctx, "User created successfully", zap.String("open_id", req.OpenID))
	}

	HandleWithLogging(c, responsedto.ToUserInfoResponse(user), err)
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表，支持按姓名、性别、时间范围等条件过滤
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1) minimum(1)
// @Param page_size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param name query string false "用户姓名（模糊搜索）" maxlength(50)
// @Param gender query int false "性别" Enums(100,200,300) example(100)
// @Param start_time query string false "开始时间" format(date) example("2023-01-01")
// @Param end_time query string false "结束时间" format(date) example("2023-12-31")
// @Success 200 {object} responsedto.UserListResponse "获取成功"
// @Failure 400 {object} services_internal_interfaces_http_swagger.SwaggerError "请求参数验证失败"
// @Failure 401 {object} services_internal_interfaces_http_swagger.SwaggerError "未授权访问"
// @Failure 500 {object} services_internal_interfaces_http_swagger.SwaggerError "服务器内部错误"
// @Security BearerAuth
// @Router /users [get]
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
	} else {
		logger.Info(ctx, "User list retrieved successfully",
			zap.Int64("total", total),
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize))
	}

	// 将领域实体转换为DTO
	userResponses := responsedto.ToUserListResponse(users)

	HandlePagingWithLogging(c, userResponses, req.Page, req.PageSize, total, err)
}

// GetUser 获取用户信息
// @Summary 获取用户详细信息
// @Description 根据用户ID获取用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID" example("user_123456789")
// @Success 200 {object} responsedto.UserInfoResponse "获取成功"
// @Failure 400 {object} services_internal_interfaces_http_swagger.SwaggerError "请求参数错误"
// @Failure 401 {object} services_internal_interfaces_http_swagger.SwaggerError "未授权访问"
// @Failure 404 {object} services_internal_interfaces_http_swagger.SwaggerError "用户不存在"
// @Failure 500 {object} services_internal_interfaces_http_swagger.SwaggerError "服务器内部错误"
// @Security BearerAuth
// @Router /users/{id} [get]
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
	}

	HandleWithLogging(c, responsedto.ToUserInfoResponse(userInfo), err)
}
