package response

import "github.com/gin-gonic/gin"

// CodeRegistry 业务码注册表接口
type CodeRegistry interface {
	// GetCodeInfo 获取业务码信息
	GetCodeInfo(code int) (*CodeInfo, bool)
	// RegisterCode 注册业务码
	RegisterCode(code int, info *CodeInfo)
	// RegisterCodes 批量注册业务码
	RegisterCodes(codes map[int]*CodeInfo)
	// GetCodeMessage 获取业务码对应的消息
	GetCodeMessage(code int) string
	// GetHTTPStatus 获取业务码对应的HTTP状态码
	GetHTTPStatus(code int) int
}

// ResponsePool 响应对象池接口
type ResponsePool interface {
	// GetResponse 获取Response对象
	GetResponse() *Response
	// PutResponse 归还Response对象
	PutResponse(resp *Response)
	// GetPageData 获取PageData对象
	GetPageData() *PageData
	// PutPageData 归还PageData对象
	PutPageData(pageData *PageData)
}

// ErrorMapper 错误映射接口
type ErrorMapper interface {
	// GetErrorMapping 获取错误类型对应的映射信息
	GetErrorMapping(errorType ErrorType) (*ErrorMapping, bool)
	// RegisterErrorMapping 注册错误映射
	RegisterErrorMapping(errorType ErrorType, mapping *ErrorMapping)
}

// ResponseHandler 响应处理器接口
type ResponseHandler interface {
	// HandleSuccess 处理成功响应
	HandleSuccess(c *gin.Context, data any)
	// HandleSuccessWithPaging 处理分页成功响应
	HandleSuccessWithPaging(c *gin.Context, data any, page, pageSize int, total int64)
	// HandleError 处理错误响应
	HandleError(c *gin.Context, err error)
	// HandleErrorWithCode 使用指定业务码处理错误响应
	HandleErrorWithCode(c *gin.Context, code int, message string)
	// HandleErrorWithData 处理带有额外数据的错误响应
	HandleErrorWithData(c *gin.Context, err error, data any)
	// HandleErrorWithOptions 使用选项处理错误响应
	HandleErrorWithOptions(c *gin.Context, err error, options ...ErrorHandleOption)
}

// UnifiedAPI 统一API接口
type UnifiedAPI interface {
	// Handle 统一处理函数 - 自动判断成功或错误
	Handle(c *gin.Context, data any, err error)
	// HandleWith 统一处理函数 - 支持选项配置
	HandleWith(c *gin.Context, data any, err error, options ...ErrorHandleOption)
	// HandlePaging 统一分页处理函数
	HandlePaging(c *gin.Context, data any, page, pageSize int, total int64, err error)
}

// Engine 响应引擎接口
type Engine interface {
	CodeRegistry
	ResponsePool
	ErrorMapper
	ResponseHandler
	UnifiedAPI
	
	// GetContextManager 获取上下文管理器
	GetContextManager() *ContextManager
}