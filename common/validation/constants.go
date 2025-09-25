package validation

// 错误消息常量
const (
	// ErrEmptyRequestBody 空请求体错误消息
	ErrEmptyRequestBody = "请求体参数不能为空"

	// ErrValidationFailed 验证失败错误消息
	ErrValidationFailed = "请求参数验证失败"

	// HTTP状态码
	HTTPBadRequest = 400
)

// 标签常量
const (
	LabelTag = "label"
	JSONTag  = "json"
)
