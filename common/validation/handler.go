package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"common/logger"
)

// Validatable 可验证接口
type Validatable interface {
	Validate() error
}

// ValidationError 自定义验证错误
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError 创建验证错误的辅助函数
func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

// Verify 执行绑定操作并自动处理错误
func Verify(c *gin.Context, params interface{}, bindMethod BindMethod, logger *zap.Logger, trans ut.Translator) bool {
	if err := bindMethod(c, params); err != nil {
		handleError(c, params, err, logger, trans)
		return false
	}
	return true
}

// ValidateError 处理自定义验证错误的辅助函数
func ValidateError(c *gin.Context, params interface{}, err error, logger *zap.Logger, trans ut.Translator) bool {
	if err != nil {
		handleError(c, params, err, logger, trans)
		return false
	}
	return true
}

// handleError 处理验证错误
func handleError(c *gin.Context, params interface{}, err error, zapLogger *zap.Logger, trans ut.Translator) {
	ctx := c.Request.Context()
	logger.Error(zapLogger, ctx, "invalid params",
		zap.Error(err),
		zap.String("url", c.Request.URL.String()),
		zap.Any("params", params),
	)

	switch err := err.(type) {
	case ValidationError:
		// 自定义验证错误
		respondWithError(c, HTTPBadRequest, err.Message, nil)
	case validator.ValidationErrors:
		// 验证器错误
		respondWithError(c, HTTPBadRequest, ErrValidationFailed, removeTopStruct(err.Translate(trans)))
	case *json.UnmarshalTypeError:
		// JSON类型错误
		fieldName := getFieldJSONName(params, err)
		message := buildTypeErrorMessage(fieldName, err.Type.String())
		respondWithError(c, HTTPBadRequest, message, nil)
	default:
		// 其他错误
		respondWithError(c, HTTPBadRequest, err.Error(), nil)
	}
}

// respondWithError 统一的错误响应格式
func respondWithError(c *gin.Context, code int, message string, errors interface{}) {
	response := gin.H{
		"code":    code,
		"message": message,
	}

	if errors != nil {
		response["errors"] = errors
	}

	c.JSON(code, response)
}

// removeTopStruct 移除顶层结构体名称，提取错误信息
func removeTopStruct(fields map[string]string) string {
	var errors []string
	for _, err := range fields {
		parts := strings.Split(err, "|")
		if len(parts) == 2 {
			errors = append(errors, parts[1])
		} else {
			errors = append(errors, err)
		}
	}
	return strings.Join(errors, "; ")
}

// getFieldJSONName 根据JSON字段名获取显示名称
func getFieldJSONName(params interface{}, unmarshalTypeError *json.UnmarshalTypeError) string {
	t := reflect.TypeOf(params)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByNameFunc(func(name string) bool {
		f, _ := t.FieldByName(name)
		jsonName := strings.SplitN(f.Tag.Get(JSONTag), ",", 2)[0]
		return jsonName == unmarshalTypeError.Field
	})

	if !found {
		return unmarshalTypeError.Field
	}

	// 优先使用 label 标签
	if label := field.Tag.Get(LabelTag); label != "" {
		return label
	}

	// 其次使用 json 标签
	if jsonName := strings.SplitN(field.Tag.Get(JSONTag), ",", 2)[0]; jsonName != "" {
		return jsonName
	}

	return field.Name
}

// buildTypeErrorMessage 构建类型错误消息
func buildTypeErrorMessage(fieldName, fieldType string) string {
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
		"bool":    "布尔",
		"string":  "字符串",
		"array":   "数组",
		"slice":   "数组",
		"map":     "映射",
	}

	if typeName, exists := typeMessages[fieldType]; exists {
		return fmt.Sprintf("%s字段类型不正确，应为%s类型", fieldName, typeName)
	}

	return fmt.Sprintf("%s字段类型不正确，应为%s类型", fieldName, fieldType)
}
