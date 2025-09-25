package errors

import "errors"

// 基础设施层通用错误
var (
	ErrInternalServer = errors.New("内部服务器错误")
	ErrInvalidRequest = errors.New("无效的请求")
	ErrUnauthorized   = errors.New("未授权")
	ErrForbidden      = errors.New("禁止访问")
)

// 数据库相关错误
var (
	ErrDatabaseConnection = errors.New("数据库连接失败")
	ErrRecordNotFound     = errors.New("记录不存在")
	ErrDuplicateKey       = errors.New("重复键值")
)

// 外部服务相关错误
var (
	ErrExternalServiceUnavailable = errors.New("外部服务不可用")
	ErrTimeout                    = errors.New("请求超时")
	ErrNetworkError               = errors.New("网络错误")
)
