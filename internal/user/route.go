package user

import (
	"database/sql"

	"github.com/go-chi/chi"
)

type UserRouter struct {
	DB *sql.DB
}

func (ur *UserRouter) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	handler := NewUserHandler(ur.DB)

	r.Route("/user", func(r chi.Router) {
		r.Post("/magic-link", handler.RequestMagicLink)
		r.Get("/verify", handler.VerifyToken)
	})

	return r
}
