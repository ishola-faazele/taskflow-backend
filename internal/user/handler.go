package user

import (
	"database/sql"
	"encoding/json"
	"net/http"

	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ishola-faazele/taskflow/pkg/utils/domain_errors"
)

type UserHandler struct {
	service   *UserService
	responder *domain_errors.APIResponder
}

func NewUserHandler(db *sql.DB, conn *amqp.Connection) *UserHandler {
	postgresAuthRepo := NewPostgresAuthRepository(db)
	postgresProfileRepo := NewPostgresUserProfileRepository(db)
	service := NewUserService(postgresAuthRepo, postgresProfileRepo, conn)
	responder := domain_errors.NewAPIResponder()

	return &UserHandler{
		service:   service,
		responder: responder,
	}
}

type MagicLinkRequestDTO struct {
	Email string `json:"email"`
}

type UserProfileDTO struct {
	Name string `json:"name"`
}
type VerifyTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// Creates a magic link which the user has to verify to log in
func (h *UserHandler) RequestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req MagicLinkRequestDTO
	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "INVALID_REQUEST_BODY", err)
		return
	}
	// send magic link
	if err := h.service.GetMagicLink(req.Email); err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_TO_SEND_MAGIC_LINK", err)
		return
	}

	h.responder.NoContent(w)
}

// verifies token embedded in the magic link
func (h UserHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	// get token from query param
	token := r.URL.Query().Get("token")
	// verify token and return access and refresh tokens
	access, refresh, verifyErr := h.service.VerifyToken(token)
	if verifyErr != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "ERROR_VERIFYING_TOKEN", verifyErr)
		return
	}
	// encode access token
	data := VerifyTokenResponse{
		AccessToken: access,
	}
	// embed refresh token in http cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
	})
	h.responder.Success(w, r, http.StatusOK, "TOKEN_SUCCESSFULLY_VERIFIED", data)
}

// Returns new access and refresh tokens while invalidating the old refresh token
func (h UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		h.responder.Error(w, r, http.StatusUnauthorized, "REFRESH_TOKEN_NOT_FOUND_IN_COOKIE", err)
		return
	}
	// verify refresh token and return new access and refresh tokens
	access, refresh, token_err := h.service.RefreshToken(cookie.Value)
	if token_err != nil {
		h.responder.Error(w, r, http.StatusUnauthorized, "ERROR_REFRESHING_TOKEN", token_err)
		return
	}
	// encode access token
	data := VerifyTokenResponse{
		AccessToken: access,
	}
	// embed new refresh token in http cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
	})
	h.responder.Success(w, r, http.StatusOK, "TOKEN SUCCESSFULLY_REFRESHED", data)
}

// A route for users to get their own auth data
func (h UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "UNAUTHORIZED: USER_ID_NOT_FOUND_IN_CONTEXT", nil)
		return
	}
	user, err := h.service.GetByID(userID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_TO_GET_USER", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "USER_RETRIEVED_SUCCESSFULLY", user)
}

// func (h UserHandler) GetEmail(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	userID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
// 	if !ok || userID == "" {
// 		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
// 		return
// 	}
// 	user, err := h.service.GetByEmail(userID)
// 	if err != nil {
// 		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to Get User Email", err)
// 		return
// 	}
// 	h.responder.Success(w, r, http.StatusOK, "User Retrieved Successfully", user.Email)
// }

// A route for users to get their own profile data.
func (h UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "UNAUTHORIZED: USER_ID_NOT_FOUND_IN_CONTEXT", nil)
		return
	}
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_TO_GET_PROFILE", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "PROFILE_RETRIEVED_SUCCESSFULLY", profile)
}

// A route to update a user's profile
func (h UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "UNAUTHORIZED: USER_ID_NOT_FOUND_IN_CONTEXT", nil)
		return
	}

	var profile UserProfileDTO
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "INVALID_REQUEST_BODY", err)
		return
	}

	updatedProfile, err := h.service.UpdateProfile(userID, profile.Name)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_TO_UPDATE_PROFILE", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "PROFILE_UPDATED_SUCCESSFULLY", updatedProfile)
}

// A route to get a user's public profile
func (h UserHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	publicProfile, err := h.service.GetPublicProfile(userID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED_TO_GET_PUBLIC_PROFILE", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "PUBLIC_PROFILE_RETRIEVED_SUCCESSFULLY", publicProfile)
}
