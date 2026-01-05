package project

import "github.com/ishola-faazele/taskflow/internal/shared/domain_errors"

type ProjectRepository interface {
	Create(project *Project) (*Project, domain_errors.DomainError)
	GetByID(id string) (*Project, domain_errors.DomainError)
	Update(input *UpdateProjectInput, id string) (*Project, domain_errors.DomainError)
	Delete(id string) domain_errors.DomainError
	ListByWorkspace(wsID string) ([]*Project, domain_errors.DomainError)
}
