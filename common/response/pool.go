package response

import "sync"

// ResponsePool 响应对象池
type ResponsePool struct {
	responsePool sync.Pool
	pageDataPool sync.Pool
}

// NewResponsePool 创建新的响应对象池
func NewResponsePool() *ResponsePool {
	return &ResponsePool{
		responsePool: sync.Pool{
			New: func() any {
				return &Response{}
			},
		},
		pageDataPool: sync.Pool{
			New: func() any {
				return &PageData{}
			},
		},
	}
}

// GetResponse 从池中获取Response对象
func (rp *ResponsePool) GetResponse() *Response {
	resp := rp.responsePool.Get().(*Response)
	// 重置对象状态
	resp.Code = 0
	resp.Message = ""
	resp.Data = nil
	return resp
}

// PutResponse 将Response对象放回池中
func (rp *ResponsePool) PutResponse(resp *Response) {
	if resp != nil {
		rp.responsePool.Put(resp)
	}
}

// GetPageData 从池中获取PageData对象
func (rp *ResponsePool) GetPageData() *PageData {
	pageData := rp.pageDataPool.Get().(*PageData)
	// 重置对象状态
	pageData.Items = nil
	pageData.Pagination = nil
	return pageData
}

// PutPageData 将PageData对象放回池中
func (rp *ResponsePool) PutPageData(pageData *PageData) {
	if pageData != nil {
		rp.pageDataPool.Put(pageData)
	}
}
