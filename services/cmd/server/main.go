package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go.uber.org/fx"
)

func main() {
	// 添加命令行参数
	var (
		generateGraph = flag.Bool("graph", false, "Generate dependency graph and exit")
		graphOutput   = flag.String("graph-output", "dependency-graph.dot", "Output file for dependency graph")
		showModules   = flag.Bool("modules", false, "Show module information and exit")
	)
	flag.Parse()

	// 创建应用配置
	config := NewAppConfig()
	
	// 创建模块注册表
	registry := NewModuleRegistry()

	// 如果请求显示模块信息
	if *showModules {
		showModuleInfo(registry)
		return
	}

	// 创建应用容器
	app := createApplication(registry, config)

	// 检查应用初始化错误
	if err := app.Err(); err != nil {
		handleApplicationError(err)
	}

	// 如果请求生成依赖图
	if *generateGraph {
		generateDependencyGraph(registry, *graphOutput)
		return
	}

	// 启动应用容器
	fmt.Printf("Starting application in %s mode...\n", config.Environment)
	app.Run()
}

// createApplication 创建应用容器
func createApplication(registry *ModuleRegistry, config *AppConfig) *fx.App {
	modules := registry.GetModulesForConfig(config)
	
	return fx.New(modules...)
}

// handleApplicationError 处理应用初始化错误
func handleApplicationError(err error) {
	// 尝试生成依赖关系图以帮助调试
	if visualization, verr := fx.VisualizeError(err); verr == nil {
		fmt.Println("=== Dependency Graph Visualization ===")
		fmt.Println(visualization)
		fmt.Println("=====================================")
	}

	// 记录详细的错误信息并退出
	log.Fatalf("Failed to initialize application dependencies: %v", err)
}

// showModuleInfo 显示模块信息
func showModuleInfo(registry *ModuleRegistry) {
	fmt.Println("=== Application Module Information ===")
	
	modules := registry.GetModuleInfo()
	for _, module := range modules {
		fmt.Printf("%d. %s\n", module.Order, module.Name)
		fmt.Printf("   Description: %s\n", module.Description)
		fmt.Println()
	}
	
	fmt.Println("Module loading order follows Clean Architecture principles:")
	fmt.Println("- Inner layers (Domain) don't depend on outer layers")
	fmt.Println("- Dependencies flow inward toward the domain")
	fmt.Println("- Infrastructure implements domain interfaces")
}

// generateDependencyGraph 生成依赖关系图
func generateDependencyGraph(registry *ModuleRegistry, outputFile string) {
	fmt.Println("Generating dependency graph...")
	
	// 创建一个新的应用来获取依赖图
	var dotGraph string

	// 使用完整的模块集合生成依赖图
	modules := registry.GetAllModules()
	
	graphApp := fx.New(
		append(modules, 
			// 提供 DotGraph
			fx.Provide(func(graph fx.DotGraph) string {
				return string(graph)
			}),

			// 获取依赖图
			fx.Invoke(func(graph fx.DotGraph) {
				dotGraph = string(graph)
			}),
		)...,
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

	fmt.Printf("Dependency graph generated successfully: %s\n", outputFile)
	fmt.Println("\nTo visualize the graph, you can:")
	fmt.Printf("1. Install Graphviz: brew install graphviz (macOS) or apt-get install graphviz (Ubuntu)\n")
	fmt.Printf("2. Generate PNG: dot -Tpng %s -o dependency-graph.png\n", outputFile)
	fmt.Printf("3. Generate SVG: dot -Tsvg %s -o dependency-graph.svg\n", outputFile)
	fmt.Printf("4. View online: Upload %s to http://magjac.com/graphviz-visual-editor/\n", outputFile)
}
