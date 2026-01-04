// Package main Go Micro Scaffold API
//
// 微服务脚手架 API 文档
//
// 这是一个基于 Go 和 Gin 框架构建的微服务脚手架项目的 API 文档。
// 提供用户管理、认证授权等核心功能的 RESTful API 接口。
//
//	@title						Go Micro Scaffold API
//	@version					1.0.0
//	@description				微服务脚手架 API 文档
//	@termsOfService				http://swagger.io/terms/
//
//	@contact.name				API Support
//	@contact.email				support@example.com
//	@contact.url				http://www.example.com/support
//
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//
//	@host						localhost:8080
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Token，格式：Bearer {token}
//
//	@schemes					http https
//	@produce					json
//	@consumes					json
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go.uber.org/fx"

	commonDI "common/di"
	"user-services/internal/application"
	"user-services/internal/domain/user"
	"user-services/internal/infrastructure"
	"user-services/internal/interfaces/http"
)

func main() {
	// 添加命令行参数
	var (
		generateGraph = flag.Bool("graph", false, "Generate dependency graph and exit")
		graphOutput   = flag.String("graph-output", "dependency-graph.dot", "Output file for dependency graph")
	)
	flag.Parse()

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

		// 接口模块
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

	// 如果请求生成依赖图
	if *generateGraph {
		generateDependencyGraph(app, *graphOutput)
		return
	}

	// 启动应用容器
	app.Run()
}

// generateDependencyGraph 生成依赖关系图
func generateDependencyGraph(app *fx.App, outputFile string) {
	// 创建一个新的应用来获取依赖图
	var dotGraph string

	graphApp := fx.New(
		// 使用common库的Web模块（包括HTTP模块）
		commonDI.GetWebModules(),

		// 领域模块
		user.DomainModule,

		// 应用模块
		application.ApplicationModule,

		// 基础设施模块
		infrastructure.InfrastructureModule,

		// 接口模块（包含新的HTTP服务器实现）
		http.InterfaceModuleFinal,

		// 提供 DotGraph
		fx.Provide(func(graph fx.DotGraph) string {
			return string(graph)
		}),

		// 获取依赖图
		fx.Invoke(func(graph fx.DotGraph) {
			dotGraph = string(graph)
		}),
	)

	if err := graphApp.Err(); err != nil {
		log.Fatalf("Failed to generate dependency graph: %v", err)
	}

	// 写入文件
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create graph output file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(dotGraph); err != nil {
		log.Fatalf("Failed to write dependency graph: %v", err)
	}

	/*
		# 安装 Graphviz 后运行 (转为png)
		dot -Tpng dependency-graph.dot -o graph.png
	*/
	fmt.Printf("Dependency graph generated successfully: %s\n", outputFile)
	fmt.Println("To visualize the graph, you can:")
	fmt.Printf("1. Install Graphviz: brew install graphviz (macOS) or apt-get install graphviz (Ubuntu)\n")
	fmt.Printf("2. Generate PNG: dot -Tpng %s -o dependency-graph.png\n", outputFile)
	fmt.Printf("3. Generate SVG: dot -Tsvg %s -o dependency-graph.svg\n", outputFile)
	fmt.Printf("4. View online: Upload %s to http://magjac.com/graphviz-visual-editor/\n", outputFile)
}
