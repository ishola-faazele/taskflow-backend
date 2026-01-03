package workspace

import (
	"database/sql"

	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, DB *sql.DB) {
	dm := domain_middleware.NewDomainMiddleware()
	handler := NewWorkspaceHandler(DB)
	r.Use(dm.Authenticate)

	r.Post("/", handler.CreateWorkspace)
	r.Get("/mine", handler.ListWorkspaces)
	r.Put("/{id}", handler.UpdateWorkspace)
	r.Get("/{id}", handler.GetWorkspace)
	r.Delete("/{id}", handler.DeleteWorkspace)
}
