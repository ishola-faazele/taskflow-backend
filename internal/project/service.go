package project

import (
	"time"

	"github.com/google/uuid"
	"github.com/ishola-faazele/taskflow/pkg/utils/domain_errors"
)

type ProjectService struct {
	projectRepo ProjectRepository
}

func NewProjectService(pjRepo ProjectRepository) *ProjectService {
	return &ProjectService{
		projectRepo: pjRepo,
	}
}

func (pjs *ProjectService) Create(name, desc, ws_id, creator string) (*Project, error) {
	project := &Project{
		ID:          uuid.NewString(),
		Name:        name,
		Description: desc,
		WorkspaceID: ws_id,
		Creator:     creator,
		CreatedAt:   time.Now().UTC(),
	}
	return pjs.projectRepo.Create(project)
}

func (pjs *ProjectService) GetByID(id string) (*Project, error) {
	// validate id
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("id", id, "PROJECT ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.GetByID(id)
}

func (pjs *ProjectService) Update(input *UpdateProjectInput, id string) (*Project, error) {
	// validate input
	if input.Name == nil && input.Description == nil {
		return nil, domain_errors.NewValidationError("input", "NO FIELDS TO UPDATE")
	}
	// validate id
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("id", id, "PROJECT ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.Update(input, id)
}

func (pjs *ProjectService) Delete(id string) error {
	// validate id
	if err := uuid.Validate(id); err != nil {
		return domain_errors.NewValidationErrorWithValue("id", id, "PROJECT ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.Delete(id)
}
func (pjs *ProjectService) ListByWorkspace(wsID, requester string) ([]*Project, error) {
	// validate wsID
	if err := uuid.Validate(wsID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("workspace_id", wsID, "WORKSPACE ID IS NOT A VALID UUID")
	}
	// check requester membership in workspace could be added here
	return pjs.projectRepo.ListByWorkspace(wsID)
}

// ============================================================================
// TASK METHODS
// ============================================================================

func (pjs *ProjectService) CreateTask(input *CreateTaskInput) (*Task, domain_errors.DomainError) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}
	task := &Task{
		ID:          uuid.NewString(),
		ProjectID:   input.ProjectID,
		ParentID:    &input.ParentID,
		Name:        input.Name,
		Description: input.Description,
		Creator:     input.Creator,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	var parentID *string
	if input.ParentID != "" {
		parentID = &input.ParentID
	}
	task.ParentID = parentID
	return pjs.projectRepo.CreateTask(task)
}
func (pjs *ProjectService) GetTaskByID(id string) (*Task, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("id", id, "TASK ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.GetTaskByID(id)
}

func (pjs *ProjectService) UpdateTask(input *UpdateTaskInput, id string) (*Task, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("task_id", id, "TASK ID IS NOT A VALID UUID")
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return pjs.projectRepo.UpdateTask(input, id)
}

func (pjs *ProjectService) DeleteTask(id string) domain_errors.DomainError {
	if err := uuid.Validate(id); err != nil {
		return domain_errors.NewValidationErrorWithValue("task_id", id, "TASK ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.DeleteTask(id)
}

// Tree queries

// Returns the immediate children of a parent node
func (pjs *ProjectService) ListSubtasks(parentID string) ([]*Task, domain_errors.DomainError) {
	if err := uuid.Validate(parentID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("parent_id", parentID, "PARENT ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.ListSubtasks(parentID)
}

// Returns a task and all its descendants
func (pjs *ProjectService) GetTaskTree(rootID string) (*TaskTree, domain_errors.DomainError) {
	if err := uuid.Validate(rootID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("root_id", rootID, "ROOT ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.GetTaskTree(rootID)
}

// Returns a task and its immediate childres in a tree
func (pjs *ProjectService) GetTaskWithChildren(id string) (*TaskWithChildren, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("task_id", id, "TASK ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.GetTaskWithChildren(id)
}

// Returns all tasks without parents
func (pjs *ProjectService) GetRootTasks(projectID string) ([]*Task, domain_errors.DomainError) {
	if err := uuid.Validate(projectID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("project_id", projectID, "PROJECT ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.GetRootTasks(projectID)
}

// Project queries

// Returns all tasks in a project in a flatlist
func (pjs *ProjectService) ListTasksByProject(projectID string) ([]*Task, domain_errors.DomainError) {
	if err := uuid.Validate(projectID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("project_id", projectID, "PROJECT ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.ListTasksByProject(projectID)
}

// Returns all tasks in a project in a tree
func (pjs *ProjectService) GetProjectTaskTree(projectID string) ([]*TaskTree, domain_errors.DomainError) {
	if err := uuid.Validate(projectID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("project_id", projectID, "PROJECT ID IS NOT A VALID UUID")

	}
	return pjs.projectRepo.GetProjectTaskTree(projectID)
}

// Utility

func (pjs *ProjectService) GetTaskDepth(id string) (int, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return 0, domain_errors.NewValidationErrorWithValue("task_id", id, "TASK ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.GetTaskDepth(id)
}

func (pjs *ProjectService) CountSubtasks(parentID string) (int, domain_errors.DomainError) {
	if err := uuid.Validate(parentID); err != nil {
		return 0, domain_errors.NewValidationErrorWithValue("task_id", parentID, "TASK ID IS NOT A VALID UUID")
	}
	return pjs.projectRepo.CountSubtasks(parentID)
}
