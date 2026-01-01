package workspace

type Workspace struct {
	ID   string
	Name string
	OwnerID string
}

type Membership struct {
	UserID         string
	OrganizationID string
	Role           Role
}

type Invitation struct {
	Email          string
	OrganizationID string
	Token          string
}

type Role string
const (
	RoleMember Role = "member"
	RoleAdmin Role = "admin"
	RoleOwner Role = "owner"
)
