package workspace

import (
	"github.com/go-chi/chi/v5"
	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"
	"github.com/ishola-faazele/taskflow/internal/shared"
)

func RegisterRoutes(r chi.Router, as *shared.AppState) {
	dm := domain_middleware.NewDomainMiddleware()
	handler := NewWorkspaceHandler(as.DB, as.AmqpConn)
	r.Use(dm.Authenticate)

	// Workspace routes
	r.Post("/", handler.CreateWorkspace)
	r.Get("/mine", handler.ListWorkspaces)
	r.Put("/{id}", handler.UpdateWorkspace)
	r.Get("/{id}", handler.GetWorkspace)
	r.Delete("/{id}", handler.DeleteWorkspace)

	// Invitation routes
	r.Post("/invitation", handler.CreateInvitation)
	r.Get("/invitation", handler.ListWorkspaceInvitations)
	r.Get("/invitation/{id}", handler.GetInvitation)
	r.Delete("/invitation/{id}", handler.DeleteInvitation)

	// Membership routes
	r.Get("/membership/add", handler.AddMembership)
	r.Get("/membership", handler.ListWorkspaceMembers)
	r.Post("/membership/remove", handler.RemoveMembership)
}
