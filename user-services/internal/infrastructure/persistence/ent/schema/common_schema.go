package schema

import (
	schema "common/schema/common"
)

// 这里是公共的Schema,多个服务可以引用
type CommonSchema struct {
	schema.BaseSchema
}
