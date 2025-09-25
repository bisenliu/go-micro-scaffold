package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
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

// ResponseData 通用响应数据结构
type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// RequestLogMiddleware 请求日志中间件，记录请求信息并包含用户ID和真实业务状态码
func RequestLogMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()
		ctx := c.Request.Context()

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

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
		var userID interface{}
		var userIDExists bool
		if userID, userIDExists = c.Get("userID"); !userIDExists {
			userID = "anonymous"
		}

		// 尝试解析响应体获取业务状态码
		businessCode := httpStatus // 默认使用HTTP状态码
		var responseData ResponseData
		if err := json.Unmarshal(blw.body.Bytes(), &responseData); err == nil {
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

		// 根据业务状态码选择日志级别
		if businessCode != 200 && businessCode != 0 { // 0表示未解析到业务状态码
			logger.Error(zapLogger, ctx, "Request completed with error", logFields...)
		} else {
			logger.Info(zapLogger, ctx, "Request completed successfully", logFields...)
		}
	}
}

// RequestLogWithDetailsMiddleware 详细请求日志中间件，包含请求体和响应体（谨慎使用，可能包含敏感信息）
func RequestLogWithDetailsMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()
		ctx := c.Request.Context()

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 读取请求体（如果需要）
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = c.GetRawData()
			// 重新设置请求体，以便后续处理器可以读取
			c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))
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
		var userID interface{}
		var userIDExists bool
		if userID, userIDExists = c.Get("userID"); !userIDExists {
			userID = "anonymous"
		}

		// 尝试解析响应体获取业务状态码
		businessCode := httpStatus
		var responseData ResponseData
		responseBody := blw.body.Bytes()
		if err := json.Unmarshal(responseBody, &responseData); err == nil {
			businessCode = responseData.Code
		}

		// 构建详细日志字段
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
			zap.ByteString("request_body", requestBody),
			zap.ByteString("response_body", responseBody),
		}

		// 根据业务状态码选择日志级别
		if businessCode != 200 && businessCode != 0 {
			logger.Error(zapLogger, ctx, "Detailed request completed with error", logFields...)
		} else {
			logger.Info(zapLogger, ctx, "Detailed request completed successfully", logFields...)
		}
	}
}

// ErrorLogMiddleware 错误日志中间件
func ErrorLogMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		ctx := c.Request.Context()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error(zapLogger, ctx, "Request error",
					zap.Error(err.Err),
					zap.Uint("type", uint(err.Type)),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method))
			}
		}
	}
}
