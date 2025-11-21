package response

import (
	"net/http"
	"sync"
)

// codeRegistry 业务码注册表实现
type codeRegistry struct {
	codes map[int]*CodeInfo
	mu    sync.RWMutex
}

// newCodeRegistry 创建新的业务码注册表
func newCodeRegistry() CodeRegistry {
	registry := &codeRegistry{
		codes: make(map[int]*CodeInfo),
	}
	
	// 初始化默认业务码
	registry.initDefaultCodes()
	
	return registry
}

// initDefaultCodes 初始化默认业务码
func (cr *codeRegistry) initDefaultCodes() {
	for code, info := range CodeInfoMap {
		cr.codes[code] = info
	}
}

// GetCodeInfo 获取业务码信息
func (cr *codeRegistry) GetCodeInfo(code int) (*CodeInfo, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	info, exists := cr.codes[code]
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

// GetCodeMessage 获取业务码对应的消息
func (cr *codeRegistry) GetCodeMessage(code int) string {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	if info, exists := cr.codes[code]; exists {
		return info.Message
	}
	return "未知错误"
}

// GetHTTPStatus 获取业务码对应的HTTP状态码
func (cr *codeRegistry) GetHTTPStatus(code int) int {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	if info, exists := cr.codes[code]; exists {
		return info.HTTPStatus
	}
	return http.StatusInternalServerError
}

// RegisterCode 注册新的业务码
func (cr *codeRegistry) RegisterCode(code int, info *CodeInfo) {
	if info == nil {
		return
	}

	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.codes[code] = &CodeInfo{
		Code:       info.Code,
		Message:    info.Message,
		HTTPStatus: info.HTTPStatus,
		Category:   info.Category,
	}
}

// RegisterCodes 批量注册业务码
func (cr *codeRegistry) RegisterCodes(codes map[int]*CodeInfo) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	for code, info := range codes {
		if info != nil {
			cr.codes[code] = &CodeInfo{
				Code:       info.Code,
				Message:    info.Message,
				HTTPStatus: info.HTTPStatus,
				Category:   info.Category,
			}
		}
	}
}