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
func NewNotFoundError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeNotFound, message, rootCause)
}

// NewValidationError 创建验证失败错误
func NewValidationError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeValidationFailed, message, rootCause)
}

// NewAlreadyExistsError 创建资源已存在错误
func NewAlreadyExistsError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeAlreadyExists, message, rootCause)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeUnauthorized, message, rootCause)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeForbidden, message, rootCause)
}

// NewBusinessRuleViolationError 创建业务规则违反错误
func NewBusinessRuleViolationError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeBusinessRuleViolation, message, rootCause)
}

// NewInvalidDataError 创建无效数据错误
func NewInvalidDataError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeInvalidData, message, rootCause)
}

// NewInternalServerError 创建内部服务器错误
func NewInternalServerError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeInternalServer, message, rootCause)
}

// NewDatabaseConnectionError 创建数据库连接错误
func NewDatabaseConnectionError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeDatabaseConnection, message, rootCause)
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeTimeout, message, rootCause)
}

// NewNetworkError 创建网络错误
func NewNetworkError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeNetworkError, message, rootCause)
}

// NewRecordNotFoundError 创建记录不存在错误
func NewRecordNotFoundError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeRecordNotFound, message, rootCause)
}

// NewDuplicateKeyError 创建重复键值错误
func NewDuplicateKeyError(message string, cause ...error) *DomainError {
	var rootCause error
	if len(cause) > 0 {
		rootCause = cause[0]
	}
	return NewDomainErrorWithCause(response.ErrorTypeDuplicateKey, message, rootCause)
}
