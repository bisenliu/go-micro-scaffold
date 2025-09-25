package commandhandler

import (
	"context"

	command "services/internal/application/command/user"
	"services/internal/domain/user/entity"
	"services/internal/domain/user/repository"
	"services/internal/domain/user/service"
)

// UserCommandHandler 用户命令处理器
type UserCommandHandler struct {
	userRepo          repository.UserRepository
	userDomainService *service.UserDomainService
}

// NewUserCommandHandler 创建用户命令处理器
func NewUserCommandHandler(
	userRepo repository.UserRepository,
	userDomainService *service.UserDomainService,
) *UserCommandHandler {
	return &UserCommandHandler{
		userRepo:          userRepo,
		userDomainService: userDomainService,
	}
}

// HandleCreateUser 处理创建用户命令
func (h *UserCommandHandler) HandleCreateUser(ctx context.Context, cmd *command.CreateUserCommand) (*entity.User, error) {

	return nil, nil
}
