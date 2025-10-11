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

	// 尝试将错误解包为自定义的 DomainError 类型
	if errors.As(err, &domainErr) {
		// 获取 DomainError 包装的底层错误和用户友好的消息
		baseErr := domainErr.Unwrap()
		message := domainErr.Message

		switch {
		case errors.Is(baseErr, domainerrors.ErrNotFound):
			response.NotFound(c, message)
			return

		case errors.Is(baseErr, domainerrors.ErrValidationFailed):
			response.ValidationError(c, message, nil)
			return

		// 将所有导致 400 Bad Request 的错误合并处理
		case errors.Is(baseErr, domainerrors.ErrAlreadyExists),
			errors.Is(baseErr, domainerrors.ErrInvalidData):
			response.BadRequest(c, message)
			return
		}
	}

	// 非 DomainError 类型，或 DomainError 中包含的底层错误不是预期的错误类型，
	// 都作为未预期的内部服务器错误处理，并对外部隐藏详细错误信息。
	// 实际的 err.Error() 应该在更上层或中间件中记录到日志。
	response.InternalServerError(c, "服务器发生未知错误")
}
