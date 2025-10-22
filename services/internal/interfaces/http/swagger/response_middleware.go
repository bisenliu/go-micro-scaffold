package swagger

import (
	"bytes"
	"encoding/json"
	"net/http"

	"common/response"

	"github.com/gin-gonic/gin"
)

// SwaggerResponseMiddleware Swagger响应格式中间件
// 将系统内部响应格式转换为Swagger文档中定义的标准格式
type SwaggerResponseMiddleware struct {
	enabled bool
}

// NewSwaggerResponseMiddleware 创建Swagger响应中间件
func NewSwaggerResponseMiddleware(enabled bool) *SwaggerResponseMiddleware {
	return &SwaggerResponseMiddleware{
		enabled: enabled,
	}
}

// Middleware 中间件函数
func (m *SwaggerResponseMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.enabled {
			c.Next()
			return
		}

		// 只对API路径应用转换
		if !m.shouldApplyConversion(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 创建响应写入器包装器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// 继续处理请求
		c.Next()

		// 转换响应格式
		m.convertResponse(c, writer)
	}
}

// shouldApplyConversion 判断是否应该应用响应转换
func (m *SwaggerResponseMiddleware) shouldApplyConversion(path string) bool {
	// 只对API路径应用转换，排除Swagger UI路径
	if len(path) >= 4 && path[:4] == "/api" {
		return true
	}
	if len(path) >= 8 && path[:8] == "/swagger" {
		return false
	}
	if path == "/health" {
		return true
	}
	return false
}

// convertResponse 转换响应格式
func (m *SwaggerResponseMiddleware) convertResponse(c *gin.Context, writer *responseWriter) {
	statusCode := writer.Status()

	// 只转换错误响应（4xx, 5xx）
	if statusCode < 400 {
		// 成功响应保持原格式
		c.Writer = writer.ResponseWriter
		c.Writer.WriteHeader(statusCode)
		c.Writer.Write(writer.body.Bytes())
		return
	}

	// 解析原始响应
	var originalResp response.Response
	if err := json.Unmarshal(writer.body.Bytes(), &originalResp); err != nil {
		// 如果解析失败，返回原始响应
		c.Writer = writer.ResponseWriter
		c.Writer.WriteHeader(statusCode)
		c.Writer.Write(writer.body.Bytes())
		return
	}

	// 转换为Swagger格式
	swaggerResp := m.convertToSwaggerFormat(originalResp, statusCode)

	// 序列化并返回
	swaggerBytes, err := json.Marshal(swaggerResp)
	if err != nil {
		// 如果序列化失败，返回原始响应
		c.Writer = writer.ResponseWriter
		c.Writer.WriteHeader(statusCode)
		c.Writer.Write(writer.body.Bytes())
		return
	}

	// 设置正确的Content-Type和返回转换后的响应
	c.Writer = writer.ResponseWriter
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(statusCode)
	c.Writer.Write(swaggerBytes)
}

// convertToSwaggerFormat 将内部响应格式转换为Swagger格式
func (m *SwaggerResponseMiddleware) convertToSwaggerFormat(resp response.Response, statusCode int) interface{} {
	switch statusCode {
	case http.StatusBadRequest:
		if resp.Code == response.CodeValidation {
			return &ValidationErrorResponse{
				Error:   "Validation Failed",
				Message: resp.Message,
				Code:    statusCode,
				Details: ValidationErrorDetails{
					Fields: []FieldError{
						{
							Field:   "request",
							Message: resp.Message,
							Value:   "",
						},
					},
				},
			}
		}
		return &ErrorResponse{
			Error:   "Bad Request",
			Message: resp.Message,
			Code:    statusCode,
		}
	case http.StatusUnauthorized:
		return &UnauthorizedErrorResponse{
			Error:   "Unauthorized",
			Message: resp.Message,
			Code:    statusCode,
		}
	case http.StatusForbidden:
		return &ForbiddenErrorResponse{
			Error:   "Forbidden",
			Message: resp.Message,
			Code:    statusCode,
		}
	case http.StatusNotFound:
		return &NotFoundErrorResponse{
			Error:   "Not Found",
			Message: resp.Message,
			Code:    statusCode,
		}
	case http.StatusConflict:
		return &ConflictErrorResponse{
			Error:   "Conflict",
			Message: resp.Message,
			Code:    statusCode,
		}
	case http.StatusInternalServerError:
		return &InternalServerErrorResponse{
			Error:   "Internal Server Error",
			Message: resp.Message,
			Code:    statusCode,
		}
	case http.StatusServiceUnavailable:
		return &ServiceUnavailableErrorResponse{
			Error:   "Service Unavailable",
			Message: resp.Message,
			Code:    statusCode,
		}
	default:
		return &ErrorResponse{
			Error:   "Error",
			Message: resp.Message,
			Code:    statusCode,
		}
	}
}

// responseWriter 响应写入器包装器
type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

// Write 写入响应体
func (w *responseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

// WriteHeader 写入响应头
func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

// Status 获取状态码
func (w *responseWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

// SwaggerResponseConfig Swagger响应配置
type SwaggerResponseConfig struct {
	Enabled      bool     `yaml:"enabled"`       // 是否启用响应格式转换
	ApiPaths     []string `yaml:"api_paths"`     // 需要转换的API路径前缀
	ExcludePaths []string `yaml:"exclude_paths"` // 排除的路径前缀
}

// DefaultSwaggerResponseConfig 默认Swagger响应配置
func DefaultSwaggerResponseConfig() *SwaggerResponseConfig {
	return &SwaggerResponseConfig{
		Enabled:      true,
		ApiPaths:     []string{"/api", "/health"},
		ExcludePaths: []string{"/swagger", "/docs"},
	}
}

// CreateSwaggerResponseMiddleware 创建Swagger响应中间件
func CreateSwaggerResponseMiddleware(config *SwaggerResponseConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultSwaggerResponseConfig()
	}

	middleware := NewSwaggerResponseMiddleware(config.Enabled)
	return middleware.Middleware()
}

// 全局中间件实例
var defaultResponseMiddleware = NewSwaggerResponseMiddleware(true)

// GetDefaultSwaggerResponseMiddleware 获取默认Swagger响应中间件
func GetDefaultSwaggerResponseMiddleware() gin.HandlerFunc {
	return defaultResponseMiddleware.Middleware()
}

// SetSwaggerResponseMiddlewareEnabled 设置Swagger响应中间件启用状态
func SetSwaggerResponseMiddlewareEnabled(enabled bool) {
	defaultResponseMiddleware.enabled = enabled
}
