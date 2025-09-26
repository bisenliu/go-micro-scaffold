package valueobject

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

// ID 领域对象唯一标识
type ID struct {
	value int64
}

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(fmt.Sprintf("Failed to create snowflake node: %v", err))
	}
}

// NewID 创建新的ID
func NewID() ID {
	return ID{value: node.Generate().Int64()}
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
