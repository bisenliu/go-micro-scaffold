package response

import (
	"sync"
)

// EngineManager 引擎管理器
type EngineManager struct {
	defaultEngine Engine
	mu            sync.RWMutex
}

// NewEngineManager 创建新的引擎管理器
func NewEngineManager() *EngineManager {
	return &EngineManager{
		defaultEngine: NewResponseEngine(),
	}
}

// GetDefaultEngine 获取默认引擎
func (em *EngineManager) GetDefaultEngine() Engine {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.defaultEngine
}

// SetDefaultEngine 设置默认引擎
func (em *EngineManager) SetDefaultEngine(engine Engine) {
	if engine == nil {
		return
	}
	
	em.mu.Lock()
	defer em.mu.Unlock()
	em.defaultEngine = engine
}

// 全局引擎管理器实例
var globalEngineManager = NewEngineManager()

// GetGlobalEngineManager 获取全局引擎管理器
func GetGlobalEngineManager() *EngineManager {
	return globalEngineManager
}

// SetGlobalEngineManager 设置全局引擎管理器（主要用于测试）
func SetGlobalEngineManager(manager *EngineManager) {
	if manager != nil {
		globalEngineManager = manager
	}
}