package idgen

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"

	"common/config"
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
func NewSnowflakeGenerator(cfg *config.Config) (Generator, error) {
	// 解析配置中的起始时间
	startTime, err := time.Parse(time.DateOnly, cfg.SnowFlake.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid snowflake start time format: %v", err)
	}

	// 设置snowflake起始时间
	snowflake.Epoch = startTime.UnixNano() / 1000000

	// 创建snowflake节点
	node, err := snowflake.NewNode(cfg.SnowFlake.MachineID)
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
