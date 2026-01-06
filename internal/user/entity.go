package user

import "time"

type Auth struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
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
