package netutil

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// XForwardedFor X-Forwarded-For头
	XForwardedFor = "X-Forwarded-For"
	// XRealIP X-Real-IP头
	XRealIP = "X-Real-IP"
	// XClientIP X-Client-IP头
	XClientIP = "X-Client-IP"
)

// GetClientIP 获取客户端真实IP
// 它按以下顺序检查头：X-Forwarded-For, X-Real-IP, X-Client-IP,
// 最后回退到 gin 的 c.ClientIP()。
func GetClientIP(c *gin.Context) string {
	// 1. 优先检查X-Forwarded-For头
	xForwardedFor := c.GetHeader(XForwardedFor)
	if xForwardedFor != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 2. 检查X-Real-IP头
	xRealIP := c.GetHeader(XRealIP)
	if xRealIP != "" {
		return xRealIP
	}

	// 3. 检查X-Client-IP头
	xClientIP := c.GetHeader(XClientIP)
	if xClientIP != "" {
		return xClientIP
	}

	// 4. 使用gin内置方法获取IP
	return c.ClientIP()
}
