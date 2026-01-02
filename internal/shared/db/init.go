package db

import (
	"database/sql"
	"fmt"
	"github.com/ishola-faazele/taskflow/internal/shared/logger"
)

// MigrationManager handles database schema migrations
type MigrationManager struct {
	db *sql.DB
	logger *logger.StdLogger
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	logger := logger.NewStdLogger()
	return &MigrationManager{db: db, logger: logger}
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

// createWorkspaceTable creates the workspace table with indices
func (m *MigrationManager) createWorkspaceTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS workspace (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			owner_id VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create workspace table: %w", err)
	}
	
	// Create index on owner_id for faster lookups
	indexQuery := `
		CREATE INDEX IF NOT EXISTS idx_workspace_owner_id 
		ON workspace(owner_id)
	`
	
	_, err = m.db.Exec(indexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on workspace.owner_id: %w", err)
	}
	
	// Create index on name for potential search functionality
	nameIndexQuery := `
		CREATE INDEX IF NOT EXISTS idx_workspace_name 
		ON workspace(name)
	`
	
	_, err = m.db.Exec(nameIndexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on workspace.name: %w", err)
	}
	
	m.logger.Info("Created workspace table with indices")
	return nil
}

// createMembershipTable creates the membership table with indices
func (m *MigrationManager) createMembershipTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS membership (
			user_id VARCHAR(255) NOT NULL,
			organization_id VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, organization_id),
			CONSTRAINT fk_membership_workspace 
				FOREIGN KEY (organization_id) 
				REFERENCES workspace(id) 
				ON DELETE CASCADE
		)
	`
	
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create membership table: %w", err)
	}
	
	// Create index on organization_id for faster lookups
	indexQuery := `
		CREATE INDEX IF NOT EXISTS idx_membership_organization_id 
		ON membership(organization_id)
	`
	
	_, err = m.db.Exec(indexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on membership.organization_id: %w", err)
	}
	
	// Create index on user_id for reverse lookups
	userIndexQuery := `
		CREATE INDEX IF NOT EXISTS idx_membership_user_id 
		ON membership(user_id)
	`
	
	_, err = m.db.Exec(userIndexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on membership.user_id: %w", err)
	}
	
	// Create index on role for role-based queries
	roleIndexQuery := `
		CREATE INDEX IF NOT EXISTS idx_membership_role 
		ON membership(role)
	`
	
	_, err = m.db.Exec(roleIndexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on membership.role: %w", err)
	}
	
	m.logger.Info("Created membership table with indices")
	return nil
}

// createInvitationTable creates the invitation table with indices
func (m *MigrationManager) createInvitationTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS invitation (
			token VARCHAR(255) PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			organization_id VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP,
			CONSTRAINT fk_invitation_workspace 
				FOREIGN KEY (organization_id) 
				REFERENCES workspace(id) 
				ON DELETE CASCADE
		)
	`
	
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create invitation table: %w", err)
	}
	
	// Create index on email for lookups
	emailIndexQuery := `
		CREATE INDEX IF NOT EXISTS idx_invitation_email 
		ON invitation(email)
	`
	
	_, err = m.db.Exec(emailIndexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on invitation.email: %w", err)
	}
	
	// Create index on organization_id
	orgIndexQuery := `
		CREATE INDEX IF NOT EXISTS idx_invitation_organization_id 
		ON invitation(organization_id)
	`
	
	_, err = m.db.Exec(orgIndexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on invitation.organization_id: %w", err)
	}
	
	// Create index on expires_at for cleanup queries
	expiresIndexQuery := `
		CREATE INDEX IF NOT EXISTS idx_invitation_expires_at 
		ON invitation(expires_at)
	`
	
	_, err = m.db.Exec(expiresIndexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on invitation.expires_at: %w", err)
	}
	
	m.logger.Info("Created invitation table with indices")
	return nil
}

// EnsureTablesExist checks if tables exist and creates them if they don't
func (m *MigrationManager) EnsureTablesExist() error {
	m.logger.Info("Checking database schema...")
	
	// Check and create workspace table first (since it's referenced by foreign keys)
	exists, err := m.tableExists("workspace")
	if err != nil {
		return err
	}
	
	if !exists {
		m.logger.Info("workspace table does not exist, creating...")
		if err := m.createWorkspaceTable(); err != nil {
			return err
		}
	} else {
		m.logger.Info("workspace table already exists")
	}
	
	// Check and create membership table
	exists, err = m.tableExists("membership")
	if err != nil {
		return err
	}
	
	if !exists {
		m.logger.Info("membership table does not exist, creating...")
		if err := m.createMembershipTable(); err != nil {
			return err
		}
	} else {
		m.logger.Info("membership table already exists")
	}
	
	// Check and create invitation table
	exists, err = m.tableExists("invitation")
	if err != nil {
		return err
	}
	
	if !exists {
		m.logger.Info("invitation table does not exist, creating...")
		if err := m.createInvitationTable(); err != nil {
			return err
		}
	} else {
		m.logger.Info("invitation table already exists")
	}
	
	m.logger.Info("Database schema is ready")
	return nil
}

// DropAllTables drops all workspace-related tables (useful for testing)
func (m *MigrationManager) DropAllTables() error {
	m.logger.Warn("Dropping all tables...")
	
	queries := []string{
		"DROP TABLE IF EXISTS invitation CASCADE",
		"DROP TABLE IF EXISTS membership CASCADE",
		"DROP TABLE IF EXISTS workspace CASCADE",
	}
	
	for _, query := range queries {
		_, err := m.db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to drop tables: %w", err)
		}
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

// Example usage function
// func ExampleUsage() {
// 	// Connect to database
// 	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/dbname?sslmode=disable")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
	
// 	// Create migration manager
// 	migrationMgr := NewMigrationManager(db)
	
// 	// Ensure tables exist
// 	if err := migrationMgr.EnsureTablesExist(); err != nil {
// 		log.Fatal(err)
// 	}
	
// 	// Now you can use the repositories
// 	workspaceRepo := NewPostgresWorkspaceRepository(db)
// 	membershipRepo := NewPostgresMembershipRepository(db)
// 	invitationRepo := NewPostgresInvitationRepository(db)
	
// 	_ = workspaceRepo
// 	_ = membershipRepo
// 	_ = invitationRepo
// }