package project

import (
	"database/sql"

	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, DB *sql.DB) {
	dm := domain_middleware.NewDomainMiddleware()
	handler := NewProjectHandler(DB)
	r.Use(dm.Authenticate)

	// Project routes
	r.Post("/", handler.CreateProject)
	r.Get("/", handler.ListProjectsByWorkspace)
	r.Get("/{id}", handler.GetProject)
	r.Put("/{id}", handler.UpdateProject)
	r.Delete("/{id}", handler.DeleteProject)
}
