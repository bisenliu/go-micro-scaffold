package errors

import "errors"

// 应用层命令处理错误
var (
	ErrCommandValidation = errors.New("命令验证失败")
	ErrCommandExecution  = errors.New("命令执行失败")
	ErrQueryExecution    = errors.New("查询执行失败")
)

// 应用层业务流程错误
var (
	ErrBusinessRuleViolation = errors.New("业务规则违反")
	ErrConcurrencyConflict   = errors.New("并发冲突")
	ErrResourceLocked        = errors.New("资源已锁定")
)

// 应用层验证错误
var (
	ErrInvalidInput     = errors.New("无效输入")
	ErrMissingRequired  = errors.New("必填字段缺失")
	ErrValidationFailed = errors.New("验证失败")
)
