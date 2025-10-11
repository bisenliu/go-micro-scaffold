package queryhandler

import (
	"context"

	"common/databases/redis"
	"services/internal/application/query/user"
	"services/internal/domain/user/entity"
	"services/internal/domain/user/repository"
)

// UserQueryHandler 用户查询处理器
type UserQueryHandler struct {
	userRepo    repository.UserRepository
	redisClient *redis.RedisClient
}

// NewUserQueryHandler 创建用户查询处理器
func NewUserQueryHandler(
	userRepo repository.UserRepository,
	redisClient *redis.RedisClient,
) *UserQueryHandler {
	return &UserQueryHandler{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

// HandleListUsers 处理用户列表查询
func (h *UserQueryHandler) HandleListUsers(ctx context.Context, query *user.ListUsersQuery) ([]*entity.User, int64, error) {
	// 计算偏移量
	offset := (query.Page - 1) * query.PageSize

	// 构建过滤条件
	filter := &repository.UserListFilter{
		Name:      query.Name,
		Gender:    query.Gender,
		StartTime: query.StartTime,
		EndTime:   query.EndTime,
	}

	// 调用仓储层查询
	users, total, err := h.userRepo.ListWithFilter(ctx, filter, offset, query.PageSize)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// HandleGetUser 处理获取用户查询
func (h *UserQueryHandler) HandleGetUser(ctx context.Context, query *user.GetUserQuery) (*entity.User, error) {

	return h.userRepo.GetByID(ctx, query.ID)
}

