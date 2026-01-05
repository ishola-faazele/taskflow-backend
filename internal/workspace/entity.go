package workspace

import "time"

type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Invitation struct {
	ID           string    `json:"id"`
	WorkspaceID  string    `json:"workspace_id"`
	InviterID    string    `json:"inviter_id"`
	InviteeEmail string    `json:"invitee_email"`
	InviteeID    string    `json:"invitee_id"`
	Role         Role      `json:"role"`
	IsValid      bool      `json:"is_valid"`
	CreatedAt    time.Time `json:"created_at"`
}

type Membership struct {
	UserID      string
	WorkspaceID string
	Role        Role
	CreatedAt   time.Time
}

type Role string

const (
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
	RoleOwner  Role = "owner"
)
