package mapper

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

// BaseMapper 基础映射器
// 提供一些通用的映射功能
type BaseMapper struct{}

// NewBaseMapper 创建基础映射器
func NewBaseMapper() *BaseMapper {
	return &BaseMapper{}
}

// SafeStringToUUID 安全地将字符串转换为UUID
func (m *BaseMapper) SafeStringToUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, nil
	}
	return uuid.Parse(s)
}

// SafeUUIDToString 安全地将UUID转换为字符串
func (m *BaseMapper) SafeUUIDToString(id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}
	return id.String()
}

// SafeTimeToUnixMilli 安全地将时间转换为毫秒时间戳
func (m *BaseMapper) SafeTimeToUnixMilli(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixMilli()
}

// SafeUnixMilliToTime 安全地将毫秒时间戳转换为时间
func (m *BaseMapper) SafeUnixMilliToTime(timestamp int64) time.Time {
	if timestamp == 0 {
		return time.Time{}
	}
	return time.UnixMilli(timestamp)
}

// IsNil 检查接口是否为nil
func (m *BaseMapper) IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}

// SafeStringPointerValue 安全地获取字符串指针的值
func (m *BaseMapper) SafeStringPointerValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// SafeIntPointerValue 安全地获取整数指针的值
func (m *BaseMapper) SafeIntPointerValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// SafeStringValuePointer 安全地创建字符串值的指针
func (m *BaseMapper) SafeStringValuePointer(value string) *string {
	return &value
}

// SafeIntValuePointer 安全地创建整数值的指针
func (m *BaseMapper) SafeIntValuePointer(value int) *int {
	return &value
}

// CopyNonZeroFields 复制非零值字段
// 这是一个简单的实现，实际项目中可能需要更复杂的逻辑
func (m *BaseMapper) CopyNonZeroFields(dest, src interface{}) {
	destValue := reflect.ValueOf(dest)
	srcValue := reflect.ValueOf(src)
	
	if destValue.Kind() != reflect.Ptr || srcValue.Kind() != reflect.Ptr {
		return
	}
	
	destElem := destValue.Elem()
	srcElem := srcValue.Elem()
	
	if destElem.Type() != srcElem.Type() {
		return
	}
	
	for i := 0; i < srcElem.NumField(); i++ {
		srcField := srcElem.Field(i)
		destField := destElem.Field(i)
		
		if !srcField.IsZero() && destField.CanSet() {
			destField.Set(srcField)
		}
	}
}