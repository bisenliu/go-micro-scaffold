package swagger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"common/response"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// ValidationHelper Swagger验证辅助工具
type ValidationHelper struct {
	adapter *ResponseAdapter
}

// NewValidationHelper 创建验证辅助工具
func NewValidationHelper() *ValidationHelper {
	return &ValidationHelper{
		adapter: NewResponseAdapter(),
	}
}

// HandleValidationError 处理验证错误并返回Swagger格式响应
func (h *ValidationHelper) HandleValidationError(c *gin.Context, params interface{}, err error, trans ut.Translator) {
	switch err := err.(type) {
	case validator.ValidationErrors:
		fieldErrors := h.convertValidatorErrors(err, trans, params)
		validationErr := CreateValidationErrorWithFields("Request validation failed", fieldErrors)
		h.adapter.AdaptErrorResponse(c, validationErr)
	case *json.UnmarshalTypeError:
		fieldError := h.convertUnmarshalTypeError(err, params)
		validationErr := CreateValidationErrorWithFields("Request validation failed", []FieldError{fieldError})
		h.adapter.AdaptErrorResponse(c, validationErr)
	default:
		fieldError := FieldError{
			Field:   "request",
			Message: err.Error(),
			Value:   "",
		}
		validationErr := CreateValidationErrorWithFields("Request validation failed", []FieldError{fieldError})
		h.adapter.AdaptErrorResponse(c, validationErr)
	}
}

// convertValidatorErrors 转换验证器错误为Swagger字段错误格式
func (h *ValidationHelper) convertValidatorErrors(validationErrors validator.ValidationErrors, trans ut.Translator, params interface{}) []FieldError {
	var fieldErrors []FieldError

	for _, err := range validationErrors {
		fieldError := FieldError{
			Field:   h.getJSONFieldName(params, err.Field()),
			Message: err.Translate(trans),
			Value:   fmt.Sprintf("%v", err.Value()),
		}
		fieldErrors = append(fieldErrors, fieldError)
	}

	return fieldErrors
}

// convertUnmarshalTypeError 转换JSON类型错误为Swagger字段错误格式
func (h *ValidationHelper) convertUnmarshalTypeError(err *json.UnmarshalTypeError, params interface{}) FieldError {
	fieldName := h.getJSONFieldName(params, err.Field)
	message := h.buildTypeErrorMessage(fieldName, err.Type.String())

	return FieldError{
		Field:   fieldName,
		Message: message,
		Value:   fmt.Sprintf("%v", err.Value),
	}
}

// getJSONFieldName 获取字段的JSON名称
func (h *ValidationHelper) getJSONFieldName(params interface{}, fieldName string) string {
	t := reflect.TypeOf(params)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 尝试通过字段名查找
	if field, found := t.FieldByName(fieldName); found {
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			jsonName := strings.SplitN(jsonTag, ",", 2)[0]
			if jsonName != "" && jsonName != "-" {
				return jsonName
			}
		}
		return fieldName
	}

	// 如果直接查找失败，尝试通过JSON标签查找
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			jsonName := strings.SplitN(jsonTag, ",", 2)[0]
			if jsonName == fieldName {
				return jsonName
			}
		}
	}

	return fieldName
}

// buildTypeErrorMessage 构建类型错误消息
func (h *ValidationHelper) buildTypeErrorMessage(fieldName, fieldType string) string {
	typeMessages := map[string]string{
		"int":     "整数",
		"int8":    "整数",
		"int16":   "整数",
		"int32":   "整数",
		"int64":   "整数",
		"uint":    "正整数",
		"uint8":   "正整数",
		"uint16":  "正整数",
		"uint32":  "正整数",
		"uint64":  "正整数",
		"float32": "浮点数",
		"float64": "浮点数",
		"bool":    "布尔值",
		"string":  "字符串",
		"array":   "数组",
		"slice":   "数组",
		"map":     "对象",
	}

	if typeName, exists := typeMessages[fieldType]; exists {
		return fmt.Sprintf("%s字段类型不正确，应为%s类型", fieldName, typeName)
	}

	return fmt.Sprintf("%s字段类型不正确，应为%s类型", fieldName, fieldType)
}

// CreateBusinessError 创建业务错误
func (h *ValidationHelper) CreateBusinessError(message string, code int) *response.DomainError {
	switch code {
	case 400:
		return response.NewValidationError(message)
	case 401:
		return response.NewUnauthorizedError(message)
	case 403:
		return response.NewForbiddenError(message)
	case 404:
		return response.NewNotFoundError(message)
	case 409:
		return response.NewAlreadyExistsError(message)
	case 500:
		return response.NewInternalServerError(message)
	default:
		return response.NewBusinessRuleViolationError(message)
	}
}

// SwaggerValidationAdapter Swagger验证适配器
// 用于将现有的验证逻辑适配到Swagger格式
type SwaggerValidationAdapter struct {
	helper *ValidationHelper
}

// NewSwaggerValidationAdapter 创建Swagger验证适配器
func NewSwaggerValidationAdapter() *SwaggerValidationAdapter {
	return &SwaggerValidationAdapter{
		helper: NewValidationHelper(),
	}
}

// AdaptValidationError 适配验证错误
func (a *SwaggerValidationAdapter) AdaptValidationError(c *gin.Context, params interface{}, err error, trans ut.Translator) {
	a.helper.HandleValidationError(c, params, err, trans)
}

// CreateFieldError 创建字段错误的便捷函数
func CreateSwaggerFieldError(field, message, value string) FieldError {
	return FieldError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// CreateMultipleFieldErrors 创建多个字段错误
func CreateMultipleFieldErrors(errors map[string]string) []FieldError {
	var fieldErrors []FieldError
	for field, message := range errors {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   field,
			Message: message,
			Value:   "",
		})
	}
	return fieldErrors
}

// ValidateAndHandleSwagger 验证并处理Swagger格式错误的便捷函数
func ValidateAndHandleSwagger(c *gin.Context, params interface{}, bindFunc func() error) bool {
	if err := bindFunc(); err != nil {
		helper := NewValidationHelper()
		helper.HandleValidationError(c, params, err, nil)
		return false
	}
	return true
}

// 全局验证辅助工具实例
var defaultValidationHelper = NewValidationHelper()

// GetDefaultValidationHelper 获取默认验证辅助工具
func GetDefaultValidationHelper() *ValidationHelper {
	return defaultValidationHelper
}

// HandleSwaggerValidationError 全局便捷函数
func HandleSwaggerValidationError(c *gin.Context, params interface{}, err error, trans ut.Translator) {
	defaultValidationHelper.HandleValidationError(c, params, err, trans)
}
