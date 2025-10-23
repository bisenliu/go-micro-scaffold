package handler

import (
	"common/logger"
	"common/response"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandleWithLogging 使用新的统一API处理响应，并记录DomainError的上下文信息
func HandleWithLogging(c *gin.Context, data any, err error) {
	// 如果有错误且是DomainError，记录详细的上下文信息
	if err != nil {
		if domainErr, ok := err.(*response.DomainError); ok {
			logDomainError(c, domainErr)
		}
	}

	// 使用新的统一API处理响应
	response.Handle(c, data, err)
}

// HandlePagingWithLogging 使用新的统一分页API处理响应，并记录DomainError的上下文信息
func HandlePagingWithLogging(c *gin.Context, data any, page, pageSize int, total int64, err error) {
	// 如果有错误且是DomainError，记录详细的上下文信息
	if err != nil {
		if domainErr, ok := err.(*response.DomainError); ok {
			logDomainError(c, domainErr)
		}
	}

	// 使用新的统一分页API处理响应
	response.HandlePaging(c, data, page, pageSize, total, err)
}

// logDomainError 记录领域错误的详细信息
func logDomainError(c *gin.Context, domainErr *response.DomainError) {
	// 构建日志字段
	logFields := []zap.Field{
		zap.String("error_type", fmt.Sprintf("%d", domainErr.Type)),
		zap.String("message", domainErr.Message),
	}

	// 安全地添加请求信息
	if c.Request != nil && c.Request.URL != nil {
		logFields = append(logFields, zap.String("path", c.Request.URL.Path))
		logFields = append(logFields, zap.String("method", c.Request.Method))
	}

	// 添加上下文字段
	if domainErr.HasContext() {
		context := domainErr.GetContext()
		for key, value := range context {
			logFields = append(logFields, zap.Any(key, value))
		}
	}

	// 添加底层错误
	if domainErr.Cause != nil {
		logFields = append(logFields, zap.String("underlying_error", domainErr.Cause.Error()))
	}

	logger.Error(c.Request.Context(), "domain error occurred", logFields...)
}

// HandleError 处理错误响应
func HandleError(c *gin.Context, err error) {
	HandleWithLogging(c, nil, err)
}

// HandleErrorWithCode 使用指定错误码处理错误
func HandleErrorWithCode(c *gin.Context, code int, message string) {
	err := response.CreateError(response.ErrorTypeBusinessRuleViolation, message)
	response.HandleWith(c, nil, err, response.WithCode(code))
}

// HandleErrorWithData 处理带有额外数据的错误
func HandleErrorWithData(c *gin.Context, err error, data any) {
	response.HandleWith(c, nil, err, response.WithData(data))
}
