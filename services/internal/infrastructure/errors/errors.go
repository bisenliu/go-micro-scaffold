package errors

import (
	"common/response"
)

// 基础设施层通用错误
var (
	ErrInternalServer = response.NewInternalServerError("内部服务器错误")
	ErrInvalidRequest = response.NewDomainError(response.ErrorTypeInvalidRequest, "无效的请求")
	ErrUnauthorized   = response.NewUnauthorizedError("未授权")
	ErrForbidden      = response.NewForbiddenError("禁止访问")
)

// 数据库相关错误
var (
	ErrDatabaseConnection = response.NewDatabaseConnectionError("数据库连接失败")
	ErrRecordNotFound     = response.NewRecordNotFoundError("记录不存在")
	ErrDuplicateKey       = response.NewDuplicateKeyError("重复键值")
)

// 外部服务相关错误
var (
	ErrExternalServiceUnavailable = response.NewDomainError(response.ErrorTypeExternalServiceUnavailable, "外部服务不可用")
	ErrTimeout                    = response.NewTimeoutError("请求超时")
	ErrNetworkError               = response.NewNetworkError("网络错误")
)
