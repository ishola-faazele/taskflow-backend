package project

import (
	"database/sql"

	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
)

type PostgresProjectRepository struct {
	db *sql.DB
}

func NewPostgresProjectRepository(db *sql.DB) *PostgresProjectRepository {
	return &PostgresProjectRepository{db: db}
}
func (r *PostgresProjectRepository) Create(project *Project) (*Project, domain_errors.DomainError) {
	query := `INSERT INTO project (id, name, description, workspace_id, creator, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, description, workspace_id, creator, created_at`
	row := r.db.QueryRow(query, project.ID, project.Name, project.Description, project.WorkspaceID, project.Creator, project.CreatedAt)
	var createdProject Project
	err := row.Scan(&createdProject.ID, &createdProject.Name, &createdProject.Description, &createdProject.WorkspaceID, &createdProject.Creator, &createdProject.CreatedAt)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("project creation", err)
	}
	return &createdProject, nil
}

func (r *PostgresProjectRepository) GetByID(id string) (*Project, domain_errors.DomainError) {
	query := `SELECT id, name, description, workspace_id, creator, created_at FROM project WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var project Project
	err := row.Scan(&project.ID, &project.Name, &project.Description, &project.WorkspaceID, &project.Creator, &project.CreatedAt)
	if err != nil {
		return nil, domain_errors.NewNotFoundError("Project", id)
	}
	return &project, nil
}

func (r *PostgresProjectRepository) Update(input *UpdateProjectInput, id string) (*Project, domain_errors.DomainError) {
	query := `UPDATE project SET name = COALESCE($1, name), description = COALESCE($2, description) WHERE id = $3 RETURNING id, name, description, workspace_id, creator, created_at`
	row := r.db.QueryRow(query, input.Name, input.Description, id)
	var updatedProject Project
	err := row.Scan(&updatedProject.ID, &updatedProject.Name, &updatedProject.Description, &updatedProject.WorkspaceID, &updatedProject.Creator, &updatedProject.CreatedAt)
	if err != nil {
		return nil, domain_errors.NewNotFoundError("Project", id)
	}
	return &updatedProject, nil
}

func (r *PostgresProjectRepository) Delete(id string) domain_errors.DomainError {
	query := `DELETE FROM project WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return domain_errors.NewDatabaseError("project deletion", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain_errors.NewDatabaseError("project deletion", err)
	}
	if rowsAffected == 0 {
		return domain_errors.NewNotFoundError("Project", id)
	}
	return nil
}

func (r *PostgresProjectRepository) ListByWorkspace(wsID string) ([]*Project, domain_errors.DomainError) {
	query := `SELECT id, name, description, workspace_id, creator, created_at FROM project WHERE workspace_id = $1`
	rows, err := r.db.Query(query, wsID)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("project listing", err)
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.WorkspaceID, &project.Creator, &project.CreatedAt); err != nil {
			return nil, domain_errors.NewDatabaseError("project listing", err)
		}
		projects = append(projects, &project)
	}
	if err := rows.Err(); err != nil {
		return nil, domain_errors.NewDatabaseError("project listing", err)
	}
	return projects, nil
}
