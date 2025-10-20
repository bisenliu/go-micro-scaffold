package response

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// ResponseEngine 简化的响应引擎
// 集成所有功能组件，提供统一的响应处理能力
type ResponseEngine struct {
	// 内置组件
	codeRegistry   map[int]*CodeInfo
	responsePool   *sync.Pool
	pageDataPool   *sync.Pool
	contextManager *ContextManager
	errorFactory   ErrorFactory // 错误工厂依赖
	errorMappings  map[ErrorType]*ErrorMapping
	mu             sync.RWMutex
}

// NewResponseEngine 创建新的响应引擎
func NewResponseEngine() *ResponseEngine {
	return NewResponseEngineWithFactory(NewErrorFactory())
}

// NewResponseEngineWithFactory 使用指定的ErrorFactory创建响应引擎
func NewResponseEngineWithFactory(factory ErrorFactory) *ResponseEngine {
	engine := &ResponseEngine{
		codeRegistry:   make(map[int]*CodeInfo),
		contextManager: NewContextManager(),
		errorFactory:   factory,
		errorMappings:  make(map[ErrorType]*ErrorMapping),
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

	// 初始化默认业务码
	engine.initDefaultCodes()
	// 初始化错误映射
	engine.initErrorMappings()

	return engine
}

// initDefaultCodes 初始化默认业务码
func (re *ResponseEngine) initDefaultCodes() {
	for code, info := range CodeInfoMap {
		re.codeRegistry[code] = info
	}
}

// initErrorMappings 初始化错误映射
func (re *ResponseEngine) initErrorMappings() {
	re.errorMappings[ErrorTypeNotFound] = &ErrorMapping{
		BusinessCode:   CodeNotFound,
		HTTPStatus:     http.StatusNotFound,
		DefaultMessage: "资源不存在",
	}
	re.errorMappings[ErrorTypeValidationFailed] = &ErrorMapping{
		BusinessCode:   CodeValidation,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "验证失败",
	}
	re.errorMappings[ErrorTypeAlreadyExists] = &ErrorMapping{
		BusinessCode:   CodeAlreadyExists,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "资源已存在",
	}
	re.errorMappings[ErrorTypeUnauthorized] = &ErrorMapping{
		BusinessCode:   CodeUnauthorized,
		HTTPStatus:     http.StatusUnauthorized,
		DefaultMessage: "未授权访问",
	}
	re.errorMappings[ErrorTypeForbidden] = &ErrorMapping{
		BusinessCode:   CodeForbidden,
		HTTPStatus:     http.StatusForbidden,
		DefaultMessage: "禁止访问",
	}
	re.errorMappings[ErrorTypeBusinessRuleViolation] = &ErrorMapping{
		BusinessCode:   CodeBusinessError,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "业务规则违反",
	}
	re.errorMappings[ErrorTypeConcurrencyConflict] = &ErrorMapping{
		BusinessCode:   CodeConflict,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "并发冲突",
	}
	re.errorMappings[ErrorTypeResourceLocked] = &ErrorMapping{
		BusinessCode:   CodeConflict,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "资源已锁定",
	}
	re.errorMappings[ErrorTypeInvalidData] = &ErrorMapping{
		BusinessCode:   CodeValidation,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "无效的数据",
	}
	re.errorMappings[ErrorTypeCommandValidation] = &ErrorMapping{
		BusinessCode:   CodeValidation,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "命令验证失败",
	}
	re.errorMappings[ErrorTypeCommandExecution] = &ErrorMapping{
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "命令执行失败",
	}
	re.errorMappings[ErrorTypeQueryExecution] = &ErrorMapping{
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "查询执行失败",
	}
	re.errorMappings[ErrorTypeInternalServer] = &ErrorMapping{
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "内部服务器错误",
	}
	re.errorMappings[ErrorTypeInvalidRequest] = &ErrorMapping{
		BusinessCode:   CodeBadRequest,
		HTTPStatus:     http.StatusBadRequest,
		DefaultMessage: "无效的请求",
	}
	re.errorMappings[ErrorTypeDatabaseConnection] = &ErrorMapping{
		BusinessCode:   CodeInternalError,
		HTTPStatus:     http.StatusInternalServerError,
		DefaultMessage: "数据库连接失败",
	}
	re.errorMappings[ErrorTypeRecordNotFound] = &ErrorMapping{
		BusinessCode:   CodeNotFound,
		HTTPStatus:     http.StatusNotFound,
		DefaultMessage: "记录不存在",
	}
	re.errorMappings[ErrorTypeDuplicateKey] = &ErrorMapping{
		BusinessCode:   CodeAlreadyExists,
		HTTPStatus:     http.StatusConflict,
		DefaultMessage: "重复键值",
	}
	re.errorMappings[ErrorTypeExternalServiceUnavailable] = &ErrorMapping{
		BusinessCode:   CodeThirdParty,
		HTTPStatus:     http.StatusBadGateway,
		DefaultMessage: "外部服务不可用",
	}
	re.errorMappings[ErrorTypeTimeout] = &ErrorMapping{
		BusinessCode:   CodeTimeout,
		HTTPStatus:     http.StatusRequestTimeout,
		DefaultMessage: "请求超时",
	}
	re.errorMappings[ErrorTypeNetworkError] = &ErrorMapping{
		BusinessCode:   CodeThirdParty,
		HTTPStatus:     http.StatusBadGateway,
		DefaultMessage: "网络错误",
	}
}

// === 内置对象池方法 ===

// getResponse 从池中获取Response对象
func (re *ResponseEngine) getResponse() *Response {
	resp := re.responsePool.Get().(*Response)
	// 重置对象状态
	resp.Code = 0
	resp.Message = ""
	resp.Data = nil
	return resp
}

// putResponse 将Response对象放回池中
func (re *ResponseEngine) putResponse(resp *Response) {
	if resp != nil {
		re.responsePool.Put(resp)
	}
}

// getPageData 从池中获取PageData对象
func (re *ResponseEngine) getPageData() *PageData {
	pageData := re.pageDataPool.Get().(*PageData)
	// 重置对象状态
	pageData.Items = nil
	pageData.Pagination = nil
	return pageData
}

// putPageData 将PageData对象放回池中
func (re *ResponseEngine) putPageData(pageData *PageData) {
	if pageData != nil {
		re.pageDataPool.Put(pageData)
	}
}

// === 内置业务码注册表方法 ===

// GetCodeInfo 获取业务码信息
func (re *ResponseEngine) GetCodeInfo(code int) (*CodeInfo, bool) {
	re.mu.RLock()
	defer re.mu.RUnlock()

	info, exists := re.codeRegistry[code]
	if !exists {
		return nil, false
	}

	// 返回副本以防止外部修改
	return &CodeInfo{
		Code:       info.Code,
		Message:    info.Message,
		HTTPStatus: info.HTTPStatus,
	}, true
}

// GetCodeMessage 获取业务码对应的消息
func (re *ResponseEngine) GetCodeMessage(code int) string {
	re.mu.RLock()
	defer re.mu.RUnlock()

	if info, exists := re.codeRegistry[code]; exists {
		return info.Message
	}
	return "未知错误"
}

// GetHTTPStatus 获取业务码对应的HTTP状态码
func (re *ResponseEngine) GetHTTPStatus(code int) int {
	re.mu.RLock()
	defer re.mu.RUnlock()

	if info, exists := re.codeRegistry[code]; exists {
		return info.HTTPStatus
	}
	return http.StatusInternalServerError
}

// RegisterCode 注册新的业务码
func (re *ResponseEngine) RegisterCode(code int, info *CodeInfo) {
	if info == nil {
		return
	}

	re.mu.Lock()
	defer re.mu.Unlock()

	re.codeRegistry[code] = &CodeInfo{
		Code:       info.Code,
		Message:    info.Message,
		HTTPStatus: info.HTTPStatus,
	}
}

// === 内置错误映射方法 ===

// getErrorMapping 获取错误类型对应的映射信息
func (re *ResponseEngine) getErrorMapping(errorType ErrorType) (*ErrorMapping, bool) {
	re.mu.RLock()
	defer re.mu.RUnlock()

	mapping, exists := re.errorMappings[errorType]
	return mapping, exists
}

// 全局响应引擎实例
var defaultEngine = NewResponseEngine()

// HandleSuccess 处理成功响应
func (re *ResponseEngine) HandleSuccess(c *gin.Context, data any) {
	resp := re.getResponse()
	defer re.putResponse(resp)

	resp.Code = CodeSuccess
	resp.Message = re.GetCodeMessage(CodeSuccess)
	resp.Data = data

	c.JSON(http.StatusOK, resp)
}

// HandleSuccessWithPaging 处理分页成功响应
func (re *ResponseEngine) HandleSuccessWithPaging(c *gin.Context, data any, page, pageSize int, total int64) {
	resp := re.getResponse()
	defer re.putResponse(resp)

	pageData := re.getPageData()
	defer re.putPageData(pageData)

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
	resp.Message = re.GetCodeMessage(CodeSuccess)
	resp.Data = pageData

	c.JSON(http.StatusOK, resp)
}

// HandleError 处理错误响应
func (re *ResponseEngine) HandleError(c *gin.Context, err error) {
	result := re.handleError(err, nil)
	re.sendErrorResponse(c, result)
}

// HandleErrorWithCode 使用指定业务码处理错误响应
func (re *ResponseEngine) HandleErrorWithCode(c *gin.Context, code int, message string) {
	options := &ErrorHandleConfig{
		Code:    &code,
		Message: message,
	}
	result := re.handleError(nil, options)
	re.sendErrorResponse(c, result)
}

// HandleErrorWithData 处理带有额外数据的错误响应
func (re *ResponseEngine) HandleErrorWithData(c *gin.Context, err error, data any) {
	options := &ErrorHandleConfig{
		Data: data,
	}
	result := re.handleError(err, options)
	re.sendErrorResponse(c, result)
}

// HandleErrorWithOptions 使用选项处理错误响应
func (re *ResponseEngine) HandleErrorWithOptions(c *gin.Context, err error, options ...ErrorHandleOption) {
	config := &ErrorHandleConfig{}
	for _, option := range options {
		option(config)
	}
	result := re.handleError(err, config)
	re.sendErrorResponse(c, result)
}

// handleError 内置错误处理逻辑
func (re *ResponseEngine) handleError(err error, config *ErrorHandleConfig) *ErrorResult {
	if config == nil {
		config = &ErrorHandleConfig{}
	}

	// 如果没有错误且没有指定业务码，返回成功
	if err == nil && config.Code == nil {
		return &ErrorResult{
			Code:       CodeSuccess,
			Message:    "操作成功",
			HTTPStatus: http.StatusOK,
			Data:       config.Data,
		}
	}

	// 如果指定了业务码，使用业务码处理
	if config.Code != nil {
		return re.handleWithCode(*config.Code, config.Message, config.Data)
	}

	// 处理错误对象
	return re.handleErrorObject(err, config)
}

// handleWithCode 使用指定业务码处理
func (re *ResponseEngine) handleWithCode(code int, message string, data any) *ErrorResult {
	// 查找代码信息
	if codeInfo, exists := re.GetCodeInfo(code); exists {
		finalMessage := message
		if finalMessage == "" {
			finalMessage = codeInfo.Message
		}
		return &ErrorResult{
			Code:       code,
			Message:    finalMessage,
			HTTPStatus: codeInfo.HTTPStatus,
			Data:       data,
		}
	}

	// 如果代码信息不存在，使用默认处理
	finalMessage := message
	if finalMessage == "" {
		finalMessage = "未知错误"
	}
	return &ErrorResult{
		Code:       code,
		Message:    finalMessage,
		HTTPStatus: http.StatusInternalServerError,
		Data:       data,
	}
}

// handleErrorObject 处理错误对象
func (re *ResponseEngine) handleErrorObject(err error, config *ErrorHandleConfig) *ErrorResult {
	if err == nil {
		return &ErrorResult{
			Code:       CodeInternalError,
			Message:    "内部错误：空错误对象",
			HTTPStatus: http.StatusInternalServerError,
			Data:       config.Data,
		}
	}

	// 检查是否为DomainError
	if domainErr, ok := err.(*DomainError); ok {
		return re.handleDomainError(domainErr, config)
	}

	// 默认处理为内部服务器错误
	message := config.Message
	if message == "" {
		message = err.Error()
	}

	return &ErrorResult{
		Code:       CodeInternalError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Data:       config.Data,
	}
}

// handleDomainError 处理领域错误
func (re *ResponseEngine) handleDomainError(err *DomainError, config *ErrorHandleConfig) *ErrorResult {
	// 查找错误映射
	if mapping, exists := re.getErrorMapping(err.Type); exists {
		message := config.Message
		if message == "" {
			message = err.Message
		}
		if message == "" {
			message = mapping.DefaultMessage
		}

		return &ErrorResult{
			Code:       mapping.BusinessCode,
			Message:    message,
			HTTPStatus: mapping.HTTPStatus,
			Data:       config.Data,
		}
	}

	// 如果没有找到映射，使用默认处理
	message := config.Message
	if message == "" {
		message = err.Message
	}
	if message == "" {
		message = "业务处理失败"
	}

	return &ErrorResult{
		Code:       CodeBusinessError,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Data:       config.Data,
	}
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

// sendErrorResponse 发送错误响应
func (re *ResponseEngine) sendErrorResponse(c *gin.Context, result *ErrorResult) {
	resp := re.getResponse()
	defer re.putResponse(resp)

	resp.Code = result.Code
	resp.Message = result.Message
	resp.Data = result.Data

	c.JSON(result.HTTPStatus, resp)
}

// === 引擎管理方法 ===

// GetContextManager 获取上下文管理器
func (re *ResponseEngine) GetContextManager() *ContextManager {
	return re.contextManager
}

// RegisterCodes 批量注册业务码
func (re *ResponseEngine) RegisterCodes(codes map[int]*CodeInfo) {
	re.mu.Lock()
	defer re.mu.Unlock()

	for code, info := range codes {
		if info != nil {
			re.codeRegistry[code] = &CodeInfo{
				Code:       info.Code,
				Message:    info.Message,
				HTTPStatus: info.HTTPStatus,
			}
		}
	}
}

// === 统一的全局API函数 ===

// Handle 统一处理函数 - 自动判断成功或错误
func Handle(c *gin.Context, data any, err error) {
	defaultEngine.Handle(c, data, err)
}

// HandleWith 统一处理函数 - 支持选项配置
func HandleWith(c *gin.Context, data any, err error, options ...ErrorHandleOption) {
	defaultEngine.HandleWith(c, data, err, options...)
}

// HandlePaging 统一分页处理函数
func HandlePaging(c *gin.Context, data any, page, pageSize int, total int64, err error) {
	defaultEngine.HandlePaging(c, data, page, pageSize, total, err)
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
