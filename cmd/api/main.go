package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ishola-faazele/taskflow/internal/project"
	shared "github.com/ishola-faazele/taskflow/internal/shared"
	"github.com/ishola-faazele/taskflow/internal/user"
	workspace "github.com/ishola-faazele/taskflow/internal/workspace/http"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Warning: .env file not found, using system environment variables")
	}
	appState := shared.NewAppState()
	defer appState.Clean()

	// mount routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiRouter := chi.NewRouter()
	apiRouter.Route("/workspace", func(r chi.Router) {
		workspace.RegisterRoutes(r, appState)
	})
	apiRouter.Route("/user", func(r chi.Router) {
		user.RegisterRoutes(r, appState)
	})
	apiRouter.Route("/workspace/{ws_id}/project", func(r chi.Router) {
		project.RegisterProjectRoutes(r, appState.DB)
	})
	apiRouter.Route("/workspace/{ws_id}/task", func(r chi.Router) {
		project.RegisterTaskRoutes(r, appState.DB)
	})

	r.Mount("/api", apiRouter)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})
	log.Println("Server is Running")
	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}
}
