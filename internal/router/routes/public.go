package routes

import "net/http"

func Public(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from public route!"))
}
