package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"common/config"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware(cfg config.ServerConfig) gin.HandlerFunc {
	// 如果未启用，返回一个空操作的中间件
	if !cfg.EnableCORS {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, Identification")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
