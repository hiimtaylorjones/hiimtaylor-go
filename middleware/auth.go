package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager

func SetSessionManager(sm *scs.SessionManager) {
	sessionManager = sm
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminID := sessionManager.GetString(r.Context(), "admin_id")
		if adminID == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}