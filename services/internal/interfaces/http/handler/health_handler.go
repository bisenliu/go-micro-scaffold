package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/config"
	"common/databases/mysql"
	"common/databases/redis"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	mysqlClient *mysql.EntClient
	redisClient *redis.RedisClient
	config      *config.Config
	logger      *zap.Logger
}

// Ensure HealthHandler implements Handler interface
var _ Handler = (*HealthHandler)(nil)

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(
	mysqlClient *mysql.EntClient,
	redisClient *redis.RedisClient,
	config *config.Config,
	logger *zap.Logger,
) *HealthHandler {
	return &HealthHandler{
		mysqlClient: mysqlClient,
		redisClient: redisClient,
		config:      config,
		logger:      logger,
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}

// Health 健康检查
func (h *HealthHandler) Health(c *gin.Context) {
	ctx := c.Request.Context()

	response := &HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services:  make(map[string]string),
	}

	// 检查数据库连接
	if err := h.checkDatabase(ctx); err != nil {
		response.Services["database"] = "unhealthy: " + err.Error()
		response.Status = "unhealthy"
		h.logger.Error("Database health check failed", zap.Error(err))
	} else {
		response.Services["database"] = "healthy"
	}

	// 检查Redis连接
	if err := h.checkRedis(ctx); err != nil {
		response.Services["redis"] = "unhealthy: " + err.Error()
		response.Status = "unhealthy"
		h.logger.Error("Redis health check failed", zap.Error(err))
	} else {
		response.Services["redis"] = "healthy"
	}

	statusCode := http.StatusOK
	if response.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// checkDatabase 检查数据库连接
func (h *HealthHandler) checkDatabase(ctx context.Context) error {
	return h.mysqlClient.DB().PingContext(ctx)
}

// checkRedis 检查Redis连接
func (h *HealthHandler) checkRedis(ctx context.Context) error {
	return h.redisClient.Ping(ctx).Err()
}
