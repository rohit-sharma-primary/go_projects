package middleware

import (
	"net/http"
	"strings"

	"server/internal/utils"
)

const SecretToken = "afifdosa"

func HandleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		token, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || strings.TrimSpace(token) == "" || token != SecretToken {
			utils.WriteError(w, http.StatusUnauthorized, "missing or invalid bearer token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
