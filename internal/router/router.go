package router

import (
	"net/http"

	"github.com/renniemaharaj/project-list-go/internal/dashboard"
	"github.com/renniemaharaj/project-list-go/internal/meta"
	cors "github.com/renniemaharaj/project-list-go/internal/middleware"
	"github.com/renniemaharaj/project-list-go/internal/project"
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
		r.Route("/project", project.ProjectHandler)
	})

	r.Group(func(r chi.Router) {
		r.Route("/meta", meta.Meta)
	})

	r.Group(func(r chi.Router) {
		// r.Use(auth.FirebaseAuth)
		r.Route("/dashboard", dashboard.Dashboard)
	})

	return r
}
