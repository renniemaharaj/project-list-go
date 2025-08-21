package auth

import (
	"context"
	"net/http"
)

type contextKey string

const UserKey = contextKey("firebaseUser")

var (
// l = logger.New().Prefix("Auth")
)

// FirebaseAuth middleware validates Firebase token in the `token` query param
func FirebaseAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			http.Error(w, "Missing token query parameter", http.StatusUnauthorized)
			return
		}

		//Authentication would be done before this
		next.ServeHTTP(w, r.WithContext(context.Background()))
	})
}
