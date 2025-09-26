package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/response"
)

const (
	// XForwardedFor X-Forwarded-For头
	XForwardedFor = "X-Forwarded-For"
	// XRealIP X-Real-IP头
	XRealIP = "X-Real-IP"
)

// IPWhitelistMiddleware IP白名单中间件
func IPWhitelistMiddleware(allowedIPs []string, zapLogger *zap.Logger) gin.HandlerFunc {
	// 解析允许的IP地址和CIDR
	var allowedSingleIPs []net.IP
	var allowedCIDRs []*net.IPNet

	for _, ipStr := range allowedIPs {
		// 尝试解析为CIDR
		if _, cidr, err := net.ParseCIDR(ipStr); err == nil {
			allowedCIDRs = append(allowedCIDRs, cidr)
		} else {
			// 尝试解析为单个IP
			if ip := net.ParseIP(ipStr); ip != nil {
				allowedSingleIPs = append(allowedSingleIPs, ip)
			} else {
				zapLogger.Warn("Invalid IP or CIDR in whitelist", zap.String("ip", ipStr))
			}
		}
	}

	zapLogger.Info("IP whitelist middleware initialized",
		zap.Int("single_ips", len(allowedSingleIPs)),
		zap.Int("cidr_ranges", len(allowedCIDRs)))

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 获取客户端IP
		clientIP := getClientIP(c)
		if clientIP == "" {
			logger.Warn(ctx, "Unable to determine client IP")
			response.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		// 解析IP地址
		ip := net.ParseIP(clientIP)
		if ip == nil {
			logger.Warn(ctx, "Invalid client IP",
				zap.String("client_ip", clientIP))
			response.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		// 检查是否在白名单中
		allowed := false

		// 检查单个IP
		for _, allowedIP := range allowedSingleIPs {
			if ip.Equal(allowedIP) {
				allowed = true
				break
			}
		}

		// 检查CIDR范围
		if !allowed {
			for _, cidr := range allowedCIDRs {
				if cidr.Contains(ip) {
					allowed = true
					break
				}
			}
		}

		if !allowed {
			logger.Warn(ctx, "IP not in whitelist",
				zap.String("client_ip", clientIP))
			response.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		logger.Debug(ctx, "IP whitelist check passed",
			zap.String("client_ip", clientIP))

		c.Next()
	}
}

// getClientIP 获取客户端真实IP
func getClientIP(c *gin.Context) string {
	// 优先检查X-Forwarded-For头
	xForwardedFor := c.GetHeader(XForwardedFor)
	if xForwardedFor != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 检查X-Real-IP头
	xRealIP := c.GetHeader(XRealIP)
	if xRealIP != "" {
		return xRealIP
	}

	// 使用gin内置方法获取IP
	return c.ClientIP()
}
