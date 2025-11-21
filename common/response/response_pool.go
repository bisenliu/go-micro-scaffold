package response

import (
	"sync"
)

// responsePool 响应对象池实现
type responsePool struct {
	responsePool *sync.Pool
	pageDataPool *sync.Pool
}

// newResponsePool 创建新的响应对象池
func newResponsePool() ResponsePool {
	return &responsePool{
		responsePool: &sync.Pool{
			New: func() any {
				return &Response{}
			},
		},
		pageDataPool: &sync.Pool{
			New: func() any {
				return &PageData{}
			},
		},
	}
}

// GetResponse 从池中获取Response对象
func (rp *responsePool) GetResponse() *Response {
	resp := rp.responsePool.Get().(*Response)
	// 重置对象状态
	resp.Code = 0
	resp.Message = ""
	resp.Data = nil
	return resp
}

// PutResponse 将Response对象放回池中
func (rp *responsePool) PutResponse(resp *Response) {
	if resp != nil {
		// 清理对象状态，避免内存泄漏
		resp.Code = 0
		resp.Message = ""
		resp.Data = nil
		rp.responsePool.Put(resp)
	}
}

// GetPageData 从池中获取PageData对象
func (rp *responsePool) GetPageData() *PageData {
	pageData := rp.pageDataPool.Get().(*PageData)
	// 重置对象状态
	pageData.Items = nil
	pageData.Pagination = nil
	return pageData
}

// PutPageData 将PageData对象放回池中
func (rp *responsePool) PutPageData(pageData *PageData) {
	if pageData != nil {
		// 清理对象状态，避免内存泄漏
		pageData.Items = nil
		pageData.Pagination = nil
		rp.pageDataPool.Put(pageData)
	}
}