package idgen

import (
	"fmt"
)

// ID 领域对象唯一标识
type ID struct {
	value int64
}

// NewIDFromInt64 从int64创建ID
func NewIDFromInt64(value int64) ID {
	return ID{value: value}
}

// Value 获取ID值
func (id ID) Value() int64 {
	return id.value
}

// String ID的字符串表示
func (id ID) String() string {
	return fmt.Sprintf("%d", id.value)
}

// IsZero 判断是否为零值
func (id ID) IsZero() bool {
	return id.value == 0
}

// Equals 判断两个ID是否相等
func (id ID) Equals(other ID) bool {
	return id.value == other.value
}
