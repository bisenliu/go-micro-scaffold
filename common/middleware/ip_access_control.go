package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"common/pkg/contextutil"
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
		// 从Context中获取预解析的IP
		ipVal, exists := c.Get(contextutil.ClientParsedIPContextKey)
		if !exists {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		ip, ok := ipVal.(net.IP)
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
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
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

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
		// 从Context中获取预解析的IP
		ipVal, exists := c.Get(contextutil.ClientParsedIPContextKey)
		if !exists {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		ip, ok := ipVal.(net.IP)
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !isPrivateIP(ip, privateIPBlocks) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

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