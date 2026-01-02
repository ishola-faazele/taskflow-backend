package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	shared_db "github.com/ishola-faazele/taskflow/internal/shared/db"
	"github.com/ishola-faazele/taskflow/internal/workspace"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
    db, err := sqlx.Connect("pgx", "user=taskflow_user password=taskflow_password dbname=taskflow_db sslmode=disable port=5432 host=localhost")
    
    if err != nil {
        log.Fatalln(err)
    }
    defer db.Close()
    migrationMgr := shared_db.NewMigrationManager(db.DB)
    if err := migrationMgr.EnsureTablesExist(); err != nil {
        log.Fatalln("Failed to ensure tables exist:", err)
    }

    r := chi.NewRouter()
    r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	workspaceRouter := &workspace.WorkspaceRouter{DB: db.DB}
	r.Mount("/api", workspaceRouter.RegisterRoutes())
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
