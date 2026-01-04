package errors

import (
	"common/response"
)

// 应用层命令处理错误
var (
	ErrCommandValidation = response.CreateError(response.ErrorTypeCommandValidation, "命令验证失败")
	ErrCommandExecution  = response.CreateError(response.ErrorTypeCommandExecution, "命令执行失败")
	ErrQueryExecution    = response.CreateError(response.ErrorTypeQueryExecution, "查询执行失败")
)
