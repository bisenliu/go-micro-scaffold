package shared

import (
	"go.uber.org/fx"
)

// SharedModule 共享模块
// 提供 services 模块内部使用的共享服务和接口
var SharedModule = fx.Module("shared",
	// 暂时不提供额外的服务，直接使用 CommonServices
)