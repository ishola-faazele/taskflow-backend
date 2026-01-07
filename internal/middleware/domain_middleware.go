package middleware

import (
	"github.com/ishola-faazele/taskflow/internal/utils/jwt"
	workspace "github.com/ishola-faazele/taskflow/internal/workspace/service"
	"github.com/ishola-faazele/taskflow/pkg/utils/domain_errors"
)

type DomainMiddleware struct {
	jwt              *jwt.JWTUtils
	responder        *domain_errors.APIResponder
	WorkspaceService *workspace.WorkspaceService
}

func NewDomainMiddleware() *DomainMiddleware {
	responder := domain_errors.NewAPIResponder()
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
