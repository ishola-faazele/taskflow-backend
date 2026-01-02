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
	r.Route("/workspace", func(r chi.Router) {
		r.Post("/", handler.CreateWorkspace)
		r.Get("/mine", handler.ListWorkspaces)
		r.Put("/update", handler.UpdateWorkspace)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetWorkspace)
			r.Delete("/", handler.DeleteWorkspace)
		})
	})

  	return r
  	
}
