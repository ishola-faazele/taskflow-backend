package user

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
	"github.com/ishola-faazele/taskflow/internal/shared/jwt"
	"github.com/ishola-faazele/taskflow/pkg/utils"
)

type UserService struct {
	authRepo     AuthRepository
	profileRepo  UserProfileRepository
	jwtUtil      *jwt.JWTUtils
	emailService *utils.EmailService
}

func NewUserService(authRepo AuthRepository, profileRepo UserProfileRepository) *UserService {
	jwtUtil := jwt.NewJWTUtils(jwt.DefaultTokenConfig())
	emailConfig := utils.EmailConfig{
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    "587",
		SenderEmail: os.Getenv("SMTP_USER"),
		SenderName:  "TaskFlow Support",
		AppPassword: os.Getenv("SMTP_PASS"),
		FrontendURL: "http://localhost:3000",
	}
	emailService := utils.NewEmailService(emailConfig)
	return &UserService{
		authRepo:     authRepo,
		profileRepo:  profileRepo,
		jwtUtil:      jwtUtil,
		emailService: emailService,
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
		if _, err = us.authRepo.Create(new_user); err != nil {
			return err
		}
	}

	// create token using auth as claim
	authToken, token_err := us.jwtUtil.GenerateAuthToken(user.ID, email)
	if token_err != nil {
		return domain_errors.NewInternalError("FAILED TO GENERATE TOKEN", token_err)
	}
	// send email with magic link
	if err := us.emailService.SendMagicLink(email, authToken, "/api/user/verify?token="); err != nil {
		return domain_errors.NewInternalError("failed to send magic link email", err)
	}
	return nil
}

func (us *UserService) VerifyToken(token string) (string, string, domain_errors.DomainError) {
	// check if token is signed and is valid purpose(login)
	claims, parseErr := us.jwtUtil.ParseUserToken(token)
	if parseErr != nil {
		return "", "", domain_errors.NewInternalError("ERROR PARSING TOKEN", parseErr)
	}
	if claims.Purpose != jwt.PurposeAuth {
		return "", "", domain_errors.NewUnauthorizedError("INVALID TOKEN PURPOSE")
	}

	// issue new tokens for access and refresh
	access, refresh, tokenErr := us.jwtUtil.GenerateTokenPair(claims.UserID, claims.Email)
	if tokenErr != nil {
		return "", "", domain_errors.NewInternalError("FAILED TO GENERATE ACCESS AND REFRESH TOKENS", tokenErr)
	}
	return access, refresh, nil
}

func (us *UserService) RefreshToken(refreshToken string) (string, string, domain_errors.DomainError) {
	// validate refresh token
	claims, parseErr := us.jwtUtil.ParseUserToken(refreshToken)
	if parseErr != nil {
		return "", "", domain_errors.NewInternalError("ERROR PARSING REFRESH TOKEN", parseErr)
	}
	if claims.Purpose != jwt.PurposeRefresh {
		return "", "", domain_errors.NewUnauthorizedError("INVALID REFRESH TOKEN PURPOSE")
	}
	// check if token is valid
	token_hash := utils.HashToken(refreshToken)
	if valid, err := us.authRepo.IsTokenValid(token_hash); !valid || err != nil {
		return "", "", domain_errors.NewUnauthorizedError("TOKEN IS INVALID OR EXPIRED")
	}
	invalidToken := &InvalidToken{
		TokenHash:     token_hash,
		UserID:        claims.UserID,
		ExpiresAt:     claims.ExpiresAt.Time,
		InvalidatedAt: time.Now().UTC(),
	}
	if err := us.authRepo.InvalidateToken(invalidToken); err != nil {
		return "", "", domain_errors.NewInternalError("ERROR INVALIDATING REFRESH TOKEN", err)
	}
	access, refresh, token_err := us.jwtUtil.GenerateTokenPair(claims.UserID, claims.Email)
	if token_err != nil {
		return "", "", domain_errors.NewInternalError("FAILED TO GENERATE ACCESS AND REFRESH TOKENS", token_err)
	}
	return access, refresh, nil
}

func (us *UserService) GetByID(id string) (*Auth, domain_errors.DomainError) {
	return us.authRepo.GetByID(id)
}

func (us *UserService) GetByEmail(email string) (*Auth, domain_errors.DomainError) {
	// i need to validate email
	return us.authRepo.GetByEmail(email)
}

func (us *UserService) GetProfile(id string) (*UserProfile, domain_errors.DomainError) {
	return us.profileRepo.GetProfile(id)
}

func (us UserService) UpdateProfile(userID, name string) (*UserProfile, domain_errors.DomainError) {
	return us.profileRepo.UpdateProfile(userID, name)
}

func (us UserService) GetPublicProfile(id string) (*PublicProfile, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("id", id, "ID IS NOT A VALID UUID")
	}
	return us.profileRepo.GetPublicProfile(id)
}


