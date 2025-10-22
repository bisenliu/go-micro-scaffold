# Swagger 错误处理文档

本文档说明应格式。

## 概述

为了确保 API 。

## 错误响应格式

### 标准错误响应

```go
type ErrorResponse struct {
    Error   string      `json:"error" example:"Bad Rst"`
    Message string      `json:"message" example:"Inval
    Code    int         `json:"code" example:"400"`
    Details interface{} `json:"details,omitempty"`
}
```

### 验证错误响应

```go
type ValidationErrorResponse struct {
    Error   string                 
    Message string                 ailed"`
    Code    int                   "`
    Details ValidationErrorDetails `json:ails"`
}

type ValidationErrt {
"`
}

type FieldError struct {

    M
    Value   string
}
```

### 其他专用错误响应

- `UnauthorizedErrorResponse` - 401 未授权
- `ForbiddenErrorResponse` - 403 禁止访问
- `NotFoundErrorResponse` - 404 资源不存在
- `ConflictErrorResponse` - 409 资源冲突
- `InternalServerErrorResponse` - 500 服务器内部错误
- `ServiceUnavailableErrorResponse` - 503 服务不可用

## 使用方法

### 1. 更新 Swag注释

在 H类型：

```go
建用户
// @Summary 创建新用户

// @T管理
// @Accept json

// @Param户请求"
// @Success 200 {object} responsedto.UserInfoResponse "
// @Failure 400 {object} swagger.ValidationErrorResponse "请求参数验证失败"
// @Failure 409 {object} swagger.ConflictErrorResponse "用户已存在"
// @Failure 500 {object} swagger.InternalServerErrorResponse "服务器内部错误"
// rAuth
[post]
func (h *
现代码...
}
```

处理方式

#### 方式一：使用现有的统一响应系统（推荐）

现有的适当的格式：

```go
) {
    var rRequest
) {
     
    }

    user, err := h.commandHandler.HandleCreateUser(ctx, command)
    

    HandleWithLogging(c, responsedto.ToUserInfoR)
}
```

#### 方式

如果需要确保响应格式完全符合

```go
import "services/internawagger"

{
    var req requestt

        // 使用 Swagger 格式处理验证错
        swagger.HandleSwaggerValidationError)
        return
   }

nd)
    
    // 使用 Swagger 格式处理响应
    swagger.HandleWithSwagger, err)
}
```

#### 方式三：创建自定义验证错误

```go
import "services/internal/interfaces/

func (h *UserHandler) CreateU
    // 自定义验证逻辑
    if req.Name == "" {
        fieldErrors := []swagger.FieldEr{
            swagger.CreateSwaggerField,

     dErrors)

        return


    // 业务逻辑...
}
```

### 3. 创建业务错误

```go
import (
    "common/response"
    "services/internal/interfaces/http/swagger"
)

func (h *UserHandler) G
    userID := c.Param("id")
    
    user, err })
 
   

            ")
ndErr)
       n
   }
        
        // 其他错误
        swagge)
        return
    }
    
    swagger.HandleWithSw
}
```

## 错误类型映射

系统会自动将内部错误类型映射到相应的 HTTP 状态码和 Swagger 错误响应格

| 内部错误类 响应类型 |
|----|
| `
|
| `。点可以逐步迁移有端错误处理，现ger 格式的 端点中使用 SwagPI的 A

建议在新开发更容易更好，调试开发体验
4. 可维护错误处理逻辑统一且应
3. 地解析错误响可靠
2. 客户端能够 文档与实际响应格式一致
1. API格式，可以确保：
wagger 错误响应化的 S用标准结

通过使 总
```

##
}xt()
    } c.Ne  }
          
     return          Abort()
      c.r)
       eril,, ngerFormat(cwagthSleWigger.Hand       swa    
 "缺少认证令牌")rror(rizedEnauthowUe.Nens= respoerr :          
   == "" {    if tokenion")
    thorizatader("Aun := c.GetHetoke{
        .Context) nc(c *gin  return furFunc {
  in.Handle() giddleware AuthM
func```gor 响应适配器：

使用 Swagge中可以直接

A: 在中间件错误格式？ger 中使用 Swag何在中间件# Q: 如`

##)
``(c, nil, erraggerFormatandleWithSwr.Hors)
swaggeieldErr", fs("请求参数验证失败WithFielddationErrorteVali.Crea= swagger
err :"),
}id@inval"", 不正确", "邮箱格式"emailrFieldError(waggeger.CreateS
    swag, "")不能为空",", "姓名amerror("ndErFielSwagger.Createswaggerror{
    r.FieldEges := []swagfieldError
``go` 数组中：

`Fieldsetails.段错误包含在 `D 格式，将多个字se`nErrorResponValidatio: 使用 `
A错误？
Q: 如何处理复杂的验证
### 致。
档完全一保响应格式与文用于确主要处理是可选的，的 Swagger 格式正常工作。新可以e` 函数仍然e.Handl 和 `responsLogging`ndleWith不需要。现有的 `Ha改吗？

A: 全部修处理需要Q: 现有的错误常见问题

### `

## ``
ds)
}.FiellsairResp.Det errot,pty(rt.NotEm assede)
   p.CoorRes 400, err.Equal(t,ertassError)
     errorResp.iled",idation Faual(t, "Valassert.Eq   
    rResp)
 , &erroBodyp.arshal(res  json.Unmsponse
  ationErrorRegger.Validsp swa errorRe式
    var应格
    // 验证响   
 quest()nvalidRendI:= se    resp / 发送无效请求
T) {
    / *testing.nError(talidatio_VeUserstCreat
func Te```go和内容：

试中验证错误响应的格式应

确保在测## 5. 测试错误响```

#"),
}
数错误", "t", "参or("requesFieldErrger.CreateSwag  swagger{
  orFieldErrwagger.rrors := []s消息
fieldE
// 避免模糊的错误123"),
}
 "少6位",长度至 "密码ord",asswror("pErieldaggerFCreateSw swagger.l"),
   valid-emai确", "in", "邮箱格式不正"emailldError(aggerFieger.CreateSw  swagError{
  gger.Fieldwa:= []sldErrors 的错误消息
fie```go
// 好理解问题所在：

具体，帮助客户端
错误消息应该清晰、有意义的错误信息
提供### 4. 义。

r 注释中有相应的定gewag S应都在的错误响端点中，确保所有可能PI 一个 A

在同### 3. 保持一致性e{}`。

rfacstring]inte不是通用的 `map[应类型，而释中使用预定义的错误响在 Swagger 注类型

用标准错误响应 使 2.

###r"
```aggeces/http/swal/interfas/internce"serviport go
im加导入：

```件中添r 文Handleagger 错误处理的 要使用 Swger 包

在需 1. 导入 Swag实践

###
## 最佳|
Response` bleErrorUnavailaice03 | `Servble` | 5UnavailaalServiceTypeExtern| `Errorse` |
poneserverErrorRalS00 | `InternlServer` | 5ternaypeInrrorT` |
| `EsponselictErrorRe409 | `ConfyExists` | rTypeAlread`Erroe` |
| rrorResponsrbiddenE `Fon` | 403 |iddeypeForborT `Err` |
|sponsezedErrorRehori1 | `Unautzed` | 40horinautrorTypeUEr