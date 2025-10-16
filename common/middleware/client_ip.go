package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/pkg/contextutil"
	"common/pkg/netutil"
)

// ExtractClientIPMiddleware 提取客户端IP并解析为net.IP，存储在Context中
func ExtractClientIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIPStr := netutil.GetClientIP(c)

		if clientIPStr == "" {
			ctx := c.Request.Context()
			logger.Warn(ctx, "Unable to determine client IP", zap.String("path", c.Request.URL.Path))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		parsedIP := net.ParseIP(clientIPStr)
		if parsedIP == nil {
			ctx := c.Request.Context()
			logger.Warn(ctx, "Invalid client IP format", zap.String("client_ip_str", clientIPStr), zap.String("path", c.Request.URL.Path))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set(contextutil.ClientIPContextKey, clientIPStr)
		c.Set(contextutil.ClientParsedIPContextKey, parsedIP)
		ctx := c.Request.Context()
		logger.Debug(ctx, "Client IP extracted and parsed", zap.String("client_ip", clientIPStr))
		c.Next()
	}
}
