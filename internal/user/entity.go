package user

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateNewAuth(email string) *Auth {
	return &Auth{
		ID:        uuid.NewString(),
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}
}

type UserProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PublicProfile struct {
	ID    string `json:"id"` 
	Name  string `json:"name"`
	Email string `json:"email"`
}

type InvalidToken struct {
	TokenHash     string    `json:"token_hash"`
	UserID        string    `json:"user_id"`
	InvalidatedAt time.Time `json:"invalidated_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}
