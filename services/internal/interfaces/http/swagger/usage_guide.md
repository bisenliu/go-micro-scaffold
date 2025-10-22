# Swagger 错误处理使用指南

## 概述

本项目已实现标准化的 Swagger 错误响应格式，确保 API 文档与实际响应一致。

## 主要改进

### 1. 标准化错误响应类型

所有 Handler 的 Swagger 注释已更新为使用标准错误响应类型：

- `swagger.ValidationErrorResponse` - 参数验证失败
- `swagger.UnauthorizedErrorResponse` - 未授权访问  
- `swagger.ForbiddenErrorResponse` - 禁止访问
- `swagger.NotFoundErrorResponse` - 资源不存在
- `swagger.ConflictErrorResponse` - 资源冲突
- `swagger.InternalServerErrorResponse` - 服务器内部错误
- `swagger.ServiceUnavailableErrorResponse` - 服务不可用

### 2. 响应适配器

创建了响应适配器确保实际响应格式与 Swagger 文档一致：

```go
// 使用 Swagger 格式处理响应
swagger.HandleWithSwaggerFormat(c, data, err)

// 使用 Swagger 格式处理分页响应  
swagger.HandlePagingWithSwaggerFormat(c, data, page, pageSize, total, err)
```

### 3. 验证错误处理

提供了专门的验证错误处理工具：

```go
// 创建字段验证错误
fieldErrors := []swagger.FieldError{
    swagger.CreateSwaggerFieldError("name", "姓名不能为空", ""),
}
err := swagger.CreateValidationErrorWithFields("验证失败", fieldErrors)
```

## 使用方法

### 现有代码兼容性

现有的错误处理代码无需修改，仍然可以正常工作：

```go
// 现有代码继续有效
HandleWithLogging(c, data, err)
response.Handle(c, data, err)
```

### 新的 Swagger 格式处理

如需确保完全符合 Swagger 文档格式：

```go
import "services/internal/interfaces/http/swagger"

// 在 Handler 中使用
swagger.HandleWithSwaggerFormat(c, data, err)
```

## 文件更新说明

### Handler 文件更新

1. **user_handler.go** - 更新所有错误响应注释
2. **auth_handler.go** - 更新认证相关错误响应注释  
3. **health_handler.go** - 更新健康检查错误响应注释
4. **error_handler.go** - 添加 Swagger 格式处理函数

### 新增文件

1. **response_adapter.go** - 响应格式适配器
2. **response_middleware.go** - 响应格式中间件
3. **validation_helper.go** - 验证错误处理辅助工具

## 验证

运行以下命令验证更新：

```bash
# 生成 Swagger 文档
swag init -g cmd/server/main.go -o docs

# 检查生成的文档
cat docs/swagger.json | jq '.definitions' | grep -E "(Error|Response)"
```

## 总结

此次更新确保了：

1. ✅ 所有 API 返回标准化的错误格式
2. ✅ Swagger 注释使用正确的错误响应类型
3. ✅ 提供了便捷的错误处理工具
4. ✅ 保持了向后兼容性
5. ✅ 符合需求 1.4 和 3.4 的要求