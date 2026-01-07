package workspace

import (
	"database/sql"
	"fmt"

	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
	. "github.com/ishola-faazele/taskflow/internal/workspace/entity"
)

type PostgresWorkspaceRepository struct {
	db *sql.DB
}

type PostgresMembershipRepository struct {
	db *sql.DB
}

type PostgresInvitationRepository struct {
	db *sql.DB
}

// NewPostgresWorkspaceRepository creates a new workspace repository
func NewPostgresWorkspaceRepository(db *sql.DB) *PostgresWorkspaceRepository {
	return &PostgresWorkspaceRepository{db: db}
}

// NewPostgresMembershipRepository creates a new membership repository
func NewPostgresMembershipRepository(db *sql.DB) *PostgresMembershipRepository {
	return &PostgresMembershipRepository{db: db}
}

// NewPostgresInvitationRepository creates a new invitation repository
func NewPostgresInvitationRepository(db *sql.DB) *PostgresInvitationRepository {
	return &PostgresInvitationRepository{db: db}
}

// WorkspaceRepository implementation

func (r *PostgresWorkspaceRepository) Create(ws *Workspace) (*Workspace, domain_errors.DomainError) {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, domain_errors.NewDatabaseError("workspace creation - begin transaction", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			// Rollback failed, but transaction may have already been rolled back
		}
	}()

	// Insert the workspace
	workspaceQuery := `
		INSERT INTO workspace (id, name, owner_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, owner_id, created_at
	`

	row := tx.QueryRow(workspaceQuery, ws.ID, ws.Name, ws.OwnerID, ws.CreatedAt)

	result := &Workspace{}
	err = row.Scan(&result.ID, &result.Name, &result.OwnerID, &result.CreatedAt)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("workspace creation - insert workspace", err)
	}

	// Insert the owner membership
	membershipQuery := `
		INSERT INTO membership (user_id, workspace_id, role, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(membershipQuery, result.OwnerID, result.ID, RoleOwner, result.CreatedAt)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("workspace creation - insert membership", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, domain_errors.NewDatabaseError("workspace creation - commit transaction", err)
	}

	return result, nil
}

func (r *PostgresWorkspaceRepository) GetByID(id string) (*Workspace, domain_errors.DomainError) {
	query := `
		SELECT id, name, owner_id, created_at
		FROM workspace
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	workspace := &Workspace{}
	err := row.Scan(&workspace.ID, &workspace.Name, &workspace.OwnerID, &workspace.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_errors.NewNotFoundError("workspace", id)
		}
		return nil, domain_errors.NewDatabaseError("workspace query", err)
	}

	return workspace, nil
}

func (r *PostgresWorkspaceRepository) Update(ws *Workspace) (*Workspace, domain_errors.DomainError) {
	query := `
		UPDATE workspace
		SET name = $2
		WHERE id = $1
		RETURNING id, name, owner_id, created_at
	`

	row := r.db.QueryRow(query, ws.ID, ws.Name)

	result := &Workspace{}
	err := row.Scan(&result.ID, &result.Name, &result.OwnerID, &result.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_errors.NewNotFoundError("workspace", ws.ID)
		}
		return nil, domain_errors.NewDatabaseError("workspace update", err)
	}

	return result, nil
}

func (r *PostgresWorkspaceRepository) Delete(id string) domain_errors.DomainError {
	query := `DELETE FROM workspace WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return domain_errors.NewDatabaseError("workspace deletion", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return domain_errors.NewDatabaseError("workspace deletion", err)
	}

	if rows == 0 {
		return domain_errors.NewNotFoundError("workspace", id)
	}

	return nil
}

func (r *PostgresWorkspaceRepository) ListByOwner(ownerID string) ([]*Workspace, domain_errors.DomainError) {
	query := `
		SELECT id, name, owner_id, created_at
		FROM workspace
		WHERE owner_id = $1
		ORDER BY name
	`

	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("workspace query", err)
	}
	defer rows.Close()

	var workspaces []*Workspace
	for rows.Next() {
		workspace := &Workspace{}
		err := rows.Scan(&workspace.ID, &workspace.Name, &workspace.OwnerID, &workspace.CreatedAt)
		if err != nil {
			return nil, domain_errors.NewDatabaseError("workspace query", err)
		}
		workspaces = append(workspaces, workspace)
	}

	if err = rows.Err(); err != nil {
		return nil, domain_errors.NewDatabaseError("workspace query", err)
	}

	return workspaces, nil
}

// MembershipRepository implementation

func (r *PostgresMembershipRepository) Add(membership *Membership) (*Membership, error) {
	query := `
		INSERT INTO membership (user_id, workspace_id, role, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id, workspace_id, role, created_at
	`

	row := r.db.QueryRow(query, membership.UserID, membership.WorkspaceID, membership.Role, membership.CreatedAt)

	result := &Membership{}
	err := row.Scan(&result.UserID, &result.WorkspaceID, &result.Role, &result.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add membership: %w", err)
	}

	return result, nil
}

func (r *PostgresMembershipRepository) Remove(userID, organizationID string) error {
	query := `
		DELETE FROM membership
		WHERE user_id = $1 AND workspace_id = $2
	`

	result, err := r.db.Exec(query, userID, organizationID)
	if err != nil {
		return fmt.Errorf("failed to remove membership: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("membership not found")
	}

	return nil
}

func (r *PostgresMembershipRepository) ListByWorkspace(workspaceID string) ([]*Membership, error) {
	query := `
		SELECT user_id, workspace_id, role, created_at
		FROM membership
		WHERE workspace_id = $1
		ORDER BY user_id
	`

	rows, err := r.db.Query(query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list memberships: %w", err)
	}
	defer rows.Close()

	var memberships []*Membership
	for rows.Next() {
		membership := &Membership{}
		err := rows.Scan(&membership.UserID, &membership.WorkspaceID, &membership.Role, &membership.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan membership: %w", err)
		}
		memberships = append(memberships, membership)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating memberships: %w", err)
	}

	return memberships, nil
}

func (r *PostgresMembershipRepository) IsMember(userID, workspaceID string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM membership
		WHERE user_id = $1 AND workspace_id = $2
	`

	var count int
	err := r.db.QueryRow(query, userID, workspaceID).Scan(&count)
	if err != nil {
		return false, domain_errors.NewDatabaseError("Check Membership", err)
	}

	return count > 0, nil
}

// InvitationRepository implementation

func (r *PostgresInvitationRepository) Create(invitation *Invitation) (*Invitation, domain_errors.DomainError) {
	query := `
		INSERT INTO invitation (id, invitee_id, invitee_email, inviter_id, workspace_id, role, is_valid, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, invitee_id, invitee_email, inviter_id, workspace_id, role, is_valid, created_at
	`

	row := r.db.QueryRow(query, invitation.ID, invitation.InviteeID, invitation.InviteeEmail, invitation.InviterID, invitation.WorkspaceID, invitation.Role, invitation.IsValid, invitation.CreatedAt)

	result := &Invitation{}
	err := row.Scan(&result.ID, &result.InviteeID, &result.InviteeEmail, &result.InviterID, &result.WorkspaceID, &result.Role, &result.IsValid, &result.CreatedAt)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("FAILED CREATING INVITATION", err)
	}

	return result, nil
}

func (r *PostgresInvitationRepository) GetByID(id string) (*Invitation, domain_errors.DomainError) {
	query := `
		SELECT id, invitee_id, invitee_email, inviter_id, workspace_id, role, is_valid, created_at
		FROM invitation
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	result := &Invitation{}
	err := row.Scan(&result.ID, &result.InviteeID, &result.InviteeEmail, &result.InviterID, &result.WorkspaceID, &result.Role, &result.IsValid, &result.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_errors.NewNotFoundError("Invitation", id)
		}
		return nil, domain_errors.NewDatabaseError("FAILED GETTING INVITATION", err)
	}

	return result, nil
}

func (r *PostgresInvitationRepository) DeleteInvitation(id string) domain_errors.DomainError {
	query := `DELETE FROM invitation WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return domain_errors.NewDatabaseError("FAILED DELETING INVITATION", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return domain_errors.NewDatabaseError("FAILED GETTING ROWS", err)
	}

	if rows == 0 {
		return domain_errors.NewNotFoundError("Invitation", id)
	}

	return nil
}

func (r *PostgresInvitationRepository) ListInvitationToWorkspace(workspace_id string) ([]*Invitation, domain_errors.DomainError) {
	query := `
		SELECT id, invitee_id, invitee_email, inviter_id, workspace_id, role, is_valid, created_at
		FROM invitation
		WHERE workspace_id = $1
	`

	rows, err := r.db.Query(query, workspace_id)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("FAILED GETTING INVITATIONS", err)
	}
	defer rows.Close()

	var results []*Invitation
	for rows.Next() {
		invitation := &Invitation{}
		err := rows.Scan(&invitation.ID, &invitation.InviteeID, &invitation.InviteeEmail, &invitation.InviterID, &invitation.WorkspaceID, &invitation.Role, &invitation.IsValid, &invitation.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, domain_errors.NewNotFoundError("Invitations for Workspace", workspace_id)
			}
			return nil, domain_errors.NewDatabaseError("FAILED GETTING INVITATIONS", err)
		}
		results = append(results, invitation)
	}
	return results, nil
}
