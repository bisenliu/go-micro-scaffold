package idgen

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"

	"common/interfaces"
)

// Generator ID生成器接口
type Generator interface {
	NewID() ID
	NewIDFromInt64(value int64) ID
}

// snowflakeGenerator snowflake ID生成器实现
type snowflakeGenerator struct {
	node *snowflake.Node
}

// NewSnowflakeGenerator 创建snowflake ID生成器
func NewSnowflakeGenerator(configProvider interfaces.ConfigProvider) (Generator, error) {
	// 使用默认配置，后续可以从配置中获取
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	machineID := int64(1)

	// 设置snowflake起始时间
	snowflake.Epoch = startTime.UnixNano() / 1000000

	// 创建snowflake节点
	node, err := snowflake.NewNode(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake node: %v", err)
	}

	return &snowflakeGenerator{
		node: node,
	}, nil
}

// NewID 创建新的ID
func (g *snowflakeGenerator) NewID() ID {
	return ID{value: g.node.Generate().Int64()}
}

// NewIDFromInt64 从int64创建ID
func (g *snowflakeGenerator) NewIDFromInt64(value int64) ID {
	return ID{value: value}
}
