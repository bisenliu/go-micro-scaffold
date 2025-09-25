package http

import (
	"common/config"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server HTTP服务器
type Server struct {
	*http.Server
	engine *gin.Engine
}

// ServerParams 服务器依赖参数
type ServerParams struct {
	fx.In
	Config *config.Config
	Logger *zap.Logger
}

// NewServer 创建新的HTTP服务器
func NewServer(params ServerParams) *Server {
	// 设置Gin模式
	gin.SetMode(params.Config.Server.Mode)

	// 创建Gin引擎
	engine := gin.New()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    ":" + params.Config.Server.Port,
		Handler: engine,
	}

	return &Server{
		Server: server,
		engine: engine,
	}
}

// GetEngine 获取Gin引擎
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

// Start 启动服务器
func (s *Server) Start(ctx context.Context, logger *zap.Logger) error {
	logger.Info("Starting HTTP server", zap.String("addr", s.Addr))

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop(ctx context.Context, logger *zap.Logger) error {
	logger.Info("Stopping HTTP server")
	return s.Shutdown(ctx)
}

// InvokeServerLifecycle 启动和停止服务器的生命周期
func InvokeServerLifecycle(lc fx.Lifecycle, server *Server, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start(ctx, logger)
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop(ctx, logger)
		},
	})
}

// Module FX模块
var Module = fx.Module("http",
	fx.Provide(NewServer),
	fx.Provide(func(s *Server) *gin.Engine {
		return s.GetEngine()
	}),
	fx.Invoke(InvokeServerLifecycle),
)
