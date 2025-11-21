package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// responseHandler 响应处理器实现
type responseHandler struct {
	codeRegistry CodeRegistry
	responsePool ResponsePool
	errorMapper  ErrorMapper
}

// newResponseHandler 创建新的响应处理器
func newResponseHandler(codeRegistry CodeRegistry, responsePool ResponsePool, errorMapper ErrorMapper) ResponseHandler {
	return &responseHandler{
		codeRegistry: codeRegistry,
		responsePool: responsePool,
		errorMapper:  errorMapper,
	}
}

// HandleSuccess 处理成功响应
func (rh *responseHandler) HandleSuccess(c *gin.Context, data any) {
	resp := rh.responsePool.GetResponse()
	defer rh.responsePool.PutResponse(resp)

	resp.Code = CodeSuccess
	resp.Message = rh.codeRegistry.GetCodeMessage(CodeSuccess)
	resp.Data = data

	c.JSON(http.StatusOK, resp)
}

// HandleSuccessWithPaging 处理分页成功响应
func (rh *responseHandler) HandleSuccessWithPaging(c *gin.Context, data any, page, pageSize int, total int64) {
	resp := rh.responsePool.GetResponse()
	defer rh.responsePool.PutResponse(resp)

	pageData := rh.responsePool.GetPageData()
	defer rh.responsePool.PutPageData(pageData)

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
	resp.Message = rh.codeRegistry.GetCodeMessage(CodeSuccess)
	resp.Data = pageData

	c.JSON(http.StatusOK, resp)
}

// HandleError 处理错误响应
func (rh *responseHandler) HandleError(c *gin.Context, err error) {
	result := rh.handleError(err, nil)
	rh.sendErrorResponse(c, result)
}

// HandleErrorWithCode 使用指定业务码处理错误响应
func (rh *responseHandler) HandleErrorWithCode(c *gin.Context, code int, message string) {
	options := &ErrorHandleConfig{
		Code:    &code,
		Message: message,
	}
	result := rh.handleError(nil, options)
	rh.sendErrorResponse(c, result)
}

// HandleErrorWithData 处理带有额外数据的错误响应
func (rh *responseHandler) HandleErrorWithData(c *gin.Context, err error, data any) {
	options := &ErrorHandleConfig{
		Data: data,
	}
	result := rh.handleError(err, options)
	rh.sendErrorResponse(c, result)
}

// HandleErrorWithOptions 使用选项处理错误响应
func (rh *responseHandler) HandleErrorWithOptions(c *gin.Context, err error, options ...ErrorHandleOption) {
	config := &ErrorHandleConfig{}
	for _, option := range options {
		option(config)
	}
	result := rh.handleError(err, config)
	rh.sendErrorResponse(c, result)
}

// handleError 内置错误处理逻辑
func (rh *responseHandler) handleError(err error, config *ErrorHandleConfig) *ErrorResult {
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
		return rh.handleWithCode(*config.Code, config.Message, config.Data)
	}

	// 处理错误对象
	return rh.handleErrorObject(err, config)
}

// handleWithCode 使用指定业务码处理
func (rh *responseHandler) handleWithCode(code int, message string, data any) *ErrorResult {
	// 查找代码信息
	if codeInfo, exists := rh.codeRegistry.GetCodeInfo(code); exists {
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
func (rh *responseHandler) handleErrorObject(err error, config *ErrorHandleConfig) *ErrorResult {
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
		return rh.handleDomainError(domainErr, config)
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
func (rh *responseHandler) handleDomainError(err *DomainError, config *ErrorHandleConfig) *ErrorResult {
	// 查找错误映射
	if mapping, exists := rh.errorMapper.GetErrorMapping(err.Type); exists {
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

// sendErrorResponse 发送错误响应
func (rh *responseHandler) sendErrorResponse(c *gin.Context, result *ErrorResult) {
	resp := rh.responsePool.GetResponse()
	defer rh.responsePool.PutResponse(resp)

	resp.Code = result.Code
	resp.Message = result.Message
	resp.Data = result.Data

	c.JSON(result.HTTPStatus, resp)
}