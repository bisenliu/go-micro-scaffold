package response

import (
	"fmt"
)

// ErrorType 错误类型枚举
type ErrorType int

const (
	ErrorTypeNotFound ErrorType = iota
	ErrorTypeValidationFailed
	ErrorTypeAlreadyExists
	ErrorTypeUnauthorized
	ErrorTypeForbidden
	ErrorTypeBusinessRuleViolation
	ErrorTypeConcurrencyConflict
	ErrorTypeResourceLocked
	ErrorTypeInvalidData
	ErrorTypeCommandValidation
	ErrorTypeCommandExecution
	ErrorTypeQueryExecution
	ErrorTypeInternalServer
	ErrorTypeInvalidRequest
	ErrorTypeDatabaseConnection
	ErrorTypeRecordNotFound
	ErrorTypeDuplicateKey
	ErrorTypeExternalServiceUnavailable
	ErrorTypeTimeout
	ErrorTypeNetworkError
)

// ErrorMapping 错误映射结构
type ErrorMapping struct {
	BusinessCode   int
	HTTPStatus     int
	DefaultMessage string
}

// DomainError 领域错误结构
type DomainError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]any
}

// Error 实现error接口
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Cause.Error(), e.Message)
	}
	return e.Message
}

// Unwrap 返回被包装的底层错误
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// GetContext 安全地获取上下文信息
func (e *DomainError) GetContext() map[string]any {
	if e.Context == nil {
		return nil
	}
	// 返回上下文的副本以防止外部修改
	contextManager := GetDefaultContextManager()
	return contextManager.Copy(e.Context)
}

// GetContextValue 获取指定键的上下文值
func (e *DomainError) GetContextValue(key string) (any, bool) {
	if e.Context == nil {
		return nil, false
	}
	value, exists := e.Context[key]
	return value, exists
}

// HasContext 检查是否有上下文信息
func (e *DomainError) HasContext() bool {
	return len(e.Context) > 0
}

// WithContext 添加上下文信息
func (e *DomainError) WithContext(key string, value any) *DomainError {
	// 优化：如果当前上下文为空且新值为nil，直接返回自身避免不必要的分配
	if e.Context == nil && value == nil {
		return e
	}

	// 使用上下文管理器进行高效的上下文复制
	contextManager := GetDefaultContextManager()
	newContext := contextManager.CopyWithNew(e.Context, key, value)

	// 优化：如果新上下文与原上下文相同（都为nil），返回自身
	if newContext == nil && e.Context == nil {
		return e
	}

	// 创建新实例，避免修改原始错误
	return &DomainError{
		Type:    e.Type,
		Message: e.Message,
		Cause:   e.Cause,
		Context: newContext,
	}
}

// WithContextMap 批量添加上下文信息
func (e *DomainError) WithContextMap(contextMap map[string]any) *DomainError {
	// 优化：如果当前上下文和新上下文都为空，直接返回自身
	if e.Context == nil && len(contextMap) == 0 {
		return e
	}

	// 使用上下文管理器进行高效的上下文复制
	contextManager := GetDefaultContextManager()
	newContext := contextManager.CopyWithMap(e.Context, contextMap)

	// 优化：如果新上下文与原上下文相同（都为nil），返回自身
	if newContext == nil && e.Context == nil {
		return e
	}

	// 创建新实例，避免修改原始错误
	return &DomainError{
		Type:    e.Type,
		Message: e.Message,
		Cause:   e.Cause,
		Context: newContext,
	}
}
