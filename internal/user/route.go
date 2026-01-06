package user

import (
	"database/sql"

	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
}

func RegisterRoutes(r chi.Router, DB *sql.DB) {
	dm := domain_middleware.NewDomainMiddleware()
	handler := NewUserHandler(DB)

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
