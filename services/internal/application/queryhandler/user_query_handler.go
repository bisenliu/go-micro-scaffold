package queryhandler

import (
	"common/cache"
	"services/internal/domain/user/repository"
	"time"
)

const DefaultExpiration = 30 * time.Minute

// UserQueryHandler 用户查询处理器
type UserQueryHandler struct {
	userRepo    repository.UserRepository
	redisClient *cache.RedisClient
}

// NewUserQueryHandler 创建用户查询处理器
func NewUserQueryHandler(
	userRepo repository.UserRepository,
	redisClient *cache.RedisClient,
) *UserQueryHandler {
	return &UserQueryHandler{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

