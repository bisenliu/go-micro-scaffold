package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/config"
	"common/databases/redis"
	"common/logger"
	"services/internal/infrastructure/persistence"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	dbProvider  *persistence.DatabaseProvider
	redisClient *redis.RedisClient
	config      *config.Config
}

// Ensure HealthHandler implements Handler interface
var _ Handler = (*HealthHandler)(nil)

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(
	dbProvider *persistence.DatabaseProvider,
	redisClient *redis.RedisClient,
	config *config.Config,
) *HealthHandler {
	return &HealthHandler{
		dbProvider:  dbProvider,
		redisClient: redisClient,
		config:      config,
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status" example:"healthy" description:"系统整体健康状态：healthy-健康，unhealthy-不健康"`
	Timestamp time.Time         `json:"timestamp" example:"2023-01-01T12:00:00Z" description:"检查时间戳"`
	Version   string            `json:"version" example:"1.0.0" description:"应用版本号"`
	Services  map[string]string `json:"services" description:"各个服务组件的健康状态"`
}

// Health 健康检查
// @Summary 系统健康检查
// @Description 检查系统各个组件的健康状态，包括数据库、Redis等服务的连接状态
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=HealthResponse} "系统健康"
// @Failure 503 {object} response.Response "系统不健康"
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	ctx := c.Request.Context()

	responseData := &HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services:  make(map[string]string),
	}

	// 检查数据库连接
	if err := h.checkDatabase(ctx); err != nil {
		responseData.Services["database"] = "unhealthy: " + err.Error()
		responseData.Status = "unhealthy"
		logger.Error(ctx, "Database health check failed", zap.Error(err))
	} else {
		responseData.Services["database"] = "healthy"
	}

	// 检查Redis连接
	if err := h.checkRedis(ctx); err != nil {
		responseData.Services["redis"] = "unhealthy: " + err.Error()
		responseData.Status = "unhealthy"
		logger.Error(ctx, "Redis health check failed", zap.Error(err))
	} else {
		responseData.Services["redis"] = "healthy"
	}

	statusCode := http.StatusOK
	if responseData.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, responseData)
}

// checkDatabase 检查数据库连接
func (h *HealthHandler) checkDatabase(ctx context.Context) error {
	client, err := h.dbProvider.GetHealthCheckClient()
	if err != nil {
		return err
	}
	return client.Ping(ctx)
}

// checkRedis 检查Redis连接
func (h *HealthHandler) checkRedis(ctx context.Context) error {
	return h.redisClient.Ping(ctx).Err()
}
