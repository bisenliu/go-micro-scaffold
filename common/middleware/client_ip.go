package middleware

import (
	"net"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/pkg/contextutil"
	"common/pkg/netutil"
	"common/response"
)

// ExtractClientIPMiddleware 提取客户端IP并解析为net.IP，存储在Context中
func ExtractClientIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		clientIPStr := netutil.GetClientIP(c)

		if clientIPStr == "" {
			logger.Warn(ctx, "Unable to determine client IP", zap.String("path", c.Request.URL.Path))
			response.Forbidden(c, "Access denied: Unable to determine client IP")
			c.Abort()
			return
		}

		parsedIP := net.ParseIP(clientIPStr)
		if parsedIP == nil {
			logger.Warn(ctx, "Invalid client IP format", zap.String("client_ip_str", clientIPStr), zap.String("path", c.Request.URL.Path))
			response.Forbidden(c, "Access denied: Invalid client IP format")
			c.Abort()
			return
		}

		c.Set(contextutil.ClientIPContextKey, clientIPStr)
		c.Set(contextutil.ClientParsedIPContextKey, parsedIP)
		logger.Debug(ctx, "Client IP extracted and parsed", zap.String("client_ip", clientIPStr))
		c.Next()
	}
}
