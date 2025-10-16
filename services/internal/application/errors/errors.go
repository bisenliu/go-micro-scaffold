package errors

import (
	"common/response"
	domainerrors "services/internal/domain/shared/errors"
)

// 应用层命令处理错误
var (
	ErrCommandValidation = domainerrors.NewDomainError(response.ErrorTypeCommandValidation, "命令验证失败")
	ErrCommandExecution  = domainerrors.NewDomainError(response.ErrorTypeCommandExecution, "命令执行失败")
	ErrQueryExecution    = domainerrors.NewDomainError(response.ErrorTypeQueryExecution, "查询执行失败")
)
