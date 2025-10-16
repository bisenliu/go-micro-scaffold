package errors

import (
	"common/response"
)

// DomainError 是领域错误的别名
type DomainError = response.DomainError

// NewDomainError 创建新的领域错误
func NewDomainError(errorType response.ErrorType, message string) *DomainError {
	return response.NewDomainError(errorType, message)
}

// NewDomainErrorWithCause 创建带有原因的领域错误
func NewDomainErrorWithCause(errorType response.ErrorType, message string, cause error) *DomainError {
	return response.NewDomainErrorWithCause(errorType, message, cause)
}

// 便捷函数：创建特定类型的领域错误

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeNotFound, message)
}

// NewValidationError 创建验证失败错误
func NewValidationError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeValidationFailed, message)
}

// NewAlreadyExistsError 创建资源已存在错误
func NewAlreadyExistsError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeAlreadyExists, message)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeUnauthorized, message)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeForbidden, message)
}

// NewBusinessRuleViolationError 创建业务规则违反错误
func NewBusinessRuleViolationError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeBusinessRuleViolation, message)
}

// NewConcurrencyConflictError 创建并发冲突错误
func NewConcurrencyConflictError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeConcurrencyConflict, message)
}

// NewResourceLockedError 创建资源已锁定错误
func NewResourceLockedError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeResourceLocked, message)
}

// NewInvalidDataError 创建无效数据错误
func NewInvalidDataError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeInvalidData, message)
}

// NewInternalServerError 创建内部服务器错误
func NewInternalServerError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeInternalServer, message)
}

// NewDatabaseConnectionError 创建数据库连接错误
func NewDatabaseConnectionError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeDatabaseConnection, message)
}

// NewRecordNotFoundError 创建记录不存在错误
func NewRecordNotFoundError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeRecordNotFound, message)
}

// NewDuplicateKeyError 创建重复键值错误
func NewDuplicateKeyError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeDuplicateKey, message)
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeTimeout, message)
}

// NewNetworkError 创建网络错误
func NewNetworkError(message string) *DomainError {
	return NewDomainError(response.ErrorTypeNetworkError, message)
}
