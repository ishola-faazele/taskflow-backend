package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTUtils struct {
	SecretKey string
	Issuer    string
	Config    TokenConfig
}

type TokenPurpose string

const (
	PurposeAccess     TokenPurpose = "access"
	PurposeAuth       TokenPurpose = "auth"
	PurposeRefresh    TokenPurpose = "refresh"
	PurposeInvitation TokenPurpose = "invitation"
)

// TokenConfig holds expiration durations for different token purposes
type TokenConfig struct {
	AccessTokenDuration     time.Duration
	AuthTokenDuration       time.Duration
	RefreshTokenDuration    time.Duration
	InvitationTokenDuration time.Duration
}

// BaseClaims contains common fields for all token types
type BaseClaims struct {
	Purpose TokenPurpose `json:"purpose"`
	jwt.RegisteredClaims
}

// newBaseRegisteredClaims creates the standard JWT claims based on purpose
func (j *JWTUtils) newBaseRegisteredClaims(purpose TokenPurpose, subject string) jwt.RegisteredClaims {
	now := time.Now()
	var expiresAt time.Time

	// Set expiration based on purpose
	switch purpose {
	case PurposeAccess:
		expiresAt = now.Add(j.Config.AccessTokenDuration)
	case PurposeAuth:
		expiresAt = now.Add(j.Config.AuthTokenDuration)
	case PurposeRefresh:
		expiresAt = now.Add(j.Config.RefreshTokenDuration)
	case PurposeInvitation:
		expiresAt = now.Add(j.Config.InvitationTokenDuration)
	default:
		expiresAt = now.Add(j.Config.AccessTokenDuration)
	}

	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    j.Issuer,
		Subject:   subject,
		ID:        uuid.NewString(),
		Audience:  []string{string(purpose)},
	}
}

// UserClaims for authentication tokens
type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	BaseClaims
}

// InvitationClaims for invitation tokens
type InvitationClaims struct {
	InvitationID string `json:"invitation_id"`
	WorkspaceID  string `json:"workspace_id"`
	InviterID    string `json:"inviter_id"`
	InviteeEmail string `json:"invitee_email"`
	InviteeID    string `json:"invitee_id,omitempty"`
	Role         string `json:"role"`
	BaseClaims
}

// DefaultTokenConfig returns sensible defaults for token durations
func DefaultTokenConfig() TokenConfig {
	return TokenConfig{
		AccessTokenDuration:     15 * 4 * 24 * 365 * time.Minute, // Short-lived for security
		AuthTokenDuration:       15 * time.Minute,                // Medium-lived for authentication flows
		RefreshTokenDuration:    7 * 24 * time.Hour,              // Long-lived (7 days)
		InvitationTokenDuration: 24 * time.Hour,                  // Short-lived for invitation links
	}
}

// NewJWTUtils creates a new JWTUtils with the provided secret key and issuer
func NewJWTUtils(config TokenConfig) *JWTUtils {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	issuer := os.Getenv("JWT_ISSUER")
	if config == (TokenConfig{}) {
		config = DefaultTokenConfig()
	}
	return &JWTUtils{
		SecretKey: secretKey,
		Issuer:    issuer,
		Config:    config,
	}
}

// GenerateToken generates a token with any claims type
func (j *JWTUtils) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ParseToken parses a token with specific claims type
func ParseToken[T jwt.Claims](tokenStr string, secretKey string, claimsPtr T) (T, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claimsPtr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		var zero T
		return zero, err
	}

	if claims, ok := token.Claims.(T); ok && token.Valid {
		return claims, nil
	}
	var zero T
	return zero, jwt.ErrTokenInvalidClaims
}

// ParseUserToken parses a user authentication token
func (j *JWTUtils) ParseUserToken(tokenStr string) (*UserClaims, error) {
	return ParseToken(tokenStr, j.SecretKey, &UserClaims{})
}

// ParseInvitationToken parses an invitation token
func (j *JWTUtils) ParseInvitationToken(tokenStr string) (*InvitationClaims, error) {
	return ParseToken(tokenStr, j.SecretKey, &InvitationClaims{})
}

// NewUserClaims creates user claims with purpose-based expiration
func (j *JWTUtils) NewUserClaims(userID, email string, purpose TokenPurpose) *UserClaims {
	return &UserClaims{
		UserID: userID,
		Email:  email,
		BaseClaims: BaseClaims{
			Purpose:          purpose,
			RegisteredClaims: j.newBaseRegisteredClaims(purpose, userID),
		},
	}
}

// NewInvitationClaims creates invitation claims
func (j *JWTUtils) NewInvitationClaims(invitationID, workspaceID, inviterID, inviteeEmail, inviteeID, role string) *InvitationClaims {
	return &InvitationClaims{
		InvitationID: invitationID,
		WorkspaceID:  workspaceID,
		InviterID:    inviterID,
		InviteeEmail: inviteeEmail,
		InviteeID:    inviteeID,
		Role:         role,
		BaseClaims: BaseClaims{
			Purpose:          PurposeInvitation,
			RegisteredClaims: j.newBaseRegisteredClaims(PurposeInvitation, inviteeEmail),
		},
	}
}

// IsExpired checks if the token is expired
func (b *BaseClaims) IsExpired() bool {
	if b.ExpiresAt == nil {
		return true
	}
	return b.ExpiresAt.Before(time.Now())
}

// IsValid checks if the token is valid (not expired and not before NotBefore time)
func (b *BaseClaims) IsValid() bool {
	now := time.Now()

	if b.ExpiresAt != nil && b.ExpiresAt.Before(now) {
		return false
	}

	if b.NotBefore != nil && b.NotBefore.After(now) {
		return false
	}

	return true
}

// TimeUntilExpiry returns the duration until token expiration
func (b *BaseClaims) TimeUntilExpiry() time.Duration {
	if b.ExpiresAt == nil {
		return 0
	}
	return time.Until(b.ExpiresAt.Time)
}

// Helper methods for JWTUtils to simplify token generation

// GenerateAccessToken creates an access token
func (j *JWTUtils) GenerateAccessToken(userID, email string) (string, error) {
	claims := j.NewUserClaims(userID, email, PurposeAccess)
	return j.GenerateToken(claims)
}

// GenerateAuthToken creates an authentication token
func (j *JWTUtils) GenerateAuthToken(userID, email string) (string, error) {
	claims := j.NewUserClaims(userID, email, PurposeAuth)
	return j.GenerateToken(claims)
}

// GenerateRefreshToken creates a refresh token
func (j *JWTUtils) GenerateRefreshToken(userID, email string) (string, error) {
	claims := j.NewUserClaims(userID, email, PurposeRefresh)
	return j.GenerateToken(claims)
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTUtils) GenerateTokenPair(userID, email string) (accessToken, refreshToken string, err error) {
	accessToken, err = j.GenerateAccessToken(userID, email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = j.GenerateRefreshToken(userID, email)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// GenerateInvitationToken creates an invitation token
func (j *JWTUtils) GenerateInvitationToken(invitationID, workspaceID, inviterID, inviteeEmail, inviteeID, role string) (string, error) {
	claims := j.NewInvitationClaims(invitationID, workspaceID, inviterID, inviteeEmail, inviteeID, role)
	return j.GenerateToken(claims)
}
