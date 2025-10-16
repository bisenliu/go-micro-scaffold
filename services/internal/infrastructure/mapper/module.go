package mapper

import (
	"go.uber.org/fx"
)

// Module 映射器模块
var Module = fx.Module("mapper",
	fx.Provide(
		// 用户映射器
		fx.Annotate(
			NewUserMapper,
			fx.As(new(UserMapper)),
		),
		
		// 用户映射器辅助类
		NewUserMapperHelper,
		
		// 映射器工厂
		NewMapperFactory,
		
		// 映射器工厂提供者
		fx.Annotate(
			NewMapperFactoryProvider,
			fx.As(new(MapperFactoryProvider)),
		),
	),
)