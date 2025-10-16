package http

import (
	"common/interfaces"
	commonMiddleware "common/middleware"
	service "services/internal/application/service"
	"services/internal/interfaces/http/routes"
)

// NewCasbinMiddleware 创建 Casbin 中间件的 Provider
func NewCasbinMiddleware(permissionService service.PermissionServiceInterface) routes.CasbinMiddleware {
	return routes.CasbinMiddleware(commonMiddleware.CasbinMiddleware(permissionService.Enforce))
}

// NewAuthMiddleware 创建 Auth 中间件的 Provider
func NewAuthMiddleware(jwtService interfaces.JWTService, configProvider interfaces.ConfigProvider, logger interfaces.Logger) routes.AuthMiddleware {
	authConfig := configProvider.GetAuthConfig()
	return routes.AuthMiddleware(commonMiddleware.AuthMiddleware(jwtService, authConfig, logger))
}