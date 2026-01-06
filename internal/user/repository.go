package user

import "github.com/ishola-faazele/taskflow/internal/shared/domain_errors"

type AuthRepository interface {
	Create(auth *Auth) (*Auth, domain_errors.DomainError)
	GetByID(id string) (*Auth, domain_errors.DomainError)
	GetByEmail(email string) (*Auth, domain_errors.DomainError)
	IsTokenValid(token_hash string) (bool, domain_errors.DomainError)
	InvalidateToken(token_hash *InvalidToken) domain_errors.DomainError
}

type UserProfileRepository interface {
	GetProfile(id string) (*UserProfile, domain_errors.DomainError)
	UpdateProfile(userID, name string) (*UserProfile, domain_errors.DomainError)
	GetPublicProfile(id string) (*PublicProfile, domain_errors.DomainError)
}
