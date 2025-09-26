package queryhandler

import (
	"common/databases/redis"
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
