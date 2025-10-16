package errors

import (
	"common/response"
	domainerrors "services/internal/domain/shared/errors"
)

// 基础设施层通用错误
var (
	ErrInternalServer = domainerrors.NewInternalServerError("内部服务器错误")
	ErrInvalidRequest = domainerrors.NewDomainError(response.ErrorTypeInvalidRequest, "无效的请求")
	ErrUnauthorized   = domainerrors.NewUnauthorizedError("未授权")
	ErrForbidden      = domainerrors.NewForbiddenError("禁止访问")
)

// 数据库相关错误
var (
	ErrDatabaseConnection = domainerrors.NewDatabaseConnectionError("数据库连接失败")
	ErrRecordNotFound     = domainerrors.NewRecordNotFoundError("记录不存在")
	ErrDuplicateKey       = domainerrors.NewDuplicateKeyError("重复键值")
)

// 外部服务相关错误
var (
	ErrExternalServiceUnavailable = domainerrors.NewDomainError(response.ErrorTypeExternalServiceUnavailable, "外部服务不可用")
	ErrTimeout                    = domainerrors.NewTimeoutError("请求超时")
	ErrNetworkError               = domainerrors.NewNetworkError("网络错误")
)
