package workspace

import (
	"github.com/google/uuid"
	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
)


type WorkspaceService struct {
	workspaceRepo   WorkspaceRepository
	membershipRepo  MembershipRepository
	invitationRepo  InvitationRepository
}

func (s *WorkspaceService) CreateWorkspace(name, ownerID string) (*Workspace, domain_errors.DomainError) {
	// check for valid name and ownerID
	if name == "" {
		return nil, domain_errors.NewValidationErrorWithValue("name", name, "EMPTY WORKSPACE NAME")
	}

	if err := uuid.Validate(ownerID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("ownerID", ownerID, "OWNER ID IS NOT A VALID UUID")
	}

	workspace := &Workspace{
		ID:      uuid.NewString(),
		Name:    name,
		OwnerID: ownerID,
	}
	return s.workspaceRepo.Create(workspace)
}

func (s *WorkspaceService) GetWorkspaceByID(id string) (*Workspace, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("workspaceID", id, "WorkspaceID is not a valid UUID")
	}
	return s.workspaceRepo.GetByID(id)
}

func (s *WorkspaceService) UpdateWorkspace(workspace *Workspace) (*Workspace, domain_errors.DomainError) {
	if workspace.Name == "" {
		return nil, domain_errors.NewValidationErrorWithValue("name", workspace.Name, "EMPTY WORKSPACE NAME")
	}
	return s.workspaceRepo.Update(workspace)
}

func (s *WorkspaceService) DeleteWorkspace(id string) domain_errors.DomainError {
	if err := uuid.Validate(id); err != nil {
		return domain_errors.NewValidationErrorWithValue("workspaceID", id, "WorkspaceID is not a valid UUID")
	}
	return s.workspaceRepo.Delete(id)
}

func (s *WorkspaceService) ListWorkspacesByOwner(ownerID string) ([]*Workspace, domain_errors.DomainError) {
	if err := uuid.Validate(ownerID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("ownerID", ownerID, "OWNER ID IS NOT A VALID UUID")
	}
	return s.workspaceRepo.ListByOwner(ownerID)

}


// MEMBERSHIP FUNCTIONS

func (s *WorkspaceService) AddMembership(membership *Membership) (*Membership, error) {
	return s.membershipRepo.Add(membership)
}
func (s *WorkspaceService) RemoveMembership(userID, organizationID string) error {
	return s.membershipRepo.Remove(userID, organizationID)
}

func (s *WorkspaceService) ListMembershipsByOrganization(organizationID string) ([]*Membership, error) {
	return s.membershipRepo.ListByOrganization(organizationID)
}

func (s *WorkspaceService) CreateInvitation(invitation *Invitation) (*Invitation, error) {
	return s.invitationRepo.Create(invitation)
}

func (s *WorkspaceService) GetInvitationByToken(token string) (*Invitation, error) {
	return s.invitationRepo.GetByToken(token)
}

func (s *WorkspaceService) DeleteInvitation(token string) error {
	return s.invitationRepo.Delete(token)
}