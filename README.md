# Go Micro Scaffold

Go Micro Scaffold 是一个基于 Go 语言的微服务脚手架项目，采用 Clean Architecture 架构设计，集成了常用的微服务组件和最佳实践。

## 项目特点

- **Clean Architecture**: 采用领域驱动设计（DDD）和六边形架构
- **模块化设计**: 使用 Uber FX 进行依赖注入和模块管理
- **高性能**: 基于 Gin Web 框架构建
- **数据库集成**: 集成 Ent ORM，支持多数据库
- **配置管理**: 基于 Viper 的配置管理
- **日志系统**: 集成 Zap 高性能日志库
- **验证系统**: 集成验证器，支持多语言
- **中间件**: 内置常用中间件（CORS、认证、限流等）

## 项目结构

```
go-micro-scaffold/
├── common/                 # 公共库
│   ├── config/             # 配置管理
│   ├── databases/          # 数据库相关
│   │   ├── mysql/          # MySQL数据库
│   │   └── redis/          # Redis缓存
│   ├── di/                 # 依赖注入模块
│   ├── http/               # HTTP 服务
│   ├── logger/             # 日志系统
│   ├── middleware/         # 中间件
│   ├── pkg/                # 通用工具包
│   │   ├── idgen/          # ID生成器
│   │   ├── jwt/            # JWT认证
│   │   ├── timezone/       # 时区管理
│   │   └── validation/     # 验证系统
│   ├── response/           # 响应处理
│   └── schema/             # 数据库模式
├── services/               # 服务模块
│   ├── cmd/                # 命令行入口
│   │   ├── cli/            # CLI 命令
│   │   └── server/         # 服务端
│   ├── configs/            # 配置文件
│   ├── internal/           # 内部实现
│   │   ├── application/    # 应用层
│   │   ├── domain/         # 领域层
│   │   ├── infrastructure/ # 基础设施层
│   │   └── interfaces/     # 接口层
│   └── go.mod              # Go 模块定义
└── go.work                 # Go 工作区
```

## 快速开始

### 环境要求

- Go 1.24+
- MySQL 5.7+
- Redis 5.0+

### 安装依赖

```bash
# 进入项目根目录
cd go-micro-scaffold

# 初始化 Go 工作区
go work init
go work use ./services
go work use ./common

# 安装依赖
cd services
go mod tidy
```

### 配置文件

1. 复制配置文件模板：
```bash
cd services/configs
cp app.yaml.example app.yaml
```

2. 根据实际环境修改 [app.yaml](file:///Users/liubisen/Desktop/sander/Project/my/go-micro-scaffold/services/configs/app.yaml) 配置文件

### 数据库迁移

```bash
# 执行数据库迁移
cd services
go run cmd/cli/main.go migrate
```

### 启动服务

```bash
# 启动服务
cd services
go run cmd/server/main.go
```

## API 接口

### 健康检查

```bash
GET /health
GET /ping
```

### 用户相关

```bash
POST /api/v1/users    # 创建用户
```

## 项目模块说明

### 应用层 (Application Layer)

- 负责业务用例的编排
- 定义命令和查询
- 处理业务逻辑

### 领域层 (Domain Layer)

- 核心业务逻辑
- 实体、值对象、聚合根
- 领域服务和仓储接口

### 基础设施层 (Infrastructure Layer)

- 数据库实现
- 外部服务集成
- 仓储实现

### 接口层 (Interface Layer)

- HTTP 控制器
- DTO 数据传输对象
- 路由配置

## 依赖注入

项目使用 Uber FX 进行依赖注入管理，各模块通过 FX 模块进行组织和注入。

有关 Uber FX 框架核心概念的详细说明，请参考 [FX 框架指南](FX_FRAMEWORK_GUIDE.md)。

## 配置说明

项目支持丰富的配置选项，详细配置说明请参考 [app.yaml.example](file:///Users/liubisen/Desktop/sander/Project/my/go-micro-scaffold/services/configs/app.yaml.example) 文件。

## 日志系统

项目集成了 Zap 日志库，支持结构化日志输出和日志级别控制。

## 验证系统

项目集成了验证器，支持请求参数验证和多语言错误提示。

## 中间件

- CORS 跨域支持
- 认证中间件
- 限流中间件
- 请求日志中间件
- IP 白名单中间件
- Recovery 中间件

### 时区管理

项目提供了时区管理模块，用于全局设置应用程序的时区。该模块从配置文件中读取时区设置，如果没有配置则默认使用 "Asia/Shanghai"。

使用方法：
1. 在配置文件中添加时区设置：
```yaml
system:
  timezone: "Asia/Shanghai"  # 或其他时区，如 "America/New_York"
```

2. 时区模块会自动在应用启动时初始化（通过 Uber FX 依赖注入）：
```go
// 在 common/di/modules.go 中已经注册
var TimezoneModule = fx.Module("timezone",
    timezone.Module,
)
```

时区模块会全局设置 [time.Local](file:///Users/liubisen/Desktop/sander/Project/my/go-micro-scaffold/services/internal/infrastructure/persistence/ent/gen/mutation.go#L313-L313) 和环境变量，确保整个应用程序使用统一的时区。时区只在应用启动时初始化一次，而不是在每个请求中都设置。

## 开发指南

### 添加新功能

1. 在领域层创建实体和仓储接口
2. 在应用层创建命令/查询和处理器
3. 在接口层创建 HTTP 控制器和路由
4. 注册新模块到依赖注入容器

### 数据库操作

项目使用 Ent ORM 进行数据库操作，可通过以下方式生成代码：

```bash
cd services/internal/infrastructure/persistence/ent
go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
```

## 部署

### 构建

```bash
cd services
go build -o bin/server cmd/server/main.go
go build -o bin/cli cmd/cli/main.go
```

### 运行

```bash
# 运行服务
./bin/server

# 运行 CLI 工具
./bin/cli migrate
```

## 贡献

欢迎提交 Issue 和 Pull Request 来改进项目。

## 许可证

[MIT License](LICENSE)