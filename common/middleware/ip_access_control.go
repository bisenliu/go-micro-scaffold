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
	// XClientIP X-Client-IP头
	XClientIP = "X-Client-IP"
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

// InternalIPMiddleware 内网IP访问限制中间件
func InternalIPMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	// 定义私有IP网段
	privateIPBlocks := []*net.IPNet{}

	// 初始化私有IP网段
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918 - A类私有地址
		"172.16.0.0/12",  // RFC1918 - B类私有地址
		"192.168.0.0/16", // RFC1918 - C类私有地址
		"169.254.0.0/16", // RFC3927 - 链路本地地址
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 链路本地地址
		"fc00::/7",       // IPv6 唯一本地地址
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			zapLogger.Fatal("Failed to parse CIDR", zap.String("cidr", cidr), zap.Error(err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		clientIP := getClientIP(c)

		if !isPrivateIP(clientIP, privateIPBlocks) {
			logger.Warn(ctx, "Access denied for external IP",
				zap.String("ip", clientIP),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
			response.Forbidden(c, "Access denied: Internal network only")
			c.Abort()
			return
		}

		logger.Debug(ctx, "Internal IP access granted",
			zap.String("ip", clientIP),
			zap.String("path", c.Request.URL.Path))
		c.Next()
	}
}

// getClientIP 获取客户端真实IP
func getClientIP(c *gin.Context) string {
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

// isPrivateIP 判断IP是否为私有/内网IP
func isPrivateIP(ipStr string, privateIPBlocks []*net.IPNet) bool {
	// 解析IP地址
	ip := net.ParseIP(ipStr)
	if ip == nil {
		// 无法解析的IP地址，拒绝访问
		return false
	}

	// 检查是否为回环地址（127.0.0.1, ::1）
	if ip.IsLoopback() {
		return true
	}

	// 检查是否为链路本地地址（169.254.x.x, fe80::）
	if ip.IsLinkLocalUnicast() {
		return true
	}

	// 检查是否为链路本地多播地址
	if ip.IsLinkLocalMulticast() {
		return true
	}

	// 检查是否在私有IP网段内
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
