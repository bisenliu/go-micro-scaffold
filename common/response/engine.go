package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseEngine 优化的响应引擎
// 完全重写以使用新的优化组件，提供更高性能和更简洁的API
type ResponseEngine struct {
	// 核心组件
	errorHandler   UnifiedErrorHandler // 统一错误处理器
	errorFactory   ErrorFactory        // 错误工厂
	contextManager *ContextManager     // 上下文管理器
	codeRegistry   CodeRegistry        // 代码注册表
	pool           *ResponsePool       // 对象池

}

// ResponseEngineConfig 响应引擎配置
type ResponseEngineConfig struct {
	EnableContextPool  bool // 启用上下文对象池
	EnableLazyMapping  bool // 启用延迟映射
	MaxContextPoolSize int  // 上下文池最大大小
	DefaultHTTPStatus  int  // 默认HTTP状态码
}

// DefaultResponseEngineConfig 默认配置
func DefaultResponseEngineConfig() *ResponseEngineConfig {
	return &ResponseEngineConfig{
		EnableContextPool:  true,
		EnableLazyMapping:  true,
		MaxContextPoolSize: 100,
		DefaultHTTPStatus:  http.StatusInternalServerError,
	}
}

// NewResponseEngine 创建新的优化响应引擎
func NewResponseEngine() *ResponseEngine {
	return NewResponseEngineWithConfig(DefaultResponseEngineConfig())
}

// NewResponseEngineWithConfig 使用配置创建响应引擎
func NewResponseEngineWithConfig(config *ResponseEngineConfig) *ResponseEngine {
	// 创建核心组件
	contextManager := NewContextManager()
	codeRegistry := NewCodeRegistry()
	errorFactory := NewErrorFactoryWithContextManager(contextManager)

	// 创建错误映射器（支持延迟初始化）
	// 默认使用延迟映射器以提高启动性能
	var errorMapper ErrorMapper = NewLazyErrorMapper()

	// 创建统一错误处理器
	errorHandler := NewUnifiedErrorHandlerWithDependencies(errorMapper, errorFactory)

	return &ResponseEngine{
		errorHandler:   errorHandler,
		errorFactory:   errorFactory,
		contextManager: contextManager,
		codeRegistry:   codeRegistry,
		pool:           NewResponsePool(),
	}
}

// NewResponseEngineWithDependencies 使用指定依赖创建响应引擎
func NewResponseEngineWithDependencies(
	errorHandler UnifiedErrorHandler,
	errorFactory ErrorFactory,
	contextManager *ContextManager,
	codeRegistry CodeRegistry,
	pool *ResponsePool,
) *ResponseEngine {
	return &ResponseEngine{
		errorHandler:   errorHandler,
		errorFactory:   errorFactory,
		contextManager: contextManager,
		codeRegistry:   codeRegistry,
		pool:           pool,
	}
}

// 全局响应引擎实例
var defaultEngine = NewResponseEngine()

// HandleSuccess 处理成功响应
// 优化：使用对象池和代码注册表提高性能
func (re *ResponseEngine) HandleSuccess(c *gin.Context, data any) {
	resp := re.pool.GetResponse()
	defer re.pool.PutResponse(resp)

	resp.Code = CodeSuccess
	resp.Message = re.codeRegistry.GetMessage(CodeSuccess)
	resp.Data = data

	c.JSON(http.StatusOK, resp)
}

// HandleSuccessWithPaging 处理分页成功响应
// 优化：使用对象池减少内存分配，使用代码注册表提高查找性能
func (re *ResponseEngine) HandleSuccessWithPaging(c *gin.Context, data any, page, pageSize int, total int64) {
	resp := re.pool.GetResponse()
	defer re.pool.PutResponse(resp)

	pageData := re.pool.GetPageData()
	defer re.pool.PutPageData(pageData)

	// 计算总页数，避免除零错误
	var totalPages int
	if pageSize <= 0 {
		totalPages = 0
	} else {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	pageData.Items = data
	pageData.Pagination = &Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	resp.Code = CodeSuccess
	resp.Message = re.codeRegistry.GetMessage(CodeSuccess)
	resp.Data = pageData

	c.JSON(http.StatusOK, resp)
}

// HandleError 处理错误响应
// 优化：使用统一错误处理器和对象池
func (re *ResponseEngine) HandleError(c *gin.Context, err error) {
	result := re.errorHandler.HandleError(err)
	re.sendErrorResponse(c, result)
}

// HandleErrorWithCode 使用指定业务码处理错误响应
func (re *ResponseEngine) HandleErrorWithCode(c *gin.Context, code int, message string) {
	result := re.errorHandler.HandleError(nil, WithCode(code), WithMessage(message))
	re.sendErrorResponse(c, result)
}

// HandleErrorWithData 处理带有额外数据的错误响应
func (re *ResponseEngine) HandleErrorWithData(c *gin.Context, err error, data any) {
	result := re.errorHandler.HandleError(err, WithData(data))
	re.sendErrorResponse(c, result)
}

// HandleErrorWithOptions 使用选项处理错误响应（新的统一API）
func (re *ResponseEngine) HandleErrorWithOptions(c *gin.Context, err error, options ...ErrorHandleOption) {
	result := re.errorHandler.HandleError(err, options...)
	re.sendErrorResponse(c, result)
}

// HandleErrorAdvanced 高级错误处理（新的统一API）
func (re *ResponseEngine) HandleErrorAdvanced(c *gin.Context, err error, code *int, message string, data any) {
	var options []ErrorHandleOption

	if code != nil {
		options = append(options, WithCode(*code))
	}
	if message != "" {
		options = append(options, WithMessage(message))
	}
	if data != nil {
		options = append(options, WithData(data))
	}

	result := re.errorHandler.HandleError(err, options...)
	re.sendErrorResponse(c, result)
}

// === 新的简洁统一API ===

// Handle 新的统一处理方法 - 自动判断成功或错误
func (re *ResponseEngine) Handle(c *gin.Context, data any, err error) {
	if err != nil {
		re.HandleError(c, err)
		return
	}
	re.HandleSuccess(c, data)
}

// HandleWith 新的统一处理方法 - 支持选项配置
func (re *ResponseEngine) HandleWith(c *gin.Context, data any, err error, options ...ErrorHandleOption) {
	if err != nil {
		re.HandleErrorWithOptions(c, err, options...)
		return
	}
	re.HandleSuccess(c, data)
}

// HandlePaging 新的统一分页处理方法
func (re *ResponseEngine) HandlePaging(c *gin.Context, data any, page, pageSize int, total int64, err error) {
	if err != nil {
		re.HandleError(c, err)
		return
	}
	re.HandleSuccessWithPaging(c, data, page, pageSize, total)
}

// CreateError 使用引擎的错误工厂创建错误
func (re *ResponseEngine) CreateError(errorType ErrorType, message string, cause ...error) *DomainError {
	return re.errorFactory.Create(errorType, message, cause...)
}

// CreateErrorWithContext 使用引擎的错误工厂创建带上下文的错误
func (re *ResponseEngine) CreateErrorWithContext(errorType ErrorType, message string, context map[string]any, cause ...error) *DomainError {
	return re.errorFactory.CreateWithContext(errorType, message, context, cause...)
}

// sendErrorResponse 发送错误响应
// 优化：使用对象池减少内存分配
func (re *ResponseEngine) sendErrorResponse(c *gin.Context, result *ErrorResult) {
	resp := re.pool.GetResponse()
	defer re.pool.PutResponse(resp)

	resp.Code = result.Code
	resp.Message = result.Message
	resp.Data = result.Data

	c.JSON(result.HTTPStatus, resp)
}

// === 引擎管理方法 ===

// GetErrorFactory 获取错误工厂
func (re *ResponseEngine) GetErrorFactory() ErrorFactory {
	return re.errorFactory
}

// GetContextManager 获取上下文管理器
func (re *ResponseEngine) GetContextManager() *ContextManager {
	return re.contextManager
}

// GetCodeRegistry 获取代码注册表
func (re *ResponseEngine) GetCodeRegistry() CodeRegistry {
	return re.codeRegistry
}

// GetPool 获取对象池
func (re *ResponseEngine) GetPool() *ResponsePool {
	return re.pool
}

// RegisterCode 注册新的业务码
func (re *ResponseEngine) RegisterCode(code int, info *CodeInfo) {
	re.codeRegistry.Register(code, info)
}

// RegisterCodes 批量注册业务码
func (re *ResponseEngine) RegisterCodes(codes map[int]*CodeInfo) {
	re.codeRegistry.RegisterBatch(codes)
}

// === 引擎组件访问方法 ===

// GetErrorHandler 获取统一错误处理器
func (re *ResponseEngine) GetErrorHandler() UnifiedErrorHandler {
	return re.errorHandler
}

// === 全局便捷函数，使用默认引擎 ===

// 全局API函数

// OK 成功响应
func OK(c *gin.Context, data any) {
	defaultEngine.HandleSuccess(c, data)
}

// OKWithPaging 分页成功响应
func OKWithPaging(c *gin.Context, data any, page, pageSize int, total int64) {
	defaultEngine.HandleSuccessWithPaging(c, data, page, pageSize, total)
}

// Fail 错误响应
func Fail(c *gin.Context, err error) {
	defaultEngine.HandleError(c, err)
}

// FailWithCode 使用指定业务码的错误响应
func FailWithCode(c *gin.Context, code int, message string) {
	defaultEngine.HandleErrorWithCode(c, code, message)
}

// FailWithData 带有额外数据的错误响应
func FailWithData(c *gin.Context, err error, data any) {
	defaultEngine.HandleErrorWithData(c, err, data)
}

// FailWithOptions 使用选项处理错误响应
func FailWithOptions(c *gin.Context, err error, options ...ErrorHandleOption) {
	defaultEngine.HandleErrorWithOptions(c, err, options...)
}

// FailAdvanced 高级错误处理
func FailAdvanced(c *gin.Context, err error, code *int, message string, data any) {
	defaultEngine.HandleErrorAdvanced(c, err, code, message, data)
}

// FailWith 简化的统一错误处理函数
func FailWith(c *gin.Context, err error, options ...ErrorHandleOption) {
	defaultEngine.HandleErrorWithOptions(c, err, options...)
}

// === 新的简洁统一API全局函数 ===

// Handle 新的统一处理函数 - 自动判断成功或错误
// 这是最简洁的API，根据err是否为nil自动选择成功或错误处理
func Handle(c *gin.Context, data any, err error) {
	defaultEngine.Handle(c, data, err)
}

// HandleWith 新的统一处理函数 - 支持选项配置
// 提供更灵活的配置选项，同时保持简洁性
func HandleWith(c *gin.Context, data any, err error, options ...ErrorHandleOption) {
	defaultEngine.HandleWith(c, data, err, options...)
}

// HandlePaging 新的统一分页处理函数
// 自动处理分页成功或错误响应
func HandlePaging(c *gin.Context, data any, page, pageSize int, total int64, err error) {
	defaultEngine.HandlePaging(c, data, page, pageSize, total, err)
}

// === 错误创建便捷函数 ===
// 注意：CreateError 和 CreateErrorWithContext 函数已在 factory.go 中定义
// 这里提供引擎级别的错误创建方法

// CreateEngineError 使用默认引擎的工厂创建错误
func CreateEngineError(errorType ErrorType, message string, cause ...error) *DomainError {
	return defaultEngine.CreateError(errorType, message, cause...)
}

// CreateEngineErrorWithContext 使用默认引擎的工厂创建带上下文的错误
func CreateEngineErrorWithContext(errorType ErrorType, message string, context map[string]any, cause ...error) *DomainError {
	return defaultEngine.CreateErrorWithContext(errorType, message, context, cause...)
}

// === 引擎管理全局函数 ===

// GetDefaultEngine 获取默认响应引擎
func GetDefaultEngine() *ResponseEngine {
	return defaultEngine
}

// SetDefaultEngine 设置默认响应引擎
func SetDefaultEngine(engine *ResponseEngine) {
	if engine != nil {
		defaultEngine = engine
	}
}

// RegisterGlobalCode 在默认引擎中注册业务码
func RegisterGlobalCode(code int, info *CodeInfo) {
	defaultEngine.RegisterCode(code, info)
}

// RegisterGlobalCodes 在默认引擎中批量注册业务码
func RegisterGlobalCodes(codes map[int]*CodeInfo) {
	defaultEngine.RegisterCodes(codes)
}
