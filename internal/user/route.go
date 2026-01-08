package user

import (
	"github.com/go-chi/chi/v5"
	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"
	"github.com/ishola-faazele/taskflow/internal/shared"
)

func RegisterRoutes(r chi.Router, as *shared.AppState) {
	dm := domain_middleware.NewDomainMiddleware()
	handler := NewUserHandler(as.DB, as.AmqpConn)

	// Public routes (no authentication required)
	r.Post("/magic-link", handler.RequestMagicLink)
	r.Get("/verify", handler.VerifyToken)
	r.Get("/refresh-token", handler.RefreshToken)

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(dm.Authenticate)

		r.Get("/auth", handler.GetByID)
		r.Get("/profile", handler.GetProfile)
		r.Put("/profile", handler.UpdateProfile)
		r.Get("/profile/{id}", handler.GetPublicProfile)
	})
}
