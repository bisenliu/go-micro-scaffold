package ent

import (
	"go.uber.org/fx"

	"services/internal/infrastructure/persistence"
	"services/internal/infrastructure/persistence/ent/gen"
)

// Module Ent 模块
var Module = fx.Module("ent",
	// 提供 DatabaseProvider
	fx.Provide(persistence.NewDatabaseProvider),
	
	// 提供 gen.Client，基于 DatabaseProvider
	fx.Provide(func(provider *persistence.DatabaseProvider) (*gen.Client, error) {
		return provider.CreateEntClient()
	}),
)
