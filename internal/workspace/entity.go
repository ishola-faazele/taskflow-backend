package workspace

type Workspace struct {
	ID   string
	Name string
	OwnerID string
}

type Membership struct {
	UserID         string
	WorkspaceID string
	Role           Role
}

type Invitation struct {
	UserID         string
	WorkspaceID string
	Token          string
	IsValid 		bool
}

type Role string
const (
	RoleMember Role = "member"
	RoleAdmin Role = "admin"
	RoleOwner Role = "owner"
)
