# Swagger API 文档使用指南

## 概述

本项目已集成 Swagger API 文档功能，提供自动化的 API 文档生成和交互式 API 测试界面。本指南将详细说明如何使用和维护 Swagger 文档。

## 🚀 快速开始

### 访问 Swagger UI

启动应用后，可通过以下地址访问 Swagger UI：

```
http://localhost:8080/swagger/index.html
```

### 环境配置

Swagger 功能通过配置文件控制：

```yaml
# services/configs/app.yaml
swagger:
  enabled: true                    # 是否启用 Swagger
  title: "Go Micro Scaffold API"   # API 文档标题
  description: "微服务脚手架 API 文档"
  version: "1.0.0"
  host: "localhost:8080"
  base_path: "/api/v1"
  contact:
    name: "API Support"
    email: "support@example.com"
  license:
    name: "MIT"
    url: "https://opensource.org/licenses/MIT"
```

## 📝 API 文档注释规范

### 1. 主应用注释

在 `cmd/server/main.go` 中添加主应用信息：

```go
// @title Go Micro Scaffold API
// @version 1.0.0
// @description 微服务脚手架 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Token，格式：Bearer {token}

func main() {
    // 应用启动代码...
}
```

### 2. Handler 方法注释

#### 完整注释示例

```go
// CreateUser 创建用户
// @Summary 创建新用户
// @Description 创建一个新的用户账户，需要提供用户基本信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body requestdto.CreateUserRequest true "创建用户请求"
// @Success 200 {object} responsedto.UserInfoResponse "创建成功"
// @Failure 400 {object} swagger.ValidationErrorResponse "请求参数验证失败"
// @Failure 409 {object} swagger.ConflictErrorResponse "用户已存在"
// @Failure 500 {object} swagger.InternalServerErrorResponse "服务器内部错误"
// @Security BearerAuth
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // 实现代码...
}
```

#### 注释标签说明

| 标签 | 说明 | 示例 |
|------|------|------|
| `@Summary` | API 简短描述 | `@Summary 创建新用户` |
| `@Description` | API 详细描述 | `@Description 创建一个新的用户账户` |
| `@Tags` | API 分组标签 | `@Tags 用户管理` |
| `@Accept` | 接受的内容类型 | `@Accept json` |
| `@Produce` | 返回的内容类型 | `@Produce json` |
| `@Param` | 参数定义 | `@Param id path string true "用户ID"` |
| `@Success` | 成功响应 | `@Success 200 {object} UserResponse` |
| `@Failure` | 错误响应 | `@Failure 400 {object} ErrorResponse` |
| `@Security` | 安全认证 | `@Security BearerAuth` |
| `@Router` | 路由定义 | `@Router /api/v1/users [post]` |

### 3. 数据模型注释

#### 请求 DTO 注释

```go
// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    OpenID      string `json:"open_id" binding:"required" example:"wx_123456789" validate:"required"`      // 微信OpenID
    Name        string `json:"name" binding:"required,min=2,max=50" example:"张三" validate:"required,min=2,max=50"`     // 用户姓名，长度2-50字符
    Gender      int    `json:"gender" binding:"oneof=0 1 2" example:"1" validate:"oneof=0 1 2"`              // 性别：0-未知，1-男，2-女
    PhoneNumber string `json:"phone_number" binding:"required,phone" example:"13800138000" validate:"required,phone"` // 手机号码
    Password    string `json:"password" binding:"required,min=6" example:"password123" validate:"required,min=6"`     // 密码，最少6位
}
```

#### 响应 DTO 注释

```go
// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
    ID          string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`          // 用户唯一标识
    OpenID      string `json:"open_id" example:"wx_123456789"`                             // 微信OpenID
    Name        string `json:"name" example:"张三"`                                         // 用户姓名
    Gender      int    `json:"gender" example:"1"`                                         // 性别：0-未知，1-男，2-女
    PhoneNumber string `json:"phone_number" example:"13800138000"`                         // 手机号码
    CreatedAt   int64  `json:"created_at" example:"1640995200000"`                         // 创建时间戳（毫秒）
    UpdatedAt   int64  `json:"updated_at" example:"1640995200000"`                         // 更新时间戳（毫秒）
}
```

### 4. 错误响应注释

项目提供了标准化的错误响应类型：

```go
// 在 Handler 注释中使用标准错误类型
// @Failure 400 {object} swagger.ValidationErrorResponse "请求参数验证失败"
// @Failure 401 {object} swagger.UnauthorizedErrorResponse "未授权访问"
// @Failure 403 {object} swagger.ForbiddenErrorResponse "禁止访问"
// @Failure 404 {object} swagger.NotFoundErrorResponse "资源不存在"
// @Failure 409 {object} swagger.ConflictErrorResponse "资源冲突"
// @Failure 500 {object} swagger.InternalServerErrorResponse "服务器内部错误"
// @Failure 503 {object} swagger.ServiceUnavailableErrorResponse "服务不可用"
```

## 🔧 文档生成和维护

### 1. 生成 Swagger 文档

```bash
# 进入 services 目录
cd services

# 生成 Swagger 文档
swag init -g cmd/server/main.go -o docs

# 验证生成的文档
ls -la docs/
# 应该看到：docs.go, swagger.json, swagger.yaml
```

### 2. 验证文档完整性

使用项目提供的验证脚本：

```bash
# 运行 Swagger 文档验证
cd services
go run scripts/validate-swagger.go

# 或使用 Makefile
make validate-swagger
```

### 3. 文档更新流程

1. **修改 API 代码**：更新 Handler 方法或 DTO 结构体
2. **更新注释**：按照规范更新 Swagger 注释
3. **重新生成文档**：运行 `swag init` 命令
4. **验证文档**：检查生成的文档是否正确
5. **提交代码**：将生成的文档文件一起提交

## 🛠️ 开发最佳实践

### 1. 注释编写规范

#### ✅ 好的注释示例

```go
// GetUser 获取用户信息
// @Summary 根据用户ID获取用户详细信息
// @Description 通过用户唯一标识获取用户的详细信息，包括基本资料和状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户唯一标识" format(uuid)
// @Success 200 {object} responsedto.UserInfoResponse "获取成功"
// @Failure 400 {object} swagger.ValidationErrorResponse "请求参数格式错误"
// @Failure 404 {object} swagger.NotFoundErrorResponse "用户不存在"
// @Failure 500 {object} swagger.InternalServerErrorResponse "服务器内部错误"
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
```

#### ❌ 不好的注释示例

```go
// GetUser 获取用户
// @Summary 获取用户
// @Router /api/v1/users/{id} [get]
```

### 2. 参数类型规范

| 参数位置 | 语法 | 示例 |
|----------|------|------|
| 路径参数 | `@Param name path type required "description"` | `@Param id path string true "用户ID"` |
| 查询参数 | `@Param name query type required "description"` | `@Param page query int false "页码"` |
| 请求体 | `@Param name body type required "description"` | `@Param request body CreateUserRequest true "用户信息"` |
| 头部参数 | `@Param name header type required "description"` | `@Param Authorization header string true "JWT Token"` |

### 3. 响应格式规范

```go
// 成功响应
// @Success 200 {object} ResponseType "描述"
// @Success 201 {object} ResponseType "创建成功"

// 分页响应
// @Success 200 {object} PaginatedResponse{data=[]UserResponse} "分页数据"

// 数组响应
// @Success 200 {array} UserResponse "用户列表"

// 简单响应
// @Success 200 {string} string "操作成功"
```

### 4. 安全认证配置

```go
// 需要认证的接口
// @Security BearerAuth

// 可选认证的接口
// @Security BearerAuth || {}

// 不需要认证的接口
// 不添加 @Security 标签
```

## 🔍 测试和调试

### 1. 使用 Swagger UI 测试 API

1. **访问 Swagger UI**：打开 `http://localhost:8080/swagger/index.html`
2. **选择接口**：点击要测试的 API 端点
3. **填写参数**：根据文档填写必要的参数
4. **执行请求**：点击 "Try it out" 和 "Execute"
5. **查看响应**：检查返回的状态码和响应数据

### 2. JWT 认证测试

1. **获取 Token**：先调用登录接口获取 JWT Token
2. **设置认证**：点击页面顶部的 "Authorize" 按钮
3. **输入 Token**：在弹出框中输入 `Bearer your-jwt-token`
4. **测试接口**：现在可以测试需要认证的接口

### 3. 常见问题排查

#### 文档生成失败

```bash
# 检查语法错误
swag init -g cmd/server/main.go -o docs --parseVendor

# 查看详细错误信息
swag init -g cmd/server/main.go -o docs --parseVendor --parseDependency
```

#### 接口不显示

1. 检查 Handler 方法是否有正确的注释
2. 确认路由是否正确注册
3. 验证注释语法是否符合规范

#### 参数类型错误

1. 检查 DTO 结构体的 tag 定义
2. 确认参数类型与实际代码一致
3. 验证 binding 和 validation 标签

## 📊 监控和维护

### 1. 文档质量检查

定期检查以下项目：

- [ ] 所有公开 API 都有完整的 Swagger 注释
- [ ] 错误响应类型使用标准格式
- [ ] 参数描述清晰准确
- [ ] 示例数据真实有效
- [ ] 安全认证配置正确

### 2. 自动化检查

在 CI/CD 流程中添加文档检查：

```bash
# .github/workflows/api-docs.yml
name: API Documentation Check

on: [push, pull_request]

jobs:
  swagger-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Install swag
        run: go install github.com/swaggo/swag/cmd/swag@latest
      
      - name: Generate docs
        run: |
          cd services
          swag init -g cmd/server/main.go -o docs
      
      - name: Validate docs
        run: |
          cd services
          go run scripts/validate-swagger.go
```

### 3. 版本管理

- 在 API 有重大变更时更新版本号
- 保持文档与代码同步更新
- 记录 API 变更历史

## 🔒 安全考虑

### 1. 生产环境配置

```yaml
# 生产环境配置
swagger:
  enabled: false  # 生产环境建议禁用
```

### 2. 敏感信息保护

- 不在示例中使用真实的敏感数据
- 密码字段使用通用示例值
- API Key 和 Token 使用占位符

### 3. 访问控制

```go
// 可选：添加 IP 白名单控制
if swaggerConfig.Enabled {
    swaggerGroup := engine.Group("/swagger")
    swaggerGroup.Use(middleware.IPWhitelistMiddleware(allowedIPs))
    swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

## 📚 参考资源

- [Swagger 官方文档](https://swagger.io/docs/)
- [swaggo/swag 文档](https://github.com/swaggo/swag)
- [OpenAPI 3.0 规范](https://swagger.io/specification/)
- [Gin Swagger 集成](https://github.com/swaggo/gin-swagger)

## 🤝 贡献指南

如需改进 Swagger 文档：

1. 遵循本指南的注释规范
2. 确保文档与代码同步
3. 添加必要的测试用例
4. 更新相关文档

---

通过遵循本指南，可以确保 API 文档的质量和一致性，提升开发效率和用户体验。