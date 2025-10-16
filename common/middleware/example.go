package middleware

import (
	"github.com/gin-gonic/gin"
)

// ExampleUsage 展示如何使用新的中间件配置系统
func ExampleUsage() {
	// 这个文件仅用于展示用法，不会被编译到最终程序中
	
	// 1. 使用标准中间件链
	_ = func(provider interface{}) {
		// middlewareProvider := provider.(interfaces.MiddlewareProvider)
		// builder := middlewareProvider.GetBuilder()
		// standardMiddlewares := builder.BuildStandardChain()
		// 
		// router := gin.New()
		// router.Use(standardMiddlewares...)
	}
	
	// 2. 使用安全中间件链（包含认证）
	_ = func(provider interface{}) {
		// middlewareProvider := provider.(interfaces.MiddlewareProvider)
		// builder := middlewareProvider.GetBuilder()
		// secureMiddlewares := builder.BuildSecureChain()
		// 
		// router := gin.New()
		// router.Use(secureMiddlewares...)
	}
	
	// 3. 使用自定义配置
	_ = func(provider interface{}) {
		// config := &MiddlewareConfig{
		// 	Auth: &AuthConfig{
		// 		Enabled:   true,
		// 		Whitelist: []string{"/health", "/metrics"},
		// 	},
		// 	CORS: &CORSConfig{
		// 		Enabled: true,
		// 		AllowOrigins: []string{"https://example.com"},
		// 	},
		// 	RateLimit: &RateLimitConfig{
		// 		Enabled: true,
		// 		Rate:    100,
		// 		Burst:   10,
		// 	},
		// }
		// 
		// middlewareProvider := provider.(interfaces.MiddlewareProvider)
		// builder := middlewareProvider.GetBuilder()
		// customMiddlewares := builder.BuildCustomChain(config)
		// 
		// router := gin.New()
		// router.Use(customMiddlewares...)
	}
	
	// 4. 手动构建中间件链
	_ = func(provider interface{}) {
		// middlewareProvider := provider.(interfaces.MiddlewareProvider)
		// 
		// router := gin.New()
		// router.Use(middlewareProvider.CreateRecoveryMiddleware())
		// router.Use(middlewareProvider.CreateCORSMiddleware())
		// router.Use(middlewareProvider.CreateRequestLogMiddleware())
		// router.Use(middlewareProvider.CreateRateLimitMiddleware())
		// router.Use(middlewareProvider.CreateAuthMiddleware())
	}
}

// ChainExample 展示中间件链的使用
func ChainExample() {
	// 创建中间件链
	chain := NewMiddlewareChain()
	
	// 添加中间件
	chain.Add("recovery", gin.Recovery())
	chain.Add("logger", gin.Logger())
	
	// 获取中间件列表
	middlewares := chain.Build()
	names := chain.Names()
	
	// 应用到路由器
	router := gin.New()
	router.Use(middlewares...)
	
	// 打印中间件信息
	_ = names // ["recovery", "logger"]
	_ = chain.Length() // 2
}