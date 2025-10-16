package middleware

import (
	"github.com/gin-gonic/gin"
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	// Auth 认证中间件配置
	Auth *AuthConfig `json:"auth,omitempty"`
	
	// CORS 跨域中间件配置
	CORS *CORSConfig `json:"cors,omitempty"`
	
	// RateLimit 限流中间件配置
	RateLimit *RateLimitConfig `json:"rate_limit,omitempty"`
	
	// Recovery 恢复中间件配置
	Recovery *RecoveryConfig `json:"recovery,omitempty"`
	
	// RequestLog 请求日志中间件配置
	RequestLog *RequestLogConfig `json:"request_log,omitempty"`
	
	// IPAccess IP访问控制配置
	IPAccess *IPAccessConfig `json:"ip_access,omitempty"`
}

// AuthConfig 认证中间件配置
type AuthConfig struct {
	Enabled   bool     `json:"enabled"`
	Whitelist []string `json:"whitelist"`
}

// CORSConfig CORS中间件配置
type CORSConfig struct {
	Enabled          bool     `json:"enabled"`
	AllowOrigins     []string `json:"allow_origins"`
	AllowMethods     []string `json:"allow_methods"`
	AllowHeaders     []string `json:"allow_headers"`
	ExposeHeaders    []string `json:"expose_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
	MaxAge           int      `json:"max_age"`
}

// RateLimitConfig 限流中间件配置
type RateLimitConfig struct {
	Enabled bool `json:"enabled"`
	Rate    int  `json:"rate"`    // 每秒请求数
	Burst   int  `json:"burst"`   // 突发请求数
}

// RecoveryConfig 恢复中间件配置
type RecoveryConfig struct {
	Enabled     bool `json:"enabled"`
	LogStack    bool `json:"log_stack"`
	LogRequest  bool `json:"log_request"`
}

// RequestLogConfig 请求日志中间件配置
type RequestLogConfig struct {
	Enabled    bool     `json:"enabled"`
	SkipPaths  []string `json:"skip_paths"`
	LogBody    bool     `json:"log_body"`
	LogHeaders bool     `json:"log_headers"`
}

// IPAccessConfig IP访问控制配置
type IPAccessConfig struct {
	Enabled       bool     `json:"enabled"`
	WhitelistIPs  []string `json:"whitelist_ips"`
	InternalOnly  bool     `json:"internal_only"`
}

// MiddlewareChain 中间件链
type MiddlewareChain struct {
	middlewares []gin.HandlerFunc
	names       []string
}

// NewMiddlewareChain 创建新的中间件链
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]gin.HandlerFunc, 0),
		names:       make([]string, 0),
	}
}

// Add 添加中间件到链中
func (mc *MiddlewareChain) Add(name string, middleware gin.HandlerFunc) *MiddlewareChain {
	if middleware != nil {
		mc.middlewares = append(mc.middlewares, middleware)
		mc.names = append(mc.names, name)
	}
	return mc
}

// Build 构建中间件链
func (mc *MiddlewareChain) Build() []gin.HandlerFunc {
	return mc.middlewares
}

// Names 获取中间件名称列表
func (mc *MiddlewareChain) Names() []string {
	return mc.names
}

// Length 获取中间件数量
func (mc *MiddlewareChain) Length() int {
	return len(mc.middlewares)
}