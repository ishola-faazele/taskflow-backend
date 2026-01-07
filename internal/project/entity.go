package project

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
)

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WorkspaceID string    `json:"workspace_id"`
	Creator     string    `json:"creator"`
	CreatedAt   time.Time `json:"created_at"`
}
type UpdateProjectInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// TASK TYPES

type TaskStatus string

const (
	TaskStatusOpen     TaskStatus = "open"
	TaskStatusInReview TaskStatus = "in_review"
	TaskStatusClosed   TaskStatus = "closed"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          string       `json:"id"`
	ParentID    *string      `json:"parent_id"`
	ProjectID   string       `json:"project_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Creator     string       `json:"creator"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     time.Time    `json:"due_date"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// TaskTree represents a task with its subtasks (nested structure)
type TaskTree struct {
	Task
	Subtasks []*TaskTree `json:"subtasks,omitempty"`
	Depth    int         `json:"depth"`
}

// TaskWithChildren represents a task with immediate children only
type TaskWithChildren struct {
	Task
	Children []*Task `json:"children,omitempty"`
}

type TaskAssignment struct {
	ID        string    `json:"id"`
	Assigner  string    `json:"assigner"`
	Assignee  string    `json:"assignee"`
	TaskID    string    `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskComment struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	TaskID    string    `json:"task_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTaskInput struct {
	ProjectID   string
	ParentID    string
	Name        string
	Description string
	Creator     string
	Status      TaskStatus
	Priority    TaskPriority
	DueDate     time.Time
}

func (taskinput *CreateTaskInput) Validate() domain_errors.DomainError {
	if err := uuid.Validate(taskinput.ProjectID); err != nil {
		return domain_errors.NewValidationErrorWithValue("project_id", taskinput.ProjectID, "PROJECT ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(taskinput.ParentID); err != nil && taskinput.ParentID != "" {
		return domain_errors.NewValidationErrorWithValue("parent_id", taskinput.ParentID, "PARENT ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(taskinput.Creator); err != nil {
		return domain_errors.NewValidationErrorWithValue("creator", taskinput.Creator, "CREATOR IS NOT A VALID UUID")
	}
	if taskinput.Name == "" {
		return domain_errors.NewValidationError("name", "NAME CANNOT BE EMPTY")
	}
	if taskinput.Description == "" {
		return domain_errors.NewValidationError("description", "DESCRIPTION CANNOT BE EMPTY")
	}
	if taskinput.Status == "" {
		return domain_errors.NewValidationError("status", "STATUS CANNOT BE EMPTY")
	}
	if taskinput.Priority == "" {
		return domain_errors.NewValidationError("priority", "PRIORITY CANNOT BE EMPTY")
	}
	if taskinput.DueDate.IsZero() {
		return domain_errors.NewValidationError("due_date", "DUE DATE CANNOT BE EMPTY")
	}
	return nil
}

type UpdateTaskInput struct {
	Name        *string       `json:"name"`
	Description *string       `json:"description"`
	Status      *TaskStatus   `json:"status"`
	Priority    *TaskPriority `json:"priority"`
	DueDate     *time.Time    `json:"due_date"`
}

func (taskinput *UpdateTaskInput) Validate() domain_errors.DomainError {
	fmt.Printf("update input: +%v\n", taskinput)
	if taskinput.Name == nil || taskinput.Description == nil || taskinput.DueDate == nil || taskinput.Priority == nil || taskinput.Status == nil {
		return domain_errors.NewValidationError("input", "ONE OR MORE OF THE OF INPUTS IS EMPTY")
	}
	return nil
}
