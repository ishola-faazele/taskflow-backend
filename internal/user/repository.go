package user

import "github.com/ishola-faazele/taskflow/internal/shared/domain_errors"

type AuthRepository interface {
	Create(auth *Auth) (*Auth, domain_errors.DomainError)
	GetByID(id string) (*Auth, domain_errors.DomainError)
	GetByEmail(email string) (*Auth, domain_errors.DomainError)
}

type UserProfileRepository interface {
	GetProfile(id string) (*UserProfile, domain_errors.DomainError)
	UpdateProfile(profile *UserProfile) (*UserProfile, domain_errors.DomainError)
}
