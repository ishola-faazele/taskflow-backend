package workspace

import "time"

type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Membership struct {
	UserID      string
	WorkspaceID string
	Role        Role
}

type Invitation struct {
	UserID      string
	WorkspaceID string
	Token       string
	IsValid     bool
}

type Role string

const (
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
	RoleOwner  Role = "owner"
)
