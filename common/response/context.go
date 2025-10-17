package response

import (
	"maps"
	"sync"
)

// ContextManager 上下文管理器，提供高效的上下文操作
type ContextManager struct {
	pool sync.Pool
}

// NewContextManager 创建新的上下文管理器
func NewContextManager() *ContextManager {
	return &ContextManager{
		pool: sync.Pool{
			New: func() any {
				return make(map[string]any, 4) // 预分配4个元素的容量
			},
		},
	}
}

// GetContext 从对象池获取一个空的上下文map
func (cm *ContextManager) GetContext() map[string]any {
	ctx := cm.pool.Get().(map[string]any)
	// 清空map，但保留容量
	for k := range ctx {
		delete(ctx, k)
	}
	return ctx
}

// ReleaseContext 将上下文map归还到对象池
func (cm *ContextManager) ReleaseContext(ctx map[string]any) {
	if ctx != nil && len(ctx) <= 16 { // 只回收小容量的map，避免内存浪费
		cm.pool.Put(ctx)
	}
}

// CopyWithNew 复制现有上下文并添加新的键值对
// 如果existing为nil且value为nil，返回nil避免不必要的分配
func (cm *ContextManager) CopyWithNew(existing map[string]any, key string, value any) map[string]any {
	// 优化：如果现有上下文为空且新值也为nil，直接返回nil
	if existing == nil && value == nil {
		return nil
	}

	// 优化：如果现有上下文为空且只添加一个值，创建最小map
	if existing == nil {
		if value == nil {
			return nil
		}
		newCtx := cm.GetContext()
		newCtx[key] = value
		return newCtx
	}

	// 复制现有上下文
	newCtx := cm.GetContext()
	maps.Copy(newCtx, existing)

	// 添加或更新新值
	if value != nil {
		newCtx[key] = value
	} else {
		delete(newCtx, key) // 如果value为nil，删除该键
	}

	return newCtx
}

// CopyWithMap 复制现有上下文并批量添加新的键值对
func (cm *ContextManager) CopyWithMap(existing, additional map[string]any) map[string]any {
	// 优化：如果两个map都为空，返回nil
	if existing == nil && additional == nil {
		return nil
	}

	// 优化：如果现有上下文为空，直接复制additional
	if existing == nil {
		if len(additional) == 0 {
			return nil
		}
		newCtx := cm.GetContext()
		maps.Copy(newCtx, additional)
		return newCtx
	}

	// 优化：如果additional为空，直接复制existing
	if len(additional) == 0 {
		newCtx := cm.GetContext()
		maps.Copy(newCtx, existing)
		return newCtx
	}

	// 复制现有上下文并添加新的键值对
	newCtx := cm.GetContext()
	maps.Copy(newCtx, existing)
	maps.Copy(newCtx, additional)

	return newCtx
}

// Copy 创建上下文的副本
func (cm *ContextManager) Copy(ctx map[string]any) map[string]any {
	if len(ctx) == 0 {
		return nil
	}

	newCtx := cm.GetContext()
	maps.Copy(newCtx, ctx)
	return newCtx
}

// IsEmpty 检查上下文是否为空
func (cm *ContextManager) IsEmpty(ctx map[string]any) bool {
	return len(ctx) == 0
}

// 全局上下文管理器实例
var defaultContextManager = NewContextManager()

// GetDefaultContextManager 获取默认的上下文管理器
func GetDefaultContextManager() *ContextManager {
	return defaultContextManager
}
