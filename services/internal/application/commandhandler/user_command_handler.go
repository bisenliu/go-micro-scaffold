package commandhandler

import (
	"context"

	command "services/internal/application/command/user"
	"services/internal/domain/user/entity"
	"services/internal/domain/user/repository"
	"services/internal/domain/user/service"
	"services/internal/infrastructure/messaging"
)

// UserCommandHandler 用户命令处理器
type UserCommandHandler struct {
	userRepo          repository.UserRepository
	userDomainService *service.UserDomainService
	eventPublisher    messaging.EventPublisher
}

// NewUserCommandHandler 创建用户命令处理器
func NewUserCommandHandler(
	userRepo repository.UserRepository,
	userDomainService *service.UserDomainService,
	eventPublisher messaging.EventPublisher,
) *UserCommandHandler {
	return &UserCommandHandler{
		userRepo:          userRepo,
		userDomainService: userDomainService,
		eventPublisher:    eventPublisher,
	}
}

// HandleCreateUser 处理创建用户命令
func (h *UserCommandHandler) HandleCreateUser(ctx context.Context, cmd *command.CreateUserCommand) (*entity.User, error) {

	user, err := h.userDomainService.CreateUser(ctx, cmd.OpenID, cmd.Name, cmd.PhoneNumber, cmd.Password, cmd.Gender)
	if err != nil {
		return nil, err
	}

	// 发布用户创建事件
	if err := h.eventPublisher.PublishUserCreated(ctx, user); err != nil {
		// 记录日志但不影响主流程
		// TODO: 添加日志记录
	}

	return user, nil
}
