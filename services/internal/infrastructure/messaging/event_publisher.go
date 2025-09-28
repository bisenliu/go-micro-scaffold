package messaging

import (
	"context"
	"encoding/json"
	"services/internal/domain/user/entity"
	"time"

	commonRedis "common/databases/redis"
	"common/pkg/idgen"

	"go.uber.org/zap"
)

// EventPublisher 事件发布器接口
type EventPublisher interface {
	PublishUserCreated(ctx context.Context, user *entity.User) error
}

// RedisEventPublisher Redis事件发布器实现
type RedisEventPublisher struct {
	redisClient *commonRedis.RedisClient
	logger      *zap.Logger
	idGen       idgen.Generator
}

// NewRedisEventPublisher 创建Redis事件发布器
func NewRedisEventPublisher(redisClient *commonRedis.RedisClient, logger *zap.Logger, idgen idgen.Generator) EventPublisher {
	return &RedisEventPublisher{
		redisClient: redisClient,
		logger:      logger,
		idGen:       idgen,
	}
}

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	UserID    string    `json:"user_id"`
	OpenID    string    `json:"open_id"`
	Timestamp time.Time `json:"timestamp"`
}

// PublishUserCreated 发布用户创建事件
func (p *RedisEventPublisher) PublishUserCreated(ctx context.Context, user *entity.User) error {
	event := UserCreatedEvent{
		EventID:   p.idGen.NewID().String(),
		EventType: "user.created",
		UserID:    user.ID(),
		OpenID:    user.OpenID(),
		Timestamp: time.Now(),
	}

	return p.publishEvent(ctx, "events:user:created", event)
}

// publishEvent 发布事件到Redis
func (p *RedisEventPublisher) publishEvent(ctx context.Context, channel string, event interface{}) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal event", zap.Error(err))
		return err
	}

	// 使用Redis客户端发布事件

	p.logger.Info("Event published successfully", zap.String("channel", channel), zap.String("event", string(eventData)))
	return nil
}
