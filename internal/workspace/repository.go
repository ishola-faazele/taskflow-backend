package workspace

import "github.com/ishola-faazele/taskflow/internal/shared/domain_errors"

type WorkspaceRepository interface {
	Create(org *Workspace) (*Workspace, domain_errors.DomainError)
	GetByID(id string) (*Workspace, domain_errors.DomainError)
	Update(org *Workspace) (*Workspace, domain_errors.DomainError)
	Delete(id string) domain_errors.DomainError
	ListByOwner(ownerID string) ([]*Workspace, domain_errors.DomainError)
}

type MembershipRepository interface {
	Add(membership *Membership) (*Membership, error)
	Remove(userID, organizationID string) error
	ListByOrganization(organizationID string) ([]*Membership, error)
}

type InvitationRepository interface {
	Create(invitation *Invitation) (*Invitation, error)
	GetByToken(token string) (*Invitation, error)
	Delete(token string) error
	ListInvitationToWorkspace(workspace_id string) ([]*Invitation, error)
}