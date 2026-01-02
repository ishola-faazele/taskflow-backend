package workspace

import (
	"database/sql"

	"github.com/go-chi/chi"
)
type WorkspaceRouter struct {
	DB *sql.DB 

}
func(wr *WorkspaceRouter) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	handler, err := NewWorkspaceHandler(wr.DB)
	if err != nil {
		panic(err)
	}
	r.Route("/workspaces", func(r chi.Router) {
		r.Post("/", handler.CreateWorkspace)
		r.Get("/mine", handler.ListWorkspaces)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetWorkspace)
			r.Put("/", handler.UpdateWorkspace)
			r.Delete("/", handler.DeleteWorkspace)
		})
	})

  	return r
  	
}
