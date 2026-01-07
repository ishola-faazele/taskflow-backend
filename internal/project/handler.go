package project

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

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
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	wsID := r.PathValue("ws_id")
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
	project, err := h.service.Create(req.Name, req.Description, wsID, ownerID)
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
	wsID := r.PathValue("ws_id")
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

// ============================================================================
// TASK METHODS
// ============================================================================
type CreateTaskDTO struct {
	ProjectID   string       `json:"project_id"`
	ParentID    string       `json:"parent_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     time.Time    `json:"due_date"`
}

func (dto *CreateTaskDTO) CreateTaskInput() *CreateTaskInput {

	return &CreateTaskInput{
		ProjectID:   dto.ProjectID,
		ParentID:    dto.ParentID,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      dto.Status,
		Priority:    dto.Priority,
		DueDate:     dto.DueDate,
	}
}

func (h *ProjectHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	ctx := r.Context()
	creator, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || creator == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}

	taskInput := req.CreateTaskInput()
	taskInput.Creator = creator

	task, err := h.service.CreateTask(taskInput)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to Create Task", err)
		return
	}
	h.responder.Success(w, r, http.StatusCreated, "Task Created successfully", task)
}

func (h *ProjectHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to get task", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Task Retrieved Successfully", task)
}

func (h *ProjectHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	task, err := h.service.UpdateTask(&req, id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to get task", err)
		return
	}
	h.responder.Success(w, r, http.StatusAccepted, "Task Updated Successfully", task)
}

func (h *ProjectHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.service.DeleteTask(id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_DELETE_TASK", err)
		return
	}
	h.responder.NoContent(w)
}

// Tree queries
func (h *ProjectHandler) ListSubtasks(w http.ResponseWriter, r *http.Request) {
	parentID := r.PathValue("id")
	subtasks, err := h.service.ListSubtasks(parentID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_GET_SUBTASKS", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Subtasks Retrieved Successfully", subtasks)
}

func (h *ProjectHandler) GetTaskTree(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	taskTree, err := h.service.GetTaskTree(id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_GET_TASK_TREE", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Task Tree Retrieved Successfully", taskTree)
}

func (h *ProjectHandler) GetTaskWithChildren(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	taskWithChildren, err := h.service.GetTaskWithChildren(id)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_GET_TASK_WITH_CHILDREN", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Task With Children Retrieved Successfully", taskWithChildren)
}

func (h *ProjectHandler) GetRootTasks(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	rootTasks, err := h.service.GetRootTasks(projectID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_GET_ROOT_TASKS", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Root Tasks Retrieved Successfully", rootTasks)
}

// Project queries
// Lists all root task in a project in a flatlist
func (h *ProjectHandler) ListTasksByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	tasks, err := h.service.ListTasksByProject(projectID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_LIST_TASKS_BY_PROJECT", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Tasks Retrieved Successfully", tasks)
}

// List task in project in tree structure
func (h *ProjectHandler) GetProjectTaskTree(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	taskTree, err := h.service.GetProjectTaskTree(projectID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_LIST_TASKTREE_BY_PROJECT", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Tasks Retrieved Successfully", taskTree)
}
