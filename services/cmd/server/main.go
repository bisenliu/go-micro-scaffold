package main

import (
	"fmt"
	"log"

	"go.uber.org/fx"

	commonDI "common/di"
	"services/internal/application"
	"services/internal/domain/user"
	"services/internal/infrastructure"
	"services/internal/interfaces/http"
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

	if err := app.Err(); err != nil {
		// 尝试生成依赖关系图以帮助调试
		if visualization, verr := fx.VisualizeError(err); verr == nil {
			fmt.Println("Dependency graph visualization:")
			fmt.Println(visualization)
		}

		// 记录详细的错误信息并退出
		log.Fatalf("Failed to initialize application dependencies: %v", err)
	}

	// 启动应用容器
	app.Run()
}
