package response

// Response 统一的响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
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
	Items any `json:"items"`
	*Pagination
}

// ErrorResult 错误处理结果
type ErrorResult struct {
	Code       int
	Message    string
	HTTPStatus int
	Data       any
}

// === 错误处理选项类型 ===

// ErrorHandleOption 错误处理选项
type ErrorHandleOption func(*ErrorHandleConfig)

// ErrorHandleConfig 错误处理配置
type ErrorHandleConfig struct {
	Code    *int   // 指定的业务码
	Message string // 自定义消息
	Data    any    // 额外数据
}

// WithCode 指定业务码选项
func WithCode(code int) ErrorHandleOption {
	return func(config *ErrorHandleConfig) {
		config.Code = &code
	}
}

// WithMessage 指定自定义消息选项
func WithMessage(message string) ErrorHandleOption {
	return func(config *ErrorHandleConfig) {
		config.Message = message
	}
}

// WithData 指定额外数据选项
func WithData(data any) ErrorHandleOption {
	return func(config *ErrorHandleConfig) {
		config.Data = data
	}
}
