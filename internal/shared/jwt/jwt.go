package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTUtils struct {
	SecretKey string
	Issuer    string
	Config    TokenConfig
}

// TokenConfig holds expiration durations for different token purposes
type TokenConfig struct {
	AccessTokenDuration  time.Duration
	AuthTokenDuration    time.Duration
	RefreshTokenDuration time.Duration
}

// DefaultTokenConfig returns sensible defaults for token durations
func DefaultTokenConfig() TokenConfig {
	return TokenConfig{
		AccessTokenDuration:  15 * time.Minute,   // Short-lived for security
		AuthTokenDuration:    1 * time.Hour,      // Medium-lived for authentication flows
		RefreshTokenDuration: 7 * 24 * time.Hour, // Long-lived (7 days)
	}
}

// NewJWTUtils creates a new JWTUtils with the provided secret key and issuer
func NewJWTUtils(secretKey, issuer string, config TokenConfig) *JWTUtils {
	if config == (TokenConfig{}) { // Check if config is zero-valued
		config = DefaultTokenConfig()
	}
	return &JWTUtils{
		SecretKey: secretKey,
		Issuer:    issuer,
		Config:    config,
	}
}

type TokenPurpose string

const (
	PurposeAccess  TokenPurpose = "access"
	PurposeAuth    TokenPurpose = "auth"
	PurposeRefresh TokenPurpose = "refresh"
)

type Claims struct {
	UserID  string       `json:"user_id"`
	Email   string       `json:"email"`
	Purpose TokenPurpose `json:"purpose"`
	jwt.RegisteredClaims
}

func (j *JWTUtils) GenerateToken(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *JWTUtils) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// NewClaims creates a new Claims with purpose-based expiration and metadata
func NewClaims(userID, email, issuer string, purpose TokenPurpose, config TokenConfig) *Claims {
	now := time.Now()
	var expiresAt time.Time

	// Set expiration based on purpose
	switch purpose {
	case PurposeAccess:
		expiresAt = now.Add(config.AccessTokenDuration)
	case PurposeAuth:
		expiresAt = now.Add(config.AuthTokenDuration)
	case PurposeRefresh:
		expiresAt = now.Add(config.RefreshTokenDuration)
	default:
		// Default to access token duration
		expiresAt = now.Add(config.AccessTokenDuration)
	}

	return &Claims{
		UserID:  userID,
		Email:   email,
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    issuer,
			Subject:   userID,
			ID:        uuid.NewString(),
			Audience:  []string{string(purpose)},
		},
	}
}

// NewClaimsWithDefaults creates Claims using DefaultTokenConfig
func NewClaimsWithDefaults(userID, email, issuer string, purpose TokenPurpose) *Claims {
	return NewClaims(userID, email, issuer, purpose, DefaultTokenConfig())
}

// IsExpired checks if the token is expired
func (c *Claims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return true
	}
	return c.ExpiresAt.Before(time.Now())
}

// IsValid checks if the token is valid (not expired and not before NotBefore time)
func (c *Claims) IsValid() bool {
	now := time.Now()

	if c.ExpiresAt != nil && c.ExpiresAt.Before(now) {
		return false
	}

	if c.NotBefore != nil && c.NotBefore.After(now) {
		return false
	}

	return true
}

// TimeUntilExpiry returns the duration until token expiration
func (c *Claims) TimeUntilExpiry() time.Duration {
	if c.ExpiresAt == nil {
		return 0
	}
	return time.Until(c.ExpiresAt.Time)
}

// Helper methods for JWTUtils to simplify token generation

// GenerateAccessToken creates an access token
func (j *JWTUtils) GenerateAccessToken(userID, email string) (string, error) {
	claims := NewClaimsWithDefaults(userID, email, j.Issuer, PurposeAccess)
	return j.GenerateToken(claims)
}

// GenerateAuthToken creates an authentication token
func (j *JWTUtils) GenerateAuthToken(userID, email string) (string, error) {
	claims := NewClaimsWithDefaults(userID, email, j.Issuer, PurposeAuth)
	return j.GenerateToken(claims)
}

// GenerateRefreshToken creates a refresh token
func (j *JWTUtils) GenerateRefreshToken(userID, email string) (string, error) {
	claims := NewClaimsWithDefaults(userID, email, j.Issuer, PurposeRefresh)
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
