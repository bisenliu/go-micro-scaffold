package casbin

import "go.uber.org/fx"

// Module 提供了 Casbin Enforcer
var Module = fx.Module("casbin",
	fx.Provide(NewEnforcer),
)
