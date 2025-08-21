package routes

import (
	"net/http"
	// "backend/internal/auth"
	// "backend/internal/gorilla"
	// fbAuth "firebase.google.com/go/v4/auth"
)

func Protected(w http.ResponseWriter, r *http.Request) {
	// user := r.Context().Value(auth.UserKey).(*fbAuth.Token)
	// w.Write([]byte("Hello, " + user.UID))

	// upgrade to web socket
	// gorilla.UpgradeHandler(w, r)
}
