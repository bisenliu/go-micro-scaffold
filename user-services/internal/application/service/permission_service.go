package service

import (
	"context"

	"github.com/casbin/casbin/v2"
)

// PermissionServiceInterface 权限服务接口
type PermissionServiceInterface interface {
	// Enforce 检查权限
	Enforce(ctx context.Context, sub, obj, act string) (bool, error)
	// AddPolicy 添加策略
	AddPolicy(ctx context.Context, sub, obj, act string) (bool, error)
	// AddRoleForUser 为用户添加角色
	AddRoleForUser(ctx context.Context, user, role string) (bool, error)
}

// PermissionService 权限服务
type PermissionService struct {
	enforcer *casbin.SyncedCachedEnforcer
}

// NewPermissionService 创建权限服务
func NewPermissionService(enforcer *casbin.SyncedCachedEnforcer) PermissionServiceInterface {
	return &PermissionService{
		enforcer: enforcer,
	}
}

// Enforce 检查权限
func (s *PermissionService) Enforce(ctx context.Context, sub, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加策略
func (s *PermissionService) AddPolicy(ctx context.Context, sub, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(sub, obj, act)
}

// AddRoleForUser 为用户添加角色
func (s *PermissionService) AddRoleForUser(ctx context.Context, user, role string) (bool, error) {
	return s.enforcer.AddRoleForUser(user, role)
}
