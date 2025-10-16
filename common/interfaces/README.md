# Common 接口层

本目录定义了 common 模块对外提供的所有接口，实现了模块间的依赖倒置原则。

## 接口概览

### 核心接口

- **ConfigProvider**: 配置提供者接口，统一配置访问
- **Logger**: 日志接口，支持结构化日志和上下文传递
- **DatabaseManager**: 数据库管理器接口，支持主从分离
- **CacheManager**: 缓存管理器接口，统一缓存操作
- **Validator**: 验证器接口，数据验证和错误处理
- **MiddlewareProvider**: 中间件提供者接口，创建各种中间件
- **JWTService**: JWT服务接口，令牌管理

### 聚合接口

- **CommonServices**: 聚合所有 common 服务的接口，供外部模块使用

## 使用示例

### 在 Services 模块中使用

```go
// services/internal/shared/interfaces/common.go
package interfaces

import "common/interfaces"

// 重新导出接口，避免直接依赖 common 包
type CommonServices = interfaces.CommonServices
type ConfigProvider = interfaces.ConfigProvider
type Logger = interfaces.Logger
type DatabaseManager = interfaces.DatabaseManager
```

```go
// services/internal/application/service/user_service.go
package service

import (
    "context"
    
    "services/internal/shared/interfaces"
)

type UserService struct {
    commonServices interfaces.CommonServices
}

func NewUserService(commonServices interfaces.CommonServices) *UserService {
    return &UserService{
        commonServices: commonServices,
    }
}

func (s *UserService) CreateUser(ctx context.Context, name string) error {
    // 使用配置
    config := s.commonServices.Config()
    dbConfig := config.GetDatabaseConfig()
    
    // 使用日志
    logger := s.commonServices.Logger()
    logger.Info(ctx, "Creating user", zap.String("name", name))
    
    // 使用数据库
    db := s.commonServices.Database()
    primaryDB := db.GetPrimaryDB()
    
    // 业务逻辑...
    return nil
}
```

### 依赖注入配置

```go
// services/internal/shared/di.go
package shared

import (
    "go.uber.org/fx"
    
    "services/internal/shared/interfaces"
)

var SharedModule = fx.Module("shared",
    // 从 common 模块获取 CommonServices
    fx.Provide(
        func(commonServices interfaces.CommonServices) interfaces.CommonServices {
            return commonServices
        },
    ),
)
```

## 接口设计原则

1. **依赖倒置**: 高层模块不依赖低层模块，都依赖抽象
2. **接口隔离**: 接口职责单一，客户端不依赖不需要的接口
3. **开闭原则**: 对扩展开放，对修改封闭
4. **里氏替换**: 子类型必须能够替换其基类型

## 实现适配器

所有接口都有对应的适配器实现，将现有的具体实现适配到接口：

- `common/config/provider.go`: ConfigProvider 的实现
- `common/logger/adapter.go`: Logger 的实现  
- `common/databases/adapter.go`: DatabaseManager 的实现
- `common/services.go`: CommonServices 的聚合实现

## 扩展指南

添加新接口时，请遵循以下步骤：

1. 在对应的接口文件中定义接口
2. 创建适配器实现现有功能
3. 在 `common/services.go` 中添加到 CommonServices
4. 在 `common/di/modules.go` 中注册依赖注入
5. 更新文档和示例