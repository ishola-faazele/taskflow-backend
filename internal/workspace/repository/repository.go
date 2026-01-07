package workspace

import (
	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
	. "github.com/ishola-faazele/taskflow/internal/workspace/entity"
)

type WorkspaceRepository interface {
	Create(ws *Workspace) (*Workspace, domain_errors.DomainError)
	GetByID(id string) (*Workspace, domain_errors.DomainError)
	Update(ws *Workspace) (*Workspace, domain_errors.DomainError)
	Delete(id string) domain_errors.DomainError
	ListByOwner(ownerID string) ([]*Workspace, domain_errors.DomainError)
}

type InvitationRepository interface {
	Create(invitation *Invitation) (*Invitation, domain_errors.DomainError)
	GetByID(id string) (*Invitation, domain_errors.DomainError)
	DeleteInvitation(id string) domain_errors.DomainError
	ListInvitationToWorkspace(ws_id string) ([]*Invitation, domain_errors.DomainError)
}

type MembershipRepository interface {
	Add(membership *Membership) (*Membership, error)
	Remove(userID, workspaceID string) error
	ListByWorkspace(workspaceID string) ([]*Membership, error)
	IsMember(userID, workspaceID string) (bool, error)
}
