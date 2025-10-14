package middleware

import (
	"net"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/response"
)

// IPWhitelistMiddleware IP白名单中间件
func IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
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
			}
		}
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 从Context中获取预解析的IP
		ipVal, exists := c.Get(ClientParsedIPContextKey)
		if !exists {
			logger.Error(ctx, "Parsed client IP not found in context")
			response.Forbidden(c, "Access denied: Internal error")
			c.Abort()
			return
		}
		ip, ok := ipVal.(net.IP)
		if !ok {
			logger.Error(ctx, "Parsed client IP in context is not of type net.IP")
			response.Forbidden(c, "Access denied: Internal error")
			c.Abort()
			return
		}

		clientIPStrVal, exists := c.Get(ClientIPContextKey)
		if !exists {
			logger.Error(ctx, "Client IP string not found in context")
			response.Forbidden(c, "Access denied: Internal error")
			c.Abort()
return
		}
		clientIPStr, ok := clientIPStrVal.(string)
		if !ok {
			logger.Error(ctx, "Client IP string in context is not of type string")
			response.Forbidden(c, "Access denied: Internal error")
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
				zap.String("client_ip", clientIPStr))
			response.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		logger.Debug(ctx, "IP whitelist check passed",
			zap.String("client_ip", clientIPStr))

		c.Next()
	}
}

// InternalIPMiddleware 内网IP访问限制中间件
func InternalIPMiddleware() gin.HandlerFunc {
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
			panic("Failed to parse CIDR '" + cidr + "': " + err.Error())
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 从Context中获取预解析的IP
		ipVal, exists := c.Get(ClientParsedIPContextKey)
		if !exists {
			logger.Error(ctx, "Parsed client IP not found in context")
			response.Forbidden(c, "Access denied: Internal error")
			c.Abort()
			return
		}
		ip, ok := ipVal.(net.IP)
		if !ok {
			logger.Error(ctx, "Parsed client IP in context is not of type net.IP")
			response.Forbidden(c, "Access denied: Internal error")
			c.Abort()
			return
		}

		clientIPStrVal, exists := c.Get(ClientIPContextKey)
		if !exists {
			logger.Error(ctx, "Client IP string not found in context")
			response.Forbidden(c, "Access denied: Internal error")
			c.Abort()
			return
		}
		clientIPStr, ok := clientIPStrVal.(string)
		if !ok {
			logger.Error(ctx, "Client IP string in context is not of type string")
			response.Forbidden(c, "Access denied: Internal error")
			c.Abort()
			return
		}

		if !isPrivateIP(ip, privateIPBlocks) {
			logger.Warn(ctx, "Access denied for external IP",
				zap.String("ip", clientIPStr),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
			response.Forbidden(c, "Access denied: Internal network only")
			c.Abort()
			return
		}

		logger.Debug(ctx, "Internal IP access granted",
			zap.String("ip", clientIPStr),
			zap.String("path", c.Request.URL.Path))
		c.Next()
	}
}

// isPrivateIP 判断IP是否为私有/内网IP
func isPrivateIP(ip net.IP, privateIPBlocks []*net.IPNet) bool {
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
