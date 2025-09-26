package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// Validator 验证器结构体，封装了验证和翻译功能
type Validator struct {
	uni      *ut.UniversalTranslator // 通用翻译器
	validate *validator.Validate     // 验证器实例
	trans    ut.Translator           // 当前语言翻译器
}

// SupportedLocales 支持的语言环境
var SupportedLocales = map[string]bool{
	"zh": true,
	"en": true,
}

// NewValidator 创建一个新的验证器实例
func NewLocalizedValidator(locale string) (*Validator, error) {

	// 验证语言环境是否支持
	if !SupportedLocales[locale] {
		return nil, fmt.Errorf("unsupported locale: %s", locale)
	}

	// 获取gin框架中的验证器引擎
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("failed to get validator engine")
	}

	// 注册自定义标签名称函数，支持json和label标签
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonName := strings.SplitN(fld.Tag.Get(JSONTag), ",", 2)[0]
		label := fld.Tag.Get(LabelTag)

		// 如果json标签为"_"，则忽略该字段
		if jsonName == "_" {
			return ""
		}

		// 优先使用label标签，格式为 "jsonName|label"
		if label != "" {
			return jsonName + "|" + label
		}

		return jsonName
	})

	// 创建翻译器
	zhT := zh.New() // 中文翻译器
	enT := en.New() // 英文翻译器

	// 创建通用翻译器，第一个参数是默认语言
	uni := ut.New(enT, zhT, enT)

	// 获取指定语言的翻译器
	trans, ok := uni.GetTranslator(locale)
	if !ok {
		return nil, fmt.Errorf("failed to get translator for locale: %s", locale)
	}

	// 注册默认翻译
	if err := zhTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, fmt.Errorf("failed to register translations: %w", err)
	}

	//在校验器注册自定义的校验方法
	if err := v.RegisterValidation("enum", ValidateEnum); err != nil {
		return nil, fmt.Errorf("failed to register validation: %w", err)
	}

	//注意！因为这里会使用到trans实例
	//所以这一步注册要放到trans初始化的后面
	if err := v.RegisterTranslation(
		"enum",
		trans,
		registerTranslator("enum", "{0}不合法"),
		translate,
	); err != nil {
		return nil, fmt.Errorf("failed to register translation: %w", err)
	}

	return &Validator{
		uni:      uni,
		validate: v,
		trans:    trans,
	}, nil
}

// GetTranslator 获取验证器翻译器
func (v *Validator) GetTranslator() ut.Translator {
	return v.trans
}

// GetValidate 获取验证器实例
func (v *Validator) GetValidate() *validator.Validate {
	return v.validate
}

// Verify 执行绑定操作并自动处理错误，封装了翻译器
func (v *Validator) Verify(c *gin.Context, params interface{}, bindMethod BindMethod) bool {
	return verify(c, params, bindMethod, v.trans)
}

// ValidateError 处理自定义验证错误的辅助函数
func (v *Validator) ValidateError(c *gin.Context, params interface{}, err error) bool {
	return validateError(c, params, err, v.trans)
}

// registerTranslator为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fmt.Errorf("translator failed: %w", fe.(error)))
	}
	return msg
}
