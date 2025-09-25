package validation

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

// BindMethod 定义通用的绑定方法类型
type BindMethod func(*gin.Context, interface{}) error

// URIBindAdapter URI参数绑定适配器
func URIBindAdapter(c *gin.Context, obj interface{}) error {
	return c.ShouldBindUri(obj)
}

// JSONBindAdapter JSON请求体绑定适配器
func JSONBindAdapter(c *gin.Context, obj interface{}) error {
	body, err := c.GetRawData()
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return NewValidationError(ErrEmptyRequestBody)
	}

	// 重置请求体，以便后续处理
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return c.ShouldBindJSON(obj)
}

// QueryBindAdapter 查询参数绑定适配器
func QueryBindAdapter(c *gin.Context, obj interface{}) error {
	return c.ShouldBindQuery(obj)
}

// FormBindAdapter 表单参数绑定适配器
func FormBindAdapter(c *gin.Context, obj interface{}) error {
	return c.ShouldBind(obj)
}
