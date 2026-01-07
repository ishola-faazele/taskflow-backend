package workspace

import (
	"database/sql"
	"encoding/json"
	"net/http"

	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"
	. "github.com/ishola-faazele/taskflow/internal/workspace/db"
	. "github.com/ishola-faazele/taskflow/internal/workspace/entity"
	. "github.com/ishola-faazele/taskflow/internal/workspace/service"
	"github.com/ishola-faazele/taskflow/pkg/utils/domain_errors"
)

type WorkspaceHandler struct {
	service   *WorkspaceService
	responder *domain_errors.APIResponder
}

func NewWorkspaceHandler(db *sql.DB) *WorkspaceHandler {
	workspaceRepo := NewPostgresWorkspaceRepository(db)
	invitationRepo := NewPostgresInvitationRepository(db)
	membershipRepo := NewPostgresMembershipRepository(db)
	service := NewWorkspaceService(workspaceRepo, invitationRepo, membershipRepo)
	responder := domain_errors.NewAPIResponder()

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
type CreateInvitationRequest struct {
	InviteeID    string `json:"invitee_id"`
	InviteeEmail string `json:"invitee_email"`
	WorkspaceID  string `json:"workspace_id"`
	Role         Role   `json:"role"`
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

// CreateInvitation handles invitation creation
func (h *WorkspaceHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	var req CreateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	ctx := r.Context()
	inviterID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || inviterID == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}
	invitation, err := h.service.CreateInvitation(req.InviteeID, inviterID, req.WorkspaceID, req.InviteeEmail, req.Role)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to create invitation", err)
		return
	}

	location := "/api/workspace/invitation/" + invitation.ID
	h.responder.Created(w, r, location, invitation)
}

// GetInvitation handles fetching a single invitation
func (h *WorkspaceHandler) GetInvitation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	invitation, err := h.service.GetInvitation(id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to retrieve invitation", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Invitation retrieved successfully", invitation)
}

func (h *WorkspaceHandler) DeleteInvitation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()
	requester, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || requester == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}

	if err := h.service.DeleteInvitation(id, requester); err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to delete invitation", err)
		return
	}

	h.responder.NoContent(w)
}

func (h *WorkspaceHandler) ListWorkspaceInvitations(w http.ResponseWriter, r *http.Request) {
	wsID := r.URL.Query().Get("ws")
	if wsID == "" {
		h.responder.Error(w, r, http.StatusBadRequest, "Workspace ID is required", nil)
		return
	}
	ctx := r.Context()
	requester, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || requester == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}

	invitations, err := h.service.ListWorkspaceInvitations(wsID, requester)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to list invitations", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Invitations retrieved successfully", invitations)
}

// / MEMBERSHIP HANDLERS
type RemoveMembershipRequest struct {
	UserID      string `json:"user_id"`
	WorkspaceID string `json:"workspace_id"`
}

// AddMembership handles adding a membership to a workspace
func (h *WorkspaceHandler) AddMembership(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
	token := r.URL.Query().Get("token")
	membership, err := h.service.AddMembership(token)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to add membership", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Membership added successfully", membership)
}

func (h *WorkspaceHandler) RemoveMembership(w http.ResponseWriter, r *http.Request) {
	var req RemoveMembershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	ctx := r.Context()
	requester, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || requester == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}

	if err := h.service.RemoveMembership(req.UserID, req.WorkspaceID, requester); err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to remove membership", err)
		return
	}

	h.responder.NoContent(w)
}
func (h *WorkspaceHandler) ListWorkspaceMembers(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
	wsID := r.URL.Query().Get("ws")
	if wsID == "" {
		h.responder.Error(w, r, http.StatusBadRequest, "Workspace ID is required", nil)
		return
	}
	ctx := r.Context()
	requester, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || requester == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}
	memberships, err := h.service.ListWorkspaceMembers(wsID, requester)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to list workspace members", err)
		return
	}

	h.responder.Success(w, r, http.StatusOK, "Workspace members retrieved successfully", memberships)
}
