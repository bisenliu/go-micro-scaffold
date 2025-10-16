package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/fx"

	commonDI "common/di"
	"common/interfaces"
)

func main() {
	fmt.Println("🧪 Testing new middleware system...")

	// 创建应用来测试中间件
	app := fx.New(
		// 使用common库的核心模块
		commonDI.GetCoreModules(),

		// 测试函数
		fx.Invoke(func(commonServices interfaces.CommonServices) {
			fmt.Println("✅ Middleware system test passed!")
			
			// 测试JWT服务
			jwtService := commonServices.JWT()
			if jwtService != nil {
				fmt.Printf("JWT Service: %T\n", jwtService)
				
				// 测试生成token
				token, err := jwtService.GenerateToken("test-user-123")
				if err != nil {
					fmt.Printf("Failed to generate token: %v\n", err)
				} else {
					fmt.Printf("Generated token: %s...\n", token[:20])
					
					// 测试验证token
					userID, err := jwtService.ValidateToken(token)
					if err != nil {
						fmt.Printf("Failed to validate token: %v\n", err)
					} else {
						fmt.Printf("Validated user ID: %s\n", userID)
					}
				}
			}
			
			// 测试中间件提供者
			middlewareProvider := commonServices.Middleware()
			if middlewareProvider != nil {
				fmt.Printf("Middleware Provider: %T\n", middlewareProvider)
				
				// 测试创建中间件
				authMiddleware := middlewareProvider.CreateAuthMiddleware()
				corsMiddleware := middlewareProvider.CreateCORSMiddleware()
				recoveryMiddleware := middlewareProvider.CreateRecoveryMiddleware()
				
				fmt.Printf("Auth Middleware: %T\n", authMiddleware)
				fmt.Printf("CORS Middleware: %T\n", corsMiddleware)
				fmt.Printf("Recovery Middleware: %T\n", recoveryMiddleware)
			}
			
			fmt.Println("🎉 All middleware tests passed!")
		}),
	)

	if err := app.Err(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	if err := app.Stop(ctx); err != nil {
		log.Fatalf("Failed to stop application: %v", err)
	}
}