package user

import (
	"database/sql"

	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
)

// PostgresAuthRepository handles auth data persistence
type PostgresAuthRepository struct {
	db *sql.DB
}

// PostgresUserProfileRepository handles user profile data persistence
type PostgresUserProfileRepository struct {
	db *sql.DB
}

// NewPostgresAuthRepository creates a new auth repository
func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

// NewPostgresUserProfileRepository creates a new user profile repository
func NewPostgresUserProfileRepository(db *sql.DB) *PostgresUserProfileRepository {
	return &PostgresUserProfileRepository{db: db}
}

// Auth Repository Implementation

func (r *PostgresAuthRepository) Create(auth *Auth) (*Auth, domain_errors.DomainError) {
	// Start a transaction to ensure both auth and profile are created atomically
	tx, err := r.db.Begin()
	if err != nil {
		return nil, domain_errors.NewDatabaseError("start transaction", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Insert auth record
	authQuery := `
		INSERT INTO auth (id, email, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, email, created_at
	`

	row := tx.QueryRow(authQuery, auth.ID, auth.Email, auth.CreatedAt)

	result := &Auth{}
	err = row.Scan(&result.ID, &result.Email, &result.CreatedAt)
	if err != nil {
		return nil, domain_errors.NewDatabaseError("auth creation", err)
	}

	// Create corresponding profile with the same ID
	profileQuery := `
		INSERT INTO user_profile (id, name)
		VALUES ($1, $2)
	`

	_, err = tx.Exec(profileQuery, result.ID, "")
	if err != nil {
		return nil, domain_errors.NewDatabaseError("profile creation", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, domain_errors.NewDatabaseError("commit transaction", err)
	}

	return result, nil
}

func (r *PostgresAuthRepository) GetByID(id string) (*Auth, domain_errors.DomainError) {
	query := `
		SELECT id, email, created_at
		FROM auth
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	auth := &Auth{}
	err := row.Scan(&auth.ID, &auth.Email, &auth.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_errors.NewNotFoundError("auth", id)
		}
		return nil, domain_errors.NewDatabaseError("auth query", err)
	}

	return auth, nil
}

func (r *PostgresAuthRepository) GetByEmail(email string) (*Auth, domain_errors.DomainError) {
	query := `
		SELECT id, email, created_at
		FROM auth
		WHERE email = $1
	`

	row := r.db.QueryRow(query, email)

	auth := &Auth{}
	err := row.Scan(&auth.ID, &auth.Email, &auth.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil when user doesn't exist (not an error)
		}
		return nil, domain_errors.NewDatabaseError("auth query", err)
	}

	return auth, nil
}

// UserProfile Repository Implementation

func (r *PostgresUserProfileRepository) GetProfile(id string) (*UserProfile, domain_errors.DomainError) {
	query := `
		SELECT id, name
		FROM user_profile
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	profile := &UserProfile{}
	err := row.Scan(&profile.ID, &profile.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_errors.NewNotFoundError("user profile", id)
		}
		return nil, domain_errors.NewDatabaseError("profile query", err)
	}

	return profile, nil
}

func (r *PostgresUserProfileRepository) UpdateProfile(profile *UserProfile) (*UserProfile, domain_errors.DomainError) {
	query := `
		UPDATE user_profile
		SET name = $2
		WHERE id = $1
		RETURNING id, name
	`

	row := r.db.QueryRow(query, profile.ID, profile.Name)

	result := &UserProfile{}
	err := row.Scan(&result.ID, &result.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_errors.NewNotFoundError("user profile", profile.ID)
		}
		return nil, domain_errors.NewDatabaseError("profile update", err)
	}

	return result, nil
}
