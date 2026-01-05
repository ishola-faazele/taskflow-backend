package project

import "time"

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WorkspaceID string    `json:"workspace_id"`
	Creator     string    `json:"creator"`
	CreatedAt   time.Time `json:"created_at"`
}
type UpdateProjectInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
