package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseEngine 响应引擎
type ResponseEngine struct {
	errorHandler *ErrorHandler
	pool         *ResponsePool
}

// NewResponseEngine 创建新的响应引擎
func NewResponseEngine() *ResponseEngine {
	return &ResponseEngine{
		errorHandler: NewErrorHandler(),
		pool:         NewResponsePool(),
	}
}

// 全局响应引擎实例
var defaultEngine = NewResponseEngine()

// HandleSuccess 处理成功响应
func (re *ResponseEngine) HandleSuccess(c *gin.Context, data interface{}) {
	resp := re.pool.GetResponse()
	defer re.pool.PutResponse(resp)

	resp.Code = CodeSuccess
	resp.Message = GetCodeMessage(CodeSuccess)
	resp.Data = data

	c.JSON(http.StatusOK, resp)
}

// HandleSuccessWithPaging 处理分页成功响应
func (re *ResponseEngine) HandleSuccessWithPaging(c *gin.Context, data interface{}, page, pageSize int, total int64) {
	resp := re.pool.GetResponse()
	defer re.pool.PutResponse(resp)

	pageData := re.pool.GetPageData()
	defer re.pool.PutPageData(pageData)

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if pageSize <= 0 {
		totalPages = 0
	}

	pageData.Items = data
	pageData.Pagination = &Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	resp.Code = CodeSuccess
	resp.Message = GetCodeMessage(CodeSuccess)
	resp.Data = pageData

	c.JSON(http.StatusOK, resp)
}

// HandleError 处理错误响应
func (re *ResponseEngine) HandleError(c *gin.Context, err error) {
	result := re.errorHandler.Handle(err)
	re.sendErrorResponse(c, result)
}

// HandleErrorWithCode 使用指定业务码处理错误响应
func (re *ResponseEngine) HandleErrorWithCode(c *gin.Context, code int, message string) {
	result := re.errorHandler.HandleWithCode(code, message)
	re.sendErrorResponse(c, result)
}

// HandleErrorWithData 处理带有额外数据的错误响应
func (re *ResponseEngine) HandleErrorWithData(c *gin.Context, err error, data interface{}) {
	result := re.errorHandler.HandleWithData(err, data)
	re.sendErrorResponse(c, result)
}

// sendErrorResponse 发送错误响应
func (re *ResponseEngine) sendErrorResponse(c *gin.Context, result *ErrorResult) {
	resp := re.pool.GetResponse()
	defer re.pool.PutResponse(resp)

	resp.Code = result.Code
	resp.Message = result.Message
	resp.Data = result.Data

	c.JSON(result.HTTPStatus, resp)
}

// 全局便捷函数，使用默认引擎

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	defaultEngine.HandleSuccess(c, data)
}

// OKWithPaging 分页成功响应
func OKWithPaging(c *gin.Context, data interface{}, page, pageSize int, total int64) {
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
func FailWithData(c *gin.Context, err error, data interface{}) {
	defaultEngine.HandleErrorWithData(c, err, data)
}