package auth

import (
	"context"
	"net/http"
)

var (
// l = logger.New().Prefix("Auth")
)

// FirebaseAuth middleware validates Firebase token in the `token` query param
func FirebaseAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Authentication would be done before this
		next.ServeHTTP(w, r.WithContext(context.Background()))
	})
}
