package main

import (
	"context"
	"net/http"
	"time"

	"github.com/renniemaharaj/project-list-go/internal/cache"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/demo"
	"github.com/renniemaharaj/project-list-go/internal/router"
	"github.com/renniemaharaj/project-list-go/internal/schema"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

func main() {
	mainLogger := logger.New().Prefix("Backend")

	// allow time for postgres to initialize
	time.Sleep(5 * time.Second)
	_, err := database.Automatic.Get()
	if err != nil {
		panic(err)

	}

	// will automatically initialize tables
	if err := schema.NewRepository(database.Automatic, mainLogger).InitializeDatabaseTables(context.Background()); err != nil {
		panic(err)
	}

	if err := cache.InitializeRedis(); err != nil {
		panic(err)
	}

	demoData := true
	if demoData {
		err = demo.NewRepository(database.Automatic, mainLogger).InsertSeededDemoData(context.Background())
		if err != nil {
			logger.New().Fatal(err)
		}
	}
	// setup chi router and start server
	r := router.SetupRouter()
	go func() {
		mainLogger.Info("Starting server on :8081")
		if err := http.ListenAndServe(":8081", r); err != nil {
			mainLogger.Fatal(err)
		}
	}()

	select {}
}
