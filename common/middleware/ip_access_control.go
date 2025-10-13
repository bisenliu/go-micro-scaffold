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

	// ClientIPContextKey 是在 gin.Context 中存储客户端 IP 字符串的键
	ClientIPContextKey = "clientIP"
	// ClientParsedIPContextKey 是在 gin.Context 中存储解析后的 net.IP 对象的键
	ClientParsedIPContextKey = "clientParsedIP"
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

// ExtractClientIPMiddleware 提取客户端IP并解析为net.IP，存储在Context中
func ExtractClientIPMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		clientIPStr := getClientIP(c)

		if clientIPStr == "" {
			zapLogger.Warn("Unable to determine client IP", zap.String("path", c.Request.URL.Path))
			response.Forbidden(c, "Access denied: Unable to determine client IP")
			c.Abort()
			return
		}

		parsedIP := net.ParseIP(clientIPStr)
		if parsedIP == nil {
			zapLogger.Warn("Invalid client IP format", zap.String("client_ip_str", clientIPStr), zap.String("path", c.Request.URL.Path))
			response.Forbidden(c, "Access denied: Invalid client IP format")
			c.Abort()
			return
		}

		c.Set(ClientIPContextKey, clientIPStr)
		c.Set(ClientParsedIPContextKey, parsedIP)
		logger.Debug(ctx, "Client IP extracted and parsed", zap.String("client_ip", clientIPStr))
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
