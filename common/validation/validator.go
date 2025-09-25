package validation

import (
	"fmt"
	"reflect"
	"strings"

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
func NewValidator(locale string) (*Validator, error) {
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

// TranslateError 翻译验证错误为可读的错误信息
func (v *Validator) TranslateError(err error) map[string]string {
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = e.Translate(v.trans)
		}
	}

	return errors
}

// ValidateStruct 验证结构体
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// ValidateVar 验证单个变量
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}
