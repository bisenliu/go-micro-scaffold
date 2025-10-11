package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"common/response"
	domainerrors "services/internal/domain/shared/errors"
)

// HandleErrorResponse 集中处理从应用层返回的领域错误，并将其转换为适当的HTTP响应。
func HandleErrorResponse(c *gin.Context, err error) {
	var domainErr *domainerrors.DomainError

	// 检查错误是否是我们自定义的 DomainError 类型
	if errors.As(err, &domainErr) {
		// 如果是，我们可以访问其内部字段
		baseErr := domainErr.Unwrap()
		message := domainErr.Message // <-- 这里是纯净、不冗余的消息

		if errors.Is(baseErr, domainerrors.ErrNotFound) {
			response.NotFound(c, message)
			return
		}
		if errors.Is(baseErr, domainerrors.ErrAlreadyExists) {
			response.BadRequest(c, message)
			return
		}
		if errors.Is(baseErr, domainerrors.ErrValidationFailed) {
			response.ValidationError(c, message, nil)
			return
		}
		if errors.Is(baseErr, domainerrors.ErrInvalidData) {
			response.BadRequest(c, message)
			return
		}
	}

	// 如果不是 DomainError，或者是不需要特殊处理的 DomainError，
	// 则作为未预期的内部服务器错误处理。
	// 完整的 err.Error() 信息会进入日志，但API只暴露通用消息。
	response.InternalServerError(c, "服务器发生未知错误")
}