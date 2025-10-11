package errors

import "errors"

// 应用层命令处理错误
var (
	ErrCommandValidation = errors.New("命令验证失败")
	ErrCommandExecution  = errors.New("命令执行失败")
	ErrQueryExecution    = errors.New("查询执行失败")
)
