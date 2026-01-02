package workspace

import (
	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
)

type MockWorkspaceRepository struct {
	CreateFn      func(org *Workspace) (*Workspace, error)
	GetByIDFn     func(id string) (*Workspace, error)
	UpdateFn      func(org *Workspace) (*Workspace, error)
	DeleteFn      func(id string) error
	ListByOwnerFn func(ownerID string) ([]*Workspace, error)
}

func (m *MockWorkspaceRepository) Create(org *Workspace) (*Workspace, domain_errors.DomainError) {
	// if m.CreateFn != nil {
	// 	return m.CreateFn(org)
	// }
	return org, nil
}

func (m *MockWorkspaceRepository) GetByID(id string) (*Workspace, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
}

func (m *MockWorkspaceRepository) Update(org *Workspace) (*Workspace, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(org)
	}
	return org, nil
}

func (m *MockWorkspaceRepository) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

func (m *MockWorkspaceRepository) ListByOwner(ownerID string) ([]*Workspace, error) {
	if m.ListByOwnerFn != nil {
		return m.ListByOwnerFn(ownerID)
	}
	return []*Workspace{}, nil
}




type MockMembershipRepository struct {
	AddFn               func(membership *Membership) (*Membership, error)
	RemoveFn            func(userID, organizationID string) error
	ListByOrganizationFn func(organizationID string) ([]*Membership, error)
}

func (m *MockMembershipRepository) Add(membership *Membership) (*Membership, error) {
	if m.AddFn != nil {
		return m.AddFn(membership)
	}
	return membership, nil
}

func (m *MockMembershipRepository) Remove(userID, organizationID string) error {
	if m.RemoveFn != nil {
		return m.RemoveFn(userID, organizationID)
	}
	return nil
}

func (m *MockMembershipRepository) ListByOrganization(organizationID string) ([]*Membership, error) {
	if m.ListByOrganizationFn != nil {
		return m.ListByOrganizationFn(organizationID)
	}
	return []*Membership{}, nil
}


type MockInvitationRepository struct {
	CreateFn    func(invitation *Invitation) (*Invitation, error)
	GetByTokenFn func(token string) (*Invitation, error)
	DeleteFn    func(token string) error
	ListInvitationToWorkspaceFn func(workspace_id string) ([]*Invitation, error)
}

func (m *MockInvitationRepository) Create(invitation *Invitation) (*Invitation, error) {
	if m.CreateFn != nil {
		return m.CreateFn(invitation)
	}
	return invitation, nil
}

func (m *MockInvitationRepository) GetByToken(token string) (*Invitation, error) {
	if m.GetByTokenFn != nil {
		return m.GetByTokenFn(token)
	}
	return nil, nil
}

func (m *MockInvitationRepository) Delete(token string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(token)
	}
	return nil
}

func (m *MockInvitationRepository) ListInvitationToWorkspace(workspace_id string) ([]*Invitation, error) {
	return  nil, nil
}
