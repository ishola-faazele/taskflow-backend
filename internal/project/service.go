package project

import (
	"time"

	"github.com/google/uuid"
	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
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
