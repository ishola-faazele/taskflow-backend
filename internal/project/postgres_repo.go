package project

import (
	"database/sql"
	"fmt"
	"time"

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

// ============================================================================
// TASK METHODS
// ============================================================================

func (r *PostgresProjectRepository) CreateTask(task *Task) (*Task, domain_errors.DomainError) {
	query := `
		INSERT INTO task (id, parent_id, project_id, name, description,creator,  status, priority, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, parent_id, project_id, name, description, creator, status, priority, due_date, created_at, updated_at
	`

	row := r.db.QueryRow(
		query,
		task.ID,
		task.ParentID,
		task.ProjectID,
		task.Name,
		task.Description,
		task.Creator,
		task.Status,
		task.Priority,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	)

	result := &Task{}
	err := row.Scan(
		&result.ID,
		&result.ParentID,
		&result.ProjectID,
		&result.Name,
		&result.Description,
		&result.Creator,
		&result.Status,
		&result.Priority,
		&result.DueDate,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("task creation", err)
	}

	return result, nil
}

func (r *PostgresProjectRepository) GetTaskByID(id string) (*Task, domain_errors.DomainError) {
	query := `
		SELECT id, parent_id, project_id, name, description, creator, status, priority, due_date, created_at, updated_at
		FROM task
		WHERE id = $1
	`

	task := &Task{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.ParentID,
		&task.ProjectID,
		&task.Name,
		&task.Description,
		&task.Creator,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_errors.NewNotFoundError("task", id)
		}
		return nil, domain_errors.NewDatabaseError("task query", err)
	}

	return task, nil
}

func (r *PostgresProjectRepository) UpdateTask(input *UpdateTaskInput, id string) (*Task, domain_errors.DomainError) {
	query := `UPDATE task SET updated_at = $1`
	args := []interface{}{time.Now().UTC()}
	argIdx := 2

	if input.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIdx)
		args = append(args, *input.Name)
		argIdx++
	}
	if input.Description != nil {
		query += fmt.Sprintf(", description = $%d", argIdx)
		args = append(args, *input.Description)
		argIdx++
	}
	if input.Status != nil {
		query += fmt.Sprintf(", status = $%d", argIdx)
		args = append(args, *input.Status)
		argIdx++
	}
	if input.Priority != nil {
		query += fmt.Sprintf(", priority = $%d", argIdx)
		args = append(args, *input.Priority)
		argIdx++
	}
	if input.DueDate != nil {
		query += fmt.Sprintf(", due_date = $%d", argIdx)
		args = append(args, *input.DueDate)
		argIdx++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	query += ` RETURNING id, parent_id, project_id, name, description, creator, status, priority, due_date, created_at, updated_at`

	task := &Task{}
	err := r.db.QueryRow(query, args...).Scan(
		&task.ID,
		&task.ParentID,
		&task.ProjectID,
		&task.Name,
		&task.Description,
		&task.Creator,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_errors.NewNotFoundError("task", id)
		}
		return nil, domain_errors.NewDatabaseError("task update", err)
	}

	return task, nil
}

func (r *PostgresProjectRepository) DeleteTask(id string) domain_errors.DomainError {
	query := `DELETE FROM task WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return domain_errors.NewDatabaseError("task deletion", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return domain_errors.NewDatabaseError("task deletion", err)
	}

	if rows == 0 {
		return domain_errors.NewNotFoundError("task", id)
	}

	return nil
}

func (r *PostgresProjectRepository) ListSubtasks(parentID string) ([]*Task, domain_errors.DomainError) {
	query := `
		SELECT id, parent_id, project_id, name, description, status, creator, priority, due_date, created_at, updated_at
		FROM task
		WHERE parent_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("subtasks query", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(
			&task.ID,
			&task.ParentID,
			&task.ProjectID,
			&task.Name,
			&task.Description,
			&task.Creator,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, domain_errors.NewDatabaseError("subtasks scan", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, domain_errors.NewDatabaseError("subtasks iteration", err)
	}

	return tasks, nil
}

func (r *PostgresProjectRepository) GetTaskTree(rootID string) (*TaskTree, domain_errors.DomainError) {
	query := `
		WITH RECURSIVE task_tree AS (
			SELECT id, parent_id, project_id, name, description, creator, status, priority, due_date, created_at, updated_at, 0 as depth
			FROM task
			WHERE id = $1
			
			UNION ALL
			
			SELECT t.id, t.parent_id, t.project_id, t.name, t.description, t.creator, t.status, t.priority, t.due_date, t.created_at, t.updated_at, tt.depth + 1
			FROM task t
			INNER JOIN task_tree tt ON t.parent_id = tt.id
		)
		SELECT id, parent_id, project_id, name, description, creator, status, priority, due_date, created_at, updated_at, depth
		FROM task_tree
		ORDER BY depth, created_at
	`

	rows, err := r.db.Query(query, rootID)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("task tree query", err)
	}
	defer rows.Close()

	taskMap := make(map[string]*TaskTree)
	var rootTask *TaskTree

	for rows.Next() {
		task := &Task{}
		var depth int
		err := rows.Scan(
			&task.ID,
			&task.ParentID,
			&task.ProjectID,
			&task.Name,
			&task.Description,
			&task.Creator,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&depth,
		)
		if err != nil {
			return nil, domain_errors.NewDatabaseError("task tree scan", err)
		}

		taskTree := &TaskTree{
			Task:     *task,
			Subtasks: []*TaskTree{},
			Depth:    depth,
		}

		taskMap[task.ID] = taskTree

		if depth == 0 {
			rootTask = taskTree
		}
	}

	if err = rows.Err(); err != nil {
		return nil, domain_errors.NewDatabaseError("task tree iteration", err)
	}

	if rootTask == nil {
		return nil, domain_errors.NewNotFoundError("task", rootID)
	}

	// Build tree structure
	for _, taskTree := range taskMap {
		if taskTree.ParentID != nil {
			if parent, exists := taskMap[*taskTree.ParentID]; exists {
				parent.Subtasks = append(parent.Subtasks, taskTree)
			}
		}
	}

	return rootTask, nil
}

func (r *PostgresProjectRepository) GetTaskWithChildren(id string) (*TaskWithChildren, domain_errors.DomainError) {
	task, err := r.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	children, err := r.ListSubtasks(id)
	if err != nil {
		return nil, err
	}

	return &TaskWithChildren{
		Task:     *task,
		Children: children,
	}, nil
}

func (r *PostgresProjectRepository) GetRootTasks(projectID string) ([]*Task, domain_errors.DomainError) {
	query := `
		SELECT id, parent_id, project_id, name, description, creator, status, priority, due_date, created_at, updated_at
		FROM task
		WHERE project_id = $1 AND (parent_id IS NULL OR parent_id = '')
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("root tasks query", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(
			&task.ID,
			&task.ParentID,
			&task.ProjectID,
			&task.Name,
			&task.Description,
			&task.Creator,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, domain_errors.NewDatabaseError("root tasks scan", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
func (r *PostgresProjectRepository) ListTasksByProject(projectID string) ([]*Task, domain_errors.DomainError) {
	query := `
		SELECT id, parent_id, project_id, name, description,creator, status, priority, due_date, created_at, updated_at
		FROM task
		WHERE project_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("project tasks query", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(
			&task.ID,
			&task.ParentID,
			&task.ProjectID,
			&task.Name,
			&task.Description,
			&task.Creator,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, domain_errors.NewDatabaseError("project tasks scan", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *PostgresProjectRepository) GetProjectTaskTree(projectID string) ([]*TaskTree, domain_errors.DomainError) {
	// Get all root tasks for this project
	rootTasks, err := r.GetRootTasks(projectID)
	if err != nil {
		return nil, err
	}

	// Build tree for each root task
	var trees []*TaskTree
	for _, rootTask := range rootTasks {
		tree, err := r.GetTaskTree(rootTask.ID)
		if err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}

	return trees, nil
}

func (r *PostgresProjectRepository) GetTaskDepth(id string) (int, domain_errors.DomainError) {
	query := `
		WITH RECURSIVE task_depth AS (
			SELECT id, parent_id, 0 as depth
			FROM task
			WHERE id = $1
			
			UNION ALL
			
			SELECT t.id, t.parent_id, td.depth + 1
			FROM task t
			INNER JOIN task_depth td ON t.id = td.parent_id
		)
		SELECT MAX(depth) FROM task_depth
	`

	var depth int
	err := r.db.QueryRow(query, id).Scan(&depth)
	if err != nil {
		return 0, domain_errors.NewDatabaseError("task depth query", err)
	}

	return depth, nil
}

func (r *PostgresProjectRepository) CountSubtasks(parentID string) (int, domain_errors.DomainError) {
	query := `
		WITH RECURSIVE task_tree AS (
			SELECT id FROM task WHERE parent_id = $1
			UNION ALL
			SELECT t.id FROM task t
			INNER JOIN task_tree tt ON t.parent_id = tt.id
		)
		SELECT COUNT(*) FROM task_tree
	`

	var count int
	err := r.db.QueryRow(query, parentID).Scan(&count)
	if err != nil {
		return 0, domain_errors.NewDatabaseError("subtask count query", err)
	}

	return count, nil
}
