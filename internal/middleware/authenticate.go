package middleware

import (
	"context"
	"net/http"

	"github.com/ishola-faazele/taskflow/internal/shared/jwt"
)

type contextKey string

const UserIDKey contextKey = "UserID"
const UserEmailKey contextKey = "UserEmail"

func (dm *DomainMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from Authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			dm.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: No token provided", nil)
			return
		}
		// strip "Bearer " prefix if present

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// parse token
		claims, err := dm.jwt.ParseToken(token)
		if err != nil {
			dm.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: Invalid token", nil)
			return
		}
		// check claims validity
		if !claims.IsValid() {
			dm.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: Token expired or invalid", nil)
			return
		}

		// check claims purpose
		if claims.Purpose != jwt.PurposeAccess {
			dm.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: Invalid token purpose", nil)
			return
		}
		// set user ID in context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		r = r.WithContext(ctx)
		// proceed to next handler
		next.ServeHTTP(w, r)
	})
}
