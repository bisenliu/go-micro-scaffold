package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"common/response"
	domainerrors "services/internal/domain/shared/errors"
)

// HandleErrorResponse 处理错误响应
func HandleErrorResponse(c *gin.Context, err error) {
	var domainErr *domainerrors.DomainError

	// 检查是否为领域错误
	if errors.As(err, &domainErr) {
		response.Fail(c, domainErr)
		return
	}

	// 非领域错误，作为内部服务器错误处理
	response.FailWithCode(c, response.CodeInternalError, "服务器发生未知错误")
}