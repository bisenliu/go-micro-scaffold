package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"common/interfaces"
)

// ServerParams 定义创建HTTP服务器所需的依赖
type ServerParams struct {
	fx.In

	Engine *gin.Engine
	ConfigProvider interfaces.ConfigProvider
}

// NewServer 创建一个新的HTTP服务器实例
func NewServer(params ServerParams) *http.Server {
	// 从配置提供者中读取服务器配置
	serverConfig := params.ConfigProvider.GetServerConfig()
	
	// 使用配置中的值，如果未设置则使用默认值
	readTimeout := serverConfig.ReadTimeout
	if readTimeout == 0 {
		readTimeout = 5 * time.Second
	}

	writeTimeout := serverConfig.WriteTimeout
	if writeTimeout == 0 {
		writeTimeout = 10 * time.Second
	}

	idleTimeout := serverConfig.IdleTimeout
	if idleTimeout == 0 {
		idleTimeout = 120 * time.Second
	}

	maxHeaderBytes := serverConfig.MaxHeaderBytes
	if maxHeaderBytes == 0 {
		maxHeaderBytes = 1 << 20 // 1MB
	}

	server := &http.Server{
		Addr:           ":" + serverConfig.Port,
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
	Logger    interfaces.Logger
}

// RegisterServerLifecycle 注册服务器的启动和停止生命周期钩子
func RegisterServerLifecycle(params ServerLifecycleParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.Logger.Info(ctx, "Starting HTTP server", zap.String("addr", params.Server.Addr))

			// 在goroutine中启动服务器，避免阻塞应用启动
			go func() {
				if err := params.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					params.Logger.Error(ctx, "Failed to start HTTP server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info(ctx, "Shutting down HTTP server")

			// 在这里可以添加服务特定的清理逻辑
			// 例如：断开数据库连接、确保消息队列中的消息处理完毕等

			// 执行标准的服务器关闭流程
			return params.Server.Shutdown(ctx)
		},
	})
}
