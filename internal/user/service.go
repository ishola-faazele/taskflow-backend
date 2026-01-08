package user

import (
	"time"

	"github.com/google/uuid"
	amqp_utils "github.com/ishola-faazele/taskflow/internal/utils/amqp"
	"github.com/ishola-faazele/taskflow/internal/utils/jwt"
	"github.com/ishola-faazele/taskflow/pkg/utils"
	"github.com/ishola-faazele/taskflow/pkg/utils/domain_errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserService struct {
	authRepo    AuthRepository
	profileRepo UserProfileRepository
	conn        *amqp.Connection
	jwtUtil     *jwt.JWTUtils
}

func NewUserService(authRepo AuthRepository, profileRepo UserProfileRepository, conn *amqp.Connection) *UserService {
	jwtUtil := jwt.NewJWTUtils(jwt.DefaultTokenConfig())
	return &UserService{
		authRepo:    authRepo,
		profileRepo: profileRepo,
		jwtUtil:     jwtUtil,
		conn:        conn,
	}
}

// Creates a magic link which the user has to verify to log in
func (us *UserService) GetMagicLink(email string) domain_errors.DomainError {
	if !utils.IsValidEmail(email) {
		return domain_errors.NewValidationErrorWithValue("email", email, "INVALID_EMAIL_FORMAT")
	}
	// check if user is in db
	user, err := us.authRepo.GetByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		// create new user auth (profile is creatd as well in db implementation)
		new_user := CreateNewAuth(email)
		if user, err = us.authRepo.Create(new_user); err != nil {
			return err
		}
	}

	// create token using auth as claim
	authToken, token_err := us.jwtUtil.GenerateAuthToken(user.ID, email)
	if token_err != nil {
		return domain_errors.NewInternalError("FAILED_TO_GENERATE_TOKEN", token_err)
	}
	// send email with magic link
	emailMsg, msgErr := amqp_utils.NewMagicLinkMessage(email, authToken, "/api/user/verify?token=")
	if msgErr != nil {
		return domain_errors.NewInternalError("FAILED_TO_CREATE_EMAIL_MESSAGE", msgErr)
	}
	// publish message to queue
	ch, chErr := us.conn.Channel()
	if chErr != nil {
		return domain_errors.NewInternalError("FAILED_TO_CREATE_CHANNEL", chErr)
	}

	if err := amqp_utils.PublishEmailMessage(ch, emailMsg); err != nil {
		return domain_errors.NewInternalError("FAILED_TO_PUBLISH_EMAIL_MESSAGE", err)
	}
	return nil
}

// verifies token embedded in the magic link
func (us *UserService) VerifyToken(token string) (string, string, domain_errors.DomainError) {
	// check if token is signed and valid
	claims, parseErr := us.jwtUtil.ParseUserToken(token)
	if parseErr != nil {
		return "", "", domain_errors.NewInternalError("ERROR_PARSING_TOKEN", parseErr)
	}
	// check if token is of valid purpose
	if claims.Purpose != jwt.PurposeAuth {
		return "", "", domain_errors.NewUnauthorizedError("INVALID_TOKEN_PURPOSE")
	}

	// issue new tokens for access and refresh
	access, refresh, tokenErr := us.jwtUtil.GenerateTokenPair(claims.UserID, claims.Email)
	if tokenErr != nil {
		return "", "", domain_errors.NewInternalError("FAILED_TO_GENERATE_ACCESS_AND_REFRESH_TOKENS", tokenErr)
	}
	return access, refresh, nil
}

// Returns new access and refresh tokens while invalidating the old refresh token
func (us *UserService) RefreshToken(refreshToken string) (string, string, domain_errors.DomainError) {
	// validate refresh token
	claims, parseErr := us.jwtUtil.ParseUserToken(refreshToken)
	if parseErr != nil {
		return "", "", domain_errors.NewInternalError("ERROR_PARSING_REFRESH_TOKEN", parseErr)
	}
	if claims.Purpose != jwt.PurposeRefresh {
		return "", "", domain_errors.NewUnauthorizedError("INVALID_REFRESH_TOKEN_PURPOSE")
	}
	// check if token_hash has not been invalidated
	token_hash := utils.HashToken(refreshToken)
	if valid, err := us.authRepo.IsTokenValid(token_hash); !valid || err != nil {
		return "", "", domain_errors.NewUnauthorizedError("TOKEN_IS_INVALID_OR_EXPIRED")
	}
	// invalidate token before creating a newone
	invalidToken := &InvalidToken{
		TokenHash:     token_hash,
		UserID:        claims.UserID,
		ExpiresAt:     claims.ExpiresAt.Time,
		InvalidatedAt: time.Now().UTC(),
	}
	if err := us.authRepo.InvalidateToken(invalidToken); err != nil {
		return "", "", domain_errors.NewInternalError("ERROR_INVALIDATING_REFRESH_TOKEN", err)
	}
	// create new access and refresh tokens
	access, refresh, token_err := us.jwtUtil.GenerateTokenPair(claims.UserID, claims.Email)
	if token_err != nil {
		return "", "", domain_errors.NewInternalError("FAILED_TO_GENERATE_ACCESS_AND_REFRESH TOKENS", token_err)
	}
	return access, refresh, nil
}

// Gets a user's own auth data
func (us *UserService) GetByID(id string) (*Auth, domain_errors.DomainError) {
	return us.authRepo.GetByID(id)
}

// Gets a user's auth data using email
func (us *UserService) GetByEmail(email string) (*Auth, domain_errors.DomainError) {
	return us.authRepo.GetByEmail(email)
}

// Get's a user's profile
func (us *UserService) GetProfile(id string) (*UserProfile, domain_errors.DomainError) {
	return us.profileRepo.GetProfile(id)
}

// updates's user prpofile
func (us UserService) UpdateProfile(userID, name string) (*UserProfile, domain_errors.DomainError) {
	return us.profileRepo.UpdateProfile(userID, name)
}

func (us UserService) GetPublicProfile(id string) (*PublicProfile, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("id", id, "ID IS NOT A VALID UUID")
	}
	return us.profileRepo.GetPublicProfile(id)
}
