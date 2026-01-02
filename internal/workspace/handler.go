package workspace

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ishola-faazele/taskflow/pkg/utils"

	"github.com/go-chi/chi"
)

type WorkspaceHandler struct {
	service   *WorkspaceService
	responder *utils.APIResponder
}

func NewWorkspaceHandler(db *sql.DB) (*WorkspaceHandler, error) {
	workspaceRepo := NewPostgresWorkspaceRepository(db)
	service := &WorkspaceService{workspaceRepo: workspaceRepo}
	responder := utils.NewAPIResponder()

	return &WorkspaceHandler{
		service:   service,
		responder: responder,
	}, nil
}

type CreateWorkspaceRequest struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

// CreateWorkspace handles workspace creation
func (h *WorkspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	workspace, err := h.service.CreateWorkspace(req.Name, req.OwnerID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to create workspace", err)
		return
	}

	location := "/api/workspace/" + workspace.ID
	h.responder.Created(w, r, location, workspace)
}

// GetWorkspace handles fetching a single workspace
func (h *WorkspaceHandler) GetWorkspace(w http.ResponseWriter, r *http.Request) {
	// may have to validate id param
	id := chi.URLParam(r, "id")
	workspace, err := h.service.GetWorkspaceByID(id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to retrieve workspace", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Workspace retrieved successfully", workspace)
}

// UpdateWorkspace handles workspace updates
func (h *WorkspaceHandler) UpdateWorkspace(w http.ResponseWriter, r *http.Request) {
	// id := chi.URLParam(r, "id")

	var req *Workspace
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// change this to user making request
	updatedWorkspace, err := h.service.UpdateWorkspace(req, req.OwnerID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to update workspace", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Workspace updated successfully", updatedWorkspace)
}

// DeleteWorkspace handles workspace deletion
func (h *WorkspaceHandler) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.responder.Error(w, r, http.StatusBadRequest, "Workspace ID is required", nil)
		return
	}

	if err := h.service.DeleteWorkspace(id); err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to delete workspace", err)
		return
	}

	h.responder.NoContent(w)
}

// ListWorkspaces handles listing workspaces by owner
func (h *WorkspaceHandler) ListWorkspaces(w http.ResponseWriter, r *http.Request) {
	ownerID := r.URL.Query().Get("owner_id")
	if ownerID == "" {
		h.responder.Error(w, r, http.StatusBadRequest, "owner_id parameter is required", nil)
		return
	}

	workspaces, err := h.service.ListWorkspacesByOwner(ownerID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to list workspaces", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Workspaces retrieved successfully", workspaces)
}
