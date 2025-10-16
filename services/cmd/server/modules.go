package main

import (
	"go.uber.org/fx"

	commonDI "common/di"
	"services/internal/application"
	"services/internal/domain/user"
	"services/internal/infrastructure"
	"services/internal/interfaces/http"
	"services/internal/shared"
)

// ModuleRegistry 模块注册表
// 按依赖顺序组织所有应用模块
type ModuleRegistry struct {
	// 基础模块 - 提供核心基础设施
	CoreModules fx.Option
	
	// 共享模块 - 提供接口适配
	SharedModule fx.Option
	
	// 领域模块 - 业务逻辑核心
	DomainModules fx.Option
	
	// 应用模块 - 应用服务层
	ApplicationModule fx.Option
	
	// 基础设施模块 - 外部依赖实现
	InfrastructureModule fx.Option
	
	// 接口模块 - 对外接口层
	InterfaceModule fx.Option
}

// NewModuleRegistry 创建模块注册表
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		CoreModules:          commonDI.GetCoreModules(),
		SharedModule:         shared.SharedModule,
		DomainModules:        user.DomainModule,
		ApplicationModule:    application.ApplicationModule,
		InfrastructureModule: infrastructure.InfrastructureModule,
		InterfaceModule:      http.InterfaceModuleFinal,
	}
}

// GetAllModules 获取所有模块（按依赖顺序）
func (r *ModuleRegistry) GetAllModules() []fx.Option {
	return []fx.Option{
		r.CoreModules,          // 1. 核心基础设施（配置、日志、数据库等）
		r.SharedModule,         // 2. 共享接口适配
		r.DomainModules,        // 3. 领域模块（业务规则）
		r.ApplicationModule,    // 4. 应用服务（用例实现）
		r.InfrastructureModule, // 5. 基础设施实现（仓储、外部服务）
		r.InterfaceModule,      // 6. 接口层（HTTP、gRPC等）
	}
}

// GetModulesForConfig 根据配置获取模块
func (r *ModuleRegistry) GetModulesForConfig(config *AppConfig) []fx.Option {
	modules := r.GetCoreModules()
	
	// 根据配置条件加载接口模块
	if config.ShouldLoadHTTPModule() {
		modules = append(modules, r.InterfaceModule)
	}
	
	return modules
}

// GetCoreModules 获取核心模块（不包含接口层）
func (r *ModuleRegistry) GetCoreModules() []fx.Option {
	return []fx.Option{
		r.CoreModules,
		r.SharedModule,
		r.DomainModules,
		r.ApplicationModule,
		r.InfrastructureModule,
	}
}

// GetModulesWithoutInterface 获取除接口层外的所有模块
func (r *ModuleRegistry) GetModulesWithoutInterface() []fx.Option {
	return r.GetCoreModules()
}

// ModuleInfo 模块信息
type ModuleInfo struct {
	Name        string
	Description string
	Order       int
}

// GetModuleInfo 获取模块信息
func (r *ModuleRegistry) GetModuleInfo() []ModuleInfo {
	return []ModuleInfo{
		{
			Name:        "Core",
			Description: "核心基础设施模块（配置、日志、数据库、JWT等）",
			Order:       1,
		},
		{
			Name:        "Shared",
			Description: "共享模块（接口适配、通用服务）",
			Order:       2,
		},
		{
			Name:        "Domain",
			Description: "领域模块（业务实体、领域服务、仓储接口）",
			Order:       3,
		},
		{
			Name:        "Application",
			Description: "应用模块（用例实现、命令查询处理器）",
			Order:       4,
		},
		{
			Name:        "Infrastructure",
			Description: "基础设施模块（仓储实现、外部服务、映射器）",
			Order:       5,
		},
		{
			Name:        "Interface",
			Description: "接口模块（HTTP路由、处理器、中间件）",
			Order:       6,
		},
	}
}