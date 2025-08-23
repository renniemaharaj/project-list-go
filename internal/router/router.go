package router

import (
	"net/http"

	cors "github.com/renniemaharaj/project-list-go/internal/middleware"
	"github.com/renniemaharaj/project-list-go/internal/router/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter() http.Handler {
	r := chi.NewRouter()

	// use middlewares
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.CORS) // CORS middleware

	// public routes
	r.Group(func(r chi.Router) {
		// public
		r.Get("/public", routes.HealthCheck)
	})

	// protected routes
	r.Group(func(r chi.Router) {
		// r.Use(auth.FirebaseAuth)
		r.Route("/projects", routes.Projects)
	})

	r.Group(func(r chi.Router) {
		// r.Use(auth.FirebaseAuth)
		r.Route("/dashboards", routes.Dashboard)
	})

	return r
}
