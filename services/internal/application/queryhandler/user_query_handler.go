package queryhandler

import (
	"services/internal/domain/user/repository"

	"common/databases/redis"
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
