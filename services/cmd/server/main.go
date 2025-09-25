package main

import (
	commonDI "common/di"
	"services/internal/application"
	"services/internal/domain/user"
	"services/internal/infrastructure"
	"services/internal/interfaces/http"

	"go.uber.org/fx"
)

func main() {
	// 创建应用容器
	app := fx.New(
		// 使用common库的Web模块
		commonDI.GetWebModules(),

		// 领域模块
		user.DomainModule,

		// 应用模块
		application.ApplicationModule,

		// 基础设施模块
		infrastructure.InfrastructureModule,

		// 接口模块 - 使用重构后的路由管理
		http.InterfaceModuleFinal,
	)

	// 启动应用容器
	app.Run()
}
