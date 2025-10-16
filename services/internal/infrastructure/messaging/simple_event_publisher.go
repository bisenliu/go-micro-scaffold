package messaging

import (
	"context"

	"common/interfaces"
	"services/internal/domain/user/entity"
)

// SimpleEventPublisher 简单的事件发布器实现
// 用于开发和测试环境，不依赖外部服务
type SimpleEventPublisher struct {
	logger interfaces.Logger
}

// NewSimpleEventPublisher 创建简单事件发布器
func NewSimpleEventPublisher(logger interfaces.Logger) EventPublisher {
	return &SimpleEventPublisher{
		logger: logger,
	}
}

// PublishUserCreated 发布用户创建事件
func (p *SimpleEventPublisher) PublishUserCreated(ctx context.Context, user *entity.User) error {
	p.logger.Info(ctx, "User created event published")
	return nil
}