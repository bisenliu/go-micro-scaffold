package middleware

import (
	"github.com/gin-gonic/gin"
)

// MiddlewareBuilderAdapter 中间件构建器适配器
// 将内部的 MiddlewareBuilder 适配为接口
type MiddlewareBuilderAdapter struct {
	builder *MiddlewareBuilder
}

// BuildStandardChain 构建标准中间件链
func (a *MiddlewareBuilderAdapter) BuildStandardChain() []gin.HandlerFunc {
	chain := a.builder.BuildStandardChain()
	return chain.Build()
}

// BuildSecureChain 构建安全中间件链（包含认证）
func (a *MiddlewareBuilderAdapter) BuildSecureChain() []gin.HandlerFunc {
	chain := a.builder.BuildSecureChain()
	return chain.Build()
}

// BuildInternalChain 构建内网访问中间件链
func (a *MiddlewareBuilderAdapter) BuildInternalChain() []gin.HandlerFunc {
	chain := a.builder.BuildInternalChain()
	return chain.Build()
}

// BuildCustomChain 根据配置构建自定义中间件链
func (a *MiddlewareBuilderAdapter) BuildCustomChain(config interface{}) []gin.HandlerFunc {
	if middlewareConfig, ok := config.(*MiddlewareConfig); ok {
		chain := a.builder.BuildCustomChain(middlewareConfig)
		return chain.Build()
	}
	
	// 如果配置类型不匹配，返回标准链
	return a.BuildStandardChain()
}