package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
	"github.com/ishola-faazele/taskflow/internal/shared/jwt"
	"github.com/ishola-faazele/taskflow/pkg/utils"
)

type UserService struct {
	authRepo    AuthRepository
	profileRepo UserProfileRepository
	jwtUtil     *jwt.JWTUtils
}

func NewUserService(authRepo AuthRepository, profileRepo UserProfileRepository) *UserService {
	jwtUtil := jwt.NewJWTUtils("your-secret-key", "your-app-name", jwt.DefaultTokenConfig())
	return &UserService{
		authRepo:    authRepo,
		profileRepo: profileRepo,
		jwtUtil:     jwtUtil,
	}
}

// Creates a magic link which the user has to verify to log in
func (us *UserService) GetMagicLink(email string) domain_errors.DomainError {
	if !utils.IsValidEmail(email) {
		return domain_errors.NewValidationErrorWithValue("email", email, "INVALID EMAIL FORMAT")
	}
	// check if user is in db
	user, err := us.authRepo.GetByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		// create new user
		new_user := &Auth{
			ID:        uuid.NewString(),
			Email:     email,
			CreatedAt: time.Now().UTC(),
		}
		// profile is creatd as well
		_, err = us.authRepo.Create(new_user)
		if err != nil {
			return nil
		}
	}

	// create token using auth as claim
	authToken, token_err := us.jwtUtil.GenerateAuthToken(email, user.ID)
	if token_err != nil {
		return domain_errors.NewInternalError("FAILED TO GENERATE TOKEN", token_err)
	}
	// send email with magic link
	_ = authToken // TODO: use this token to create magic link and send email
	return nil
}

func (us *UserService) VerifyToken(token string) (string, string, domain_errors.DomainError) {
	// check if token is signed and is valid purpose(login)
	claims, err := us.jwtUtil.ParseToken(token)
	if err != nil {
		return "", "", domain_errors.NewInternalError("ERROR PARSING TOKEN", err)
	}
	if claims.Purpose != jwt.PurposeAuth {
		return "", "", domain_errors.NewUnauthorizedError("INVALID TOKEN PURPOSE")
	}

	// issue new tokens for access and refresh
	access, refresh, token_err := us.jwtUtil.GenerateTokenPair(claims.UserID, claims.Email)
	if token_err != nil {
		return "", "", domain_errors.NewInternalError("FAILED TO GENERATE ACCESS AND REFRESH TOKENS", token_err)
	}
	return access, refresh, nil
}

func (us *UserService) GetByID(id string) (*Auth, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("id", id, "ID IS NOT A VALID UUID")
	}
	return us.authRepo.GetByID(id)
}

func (us *UserService) GetByEmail(email string) (*Auth, domain_errors.DomainError) {
	// i need to validate email
	return us.authRepo.GetByEmail(email)
}

func (us *UserService) GetProfile(id string) (*UserProfile, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("id", id, "ID IS NOT A VALID UUID")
	}
	return us.profileRepo.GetProfile(id)
}

func (us UserService) UpdateProfile(profile *UserProfile, requesterID string) (*UserProfile, domain_errors.DomainError) {
	if profile.ID != requesterID {
		return nil, domain_errors.NewUnauthorizedError("REQUESTER IS NOT THE OWNER OF PROFILE")
	}
	return us.profileRepo.UpdateProfile(profile)
}
