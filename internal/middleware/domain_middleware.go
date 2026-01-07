package middleware

import (
	"github.com/ishola-faazele/taskflow/internal/shared/jwt"
	workspace "github.com/ishola-faazele/taskflow/internal/workspace/service"
	"github.com/ishola-faazele/taskflow/pkg/utils"
)

type DomainMiddleware struct {
	jwt              *jwt.JWTUtils
	responder        *utils.APIResponder
	WorkspaceService *workspace.WorkspaceService
}

func NewDomainMiddleware() *DomainMiddleware {
	responder := utils.NewAPIResponder()
	return &DomainMiddleware{
		jwt:       jwt.NewJWTUtils(jwt.DefaultTokenConfig()),
		responder: responder,
	}
}
func NewDomainMiddlewareWithWorkspace(service *workspace.WorkspaceService) *DomainMiddleware {
	mid := NewDomainMiddleware()
	mid.WorkspaceService = service
	return mid
}
