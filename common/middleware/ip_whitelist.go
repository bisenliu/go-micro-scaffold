package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
)

// InternalIPMiddleware IP白名单中间件
func InternalIPMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	// 内网IP白名单
	allowedIPs := []string{
		"127.0.0.1",
		"::1",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	// 预解析CIDR
	var allowedCIDRs []*net.IPNet
	var allowedSingleIPs []net.IP

	for _, ipStr := range allowedIPs {
		if ip := net.ParseIP(ipStr); ip != nil {
			allowedSingleIPs = append(allowedSingleIPs, ip)
		} else if _, cidr, err := net.ParseCIDR(ipStr); err == nil {
			allowedCIDRs = append(allowedCIDRs, cidr)
		}
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		clientIP := c.ClientIP()

		// 解析客户端IP
		ip := net.ParseIP(clientIP)
		if ip == nil {
			logger.Warn(zapLogger, ctx, "Invalid client IP",
				zap.String("client_ip", clientIP))
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "Access denied",
			})
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
			logger.Warn(zapLogger, ctx, "IP not in whitelist",
				zap.String("client_ip", clientIP))
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "Access denied",
			})
			c.Abort()
			return
		}

		logger.Debug(zapLogger, ctx, "IP whitelist check passed",
			zap.String("client_ip", clientIP))

		c.Next()
	}
}
