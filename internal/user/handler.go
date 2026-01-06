package user

import (
	"database/sql"
	"encoding/json"
	"net/http"

	domain_middleware "github.com/ishola-faazele/taskflow/internal/middleware"

	"github.com/ishola-faazele/taskflow/pkg/utils"
)

type UserHandler struct {
	service   *UserService
	responder *utils.APIResponder
}

func NewUserHandler(db *sql.DB) *UserHandler {
	postgresAuthRepo := NewPostgresAuthRepository(db)
	postgresProfileRepo := NewPostgresUserProfileRepository(db)
	service := NewUserService(postgresAuthRepo, postgresProfileRepo)
	responder := utils.NewAPIResponder()

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

func (h *UserHandler) RequestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req MagicLinkRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "INVALID REQUEST BODY", err)
		return
	}
	if err := h.service.GetMagicLink(req.Email); err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "FAILED TO SEND MAGIC LINK", err)
		return
	}

	h.responder.NoContent(w)
}

func (h UserHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	access, refresh, err := h.service.VerifyToken(token)
	if err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Error Verifying Token", err)
		return
	}
	data := map[string]string{
		"access_token": access,
	}
	// embed refresh token in http cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
	})
	h.responder.Success(w, r, http.StatusOK, "Token Successfully Verified", data)
}

func (h UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		h.responder.Error(w, r, http.StatusUnauthorized, "Refresh Token Missing", err)
		return
	}
	access, refresh, token_err := h.service.RefreshToken(cookie.Value)
	if token_err != nil {
		h.responder.Error(w, r, http.StatusUnauthorized, "Error Refreshing Token", token_err)
		return
	}
	data := map[string]string{
		"access_token": access,
	}
	// embed new refresh token in http cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
	})
	h.responder.Success(w, r, http.StatusOK, "Token Successfully Refreshed", data)
}

func (h UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}
	user, err := h.service.GetByID(userID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to Get User", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "User Retrieved Successfully", user)
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

func (h UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to Get Profile", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Profile Retrieved Successfully", profile)
}

func (h UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(domain_middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
		return
	}

	var profile UserProfileDTO
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	updatedProfile, err := h.service.UpdateProfile(userID, profile.Name)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to Update Profile", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Profile Updated Successfully", updatedProfile)
}

func (h UserHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	publicProfile, err := h.service.GetPublicProfile(userID)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to Get Public Profile", err)
		return
	}
	h.responder.Success(w, r, http.StatusOK, "Public Profile Retrieved Successfully", publicProfile)
}
