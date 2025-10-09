package service

import (
	"context"
)

// PermissionServiceInterface 权限服务接口
// 定义在common包中，以便common/middleware可以引用
type PermissionServiceInterface interface {
	// Enforce 检查权限
	Enforce(ctx context.Context, sub, obj, act string) (bool, error)
}
