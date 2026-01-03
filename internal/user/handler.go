package user

import (
	"database/sql"
	"encoding/json"
	"net/http"

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

func (h *UserHandler) RequestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req MagicLinkRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.Error(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	err := h.service.GetMagicLink(req.Email)
	if err != nil {
		h.responder.Error(w, r, http.StatusInternalServerError, "Failed to Send Magic Link", err)
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
