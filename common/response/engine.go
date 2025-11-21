package response

import (
	"github.com/gin-gonic/gin"
)

// ResponseEngine 简化的响应引擎
// 集成所有功能组件，提供统一的响应处理能力
type ResponseEngine struct {
	// 内置组件
	codeRegistry   CodeRegistry
	responsePool   ResponsePool
	errorMapper    ErrorMapper
	responseHandler ResponseHandler
	contextManager *ContextManager
	errorFactory   ErrorFactory // 错误工厂依赖
}

// NewResponseEngine 创建新的响应引擎
func NewResponseEngine() *ResponseEngine {
	return NewResponseEngineWithFactory(NewErrorFactory())
}

// NewResponseEngineWithFactory 使用指定的ErrorFactory创建响应引擎
func NewResponseEngineWithFactory(factory ErrorFactory) *ResponseEngine {
	codeRegistry := newCodeRegistry()
	responsePool := newResponsePool()
	errorMapper := newErrorMapper()
	
	engine := &ResponseEngine{
		codeRegistry:   codeRegistry,
		responsePool:   responsePool,
		errorMapper:    errorMapper,
		contextManager: NewContextManager(),
		errorFactory:   factory,
	}
	
	// 创建响应处理器
	engine.responseHandler = newResponseHandler(codeRegistry, responsePool, errorMapper)

	return engine
}



// === 委托给组件的方法 ===

// GetCodeInfo 获取业务码信息
func (re *ResponseEngine) GetCodeInfo(code int) (*CodeInfo, bool) {
	return re.codeRegistry.GetCodeInfo(code)
}

// GetCodeMessage 获取业务码对应的消息
func (re *ResponseEngine) GetCodeMessage(code int) string {
	return re.codeRegistry.GetCodeMessage(code)
}

// GetHTTPStatus 获取业务码对应的HTTP状态码
func (re *ResponseEngine) GetHTTPStatus(code int) int {
	return re.codeRegistry.GetHTTPStatus(code)
}

// RegisterCode 注册新的业务码
func (re *ResponseEngine) RegisterCode(code int, info *CodeInfo) {
	re.codeRegistry.RegisterCode(code, info)
}

// RegisterCodes 批量注册业务码
func (re *ResponseEngine) RegisterCodes(codes map[int]*CodeInfo) {
	re.codeRegistry.RegisterCodes(codes)
}

// GetResponse 获取Response对象
func (re *ResponseEngine) GetResponse() *Response {
	return re.responsePool.GetResponse()
}

// PutResponse 归还Response对象
func (re *ResponseEngine) PutResponse(resp *Response) {
	re.responsePool.PutResponse(resp)
}

// GetPageData 获取PageData对象
func (re *ResponseEngine) GetPageData() *PageData {
	return re.responsePool.GetPageData()
}

// PutPageData 归还PageData对象
func (re *ResponseEngine) PutPageData(pageData *PageData) {
	re.responsePool.PutPageData(pageData)
}

// GetErrorMapping 获取错误类型对应的映射信息
func (re *ResponseEngine) GetErrorMapping(errorType ErrorType) (*ErrorMapping, bool) {
	return re.errorMapper.GetErrorMapping(errorType)
}

// RegisterErrorMapping 注册错误映射
func (re *ResponseEngine) RegisterErrorMapping(errorType ErrorType, mapping *ErrorMapping) {
	re.errorMapper.RegisterErrorMapping(errorType, mapping)
}

// HandleSuccess 处理成功响应
func (re *ResponseEngine) HandleSuccess(c *gin.Context, data any) {
	re.responseHandler.HandleSuccess(c, data)
}

// HandleSuccessWithPaging 处理分页成功响应
func (re *ResponseEngine) HandleSuccessWithPaging(c *gin.Context, data any, page, pageSize int, total int64) {
	re.responseHandler.HandleSuccessWithPaging(c, data, page, pageSize, total)
}

// HandleError 处理错误响应
func (re *ResponseEngine) HandleError(c *gin.Context, err error) {
	re.responseHandler.HandleError(c, err)
}

// HandleErrorWithCode 使用指定业务码处理错误响应
func (re *ResponseEngine) HandleErrorWithCode(c *gin.Context, code int, message string) {
	re.responseHandler.HandleErrorWithCode(c, code, message)
}

// HandleErrorWithData 处理带有额外数据的错误响应
func (re *ResponseEngine) HandleErrorWithData(c *gin.Context, err error, data any) {
	re.responseHandler.HandleErrorWithData(c, err, data)
}

// HandleErrorWithOptions 使用选项处理错误响应
func (re *ResponseEngine) HandleErrorWithOptions(c *gin.Context, err error, options ...ErrorHandleOption) {
	re.responseHandler.HandleErrorWithOptions(c, err, options...)
}

// === 统一API方法 ===

// Handle 统一处理方法 - 自动判断成功或错误
func (re *ResponseEngine) Handle(c *gin.Context, data any, err error) {
	if err != nil {
		re.HandleError(c, err)
		return
	}
	re.HandleSuccess(c, data)
}

// HandleWith 统一处理方法 - 支持选项配置
func (re *ResponseEngine) HandleWith(c *gin.Context, data any, err error, options ...ErrorHandleOption) {
	if err != nil {
		re.HandleErrorWithOptions(c, err, options...)
		return
	}
	re.HandleSuccess(c, data)
}

// HandlePaging 统一分页处理方法
func (re *ResponseEngine) HandlePaging(c *gin.Context, data any, page, pageSize int, total int64, err error) {
	if err != nil {
		re.HandleError(c, err)
		return
	}
	re.HandleSuccessWithPaging(c, data, page, pageSize, total)
}

// GetContextManager 获取上下文管理器
func (re *ResponseEngine) GetContextManager() *ContextManager {
	return re.contextManager
}

// === 全局API函数 ===

// Handle 统一处理函数 - 自动判断成功或错误
func Handle(c *gin.Context, data any, err error) {
	GetDefaultEngine().Handle(c, data, err)
}

// HandleWith 统一处理函数 - 支持选项配置
func HandleWith(c *gin.Context, data any, err error, options ...ErrorHandleOption) {
	GetDefaultEngine().HandleWith(c, data, err, options...)
}

// HandlePaging 统一分页处理函数
func HandlePaging(c *gin.Context, data any, page, pageSize int, total int64, err error) {
	GetDefaultEngine().HandlePaging(c, data, page, pageSize, total, err)
}

// === 引擎管理全局函数 ===

// GetDefaultEngine 获取默认响应引擎
func GetDefaultEngine() Engine {
	return GetGlobalEngineManager().GetDefaultEngine()
}

// SetDefaultEngine 设置默认响应引擎
func SetDefaultEngine(engine Engine) {
	GetGlobalEngineManager().SetDefaultEngine(engine)
}

// RegisterGlobalCode 在默认引擎中注册业务码
func RegisterGlobalCode(code int, info *CodeInfo) {
	GetDefaultEngine().RegisterCode(code, info)
}

// RegisterGlobalCodes 在默认引擎中批量注册业务码
func RegisterGlobalCodes(codes map[int]*CodeInfo) {
	GetDefaultEngine().RegisterCodes(codes)
}
