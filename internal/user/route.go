package user

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
}

func RegisterRoutes(r chi.Router, DB *sql.DB) {
	handler := NewUserHandler(DB)
	r.Post("/magic-link", handler.RequestMagicLink)
	r.Get("/verify", handler.VerifyToken)
}
