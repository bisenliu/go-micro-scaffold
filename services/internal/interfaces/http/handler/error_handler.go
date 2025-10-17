package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"common/logger"
	"common/response"
	domainerrors "services/internal/domain/shared/errors"
)

// HandleError 处理错误响应
func HandleError(c *gin.Context, err error) {
	if domainErr, ok := err.(*domainerrors.DomainError); ok {
		// 构建日志字段
		logFields := []zap.Field{
			zap.String("error_type", fmt.Sprintf("%d", domainErr.Type)),
			zap.String("message", domainErr.Message),
		}

		// 添加上下文字段
		if domainErr.Context != nil {
			for key, value := range domainErr.Context {
				logFields = append(logFields, zap.Any(key, value))
			}
		}

		// 添加底层错误
		if domainErr.BaseErr != nil {
			logFields = append(logFields, zap.String("underlying_error", domainErr.BaseErr.Error()))
		}

		logger.Error(c.Request.Context(), "request processing failed", logFields...)

		response.Fail(c, domainErr)
		return
	}

	response.FailWithCode(c, response.CodeInternalError, "服务器发生未知错误")
}

// HandleErrorWithCode 使用指定错误码处理错误
func HandleErrorWithCode(c *gin.Context, code int, message string) {
	response.FailWithCode(c, code, message)
}

// HandleErrorWithData 处理带有额外数据的错误
func HandleErrorWithData(c *gin.Context, err error, data interface{}) {
	response.FailWithData(c, err, data)
}
