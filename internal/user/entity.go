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
