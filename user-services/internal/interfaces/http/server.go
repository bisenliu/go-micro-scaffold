package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/config"
)

// ServerParams 定义创建HTTP服务器所需的依赖
type ServerParams struct {
	fx.In

	Engine *gin.Engine
	Config *config.Config
	Logger *zap.Logger
}

// NewServer 创建一个新的HTTP服务器实例
func NewServer(params ServerParams) *http.Server {
	// 从配置文件中读取服务器配置，如果未设置则使用默认值
	readTimeout := params.Config.Server.ReadTimeout
	if readTimeout == 0 {
		readTimeout = 5 * time.Second
	}

	writeTimeout := params.Config.Server.WriteTimeout
	if writeTimeout == 0 {
		writeTimeout = 10 * time.Second
	}

	idleTimeout := params.Config.Server.IdleTimeout
	if idleTimeout == 0 {
		idleTimeout = 120 * time.Second
	}

	maxHeaderBytes := params.Config.Server.MaxHeaderBytes
	if maxHeaderBytes == 0 {
		maxHeaderBytes = 1 << 20 // 1MB
	}

	server := &http.Server{
		Addr:           ":" + params.Config.Server.Port,
		Handler:        params.Engine,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	return server
}

// ServerLifecycleParams 定义服务器生命周期管理所需的依赖
type ServerLifecycleParams struct {
	fx.In

	Server    *http.Server
	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

// RegisterServerLifecycle 注册服务器的启动和停止生命周期钩子
func RegisterServerLifecycle(params ServerLifecycleParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			params.Logger.Info("Starting HTTP server", zap.String("addr", params.Server.Addr))

			// 在goroutine中启动服务器，避免阻塞应用启动
			go func() {
				if err := params.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					params.Logger.Error("Failed to start HTTP server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Shutting down HTTP server")

			// 在这里可以添加服务特定的清理逻辑
			// 例如：断开数据库连接、确保消息队列中的消息处理完毕等

			// 执行标准的服务器关闭流程
			return params.Server.Shutdown(ctx)
		},
	})
}
