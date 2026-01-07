package middleware

import (
	"context"
	"net/http"
)

const WorkspaceIDKey contextKey = "WorkspaceID"

func (dm *DomainMiddleware) CheckMembership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsID := r.PathValue("ws_id")
		if wsID == "" {
			dm.responder.Error(w, r, http.StatusBadRequest, "WORKSPACE ID MUST BE PROVIDED", nil)
			return
		}
		ctx := r.Context()
		requester, ok := ctx.Value(UserIDKey).(string)
		if !ok || requester == "" {
			dm.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: User ID not found in context", nil)
			return
		}
		if isMember, err := dm.WorkspaceService.IsMember(requester, wsID); !isMember || err != nil {
			dm.responder.Error(w, r, http.StatusUnauthorized, "Unauthorized: USER IS NOT A MEMBER OF WORKSPACE", err)
			return
		}
		ctx = context.WithValue(r.Context(), WorkspaceIDKey, wsID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
