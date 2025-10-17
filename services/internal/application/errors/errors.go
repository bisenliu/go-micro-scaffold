package errors

import (
	"common/response"
)

// 应用层命令处理错误
var (
	ErrCommandValidation = response.NewDomainError(response.ErrorTypeCommandValidation, "命令验证失败")
	ErrCommandExecution  = response.NewDomainError(response.ErrorTypeCommandExecution, "命令执行失败")
	ErrQueryExecution    = response.NewDomainError(response.ErrorTypeQueryExecution, "查询执行失败")
)
