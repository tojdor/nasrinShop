package auth

import (
	"net/http"
	"os"
)

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wantUser := os.Getenv("ADMIN_USER")
		wantPass := os.Getenv("ADMIN_PASSWORD")

		user, pass, ok := r.BasicAuth()
		if !ok || wantUser == "" || wantPass == "" || user != wantUser || pass != wantPass {
			w.Header().Set("WWW-Authenticate", `Basic realm="NasrinShop Admin"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}