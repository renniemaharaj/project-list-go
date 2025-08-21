package router

import (
	"net/http"

	"github.com/renniemaharaj/project-list-go/internal/auth"
	cors "github.com/renniemaharaj/project-list-go/internal/middleware"
	handlers "github.com/renniemaharaj/project-list-go/internal/router/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter() http.Handler {
	r := chi.NewRouter()

	// use middlewares
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.CORS)         // CORS middleware
	r.Use(auth.FirebaseAuth) // Auth middleware

	// public routes
	r.Get("/public", handlers.Public)

	// protected routes
	r.Get("/protected", handlers.Protected)

	return r
}
