package workspace

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"
	"github.com/ishola-faazele/taskflow/pkg/utils"
)

type WorkspaceHandler struct {
	service   *WorkspaceService
	responder *utils.APIResponder
}

func NewWorkspaceHandler(db *sql.DB) *WorkspaceHandler {
	workspaceRepo := NewPostgresWorkspaceRepository(db)
	service := &WorkspaceService{workspaceRepo: workspaceRepo}
	responder := utils.NewAPIResponder()

	return &WorkspaceHandler{
		service:   service,
		responder: responder,
	}
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}
type UpdateWorkspaceRequest struct {
	Name string `json:"name"`
}
// CreateWorkspace handles workspace creation
func (h *WorkspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkspaceRequest
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
	workspace, err := h.service.CreateWorkspace(req.Name, ownerID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to create workspace", err)
		return
	}

	location := "/api/workspace/" + workspace.ID
	h.responder.Created(w, r, location, workspace)
}

// GetWorkspace handles fetching a single workspace
func (h *WorkspaceHandler) GetWorkspace(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Println("Fetching workspace with ID:", id)

	workspace, err := h.service.GetWorkspaceByID(id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to retrieve workspace", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Workspace retrieved successfully", workspace)
}

// UpdateWorkspace handles workspace updates
func (h *WorkspaceHandler) UpdateWorkspace(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req *UpdateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// change this to user making request
	requester, ok := r.Context().Value(domain_middleware.UserIDKey).(string)
	if !ok || requester == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}
	updatedWorkspace, err := h.service.UpdateWorkspace(id, req.Name, requester)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to update workspace", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Workspace updated successfully", updatedWorkspace)
}

// DeleteWorkspace handles workspace deletion
func (h *WorkspaceHandler) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.responder.Error(w, r, http.StatusBadRequest, "Workspace ID is required", nil)
		return
	}
	ctx := r.Context()
	requester, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || requester == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}

	if err := h.service.DeleteWorkspace(id, requester); err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to delete workspace", err)
		return
	}

	h.responder.NoContent(w)
}

// ListWorkspaces handles listing workspaces by owner
func (h *WorkspaceHandler) ListWorkspaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ownerID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || ownerID == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}

	workspaces, err := h.service.ListWorkspacesByOwner(ownerID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to list workspaces", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Workspaces retrieved successfully", workspaces)
}
