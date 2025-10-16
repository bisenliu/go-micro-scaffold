package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/pkg/contextutil"
	"common/response"
)

// ResponseWriter 包装gin的ResponseWriter以捕获响应体
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLogMiddleware 请求日志中间件，记录请求信息并包含用户ID和真实业务状态码
func RequestLogMiddleware() gin.HandlerFunc {
	return requestLoggerInternal(false)
}

// RequestLogWithDetailsMiddleware 详细请求日志中间件，包含请求体和响应体（谨慎使用，可能包含敏感信息）
func RequestLogWithDetailsMiddleware() gin.HandlerFunc {
	return requestLoggerInternal(true)
}

// requestLoggerInternal 是请求日志中间件的内部实现
func requestLoggerInternal(detailed bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()
		ctx := c.Request.Context()

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 优先从 context 获取真实 IP
		var clientIP string
		if ip, exists := c.Get(contextutil.ClientIPContextKey); exists {
			clientIP, _ = ip.(string)
		} else {
			clientIP = c.ClientIP() // Fallback
		}
		userAgent := c.Request.UserAgent()

		// 读取请求体（如果需要）
		var requestBody []byte
		if detailed && c.Request.Body != nil {
			var err error
			requestBody, err = c.GetRawData()
			if err == nil {
				// 重新设置请求体，以便后续处理器可以读取
				c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))
			}
		}

		// 包装ResponseWriter以捕获响应体
		blw := &ResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(startTime)
		httpStatus := c.Writer.Status()

		// 获取用户ID（如果存在）
		var userID any = "anonymous"
		if id, exists := c.Get(contextutil.UserIDKey); exists {
			userID = id
		}

		// 尝试解析响应体获取业务状态码
		responseBody := blw.body.Bytes()
		businessCode := httpStatus // 默认使用HTTP状态码
		var responseData response.Response
		if err := json.Unmarshal(responseBody, &responseData); err == nil {
			businessCode = responseData.Code
		}

		// 构建日志字段
		logFields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Int("http_status", httpStatus),
			zap.Int("business_code", businessCode),
			zap.Duration("latency", latency),
			zap.Any("user_id", userID),
		}

		// 如果是详细模式，添加请求和响应体
		if detailed {
			logFields = append(logFields,
				zap.ByteString("request_body", requestBody),
				zap.ByteString("response_body", responseBody),
			)
		}

		// 根据业务状态码选择日志级别和消息
		msg := "Request completed"
		if detailed {
			msg = "Detailed request completed"
		}

		if businessCode != response.CodeSuccess {
			logger.Error(ctx, msg+" with error", logFields...)
		} else {
			logger.Info(ctx, msg+" successfully", logFields...)
		}
	}
}