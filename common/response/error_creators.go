package response

// 便捷函数：使用统一的错误工厂

// CreateError 使用默认工厂创建领域错误
func CreateError(errorType ErrorType, message string, cause ...error) *DomainError {
	return GetDefaultErrorFactory().Create(errorType, message, cause...)
}

// CreateErrorWithContext 使用默认工厂创建带有上下文的领域错误
func CreateErrorWithContext(errorType ErrorType, message string, context map[string]any, cause ...error) *DomainError {
	return GetDefaultErrorFactory().CreateWithContext(errorType, message, context, cause...)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeNotFound, message, cause...)
}

// NewValidationError 创建验证失败错误
func NewValidationError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeValidationFailed, message, cause...)
}

// NewAlreadyExistsError 创建资源已存在错误
func NewAlreadyExistsError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeAlreadyExists, message, cause...)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeUnauthorized, message, cause...)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeForbidden, message, cause...)
}

// NewBusinessRuleViolationError 创建业务规则违反错误
func NewBusinessRuleViolationError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeBusinessRuleViolation, message, cause...)
}

// NewInvalidDataError 创建无效数据错误
func NewInvalidDataError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeInvalidData, message, cause...)
}

// NewInternalServerError 创建内部服务器错误
func NewInternalServerError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeInternalServer, message, cause...)
}

// NewDatabaseConnectionError 创建数据库连接错误
func NewDatabaseConnectionError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeDatabaseConnection, message, cause...)
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeTimeout, message, cause...)
}

// NewNetworkError 创建网络错误
func NewNetworkError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeNetworkError, message, cause...)
}

// NewRecordNotFoundError 创建记录不存在错误
func NewRecordNotFoundError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeRecordNotFound, message, cause...)
}

// NewDuplicateKeyError 创建重复键值错误
func NewDuplicateKeyError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeDuplicateKey, message, cause...)
}

// NewCommandValidationError 创建命令验证错误
func NewCommandValidationError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeCommandValidation, message, cause...)
}

// NewCommandExecutionError 创建命令执行错误
func NewCommandExecutionError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeCommandExecution, message, cause...)
}

// NewQueryExecutionError 创建查询执行错误
func NewQueryExecutionError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeQueryExecution, message, cause...)
}

// NewInvalidRequestError 创建无效请求错误
func NewInvalidRequestError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeInvalidRequest, message, cause...)
}

// NewConcurrencyConflictError 创建并发冲突错误
func NewConcurrencyConflictError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeConcurrencyConflict, message, cause...)
}

// NewResourceLockedError 创建资源锁定错误
func NewResourceLockedError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeResourceLocked, message, cause...)
}

// NewExternalServiceUnavailableError 创建外部服务不可用错误
func NewExternalServiceUnavailableError(message string, cause ...error) *DomainError {
	return CreateError(ErrorTypeExternalServiceUnavailable, message, cause...)
}