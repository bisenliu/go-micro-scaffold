package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"common/response"
	domainerrors "services/internal/domain/shared/errors"
)

// DomainErrorAdapter 实现 response.CustomError 接口
type DomainErrorAdapter struct {
	domainErr *domainerrors.DomainError
}

func (d *DomainErrorAdapter) Error() string {
	return d.domainErr.Message
}

func (d *DomainErrorAdapter) ErrorType() string {
	baseErr := d.domainErr.Unwrap()
	
	switch {
	case errors.Is(baseErr, domainerrors.ErrNotFound):
		return "not_found"
	case errors.Is(baseErr, domainerrors.ErrValidationFailed):
		return "validation_failed"
	case errors.Is(baseErr, domainerrors.ErrAlreadyExists):
		return "already_exists"
	case errors.Is(baseErr, domainerrors.ErrUnauthorized):
		return "unauthorized"
	case errors.Is(baseErr, domainerrors.ErrForbidden):
		return "forbidden"
	case errors.Is(baseErr, domainerrors.ErrBusinessRuleViolation):
		return "business_rule_violation"
	case errors.Is(baseErr, domainerrors.ErrConcurrencyConflict):
		return "concurrency_conflict"
	case errors.Is(baseErr, domainerrors.ErrResourceLocked):
		return "resource_locked"
	case errors.Is(baseErr, domainerrors.ErrInvalidData):
		return "invalid_data"
	default:
		return "internal_server"
	}
}

// HandleErrorResponse 集中处理从应用层返回的领域错误，并将其转换为适当的HTTP响应。
func HandleErrorResponse(c *gin.Context, err error) {
	var domainErr *domainerrors.DomainError

	// 尝试将错误解包为自定义的 DomainError 类型
	if errors.As(err, &domainErr) {
		// 创建适配器并使用新的响应系统
		adapter := &DomainErrorAdapter{domainErr: domainErr}
		response.Fail(c, adapter)
		return
	}

	// 非 DomainError 类型，作为未预期的内部服务器错误处理
	response.FailWithCode(c, response.CodeInternalError, "服务器发生未知错误")
}