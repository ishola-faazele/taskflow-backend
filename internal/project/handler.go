package project

import (
	"database/sql"
	"encoding/json"
	"net/http"

	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"
	"github.com/ishola-faazele/taskflow/pkg/utils"
)

type ProjectHandler struct {
	service   *ProjectService
	responder *utils.APIResponder
}

func NewProjectHandler(db *sql.DB) *ProjectHandler {
	service := NewProjectService(NewPostgresProjectRepository(db))
	responder := utils.NewAPIResponder()
	return &ProjectHandler{
		service:   service,
		responder: responder,
	}
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WorkspaceID string `json:"workspace_id"`
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	ctx := r.Context()
	ownerID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || ownerID == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}
	project, err := h.service.Create(req.Name, req.Description, req.WorkspaceID, ownerID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to create project", err)
		return
	}
	location := "/api/project/" + project.ID
	h.responder.Success(w, r, http.StatusCreated, location, project)
}
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := h.service.GetByID(id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to get project", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "", project)
}
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var req UpdateProjectInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	id := r.PathValue("id")
	project, err := h.service.Update(&req, id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to update project", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "", project)
}
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.service.Delete(id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to delete project", err)
		return
	}
	h.responder.NoContent(w)
}

func (h *ProjectHandler) ListProjectsByWorkspace(w http.ResponseWriter, r *http.Request) {
	wsID := r.URL.Query().Get("ws")
	ctx := r.Context()
	requester, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || requester == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}
	projects, err := h.service.ListByWorkspace(wsID, requester)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to list projects", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "", projects)
}
