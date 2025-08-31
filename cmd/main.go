package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/renniemaharaj/project-list-go/internal/cache"

	"github.com/renniemaharaj/project-list-go/internal/dashboard"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/demo"
	routes "github.com/renniemaharaj/project-list-go/internal/health"
	"github.com/renniemaharaj/project-list-go/internal/meta"
	cors "github.com/renniemaharaj/project-list-go/internal/middleware"
	"github.com/renniemaharaj/project-list-go/internal/project"
	"github.com/renniemaharaj/project-list-go/internal/schema"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

func main() {
	// main logger used by main function
	mainLogger := logger.New().Prefix("Backend")

	mainLogger.Info("Allowing time for postgres to initialize...")
	// allow time for postgres to initialize
	time.Sleep(5 * time.Second)

	mainLogger.Info("Resolving automatic database profile")
	// initialize automatic database profile
	_, err := database.Automatic.Resolve()
	if err != nil {
		panic(err)

	}

	mainLogger.Info("")
	// will automatically initialize tables
	if err := schema.NewRepository(database.Automatic, mainLogger).InitializeDatabaseTables(context.Background()); err != nil {
		panic(err)
	}

	// initialize redis
	if err := cache.InitializeRedis(); err != nil {
		panic(err)
	}

	demoData := true
	// seed demo data
	if demoData {
		err = demo.NewRepository(database.Automatic, mainLogger).GenerateInsertDemoData(context.Background())
		if err != nil {
			logger.New().Fatal(err)
		}
	}
	// setup chi router and start server
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

	// private routes
	r.Group(func(r chi.Router) {
		r.Route("/meta", meta.Meta)
		r.Route("/project", project.ProjectHandler)
		r.Route("/dashboard", dashboard.Dashboard)
	})

	// start rest server
	go func() {
		mainLogger.Info("Starting server on :8081")
		server := &http.Server{
			Addr:    ":8081",
			Handler: r, // chi router as handler
		}
		if err := server.ListenAndServe(); err != nil {
			mainLogger.Fatal(err)
		}
	}()

	select {}
}
