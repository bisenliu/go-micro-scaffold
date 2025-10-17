package response

import (
	"sync"
)

// CodeRegistry 代码注册表接口，提供统一的代码信息管理
type CodeRegistry interface {
	// GetInfo 获取业务码信息
	GetInfo(code int) (*CodeInfo, bool)
	// GetMessage 获取业务码对应的消息
	GetMessage(code int) string
	// GetHTTPStatus 获取业务码对应的HTTP状态码
	GetHTTPStatus(code int) int
	// Register 注册新的业务码信息
	Register(code int, info *CodeInfo)
	// RegisterBatch 批量注册业务码信息
	RegisterBatch(codes map[int]*CodeInfo)
	// GetAllCodes 获取所有已注册的业务码
	GetAllCodes() map[int]*CodeInfo
	// Exists 检查业务码是否存在
	Exists(code int) bool
}

// codeRegistry 代码注册表实现
type codeRegistry struct {
	codes map[int]*CodeInfo
	mu    sync.RWMutex
}

// NewCodeRegistry 创建新的代码注册表
func NewCodeRegistry() CodeRegistry {
	registry := &codeRegistry{
		codes: make(map[int]*CodeInfo),
	}

	// 初始化默认的业务码
	registry.initDefaultCodes()

	return registry
}

// initDefaultCodes 初始化默认的业务码信息
func (r *codeRegistry) initDefaultCodes() {
	// 直接使用现有的 CodeInfoMap 进行初始化
	for code, info := range CodeInfoMap {
		r.codes[code] = &CodeInfo{
			Code:       info.Code,
			Message:    info.Message,
			HTTPStatus: info.HTTPStatus,
			Category:   info.Category,
		}
	}
}

// GetInfo 获取业务码信息
func (r *codeRegistry) GetInfo(code int) (*CodeInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, exists := r.codes[code]
	if !exists {
		return nil, false
	}

	// 返回副本以防止外部修改
	return &CodeInfo{
		Code:       info.Code,
		Message:    info.Message,
		HTTPStatus: info.HTTPStatus,
		Category:   info.Category,
	}, true
}

// GetMessage 获取业务码对应的消息
func (r *codeRegistry) GetMessage(code int) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if info, exists := r.codes[code]; exists {
		return info.Message
	}
	return "未知错误"
}

// GetHTTPStatus 获取业务码对应的HTTP状态码
func (r *codeRegistry) GetHTTPStatus(code int) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if info, exists := r.codes[code]; exists {
		return info.HTTPStatus
	}
	return 500 // 默认返回内部服务器错误
}

// Register 注册新的业务码信息
func (r *codeRegistry) Register(code int, info *CodeInfo) {
	if info == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 创建副本以防止外部修改
	r.codes[code] = &CodeInfo{
		Code:       info.Code,
		Message:    info.Message,
		HTTPStatus: info.HTTPStatus,
		Category:   info.Category,
	}
}

// RegisterBatch 批量注册业务码信息
func (r *codeRegistry) RegisterBatch(codes map[int]*CodeInfo) {
	if len(codes) == 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for code, info := range codes {
		if info != nil {
			r.codes[code] = &CodeInfo{
				Code:       info.Code,
				Message:    info.Message,
				HTTPStatus: info.HTTPStatus,
				Category:   info.Category,
			}
		}
	}
}

// GetAllCodes 获取所有已注册的业务码
func (r *codeRegistry) GetAllCodes() map[int]*CodeInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[int]*CodeInfo, len(r.codes))
	for code, info := range r.codes {
		result[code] = &CodeInfo{
			Code:       info.Code,
			Message:    info.Message,
			HTTPStatus: info.HTTPStatus,
			Category:   info.Category,
		}
	}

	return result
}

// Exists 检查业务码是否存在
func (r *codeRegistry) Exists(code int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.codes[code]
	return exists
}

// 全局代码注册表实例
var defaultCodeRegistry = NewCodeRegistry()

// GetDefaultCodeRegistry 获取默认的代码注册表
func GetDefaultCodeRegistry() CodeRegistry {
	return defaultCodeRegistry
}

// 便捷函数：使用默认注册表的操作

// RegisterCode 使用默认注册表注册业务码
func RegisterCode(code int, info *CodeInfo) {
	defaultCodeRegistry.Register(code, info)
}

// RegisterCodes 使用默认注册表批量注册业务码
func RegisterCodes(codes map[int]*CodeInfo) {
	defaultCodeRegistry.RegisterBatch(codes)
}

// CodeExists 检查业务码是否在默认注册表中存在
func CodeExists(code int) bool {
	return defaultCodeRegistry.Exists(code)
}
