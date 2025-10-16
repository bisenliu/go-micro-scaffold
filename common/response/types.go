package response

// Response 统一的响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Pagination 分页信息结构
type Pagination struct {
	Page       int   `json:"page"`        // 当前页
	PageSize   int   `json:"page_size"`   // 每页大小
	Total      int64 `json:"total"`       // 总数量
	TotalPages int   `json:"total_pages"` // 总页数
}

// PageData 分页数据结构
type PageData struct {
	Items interface{} `json:"items"`
	*Pagination
}

// ErrorResult 错误处理结果
type ErrorResult struct {
	Code       int
	Message    string
	HTTPStatus int
	Data       interface{}
}