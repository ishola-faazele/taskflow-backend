package project

import "github.com/ishola-faazele/taskflow/internal/shared/domain_errors"

type ProjectRepository interface {
	Create(project *Project) (*Project, domain_errors.DomainError)
	GetByID(id string) (*Project, domain_errors.DomainError)
	Update(input *UpdateProjectInput, id string) (*Project, domain_errors.DomainError)
	Delete(id string) domain_errors.DomainError
	ListByWorkspace(wsID string) ([]*Project, domain_errors.DomainError)

	// methods for tasks
	// Basic CRUD
	CreateTask(task *Task) (*Task, domain_errors.DomainError)
	GetTaskByID(id string) (*Task, domain_errors.DomainError)
	UpdateTask(input *UpdateTaskInput, id string) (*Task, domain_errors.DomainError)
	DeleteTask(id string) domain_errors.DomainError

	// Tree queries
	ListSubtasks(parentID string) ([]*Task, domain_errors.DomainError)
	GetTaskTree(rootID string) (*TaskTree, domain_errors.DomainError)
	GetTaskWithChildren(id string) (*TaskWithChildren, domain_errors.DomainError)
	GetRootTasks(projectID string) ([]*Task, domain_errors.DomainError)

	// Project queries
	ListTasksByProject(projectID string) ([]*Task, domain_errors.DomainError)
	GetProjectTaskTree(projectID string) ([]*TaskTree, domain_errors.DomainError)

	// Utility
	GetTaskDepth(id string) (int, domain_errors.DomainError)
	CountSubtasks(parentID string) (int, domain_errors.DomainError)
}
