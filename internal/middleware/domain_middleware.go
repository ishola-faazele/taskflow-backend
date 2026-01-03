package middleware

import (
	"github.com/ishola-faazele/taskflow/internal/shared/jwt"
	"github.com/ishola-faazele/taskflow/pkg/utils"
)

type DomainMiddleware struct {
	jwt       *jwt.JWTUtils
	responder *utils.APIResponder
}

func NewDomainMiddleware() *DomainMiddleware {
	responder := utils.NewAPIResponder()
	return &DomainMiddleware{
		jwt:       jwt.NewJWTUtils(jwt.DefaultTokenConfig()),
		responder: responder,
	}
}
