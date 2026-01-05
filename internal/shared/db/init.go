package db

import (
	"database/sql"
	"fmt"

	"github.com/ishola-faazele/taskflow/internal/shared/logger"
)

// TableDefinition represents a database table with its schema and indices
type TableDefinition struct {
	Name         string
	CreateSQL    string
	Indices      []string
	Dependencies []string // Tables this table depends on (for foreign keys)
}

// MigrationManager handles database schema migrations
type MigrationManager struct {
	db     *sql.DB
	logger *logger.StdLogger
	tables []TableDefinition
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	logger := logger.NewStdLogger()
	mgr := &MigrationManager{
		db:     db,
		logger: logger,
		tables: []TableDefinition{},
	}

	// Register all tables
	mgr.registerWorkspaceTables()
	mgr.registerUserTables()
	mgr.registerProjectTables()
	return mgr
}

// registerWorkspaceTables registers workspace-related tables
func (m *MigrationManager) registerWorkspaceTables() {
	// Workspace table
	m.RegisterTable(TableDefinition{
		Name: "workspace",
		CreateSQL: `
			CREATE TABLE IF NOT EXISTS workspace (
				id VARCHAR(255) PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				owner_id VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`,
		Indices: []string{
			`CREATE INDEX IF NOT EXISTS idx_workspace_owner_id ON workspace(owner_id)`,
			`CREATE INDEX IF NOT EXISTS idx_workspace_name ON workspace(name)`,
		},
		Dependencies: []string{},
	})

	// Membership table
	m.RegisterTable(TableDefinition{
		Name: "membership",
		CreateSQL: `
			CREATE TABLE IF NOT EXISTS membership (
				user_id VARCHAR(255) NOT NULL,
				workspace_id VARCHAR(255) NOT NULL,
				role VARCHAR(50) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				PRIMARY KEY (user_id, workspace_id),
				CONSTRAINT fk_membership_workspace 
					FOREIGN KEY (workspace_id)
					REFERENCES workspace(id) 
					ON DELETE CASCADE,
				CONSTRAINT fk_membership_user
					FOREIGN KEY (user_id)
					REFERENCES auth(id)
					ON DELETE CASCADE
			)
		`,
		Indices: []string{
			`CREATE INDEX IF NOT EXISTS idx_membership_workspace_id ON membership(workspace_id)`,
			`CREATE INDEX IF NOT EXISTS idx_membership_user_id ON membership(user_id)`,
			`CREATE INDEX IF NOT EXISTS idx_membership_role ON membership(role)`,
		},
		Dependencies: []string{"workspace"},
	})

	// Invitation table
	m.RegisterTable(TableDefinition{
		Name: "invitation",
		CreateSQL: `
		CREATE TABLE IF NOT EXISTS invitation (
			id VARCHAR(255) PRIMARY KEY,
			invitee_id VARCHAR(255),
			invitee_email VARCHAR(255) NOT NULL,
			inviter_id VARCHAR(255) NOT NULL,
			workspace_id VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			is_valid BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_invitation_workspace 
				FOREIGN KEY (workspace_id) 
				REFERENCES workspace(id) 
				ON DELETE CASCADE,
			CONSTRAINT fk_invitation_inviter
				FOREIGN KEY (inviter_id)
				REFERENCES auth(id)
				ON DELETE CASCADE,
			CONSTRAINT fk_invitation_invitee
				FOREIGN KEY (invitee_id)
				REFERENCES auth(id)
				ON DELETE CASCADE
		)
	`,
		Indices: []string{
			`CREATE INDEX IF NOT EXISTS idx_invitation_invitee_email ON invitation(invitee_email)`,
			`CREATE INDEX IF NOT EXISTS idx_invitation_workspace_id ON invitation(workspace_id)`,
			`CREATE INDEX IF NOT EXISTS idx_invitation_inviter_id ON invitation(inviter_id)`,
			`CREATE INDEX IF NOT EXISTS idx_invitation_is_valid ON invitation(is_valid)`,
			`CREATE INDEX IF NOT EXISTS idx_invitation_invitee_id ON invitation(invitee_id)`,
		},
		Dependencies: []string{"workspace", "auth"},
	})

}

// registerUserTables registers user-related tables
func (m *MigrationManager) registerUserTables() {
	// Auth table
	m.RegisterTable(TableDefinition{
		Name: "auth",
		CreateSQL: `
			CREATE TABLE IF NOT EXISTS auth (
				id VARCHAR(255) PRIMARY KEY,
				email VARCHAR(255) NOT NULL UNIQUE,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			)
		`,
		Indices: []string{
			`CREATE INDEX IF NOT EXISTS idx_auth_email ON auth(email)`,
		},
		Dependencies: []string{},
	})

	// User Profile table
	m.RegisterTable(TableDefinition{
		Name: "user_profile",
		CreateSQL: `
			CREATE TABLE IF NOT EXISTS user_profile (
				id VARCHAR(255) PRIMARY KEY,
				name VARCHAR(255) NOT NULL DEFAULT '',
				CONSTRAINT fk_user_profile_auth
					FOREIGN KEY (id)
					REFERENCES auth(id)
					ON DELETE CASCADE
			)
		`,
		Indices:      []string{},
		Dependencies: []string{"auth"},
	})
}
func (m *MigrationManager) registerProjectTables() {
	// Project table
	m.RegisterTable(TableDefinition{
		Name: "project",
		CreateSQL: `
			CREATE TABLE IF NOT EXISTS project (
				id VARCHAR(255) PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				description TEXT,
				workspace_id VARCHAR(255) NOT NULL,
				creator VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT fk_project_workspace
					FOREIGN KEY (workspace_id)
					REFERENCES workspace(id)
					ON DELETE CASCADE,
				CONSTRAINT fk_project_creator
					FOREIGN KEY (creator)
					REFERENCES auth(id)
					ON DELETE SET NULL
			)
		`,
		Indices: []string{
			`CREATE INDEX IF NOT EXISTS idx_project_workspace_id ON project(workspace_id)`,
			`CREATE INDEX IF NOT EXISTS idx_project_creator ON project(creator)`,
			`CREATE INDEX IF NOT EXISTS idx_project_name ON project(name)`,
		},
		Dependencies: []string{"workspace", "auth"},
	})
}

// RegisterTable adds a new table definition to the migration manager
func (m *MigrationManager) RegisterTable(table TableDefinition) {
	m.tables = append(m.tables, table)
}

// tableExists checks if a table exists in the database
func (m *MigrationManager) tableExists(tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)
	`

	var exists bool
	err := m.db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if table %s exists: %w", tableName, err)
	}

	return exists, nil
}

// createTable creates a single table with its indices
func (m *MigrationManager) createTable(table TableDefinition) error {
	// Create the table
	_, err := m.db.Exec(table.CreateSQL)
	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", table.Name, err)
	}

	// Create indices
	for _, indexSQL := range table.Indices {
		_, err := m.db.Exec(indexSQL)
		if err != nil {
			return fmt.Errorf("failed to create index for table %s: %w", table.Name, err)
		}
	}

	m.logger.Info(fmt.Sprintf("Created table '%s' with %d indices", table.Name, len(table.Indices)))
	return nil
}

// sortTablesByDependencies returns tables sorted by their dependencies
// Tables without dependencies come first, then tables that depend on them
func (m *MigrationManager) sortTablesByDependencies() []TableDefinition {
	sorted := []TableDefinition{}
	processed := make(map[string]bool)

	// Helper function to recursively add tables and their dependencies
	var addTable func(string)
	addTable = func(tableName string) {
		if processed[tableName] {
			return
		}

		// Find the table definition
		var table *TableDefinition
		for i := range m.tables {
			if m.tables[i].Name == tableName {
				table = &m.tables[i]
				break
			}
		}

		if table == nil {
			return
		}

		// First add all dependencies
		for _, dep := range table.Dependencies {
			addTable(dep)
		}

		// Then add this table
		if !processed[tableName] {
			sorted = append(sorted, *table)
			processed[tableName] = true
		}
	}

	// Add all tables
	for _, table := range m.tables {
		addTable(table.Name)
	}

	return sorted
}

// EnsureTablesExist checks if tables exist and creates them if they don't
func (m *MigrationManager) EnsureTablesExist() error {
	m.logger.Info("Checking database schema...")

	// Sort tables by dependencies
	sortedTables := m.sortTablesByDependencies()

	// Create tables in dependency order
	for _, table := range sortedTables {
		exists, err := m.tableExists(table.Name)
		if err != nil {
			return err
		}

		if !exists {
			m.logger.Info(fmt.Sprintf("Table '%s' does not exist, creating...", table.Name))
			if err := m.createTable(table); err != nil {
				return err
			}
		} else {
			m.logger.Info(fmt.Sprintf("Table '%s' already exists", table.Name))
		}
	}

	m.logger.Info("Database schema is ready")
	return nil
}

// DropAllTables drops all registered tables (useful for testing)
func (m *MigrationManager) DropAllTables() error {
	m.logger.Warn("Dropping all tables...")

	// Get tables in reverse dependency order
	sortedTables := m.sortTablesByDependencies()

	// Drop in reverse order to respect foreign key constraints
	for i := len(sortedTables) - 1; i >= 0; i-- {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", sortedTables[i].Name)
		_, err := m.db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %w", sortedTables[i].Name, err)
		}
		m.logger.Info(fmt.Sprintf("Dropped table '%s'", sortedTables[i].Name))
	}

	m.logger.Info("All tables dropped")
	return nil
}

// ResetDatabase drops and recreates all tables (useful for testing)
func (m *MigrationManager) ResetDatabase() error {
	if err := m.DropAllTables(); err != nil {
		return err
	}

	return m.EnsureTablesExist()
}

// GetTableNames returns all registered table names
func (m *MigrationManager) GetTableNames() []string {
	names := make([]string, len(m.tables))
	for i, table := range m.tables {
		names[i] = table.Name
	}
	return names
}

// AddCustomTable allows adding tables at runtime (useful for plugins/extensions)
func (m *MigrationManager) AddCustomTable(table TableDefinition) {
	m.RegisterTable(table)
}
