// Package interfaces 定义了 common 模块对外提供的所有接口
// 这些接口定义了模块间的契约，确保依赖倒置原则的实施
package interfaces

// CommonServices 聚合所有 common 服务的接口
// 这个接口将被 services 模块使用，实现对 common 模块的接口依赖
type CommonServices interface {
	// Config 获取配置提供者
	Config() ConfigProvider
	
	// Logger 获取日志器
	Logger() Logger
	
	// Database 获取数据库管理器
	Database() DatabaseManager
	
	// Cache 获取缓存管理器
	Cache() CacheManager
	
	// Validator 获取验证器
	Validator() Validator
	
	// Middleware 获取中间件提供者
	Middleware() MiddlewareProvider
	
	// JWT 获取JWT服务
	JWT() JWTService
}

// ServiceProvider 服务提供者接口
// 定义了服务的生命周期管理
type ServiceProvider interface {
	// Start 启动服务
	Start() error
	
	// Stop 停止服务
	Stop() error
	
	// Health 健康检查
	Health() error
	
	// Name 获取服务名称
	Name() string
}

// HealthChecker 健康检查接口
type HealthChecker interface {
	// Check 执行健康检查
	Check() HealthStatus
}

// HealthStatus 健康状态
type HealthStatus struct {
	Service string            `json:"service"`
	Status  string            `json:"status"` // "healthy", "unhealthy", "degraded"
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}