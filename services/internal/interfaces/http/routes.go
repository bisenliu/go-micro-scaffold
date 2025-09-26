package http

import (
	"common/middleware"
	"common/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"services/internal/interfaces/http/handler"
)

// SetupRoutesFinal 最终推荐的路由设置方案
func SetupRoutesFinal(
	engine *gin.Engine,
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	zapLogger *zap.Logger,
) {
	engine.Use(middleware.LoggerMiddleware(zapLogger))

	// 1. 系统路由（无需认证）
	setupSystemRoutesFinal(engine, healthHandler, zapLogger)

	// 2. API v1 路由组（需要认证）
	v1 := engine.Group("/api/v1")
	{
		// 这里可以添加认证中间件
		// v1.Use(middlewareManager.AuthMiddleware())

		setupUserAPIRoutes(v1, userHandler, zapLogger)
		// 后续添加其他模块
		// setupOrderAPIRoutes(v1, orderHandler, zapLogger)
		// setupProductAPIRoutes(v1, productHandler, zapLogger)
		// setupPaymentAPIRoutes(v1, paymentHandler, zapLogger)
	}

	// 3. 内部API路由组（仅内网访问）
	// internal := engine.Group("/api/internal")
	// {
	// 	// 设置IP白名单中间件
	// 	// middlewareManager.SetupInternalMiddleware(internal)

	// }

	zapLogger.Info("All routes setup completed successfully")
}

// setupSystemRoutesFinal 设置系统路由
func setupSystemRoutesFinal(engine *gin.Engine, healthHandler *handler.HealthHandler, logger *zap.Logger) {
	engine.GET("/health", healthHandler.Health)
	engine.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"message": "pong"})
	})

	logger.Info("System routes registered", zap.Int("count", 2))
}

// setupUserAPIRoutes 设置用户API路由
func setupUserAPIRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, logger *zap.Logger) {
	users := rg.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
	}

	logger.Info("User API routes registered", zap.Int("count", 5))
}

// 以下是其他模块的路由设置示例，展示如何扩展到70-80个路由

// setupOrderAPIRoutes 设置订单API路由
// func setupOrderAPIRoutes(rg *gin.RouterGroup, orderHandler *handler.OrderHandler, logger *zap.Logger) {
//     orders := rg.Group("/orders")
//     {
//         orders.POST("", orderHandler.CreateOrder)           // 创建订单
//         orders.GET("/:id", orderHandler.GetOrder)           // 获取订单详情
//         orders.PUT("/:id", orderHandler.UpdateOrder)        // 更新订单
//         orders.DELETE("/:id", orderHandler.CancelOrder)     // 取消订单
//         orders.GET("", orderHandler.ListOrders)             // 订单列表
//         orders.PUT("/:id/pay", orderHandler.PayOrder)       // 支付订单
//         orders.PUT("/:id/confirm", orderHandler.ConfirmOrder) // 确认订单
//         orders.PUT("/:id/ship", orderHandler.ShipOrder)     // 发货
//         orders.PUT("/:id/complete", orderHandler.CompleteOrder) // 完成订单
//         orders.GET("/:id/status", orderHandler.GetOrderStatus) // 订单状态
//     }
//
//     logger.Info("Order API routes registered", zap.Int("count", 10))
// }

// setupProductAPIRoutes 设置产品API路由
// func setupProductAPIRoutes(rg *gin.RouterGroup, productHandler *handler.ProductHandler, logger *zap.Logger) {
//     products := rg.Group("/products")
//     {
//         products.POST("", productHandler.CreateProduct)      // 创建产品
//         products.GET("/:id", productHandler.GetProduct)      // 获取产品详情
//         products.PUT("/:id", productHandler.UpdateProduct)   // 更新产品
//         products.DELETE("/:id", productHandler.DeleteProduct) // 删除产品
//         products.GET("", productHandler.ListProducts)        // 产品列表
//         products.PUT("/:id/status", productHandler.UpdateProductStatus) // 更新产品状态
//
//         // 产品分类相关
//         categories := products.Group("/:id/categories")
//         {
//             categories.POST("", productHandler.AddProductCategory)       // 添加产品分类
//             categories.DELETE("/:categoryId", productHandler.RemoveProductCategory) // 移除产品分类
//         }
//
//         // 产品库存相关
//         inventory := products.Group("/:id/inventory")
//         {
//             inventory.GET("", productHandler.GetProductInventory)        // 获取库存
//             inventory.PUT("", productHandler.UpdateProductInventory)     // 更新库存
//             inventory.POST("/adjust", productHandler.AdjustProductInventory) // 调整库存
//         }
//
//         // 产品价格相关
//         pricing := products.Group("/:id/pricing")
//         {
//             pricing.GET("", productHandler.GetProductPricing)           // 获取价格
//             pricing.PUT("", productHandler.UpdateProductPricing)        // 更新价格
//             pricing.GET("/history", productHandler.GetProductPriceHistory) // 价格历史
//         }
//     }
//
//     logger.Info("Product API routes registered", zap.Int("count", 15))
// }

// setupPaymentAPIRoutes 设置支付API路由
// func setupPaymentAPIRoutes(rg *gin.RouterGroup, paymentHandler *handler.PaymentHandler, logger *zap.Logger) {
//     payments := rg.Group("/payments")
//     {
//         payments.POST("", paymentHandler.CreatePayment)              // 创建支付
//         payments.GET("/:id", paymentHandler.GetPayment)              // 获取支付详情
//         payments.PUT("/:id/confirm", paymentHandler.ConfirmPayment)  // 确认支付
//         payments.PUT("/:id/cancel", paymentHandler.CancelPayment)    // 取消支付
//         payments.POST("/:id/refund", paymentHandler.RefundPayment)   // 退款
//         payments.GET("", paymentHandler.ListPayments)               // 支付列表
//         payments.GET("/:id/status", paymentHandler.GetPaymentStatus) // 支付状态
//
//         // 支付方式相关
//         methods := payments.Group("/methods")
//         {
//             methods.GET("", paymentHandler.ListPaymentMethods)       // 支付方式列表
//             methods.POST("", paymentHandler.AddPaymentMethod)        // 添加支付方式
//             methods.DELETE("/:methodId", paymentHandler.RemovePaymentMethod) // 删除支付方式
//         }
//
//         // 支付回调
//         callbacks := payments.Group("/callbacks")
//         {
//             callbacks.POST("/wechat", paymentHandler.WechatCallback)  // 微信支付回调
//             callbacks.POST("/alipay", paymentHandler.AlipayCallback) // 支付宝回调
//         }
//     }
//
//     logger.Info("Payment API routes registered", zap.Int("count", 12))
// }

// setupTeamAPIRoutes 设置团队API路由
// func setupTeamAPIRoutes(rg *gin.RouterGroup, teamHandler *handler.TeamHandler, logger *zap.Logger) {
//     teams := rg.Group("/teams")
//     {
//         teams.POST("", teamHandler.CreateTeam)                    // 创建团队
//         teams.GET("/:id", teamHandler.GetTeam)                    // 获取团队详情
//         teams.PUT("/:id", teamHandler.UpdateTeam)                 // 更新团队
//         teams.DELETE("/:id", teamHandler.DeleteTeam)              // 删除团队
//         teams.GET("", teamHandler.ListTeams)                      // 团队列表
//
//         // 团队成员管理
//         members := teams.Group("/:id/members")
//         {
//             members.POST("", teamHandler.AddTeamMember)            // 添加成员
//             members.DELETE("/:memberId", teamHandler.RemoveTeamMember) // 移除成员
//             members.GET("", teamHandler.ListTeamMembers)           // 成员列表
//             members.PUT("/:memberId/role", teamHandler.UpdateMemberRole) // 更新成员角色
//         }
//
//         // 团队活动
//         activities := teams.Group("/:id/activities")
//         {
//             activities.POST("", teamHandler.CreateTeamActivity)    // 创建活动
//             activities.GET("", teamHandler.ListTeamActivities)     // 活动列表
//             activities.PUT("/:activityId", teamHandler.UpdateTeamActivity) // 更新活动
//             activities.DELETE("/:activityId", teamHandler.DeleteTeamActivity) // 删除活动
//         }
//     }
//
//     logger.Info("Team API routes registered", zap.Int("count", 13))
// }

// setupOrderInternalRoutes 设置订单内部路由
// func setupOrderInternalRoutes(rg *gin.RouterGroup, orderHandler *handler.OrderHandler, logger *zap.Logger) {
//     orders := rg.Group("/orders")
//     {
//         orders.PUT("/:id/force-cancel", orderHandler.ForceCancelOrder)   // 强制取消订单
//         orders.PUT("/:id/force-refund", orderHandler.ForceRefundOrder)   // 强制退款
//         orders.GET("/statistics", orderHandler.GetOrderStatistics)      // 订单统计
//         orders.POST("/batch-update", orderHandler.BatchUpdateOrders)     // 批量更新订单
//     }
//
//     logger.Info("Order internal routes registered", zap.Int("count", 4))
// }
