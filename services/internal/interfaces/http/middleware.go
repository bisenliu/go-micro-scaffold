package http

import (
	"common/config"
	commonMiddleware "common/middleware"
	"common/pkg/jwt"
	service "services/internal/application/service"
	"services/internal/interfaces/http/routes"
)

// NewCasbinMiddleware 创建 Casbin 中间件的 Provider
func NewCasbinMiddleware(permissionService service.PermissionServiceInterface) routes.CasbinMiddleware {
	return routes.CasbinMiddleware(commonMiddleware.CasbinMiddleware(permissionService.Enforce))
}

// NewAuthMiddleware 创建 Auth 中间件的 Provider
func NewAuthMiddleware(jwtService *jwt.JWT, config *config.Config) routes.AuthMiddleware {
	return routes.AuthMiddleware(commonMiddleware.AuthMiddleware(jwtService, config.Auth))
}