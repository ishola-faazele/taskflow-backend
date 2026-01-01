package workspace

type WorkspaceRepository interface {
	Create(org *Workspace) (*Workspace, error)
	GetByID(id string) (*Workspace, error)
	Update(org *Workspace) (*Workspace, error)
	Delete(id string) error
	ListByOwner(ownerID string) ([]*Workspace, error)
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
}