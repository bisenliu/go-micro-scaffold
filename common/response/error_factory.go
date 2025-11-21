package response

import "sync"

// ErrorFactory 统一的错误工厂接口
type ErrorFactory interface {
	// Create 创建领域错误
	Create(errorType ErrorType, message string, cause ...error) *DomainError
	// CreateWithContext 创建带有上下文的领域错误
	CreateWithContext(errorType ErrorType, message string, context map[string]any, cause ...error) *DomainError
}

// errorFactory 错误工厂实现
type errorFactory struct {
	contextManager *ContextManager
	errorPool      *sync.Pool
}

// NewErrorFactory 创建新的错误工厂
func NewErrorFactory() ErrorFactory {
	return &errorFactory{
		contextManager: GetDefaultContextManager(),
		errorPool: &sync.Pool{
			New: func() any {
				return &DomainError{}
			},
		},
	}
}

// NewErrorFactoryWithContextManager 使用指定的上下文管理器创建错误工厂
func NewErrorFactoryWithContextManager(cm *ContextManager) ErrorFactory {
	return &errorFactory{
		contextManager: cm,
		errorPool: &sync.Pool{
			New: func() any {
				return &DomainError{}
			},
		},
	}
}

// Create 创建领域错误
func (f *errorFactory) Create(errorType ErrorType, message string, cause ...error) *DomainError {
	err := f.errorPool.Get().(*DomainError)
	
	// 重置错误对象状态
	err.Type = errorType
	err.Message = message
	err.Context = nil // 延迟分配，只在需要时创建
	
	// 处理 cause 参数
	if len(cause) > 0 {
		err.Cause = cause[0]
	} else {
		err.Cause = nil
	}

	return err
}

// CreateWithContext 创建带有上下文的领域错误
func (f *errorFactory) CreateWithContext(errorType ErrorType, message string, context map[string]any, cause ...error) *DomainError {
	err := f.errorPool.Get().(*DomainError)
	
	// 重置错误对象状态
	err.Type = errorType
	err.Message = message
	
	// 处理 cause 参数
	if len(cause) > 0 {
		err.Cause = cause[0]
	} else {
		err.Cause = nil
	}

	// 使用上下文管理器复制上下文，避免不必要的分配
	if len(context) > 0 {
		err.Context = f.contextManager.Copy(context)
	} else {
		err.Context = nil
	}

	return err
}

// FactoryManager 工厂管理器
type FactoryManager struct {
	defaultFactory ErrorFactory
	mu             sync.RWMutex
}

// NewFactoryManager 创建新的工厂管理器
func NewFactoryManager() *FactoryManager {
	return &FactoryManager{
		defaultFactory: NewErrorFactory(),
	}
}

// GetDefaultFactory 获取默认工厂
func (fm *FactoryManager) GetDefaultFactory() ErrorFactory {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return fm.defaultFactory
}

// SetDefaultFactory 设置默认工厂
func (fm *FactoryManager) SetDefaultFactory(factory ErrorFactory) {
	if factory == nil {
		return
	}
	
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.defaultFactory = factory
}

// 全局工厂管理器实例
var globalFactoryManager = NewFactoryManager()

// GetGlobalFactoryManager 获取全局工厂管理器
func GetGlobalFactoryManager() *FactoryManager {
	return globalFactoryManager
}

// SetGlobalFactoryManager 设置全局工厂管理器（主要用于测试）
func SetGlobalFactoryManager(manager *FactoryManager) {
	if manager != nil {
		globalFactoryManager = manager
	}
}

// GetDefaultErrorFactory 获取默认的错误工厂
func GetDefaultErrorFactory() ErrorFactory {
	return GetGlobalFactoryManager().GetDefaultFactory()
}
