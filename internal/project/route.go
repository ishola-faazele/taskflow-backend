package project

import (
	"database/sql"

	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"
	workspace_repository "github.com/ishola-faazele/taskflow/internal/workspace/db"
	workspace_service "github.com/ishola-faazele/taskflow/internal/workspace/service"

	"github.com/go-chi/chi/v5"
)

func RegisterProjectRoutes(r chi.Router, DB *sql.DB) {
	// Middleware
	workspaceService := workspace_service.WorkspaceService{
		MembershipRepo: workspace_repository.NewPostgresMembershipRepository(DB),
	}
	dm := domain_middleware.NewDomainMiddlewareWithWorkspace(&workspaceService)
	handler := NewProjectHandler(DB)
	r.Use(dm.Authenticate)
	r.Use(dm.CheckMembership)
	// Project routes
	r.Post("/", handler.CreateProject)
	r.Get("/all", handler.ListProjectsByWorkspace)
	r.Get("/{id}", handler.GetProject)
	r.Put("/{id}", handler.UpdateProject)
	r.Delete("/{id}", handler.DeleteProject)
}

func RegisterTaskRoutes(r chi.Router, DB *sql.DB) {
	workspaceService := workspace_service.WorkspaceService{
		MembershipRepo: workspace_repository.NewPostgresMembershipRepository(DB),
	}
	dm := domain_middleware.NewDomainMiddlewareWithWorkspace(&workspaceService)
	r.Use(dm.Authenticate)
	r.Use(dm.CheckMembership)
	handler := NewProjectHandler(DB)

	// BASIC CRUD APIS
	r.Post("/", handler.CreateTask)
	r.Get("/{id}", handler.GetTaskByID)
	r.Put("/{id}", handler.UpdateTask)
	r.Delete("/{id}", handler.DeleteTask)
	
	// TREE QUERIES
	r.Get("/{id}/subtasks", handler.ListSubtasks)
	r.Get("/{id}/tree", handler.GetTaskTree)
	r.Get("/{id}/children", handler.GetTaskWithChildren)
	r.Get("/{id}/root", handler.GetRootTasks)

	// PROJECT QUERIES
	r.Get("/{id}/project_tasks", handler.ListTasksByProject)
	r.Get("/{id}/project_tree", handler.GetProjectTaskTree)
	// UTILITY
	// r.Get("/task/{id}/depth", handler.GetTaskDepth)
	// r.Get("/task/{id}/count", handler.CountSubtasks)
}
