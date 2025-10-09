package service

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	commonService "common/service"
)

// PermissionService 权限服务
type PermissionService struct {
	enforcer *casbin.SyncedCachedEnforcer
	logger   *zap.Logger
}

// 确保PermissionService实现了common/service包中的PermissionServiceInterface接口
var _ commonService.PermissionServiceInterface = (*PermissionService)(nil)

// NewPermissionService 创建权限服务实例
func NewPermissionService(enforcer *casbin.SyncedCachedEnforcer, logger *zap.Logger) *PermissionService {
	return &PermissionService{
		enforcer: enforcer,
		logger:   logger,
	}
}

// AddPolicy 添加权限策略
func (s *PermissionService) AddPolicy(ctx context.Context, sub, obj, act string) error {
	s.logger.Info("Adding policy",
		zap.String("subject", sub),
		zap.String("object", obj),
		zap.String("action", act))

	added, err := s.enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		s.logger.Error("Failed to add policy",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act),
			zap.Error(err))
		return fmt.Errorf("failed to add policy: %w", err)
	}

	if added {
		s.logger.Info("Policy added successfully",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
	} else {
		s.logger.Info("Policy already exists",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
	}

	return nil
}

// RemovePolicy 移除权限策略
func (s *PermissionService) RemovePolicy(ctx context.Context, sub, obj, act string) error {
	s.logger.Info("Removing policy",
		zap.String("subject", sub),
		zap.String("object", obj),
		zap.String("action", act))

	removed, err := s.enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		s.logger.Error("Failed to remove policy",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act),
			zap.Error(err))
		return fmt.Errorf("failed to remove policy: %w", err)
	}

	if removed {
		s.logger.Info("Policy removed successfully",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
	} else {
		s.logger.Info("Policy not found",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
	}

	return nil
}

// Enforce 检查权限
func (s *PermissionService) Enforce(ctx context.Context, sub, obj, act string) (bool, error) {
	s.logger.Debug("Enforcing policy",
		zap.String("subject", sub),
		zap.String("object", obj),
		zap.String("action", act))

	allowed, err := s.enforcer.Enforce(sub, obj, act)
	if err != nil {
		s.logger.Error("Failed to enforce policy",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act),
			zap.Error(err))
		return false, fmt.Errorf("failed to enforce policy: %w", err)
	}

	s.logger.Debug("Policy enforcement result",
		zap.String("subject", sub),
		zap.String("object", obj),
		zap.String("action", act),
		zap.Bool("allowed", allowed))

	return allowed, nil
}

// AddRoleForUser 为用户添加角色
func (s *PermissionService) AddRoleForUser(ctx context.Context, user, role string) error {
	s.logger.Info("Adding role for user",
		zap.String("user", user),
		zap.String("role", role))

	added, err := s.enforcer.AddGroupingPolicy(user, role)
	if err != nil {
		s.logger.Error("Failed to add role for user",
			zap.String("user", user),
			zap.String("role", role),
			zap.Error(err))
		return fmt.Errorf("failed to add role for user: %w", err)
	}

	if added {
		s.logger.Info("Role added for user successfully",
			zap.String("user", user),
			zap.String("role", role))
	} else {
		s.logger.Info("Role already exists for user",
			zap.String("user", user),
			zap.String("role", role))
	}

	return nil
}

// GetRolesForUser 获取用户的角色
func (s *PermissionService) GetRolesForUser(ctx context.Context, user string) ([]string, error) {
	s.logger.Debug("Getting roles for user", zap.String("user", user))

	roles, err := s.enforcer.GetRolesForUser(user)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetPermissionsForUser 获取用户的权限
func (s *PermissionService) GetPermissionsForUser(ctx context.Context, user string) [][]string {
	s.logger.Debug("Getting permissions for user", zap.String("user", user))

	permissions, err := s.enforcer.GetPermissionsForUser(user)
	if err != nil {
		s.logger.Error("Failed to get permissions for user", zap.String("user", user), zap.Error(err))
		return nil
	}
	return permissions
}
